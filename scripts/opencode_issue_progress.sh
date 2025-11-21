#!/usr/bin/env bash
set -euo pipefail

##
# 针对“已绑定代码分支的 Plane Issue”生成每日需求进度 + Code Review 汇总，存储到报告仓库并推送飞书。
# 流程：
# 1) 调用后端 /jobs/issue-progress/tasks 获取任务列表（分支、标题、描述、Lark Chat）。
# 2) 对每个分支提取“今天”的提交上下文，调用 opencode 生成结构化 JSON（遵循模板）。
# 3) 将 JSON 写入 report 仓库 (默认 https://cnb.cool/1024hub/plane-test) 的 issue-progress/daily/{date}/issue-{issue_id}.json。
# 4) 调用后端 /jobs/issue-progress/send 将结果推送到对应 Lark 群（即时预览）。
# 5) 对无分支但已绑定群聊的 Issue 发送提醒。
#
# 依赖：
# - 环境变量：CABB_API_BASE（后端地址，如 https://cabb.onrender.com）、INTEGRATION_TOKEN（Bearer）
# - 报告仓库：PROGRESS_REPO_URL/PROGRESS_BRANCH/PROGRESS_DIR，需 CNB_TOKEN 可推送。
# - 可选：TZ，默认 Asia/Shanghai。需要 git 可读完整历史（会尝试 fetch）。
# - opencode CLI（脚本会尝试安装）。
##

export TZ=${TZ:-Asia/Shanghai}
api_base="${CABB_API_BASE:-}"
integration_token="${INTEGRATION_TOKEN:-}"
model="${OPENCODE_MODEL:-opencode/grok-code}"
api_key="${OPENCODE_API_KEY:-${OPENCODE_TOKEN:-${OPENCODE_KEY:-${ZEN_API_KEY:-${OC_API_KEY:-${OPENAI_API_KEY:-""}}}}}}"
repo_root="$(pwd)"
publish_enable="${PROGRESS_PUBLISH_ENABLE:-1}"
publish_repo_url="${PROGRESS_REPO_URL:-https://cnb.cool/1024hub/plane-test}"
publish_branch="${PROGRESS_BRANCH:-main}"
publish_dir_root="${PROGRESS_DIR:-issue-progress}"
template_file="scripts/issue_progress_template.json"
target_repo_url="${TARGET_REPO_URL:-}"
target_repo_branch="${TARGET_REPO_BRANCH:-main}"
target_repo_path="${repo_root}"

if [ -z "${api_base}" ] || [ -z "${integration_token}" ]; then
  echo "[error] CABB_API_BASE or INTEGRATION_TOKEN is missing; cannot proceed." >&2
  exit 1
fi

command -v curl >/dev/null 2>&1 || { echo "[error] curl is required"; exit 1; }
command -v jq >/dev/null 2>&1 || { echo "[error] jq is required"; exit 1; }

# Ensure git history available (for current repo)
git fetch --all --prune --tags >/dev/null 2>&1 || true

# Optional: clone target repo for analysis
if [ -n "${target_repo_url}" ]; then
  target_repo_path="tmp/target-repo"
  rm -rf "${target_repo_path}"
  mkdir -p "${target_repo_path}"
  auth_hdr=""
  if [ -n "${CNB_TOKEN:-}" ]; then
    auth_hdr="Authorization: Basic $(printf "cnb:%s" "${CNB_TOKEN}" | base64 | tr -d '\n')"
  fi
  if [ -n "${auth_hdr}" ]; then
    GIT_CURL_VERBOSE=0 git -c http.extraHeader="${auth_hdr}" clone "${target_repo_url}" "${target_repo_path}" >/dev/null 2>&1 || true
  else
    git clone "${target_repo_url}" "${target_repo_path}" >/dev/null 2>&1 || true
  fi
  if [ ! -d "${target_repo_path}/.git" ]; then
    echo "[error] failed to clone target repo ${target_repo_url}" >&2
    exit 1
  fi
  git -C "${target_repo_path}" fetch --all --prune --tags >/dev/null 2>&1 || true
  git -C "${target_repo_path}" checkout "${target_repo_branch}" >/dev/null 2>&1 || true
  git -C "${target_repo_path}" config --global --add safe.directory "$(cd "${target_repo_path}" && pwd)" || true
fi

git_ctx_repo="${target_repo_path}"

start_ts=$(date -d "today 00:00:00" +'%Y-%m-%dT00:00:00')
end_ts=$(date +'%Y-%m-%dT%H:%M:%S')
label=$(date +%Y-%m-%d)

