#!/usr/bin/env bash
set -euo pipefail

##
# 项目工作汇报生成脚本（日报/周报/月报）
#
# 作用概览
# - 跨分支聚合：统计指定时间窗（昨天/上一周/上个月）内“所有分支”的提交。
# - AI 汇总优先：基于提交统计 + 采样补丁的上下文，让 opencode 生成中文结构化“按人员总结”的汇报。
# - 幂等发布：按固定命名将汇报推送到目标仓库（同一时间窗覆盖旧文件），供团队留存与查阅。
# - 控制台预览：流水线末尾输出最新汇报的前 80 行，便于快速确认质量。
#
# 设计要点
# - 可观测而不过度：不在日志泄露密钥；补丁上下文会做基础脱敏（Authorization/Bearer/sk-/CNB_TOKEN）。
# - 稳健性：CI 若为浅克隆或无 .git，会自动克隆只读镜像用于统计；git safe.directory 避免权限告警。
# - 一致性：手动/定时共享同一实现与参数，差异仅在时间窗默认值（按钮可选）。
##

export TZ=${TZ:-Asia/Shanghai}

timeframe="${REPORT_TIMEFRAME:-yesterday}"
zen_mode="${ZEN_MODE:-1}"
# prefer fully-qualified model id per docs (opencode/<model-id>)
model="${OPENCODE_MODEL:-opencode/grok-code}"
# require AI summary or fallback allowed (manual可设置 REQUIRE_AI_SUMMARY=1)
require_ai="${REQUIRE_AI_SUMMARY:-0}"
# publish settings
publish_enable="${REPORT_PUBLISH_ENABLE:-1}"
publish_repo_url="${REPORT_PUBLISH_REPO_URL:-https://cnb.cool/1024hub/plane-test}"
publish_branch="${REPORT_PUBLISH_BRANCH:-main}"
publish_dir_root="${REPORT_PUBLISH_DIR:-reports}"
template_file="scripts/template.json"
# diff context controls (manual trigger can override via web_trigger inputs)
diff_mode="${REPORT_DIFF_MODE:-stats}" # stats | sampled_patch | full_patch
char_budget="${REPORT_DIFF_CHAR_BUDGET:-200000}"
top_files="${REPORT_DIFF_TOP_FILES:-20}"
max_hunks_per_file="${REPORT_DIFF_HUNKS_PER_FILE:-3}"
max_commits="${REPORT_DIFF_MAX_COMMITS:-100}"
exclude_globs="${REPORT_EXCLUDE_GLOBS:-node_modules/**;dist/**;vendor/**;*.min.*;*.png;*.jpg;*.jpeg;*.gif;*.webp;*.svg;*.mp4;*.mov;*.zip;*.tar;*.gz;*.lock;*.class;*.bin;*.exe;*.dylib;*.so;*.dll}"
numstat_lines="${REPORT_NUMSTAT_LINES:-5000}"
report_config_from_api="${REPORT_CONFIG_FROM_API:-1}"
report_repo_list="${REPORT_REPO_LIST:-}"
cabb_api_base="${CABB_API_BASE:-}"
integration_token="${INTEGRATION_TOKEN:-}"
# API Key 检测：支持多种变量名，最终导出为 OPENCODE_API_KEY（opencode CLI 将读取这一变量）。
api_key="${OPENCODE_API_KEY:-}"
if [ -z "${api_key}" ]; then
  api_key="${OPENCODE_TOKEN:-}"
fi
if [ -z "${api_key}" ]; then
  api_key="${OPENCODE_KEY:-}"
fi
if [ -z "${api_key}" ]; then
  api_key="${OPENCODE_ZEN_API_KEY:-}"
fi
if [ -z "${api_key}" ]; then
  api_key="${ZEN_API_KEY:-}"
fi
if [ -z "${api_key}" ]; then
  api_key="${OC_API_KEY:-}"
fi
if [ -z "${api_key}" ]; then
  api_key="${OPENAI_API_KEY:-}"
fi
if [ -z "${api_key}" ]; then
  api_key="${ANTHROPIC_API_KEY:-}"
fi
if [ -z "${api_key}" ]; then
  api_key="${OPENROUTER_API_KEY:-}"
