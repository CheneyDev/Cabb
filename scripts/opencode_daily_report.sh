#!/usr/bin/env bash
set -euo pipefail

# Generate a daily Markdown report for git commits within a time window.
# Tries to use opencode (zen mode) if available; otherwise falls back
# to a native git-based summary grouped by author.

export TZ=${TZ:-Asia/Shanghai}

timeframe="${REPORT_TIMEFRAME:-yesterday}"
zen_mode="${ZEN_MODE:-1}"
# prefer fully-qualified model id per docs (opencode/<model-id>)
model="${OPENCODE_MODEL:-opencode/grok-code}"
# Prefer OPENCODE_API_KEY; fallbacks are best-effort in case secrets use different key names
api_key="${OPENCODE_API_KEY:-}"
if [ -z "${api_key}" ]; then
  api_key="${OPENCODE_TOKEN:-}"
fi
if [ -z "${api_key}" ]; then
  api_key="${OPENCODE_KEY:-}"
fi

# Determine time window and label
case "${timeframe}" in
  yesterday)
    start_ts=$(date -d "yesterday 00:00:00" +'%Y-%m-%d 00:00:00')
    end_ts=$(date -d "today 00:00:00" +'%Y-%m-%d 00:00:00')
    label=$(date -d "yesterday" +%Y-%m-%d)
    ;;
  today)
    start_ts=$(date -d "today 00:00:00" +'%Y-%m-%d 00:00:00')
    end_ts=$(date -d "tomorrow 00:00:00" +'%Y-%m-%d 00:00:00')
    label=$(date +%Y-%m-%d)
    ;;
  *)
    echo "unknown REPORT_TIMEFRAME: ${timeframe}, fallback to yesterday" >&2
    start_ts=$(date -d "yesterday 00:00:00" +'%Y-%m-%d 00:00:00')
    end_ts=$(date -d "today 00:00:00" +'%Y-%m-%d 00:00:00')
    label=$(date -d "yesterday" +%Y-%m-%d)
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

mkdir -p daily-reports tmp || true

out_file="daily-reports/daily-report-${label}.md"
log_file="tmp/git-logs-${label}.txt"

# Make git workspace safe for CI and try to ensure history availability
if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  git_root=$(git rev-parse --show-toplevel)
  git config --global --add safe.directory "${git_root}" || true
  # Best-effort: fetch full history and all branches
  # Some CI clones are shallow; widen history so we can see yesterday's commits,
  # and include non-default branches.
  (git fetch --all --prune --tags --recurse-submodules=no 2>/dev/null || true)
  (git fetch --unshallow 2>/dev/null || git fetch --depth=0 2>/dev/null || true)
fi

# Collect raw logs for the time window
if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  git log \
    --all \
    --since="${start_ts}" \
    --until="${end_ts}" \
    --date=format:'%Y-%m-%d %H:%M' \
    --pretty=format:'%H%x09%an%x09%ad%x09%s' \
    > "${log_file}" || true
else
  : > "${log_file}"
fi

has_commits=false
if [ -s "${log_file}" ]; then
  has_commits=true
fi

echo "# 项目工作日报（${label}）" > "${out_file}"
echo "" >> "${out_file}"
echo "- 时间范围：${start_ts} 至 ${end_ts} (${TZ})" >> "${out_file}"
echo "- 仓库：${repo_slug}" >> "${out_file}"
echo "" >> "${out_file}"

if ${has_commits}; then
  echo "## 总览" >> "${out_file}"
  # Count commits by author
  awk -F "\t" '{print $2}' "${log_file}" | sort | uniq -c | sort -nr \
    | awk '{c=$1; $1=""; sub(/^ /, ""); printf("- %s：%d 次提交\n", $0, c)}' >> "${out_file}"
  echo "" >> "${out_file}"
else
  echo "> 本时段暂无提交记录。" >> "${out_file}"
  echo "" >> "${out_file}"
fi

# Try opencode (zen mode) to author a polished report from raw logs
try_opencode() {
  command -v opencode >/dev/null 2>&1 || return 1

  # Avoid leaking key to logs; pass via env
  export OPENCODE_MODEL="${model}"
  if [ -n "${api_key}" ]; then
    export OPENCODE_API_KEY="${api_key}"
  fi

  prompt="请基于以下 git 提交原始记录，生成中文 Markdown 工作日报：\n\n"
  prompt+="- 范围：${start_ts} 至 ${end_ts} (${TZ})\n"
  prompt+="- 要求：按作者分组，列出每位作者的提交条数与关键变更点；保留小标题，条目精炼；最后提供整体统计与风险/待办提示。\n"

  # Use non-interactive CLI per docs: opencode run
  if opencode run -m "${OPENCODE_MODEL}" -f "${log_file}" "${prompt}" > "${out_file}.ai" 2>/dev/null; then
    return 0
  fi
  # Fallback: attach server if running (unlikely in CI)
  if opencode run --format json -m "${OPENCODE_MODEL}" -f "${log_file}" "${prompt}" > "${out_file}.ai" 2>/dev/null; then
    return 0
  fi
  return 1
}

if ${has_commits} && try_opencode; then
  echo "## AI 汇总（opencode）" >> "${out_file}"
  echo "" >> "${out_file}"
  cat "${out_file}.ai" >> "${out_file}" || true
  echo "" >> "${out_file}"
else
  # Fallback: native grouped summary
  if ${has_commits}; then
    echo "## 分作者提交明细" >> "${out_file}"
    echo "" >> "${out_file}"
    # List unique authors preserving locale
    mapfile -t authors < <(awk -F "\t" '{print $2}' "${log_file}" | sort -u)
    for author in "${authors[@]}"; do
      echo "### ${author}" >> "${out_file}"
      echo "" >> "${out_file}"
      # Short hash + subject per commit
      git log \
        --all \
        --since="${start_ts}" \
        --until="${end_ts}" \
        --author="${author}" \
        --pretty=format:'- %h %s' \
        >> "${out_file}" || true
      echo "" >> "${out_file}"
    done
  fi
fi

# Append raw logs for reference (collapsed in most viewers)
echo "<details><summary>原始提交记录</summary>" >> "${out_file}"
echo "" >> "${out_file}"
if ${has_commits}; then
  printf '````\n' >> "${out_file}"
  cat "${log_file}" >> "${out_file}"
  printf '\n````\n' >> "${out_file}"
else
  echo "(空)" >> "${out_file}"
fi
echo "" >> "${out_file}"
echo "</details>" >> "${out_file}"

echo "report generated: ${out_file}" >&2