api_get() {
  local path="$1"
  curl -sfSL -H "Authorization: Bearer ${integration_token}" "${api_base%/}${path}"
}

api_post_json() {
  local path="$1"
  curl -sfSL -H "Authorization: Bearer ${integration_token}" -H "Content-Type: application/json" -d "$2" "${api_base%/}${path}"
}

tasks_json=$(api_get "/jobs/issue-progress/tasks") || {
  echo "[error] failed to fetch issue progress tasks" >&2
  exit 1
}

branch_links_count=$(echo "${tasks_json}" | jq '.branch_links | length')
reminders_count=$(echo "${tasks_json}" | jq '.unbound_chats | length')

echo "[info] branch links: ${branch_links_count}, reminders: ${reminders_count}"

workdir="$(mktemp -d)"
out_root="reports/issue-progress/daily/${label}"
mkdir -p "${out_root}"

if [ "${branch_links_count}" -gt 0 ]; then
  echo "${tasks_json}" | jq -c '.branch_links[]' | while read -r row; do
    issue_id=$(echo "${row}" | jq -r '.plane_issue_id')
    repo=$(echo "${row}" | jq -r '.cnb_repo_id')
    branch=$(echo "${row}" | jq -r '.branch')
    title=$(echo "${row}" | jq -r '.issue_title')
    desc=$(echo "${row}" | jq -r '.issue_description')
    chat_id=$(echo "${row}" | jq -r '.lark_chat_id')

    echo "[info] processing issue=${issue_id} repo=${repo} branch=${branch}"

    out_file_local="${out_root}/issue-${issue_id}.json"
    mkdir -p "$(dirname "${out_file_local}")"
    if [ -f "${template_file}" ]; then
      cp "${template_file}" "${out_file_local}"
    else
      echo '{"date":"","issue_id":"","issue_title":"","progress_summary":{"overview":"","details":[]},"code_review_summary":{"overview":"","details":[]}}' > "${out_file_local}"
    fi

    # Ensure branch exists locally
    if ! git -C "${git_ctx_repo}" rev-parse --verify --quiet "${branch}" >/dev/null; then
      echo "[warn] branch ${branch} missing locally; fetching from origin..." >&2
      git -C "${git_ctx_repo}" fetch origin "${branch}:${branch}" >/dev/null 2>&1 || true
    fi

    ctx_file="${workdir}/ctx-${issue_id}.md"
    {
      echo "# Issue Progress Context"
      echo "- Issue: ${title}"
      echo "- Branch: ${branch}"
      echo "- Repo: ${repo}"
      echo "- Window: ${start_ts} .. ${end_ts} (${TZ})"
      echo
      echo "## Commits (today)"
      git -C "${git_ctx_repo}" log "${branch}" --since="${start_ts}" --until="${end_ts}" --date=iso-local \
        --numstat --pretty=format:'---%ncommit %H%nauthor %an <%ae>%ndate %ad%ntitle %s%nbody %b' \
        | sed -e 's/(Bearer)[[:space:]]+[A-Za-z0-9._-]\+/\1 ***REDACTED***/Ig' \
        | head -n 4000 || true
    } > "${ctx_file}"

    read -r -d '' prompt <<'EOF' || true
你是研发 TL，需基于“指定分支今日提交”生成中文汇总，输出 JSON（遵循模板字段），不要输出 markdown 代码块或额外文本。
模板字段：
{
  "date": "YYYY-MM-DD",
  "issue_id": "uuid",
  "issue_title": "文本",
  "progress_summary": {
    "overview": "面向非研发的进度概要",
    "details": [ { "topic": "模块/需求点", "content": "进展/阻塞/下一步" } ]
  },
  "code_review_summary": {
    "overview": "整体代码质量/风险",
    "details": [ { "author": "作者", "changes": "主要变更", "suggestions": "建议/风险" } ]
  }
}
规则：
- 若无提交，progress_summary.overview 写明“今日无提交，暂无进展”，code_review_summary 概要同样说明。
- 语言简洁、信息密度高，细节 3-6 条为宜。
- 严格输出 JSON（无多余文本），字段名与类型保持一致。
EOF

    out_tmp="${workdir}/out-${issue_id}.json"
    export OPENCODE_MODEL="${model}"
    [ -n "${api_key}" ] && export OPENCODE_API_KEY="${api_key}"
    if opencode run --format json "${prompt}" -m "${OPENCODE_MODEL}" -f "${ctx_file}" > "${out_tmp}" 2>"${workdir}/stderr-${issue_id}.log"; then
      if jq . "${out_tmp}" >/dev/null 2>&1; then
        jq --arg date "${label}" --arg id "${issue_id}" --arg title "${title}" '.date=$date | .issue_id=$id | .issue_title=$title' "${out_tmp}" > "${out_tmp}.tmp" && mv "${out_tmp}.tmp" "${out_tmp}"
        cp "${out_tmp}" "${out_file_local}"
        progress_text=$(cat "${out_tmp}" | jq -r '"需求进度: "+.progress_summary.overview+"\nCode Review: "+.code_review_summary.overview')
      else
        echo "[warn] opencode output is not valid JSON for issue ${issue_id}" >&2
        progress_text="今日未生成 AI 汇总（JSON 校验失败）。"
      fi
    else
      echo "[warn] opencode failed for issue ${issue_id}, falling back to no-progress note" >&2
      progress_text="今日未生成 AI 汇总（调用失败或无提交）。"
    fi

    payload=$(jq -n --arg chat "${chat_id}" --arg txt "${progress_text}" --arg date "${label}" --arg issue "${title}" '{
      chat_id: $chat,
      date: $date,
      issue_title: $issue,
      message: $txt
    }')
    api_post_json "/jobs/issue-progress/send" "${payload}" || echo "[warn] send failed for chat ${chat_id}" >&2

    # Queue for publishing
    echo "${repo_root}/${out_file_local}" >> "${workdir}/publish-list.txt"
  done