fi

# Determine time window and label
report_type="daily"
case "${timeframe}" in
  yesterday)
    start_ts=$(date -d "yesterday 00:00:00" +'%Y-%m-%d 00:00:00')
    end_ts=$(date -d "today 00:00:00" +'%Y-%m-%d 00:00:00')
    label=$(date -d "yesterday" +%Y-%m-%d)
    report_type="daily"
    ;;
  last_week)
    # 上一自然周：本周周一 00:00:00 往前推 7 天为上一周周一 00:00:00
    dow=$(date +%u) # 1..7 (Mon=1)
    this_mon_date=$(date -d "today -$((dow-1)) days" +%Y-%m-%d)
    last_mon_date=$(date -d "${this_mon_date} -7 days" +%Y-%m-%d)
    start_ts="${last_mon_date} 00:00:00"
    end_ts="${this_mon_date} 00:00:00"
    label=$(date -d "${last_mon_date}" +%Y-%m-%d)
    report_type="weekly"
    week_start_date="${last_mon_date}"
    week_end_date=$(date -d "${this_mon_date} -1 day" +%Y-%m-%d)
    ;;
  last_month)
    first_of_this_month=$(date +%Y-%m-01)
    start_month=$(date -d "${first_of_this_month} -1 month" +%Y-%m-01)
    start_ts="${start_month} 00:00:00"
    end_ts="${first_of_this_month} 00:00:00"
    label=$(date -d "${first_of_this_month} -1 month" +%Y-%m)
    report_type="monthly"
    ;;
  *)
    # default to yesterday
    start_ts=$(date -d "yesterday 00:00:00" +'%Y-%m-%d 00:00:00')
    end_ts=$(date -d "today 00:00:00" +'%Y-%m-%d 00:00:00')
    label=$(date -d "yesterday" +%Y-%m-%d)
    report_type="daily"
    ;;
esac

repo_slug="${CNB_REPO_SLUG:-}"
if [ -z "${repo_slug}" ]; then
  # Fallback to local repo dir name
  if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    repo_slug=$(basename "$(git rev-parse --show-toplevel)")
  else
    repo_slug="repo"
  fi
fi

mkdir -p reports tmp || true

out_file="reports/report-${label}.json"
ctx_file="tmp/git-context-${label}.md"

# jq is required for template updates and config parsing
if ! command -v jq >/dev/null 2>&1; then
  echo "[error] jq is required but not found" >&2
  exit 1
fi

# Utility: exclude patterns
matches_exclude() {
  local f="$1"
  local IFS=';'
  # shellcheck disable=SC2086
  for pat in ${exclude_globs}; do
    [ -z "${pat}" ] && continue
    case "$f" in
      $pat) return 0 ;;
    esac
  done
  return 1
}

# Utility: redact secrets in patches
redact() {
  sed -E \
    -e 's/(Authorization:)[[:space:]]*[^\r\n]*/\1 ***REDACTED***/Ig' \
    -e 's/(Bearer)[[:space:]]+[A-Za-z0-9._-]+/\1 ***REDACTED***/Ig' \
    -e 's/(sk-[A-Za-z0-9_-]{10,})/sk-***REDACTED***/g' \
    -e 's/(CNB_TOKEN=)[^\r\n& ]+/\1***REDACTED***/g'
}

sanitize_slug() {
  local s="$1"
  s="${s// /-}"
  s="${s//\//-}"
  s="${s//[^A-Za-z0-9._-]/-}"
  echo "$s"
}

# Pull report config from API if requested
load_report_config() {
  local resp
  if [ -z "${report_repo_list}" ] && [ "${report_config_from_api}" = "1" ] && [ -n "${cabb_api_base}" ] && [ -n "${integration_token}" ]; then
    resp=$(curl -sfSL -H "Authorization: Bearer ${integration_token}" "${cabb_api_base%/}/jobs/report/config" || true)
    if [ -n "${resp}" ]; then
      report_repo_list=$(echo "${resp}" | jq -c '.report_repos // []' 2>/dev/null)
      api_publish_repo=$(echo "${resp}" | jq -r '.output_repo_url // empty' 2>/dev/null)
      api_publish_branch=$(echo "${resp}" | jq -r '.output_branch // empty' 2>/dev/null)
      api_publish_dir=$(echo "${resp}" | jq -r '.output_dir // empty' 2>/dev/null)
      if [ -n "${api_publish_repo}" ] && [ -z "${REPORT_PUBLISH_REPO_URL:-}" ]; then
        publish_repo_url="${api_publish_repo}"
      fi
      if [ -n "${api_publish_branch}" ] && [ -z "${REPORT_PUBLISH_BRANCH:-}" ]; then
        publish_branch="${api_publish_branch}"
      fi
      if [ -n "${api_publish_dir}" ] && [ -z "${REPORT_PUBLISH_DIR:-}" ]; then
        publish_dir_root="${api_publish_dir}"
      fi
    fi
  fi
  if [ -z "${report_repo_list}" ]; then
    report_repo_list=$(printf '[{"slug":"%s","repo_url":"","branch":""}]' "${repo_slug}")
  fi
  if ! echo "${report_repo_list}" | jq empty >/dev/null 2>&1; then
    echo "[warn] invalid REPORT_REPO_LIST JSON; falling back to current repo only" >&2
    report_repo_list=$(printf '[{"slug":"%s","repo_url":"","branch":""}]' "${repo_slug}")
  fi
}

# Prepare repo locally (clone if needed) and ensure history exists
prepare_repo() {
  local slug="$1"
  local repo_url="$2"
  local branch="$3"
  local repo_path=""
  if [ -z "${repo_url}" ]; then
    if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
      repo_path=$(git rev-parse --show-toplevel)
    fi
  else
    repo_path="tmp/repos/${slug}"
    if [ ! -d "${repo_path}/.git" ]; then
      mkdir -p "tmp/repos"
      echo "[info] cloning ${repo_url} into ${repo_path}..." >&2
      if [ -n "${CNB_TOKEN:-}" ]; then
        auth_hdr="Authorization: Basic $(printf "cnb:%s" "${CNB_TOKEN}" | base64 | tr -d '\n')"
        GIT_CURL_VERBOSE=0 git -c http.extraHeader="${auth_hdr}" clone --no-tags --filter=blob:none "${repo_url}" "${repo_path}" >/dev/null 2>&1 || true
      fi
      if [ ! -d "${repo_path}/.git" ]; then
        git clone --no-tags --filter=blob:none "${repo_url}" "${repo_path}" >/dev/null 2>&1 || true
      fi
    fi
  fi
  if [ -z "${repo_path}" ] || [ ! -d "${repo_path}/.git" ]; then
    echo "[warn] repo ${slug} unavailable" >&2
    echo ""
    return
  fi
  git -C "${repo_path}" config --global --add safe.directory "$(cd "${repo_path}" && pwd)" >/dev/null 2>&1 || true
  if git -C "${repo_path}" remote get-url origin >/dev/null 2>&1; then
    (git -C "${repo_path}" remote set-branches origin "*" >/dev/null 2>&1 || true)
  fi
  (git -C "${repo_path}" fetch --all --prune --tags --recurse-submodules=no >/dev/null 2>&1 || true)
  (git -C "${repo_path}" fetch --unshallow >/dev/null 2>&1 || git -C "${repo_path}" fetch --depth=0 >/dev/null 2>&1 || true)
  if [ -n "${branch}" ]; then
    (git -C "${repo_path}" fetch origin "${branch}" --depth=0 >/dev/null 2>&1 || true)
  fi
  echo "${repo_path}"
}

has_commits=false