fi

if [ "${reminders_count}" -gt 0 ]; then
  echo "${tasks_json}" | jq -c '.unbound_chats[]' | while read -r row; do
    chat_id=$(echo "${row}" | jq -r '.lark_chat_id')
    issue_id=$(echo "${row}" | jq -r '.plane_issue_id')
    reminder="⚠️ 每日进度提醒：Issue ${issue_id} 尚未绑定代码仓库/分支，暂无法生成开发进度与 Code Review 汇报。请先在 Cabb 绑定对应分支。"
    payload=$(jq -n --arg chat "${chat_id}" --arg txt "${reminder}" --arg date "${label}" '{
      chat_id: $chat,
      date: $date,
      issue_title: "",
      message: $txt
    }')
    api_post_json "/jobs/issue-progress/send" "${payload}" || echo "[warn] reminder send failed chat=${chat_id}" >&2
  done
fi

publish_reports() {
  if [ "${publish_enable}" = "0" ] || [ "${publish_enable}" = "false" ]; then
    return 0
  fi
  if [ -z "${publish_repo_url}" ] || [ -z "${CNB_TOKEN:-}" ]; then
    echo "[warn] publish skipped: missing repo url or CNB_TOKEN" >&2
    return 0
  }
  if [ ! -s "${workdir}/publish-list.txt" ]; then
    echo "[info] nothing to publish" >&2
    return 0
  fi
  workdir_pub="tmp/issue-progress-publish"
  rm -rf "${workdir_pub}"
  mkdir -p "${workdir_pub}"
  auth_hdr="Authorization: Basic $(printf "cnb:%s" "${CNB_TOKEN}" | base64 | tr -d '\n')"
  if ! GIT_CURL_VERBOSE=0 git -c http.extraHeader="${auth_hdr}" clone "${publish_repo_url}" "${workdir_pub}" >/dev/null 2>&1; then
    echo "[warn] clone publish repo failed" >&2
    return 0
  fi
  cd "${workdir_pub}"
  git config --global --add safe.directory "$(pwd)" || true
  if git show-ref --verify --quiet "refs/heads/${publish_branch}"; then
    git checkout "${publish_branch}" >/dev/null 2>&1 || true
  else
    git checkout -b "${publish_branch}" >/dev/null 2>&1 || true
  fi
  git -c http.extraHeader="${auth_hdr}" pull --rebase origin "${publish_branch}" >/dev/null 2>&1 || true
  target_dir="${publish_dir_root}/daily/${label}"
  mkdir -p "${target_dir}"
  while read -r f; do
    [ -f "${f}" ] && cp -f "${f}" "${target_dir}/"
  done < "${workdir}/publish-list.txt"
  git add "${target_dir}" || true
  if git diff --cached --quiet >/dev/null 2>&1; then
    cd - >/dev/null 2>&1 || true
    return 0
  fi
  git config user.name "cabb-issue-progress-bot"
  git config user.email "bot@cabb.local"
  git commit -m "chore(issue-progress): ${label}" >/dev/null 2>&1 || true
  GIT_CURL_VERBOSE=0 git -c http.extraHeader="${auth_hdr}" push origin "${publish_branch}" >/dev/null 2>&1 || true
  cd - >/dev/null 2>&1 || true
}

publish_reports

echo "[info] issue progress job completed"