# Collect repo-specific commit context and patch samples
collect_repo_context() {
  local slug="$1"
  local repo_path="$2"
  local branch="$3"
  local display_name="$4"
  local repo_url="$5"
  local budget="$6"
  local title="${display_name:-$slug}"
  local revset="--all"

  echo "## Repo: ${title}" >> "${ctx_file}"
  echo "- slug: ${slug}" >> "${ctx_file}"
  if [ -n "${repo_url}" ]; then
    echo "- repo_url: ${repo_url}" >> "${ctx_file}"
  fi
  if [ -n "${branch}" ]; then
    if git -C "${repo_path}" rev-parse --verify --quiet "${branch}" >/dev/null 2>&1; then
      revset="${branch}"
    elif git -C "${repo_path}" rev-parse --verify --quiet "origin/${branch}" >/dev/null 2>&1; then
      revset="origin/${branch}"
    else
      echo "> branch ${branch} not found; skipping" >> "${ctx_file}"
      echo >> "${ctx_file}"
      return
    fi
    echo "- branch scope: ${branch}" >> "${ctx_file}"
  else
    echo "- branch scope: all branches" >> "${ctx_file}"
  fi

  local log_file="tmp/${slug}-git-logs-${label}.txt"
  if [ "${revset}" = "--all" ]; then
    git -C "${repo_path}" log --all \
      --since="${start_ts}" \
      --until="${end_ts}" \
      --date=format:'%Y-%m-%d %H:%M' \
      --pretty=format:'%H%x09%an%x09%ad%x09%s' \
      > "${log_file}" || true
  else
    git -C "${repo_path}" log "${revset}" \
      --since="${start_ts}" \
      --until="${end_ts}" \
      --date=format:'%Y-%m-%d %H:%M' \
      --pretty=format:'%H%x09%an%x09%ad%x09%s' \
      > "${log_file}" || true
  fi

  if [ ! -s "${log_file}" ]; then
    echo "> No commits found in this repo for the selected window." >> "${ctx_file}"
    echo >> "${ctx_file}"
    return
  fi

  has_commits=true

  echo "### Author counts" >> "${ctx_file}"
  awk -F "\t" '{print $2}' "${log_file}" | sed '/^$/d' | sort | uniq -c | sort -nr \
    | awk '{c=$1; $1=""; sub(/^ /, ""); printf("- %s: %d commits\n", $0, c)}' >> "${ctx_file}"
  echo >> "${ctx_file}"

  echo "### Commit details (numstat + message)" >> "${ctx_file}"
  if [ "${revset}" = "--all" ]; then
    git -C "${repo_path}" log --all \
      --since="${start_ts}" --until="${end_ts}" \
      --date=iso-local \
      --numstat \
      --pretty=format:'---%ncommit %H%nauthor %an <%ae>%ndate %ad%ntitle %s%nbody %b' \
      | head -n "${numstat_lines}" >> "${ctx_file}"
  else
    git -C "${repo_path}" log "${revset}" \
      --since="${start_ts}" --until="${end_ts}" \
      --date=iso-local \
      --numstat \
      --pretty=format:'---%ncommit %H%nauthor %an <%ae>%ndate %ad%ntitle %s%nbody %b' \
      | head -n "${numstat_lines}" >> "${ctx_file}"
  fi
  echo >> "${ctx_file}"

  if [ "${diff_mode}" = "stats" ]; then
    echo >> "${ctx_file}"
    return
  fi

  local commit_list=()
  if [ "${revset}" = "--all" ]; then
    mapfile -t commit_list < <(git -C "${repo_path}" rev-list --all --since="${start_ts}" --until="${end_ts}" --no-merges)
  else
    mapfile -t commit_list < <(git -C "${repo_path}" rev-list "${revset}" --since="${start_ts}" --until="${end_ts}" --no-merges)
  fi
  local tmp_scores="tmp/commit-scores-${slug}-${label}.txt"
  : > "${tmp_scores}"
  for sha in "${commit_list[@]}"; do
    sum=$(git -C "${repo_path}" show --numstat --format="" "$sha" | awk '{a=$1;b=$2; if(a=="-")a=0; if(b=="-")b=0; s+=a+b} END{print s+0}')
    echo "$sum $sha" >> "${tmp_scores}"
  done
  mapfile -t top_commits < <(sort -nr -k1,1 "${tmp_scores}" | awk '{print $2}' | head -n "${max_commits}")

  local total_len=0
  local patch_ctx_file="tmp/git-patch-${slug}-${label}.txt"
  : > "${patch_ctx_file}"

  for sha in "${top_commits[@]}"; do
    hdr=$(git -C "${repo_path}" show -s --format='commit %H%nauthor %an <%ae>%ndate %ad%ntitle %s%n' --date=iso-local "$sha")
    entry_file="tmp/patch-entry-${slug}-${sha}.txt"
    printf "%s\n" "$hdr" > "$entry_file"
    mapfile -t files < <(git -C "${repo_path}" show --numstat --format="" "$sha" | awk '{a=$1;b=$2;f=$3; if(a=="-")a=0; if(b=="-")b=0; print a+b, f}' | sort -nr -k1,1 | awk '{ $1=""; sub(/^ /, ""); print }')
    count=0
    for f in "${files[@]}"; do
      [ -z "$f" ] && continue
      if matches_exclude "$f"; then
        continue
      fi
      file_patch=$(git -C "${repo_path}" show --no-color --unified=3 --format="" "$sha" -- "$f" | redact)
      if [ -n "$file_patch" ]; then
        if [ "${diff_mode}" = "sampled_patch" ]; then
          file_patch=$(printf "%s\n" "$file_patch" | awk -v N="${max_hunks_per_file}" 'BEGIN{h=0} { if($0 ~ /^@@/) h++; if(h==0 || h<=N) print }')
        fi
        printf "diff -- sampled file: %s\n%s\n\n" "$f" "$file_patch" >> "$entry_file"
        count=$((count+1))
        if [ "$count" -ge "$top_files" ]; then
          break
        fi
      fi
    done
    entry_size=$(wc -c < "$entry_file" | tr -d ' ')
    if [ "$(( total_len + entry_size ))" -le "$budget" ]; then
      cat "$entry_file" >> "$patch_ctx_file"
      total_len=$(( total_len + entry_size ))
    else
      remain=$(( budget - total_len ))
      if [ "$remain" -gt 0 ]; then
        head -c "$remain" "$entry_file" >> "$patch_ctx_file" || true
        total_len=$budget
      fi
      break
    fi
  done
  if [ -s "$patch_ctx_file" ]; then
    {
      echo "### Patch samples (truncated by budget if necessary)"
      echo "Budget: chars=${budget}, top_files=${top_files}, hunks_per_file=${max_hunks_per_file}, max_commits=${max_commits}"
      cat "$patch_ctx_file"
    } >> "$ctx_file"
    echo >> "$ctx_file"
  fi
}

# Initialize output file from template
if [ -f "${template_file}" ]; then
  cp "${template_file}" "${out_file}"
  tmp=$(mktemp)
  jq --arg d "${label}" '.date = $d' "${out_file}" > "$tmp" && mv "$tmp" "${out_file}"
else
  echo "{}" > "${out_file}"
fi

load_report_config
mapfile -t repo_entries < <(echo "${report_repo_list}" | jq -c '.[]')
if [ "${#repo_entries[@]}" -eq 0 ]; then
  repo_entries=()
  repo_entries+=("$(printf '{"slug":"%s","repo_url":"","branch":""}' "${repo_slug}")")
fi

repo_count=${#repo_entries[@]}
budget_per_repo="${char_budget}"
if [ "${repo_count}" -gt 0 ] && [ "${char_budget}" -gt 0 ]; then
  budget_per_repo=$(( char_budget / repo_count ))
  if [ "${budget_per_repo}" -lt 50000 ] && [ "${char_budget}" -ge 50000 ]; then
    budget_per_repo=50000
  fi
  if [ "${budget_per_repo}" -le 0 ]; then
    budget_per_repo="${char_budget}"
  fi
fi

{
  echo "# Multi-repo Commit Context"
  echo "- Time window: ${start_ts} .. ${end_ts} (${TZ})"
  echo "- Branch (trigger): ${CNB_BRANCH:-}"
  echo "- Repos configured: ${repo_count}"
  echo
} > "${ctx_file}"

for entry in "${repo_entries[@]}"; do
  repo_url=$(echo "${entry}" | jq -r '.repo_url // ""')
  repo_branch=$(echo "${entry}" | jq -r '.branch // ""')
  repo_display=$(echo "${entry}" | jq -r '.display_name // ""')
  slug_raw=$(echo "${entry}" | jq -r '.slug // ""')
  slug=$(sanitize_slug "${slug_raw}")
  if [ -z "${slug}" ]; then
    base_slug=$(basename "${repo_url}")
    slug=$(sanitize_slug "${base_slug:-repo}")
  fi
  repo_path=$(prepare_repo "${slug}" "${repo_url}" "${repo_branch}")
  if [ -z "${repo_path}" ]; then
    echo "## Repo: ${slug}" >> "${ctx_file}"
    echo "> repository unavailable; skipped" >> "${ctx_file}"
    echo >> "${ctx_file}"
    continue
  fi
  collect_repo_context "${slug}" "${repo_path}" "${repo_branch}" "${repo_display}" "${repo_url}" "${budget_per_repo}"
done

# Try opencode to author a polished report from commit context
# 调用 opencode 生成“按人员维度总结”的结构化中文汇报
try_opencode() {
  command -v opencode >/dev/null 2>&1 || return 1

  # Avoid leaking key to logs; pass via env
  export OPENCODE_MODEL="${model}"
  if [ -n "${api_key}" ]; then
    export OPENCODE_API_KEY="${api_key}"
  fi

  # Chinese prompt to produce a structured work report from commits and numstat
  case "${timeframe}" in
    yesterday) period_label="日报" ;;
    last_week) period_label="周报" ;;
    last_month) period_label="月报" ;;
    *) period_label="汇报" ;;
  esac
  prompt="你是资深工程经理，需基于提交记录（含 numstat 以及必要时的采样补丁片段）生成 JSON 格式的${period_label}。\n"
  prompt+="时间范围：${start_ts} 至 ${end_ts}（${TZ}）。\n"
  prompt+="输入包含多个仓库的提交，请在概述与细项中标注仓库归属（建议用 [repo] 前缀），按“仓库 + 人员”维度组织。\n"
  prompt+="请读取提供的 template.json 结构，并输出填充后的 JSON 内容（不要包含 markdown 代码块标记）。\n"
  prompt+="内容要求：\n"
  prompt+="1. progress_summary: 针对非开发岗位的汇报。\n"
  prompt+="   - overview: 本期主题、亮点、重要进展（通俗易懂）。\n"
  prompt+="   - details: 数组，每项包含 { \"topic\": \"...\", \"content\": \"...\" }，按功能/业务模块划分。\n"
  prompt+="2. code_review_summary: 针对开发岗位的 Code Review 总结。\n"
  prompt+="   - overview: 代码质量、架构变动、技术风险概述。\n"
  prompt+="   - details: 数组，每项包含 { \"author\": \"...\", \"changes\": \"...\", \"suggestions\": \"...\" }，按人员汇总。\n"
  prompt+="要求：\n- 严格遵守 JSON 格式；\n- 语言简练，结论先行；\n- 若无提交，在 overview 中说明。\n"

  # 非交互模式：消息作为第一个参数，避免被解析为 -f 的文件
  # Pass template content as context too if needed, or just rely on prompt description.
  # Let's pass template as a file to opencode context? No, opencode -f takes one file.
  # We can append template to ctx_file.
  echo -e "\n\n# Template JSON Structure\n\`\`\`json" >> "${ctx_file}"
  cat "${template_file}" >> "${ctx_file}"
  echo -e "\n\`\`\`" >> "${ctx_file}"

  if opencode run --format json "${prompt}" -m "${OPENCODE_MODEL}" -f "${ctx_file}" > "${out_file}.tmp" 2>"tmp/opencode.stderr"; then
     # Validate JSON
     if jq . "${out_file}.tmp" >/dev/null 2>&1; then
       mv "${out_file}.tmp" "${out_file}"
       return 0
     fi
  fi
  # Surface concise error context without secrets
  {
    echo "[opencode] failed to generate AI summary"
    echo "- model: ${OPENCODE_MODEL}"
    if [ -n "${OPENCODE_API_KEY:-}" ]; then echo "- apiKey: present"; else echo "- apiKey: missing"; fi
    echo "- ctx_file size: $(wc -c < "${ctx_file}" 2>/dev/null || echo 0) bytes"
    if [ -s "tmp/opencode.stderr" ]; then
      echo "- stderr (first 80 lines):"
      sed -n '1,80p' "tmp/opencode.stderr"
    fi
  } >&2
  return 1
}

if try_opencode; then
  echo "AI summary generated." >&2
else
  # Fallback: native grouped summary
  if [ "${require_ai}" = "1" ]; then
    echo "## AI 汇总生成失败" >> "${out_file}"
    echo "本次手动触发配置要求生成 AI 汇总，但调用失败。请检查 OPENCODE_API_KEY 与网络后重试。" >> "${out_file}"
    exit 1
  else
    if ${has_commits}; then
      # Fallback JSON
      jq '.progress_summary.overview = "AI generation failed, but commits exist."' "${out_file}" > "${out_file}.tmp" && mv "${out_file}.tmp" "${out_file}"
    else
      jq '.progress_summary.overview = "No commits found in configured repos."' "${out_file}" > "${out_file}.tmp" && mv "${out_file}.tmp" "${out_file}"
    fi
  fi
fi

echo "report generated: ${out_file}" >&2

# 发布：推送报告到目标仓库（ai-report/{daily|weekly|monthly}/...），同一时间窗覆盖旧文件
publish_report() {
  # Require token and repo URL
  if [ -z "${publish_repo_url}" ] || [ -z "${CNB_TOKEN:-}" ]; then
    return 0
  fi

  # Compute target path and filename by type
  target_subdir="${publish_dir_root}/${report_type}"
  case "${report_type}" in
    daily)
      target_filename="report-${label}.json"
      range_text="${label}"
      ;;
    weekly)
      if [ -n "${week_start_date:-}" ] && [ -n "${week_end_date:-}" ]; then
        target_filename="report-${week_start_date}_to_${week_end_date}.json"
        range_text="${week_start_date}~${week_end_date}"
      else
        target_filename="report-${label}.json"
        range_text="${label}"
      fi
      ;;
    monthly)
      target_filename="report-${label}.json"
      range_text="${label}"
      ;;
    *)
      target_filename="report-${label}.json"
      range_text="${label}"
      ;;
  esac

  workdir="tmp/publish"
  rm -rf "${workdir}" && mkdir -p "${workdir}"
  auth_hdr="Authorization: Basic $(printf "cnb:%s" "${CNB_TOKEN}" | base64 | tr -d '\n')"
  # Clone target
  GIT_CURL_VERBOSE=0 git -c http.extraHeader="${auth_hdr}" clone "${publish_repo_url}" "${workdir}" >/dev/null 2>&1 || return 0
  cd "${workdir}"
  git config --global --add safe.directory "$(pwd)" || true
  # Checkout branch (create if missing)
  if git show-ref --verify --quiet "refs/heads/${publish_branch}"; then
    git checkout "${publish_branch}" >/dev/null 2>&1 || true
  else
    git checkout -b "${publish_branch}" >/dev/null 2>&1 || true
  fi
  # Rebase onto latest remote to avoid non-fast-forward on push
  GIT_CURL_VERBOSE=0 git -c http.extraHeader="${auth_hdr}" pull --rebase origin "${publish_branch}" >/dev/null 2>&1 || true
  mkdir -p "${target_subdir}"
  cp -f "${OLDPWD}/${out_file}" "${target_subdir}/${target_filename}"
  git config user.name "cabb-report-bot"
  git config user.email "bot@cabb.local"
  git add -A "${target_subdir}" || true
  if git diff --cached --quiet >/dev/null 2>&1; then
    # no changes
    cd - >/dev/null 2>&1 || true
    return 0
  fi
  git commit -m "chore(report): ${report_type} ${range_text}" >/dev/null 2>&1 || true
  if ! GIT_CURL_VERBOSE=0 git -c http.extraHeader="${auth_hdr}" push origin "${publish_branch}" >/dev/null 2>&1; then
    # Try one rebase cycle then push again
    GIT_CURL_VERBOSE=0 git -c http.extraHeader="${auth_hdr}" pull --rebase origin "${publish_branch}" >/dev/null 2>&1 || true
    GIT_CURL_VERBOSE=0 git -c http.extraHeader="${auth_hdr}" push origin "${publish_branch}" >/dev/null 2>&1 || true
  fi
  cd - >/dev/null 2>&1 || true
}

if [ "${publish_enable}" = "1" ] || [ "${publish_enable}" = "true" ]; then
  publish_report || true
fi
