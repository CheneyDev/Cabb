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
# require AI summary or fallback allowed (manual可设置 REQUIRE_AI_SUMMARY=1)
require_ai="${REQUIRE_AI_SUMMARY:-0}"
# diff context controls (manual trigger can override via web_trigger inputs)
diff_mode="${REPORT_DIFF_MODE:-stats}" # stats | sampled_patch | full_patch
char_budget="${REPORT_DIFF_CHAR_BUDGET:-200000}"
top_files="${REPORT_DIFF_TOP_FILES:-20}"
max_hunks_per_file="${REPORT_DIFF_HUNKS_PER_FILE:-3}"
max_commits="${REPORT_DIFF_MAX_COMMITS:-100}"
exclude_globs="${REPORT_EXCLUDE_GLOBS:-node_modules/**;dist/**;vendor/**;*.min.*;*.png;*.jpg;*.jpeg;*.gif;*.webp;*.svg;*.mp4;*.mov;*.zip;*.tar;*.gz;*.lock;*.class;*.bin;*.exe;*.dylib;*.so;*.dll}"
numstat_lines="${REPORT_NUMSTAT_LINES:-5000}"
# Prefer OPENCODE_API_KEY; fallbacks are best-effort in case secrets use different key names
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
case "${timeframe}" in
  yesterday)
    start_ts=$(date -d "yesterday 00:00:00" +'%Y-%m-%d 00:00:00')
    end_ts=$(date -d "today 00:00:00" +'%Y-%m-%d 00:00:00')
    label=$(date -d "yesterday" +%Y-%m-%d)
    ;;
  last_week)
    # Compute previous calendar week (Mon 00:00:00 to this week's Mon 00:00:00)
    dow=$(date +%u) # 1..7 (Mon=1)
    this_mon_date=$(date -d "today -$((dow-1)) days" +%Y-%m-%d)
    last_mon_date=$(date -d "${this_mon_date} -7 days" +%Y-%m-%d)
    start_ts="${last_mon_date} 00:00:00"
    end_ts="${this_mon_date} 00:00:00"
    label=$(date -d "${last_mon_date}" +%Y-%m-%d)
    ;;
  last_month)
    first_of_this_month=$(date +%Y-%m-01)
    start_month=$(date -d "${first_of_this_month} -1 month" +%Y-%m-01)
    start_ts="${start_month} 00:00:00"
    end_ts="${first_of_this_month} 00:00:00"
    label=$(date -d "${first_of_this_month} -1 month" +%Y-%m)
    ;;
  *)
    # default to yesterday
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
patch_ctx_file="tmp/git-patch-${label}.txt"

# Make git workspace safe for CI and try to ensure history availability
ensure_repo() {
  if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    return 0
  fi
  # If workspace is not a git repo (some CI snapshots omit .git), clone a fresh copy
  if [ -n "${CNB_REPO_URL_HTTPS:-}" ] && [ -n "${CNB_TOKEN:-}" ]; then
    echo "[fallback] cloning repo to tmp/clone for history..." >&2
    mkdir -p tmp/clone
    # Avoid leaking token to logs; use header auth instead of embedding in URL
    git -c http.extraHeader="Authorization: Basic $(printf "cnb:%s" "${CNB_TOKEN}" | base64 | tr -d '\n')" \
      clone --no-tags --filter=blob:none --mirror "${CNB_REPO_URL_HTTPS}" tmp/clone 2>/dev/null || {
        echo "[warn] clone via header failed, trying credential-in-url (may be masked)" >&2
        git clone --no-tags --filter=blob:none --mirror "https://cnb:${CNB_TOKEN}@${CNB_REPO_URL_HTTPS#https://}" tmp/clone 2>/dev/null || true
      }
    if [ -d tmp/clone ]; then
      # Use a worktree from the mirror to query logs
      mkdir -p tmp/work
      git --git-dir=tmp/clone --work-tree=tmp/work init 2>/dev/null || true
      (cd tmp/work && git config --local core.bare false >/dev/null 2>&1 || true)
      export GIT_DIR="$(pwd)/tmp/clone"
      export GIT_WORK_TREE="$(pwd)/tmp/work"
      return 0
    fi
  fi
  return 1
}

ensure_repo || true

if git rev-parse --is-inside-work-tree >/dev/null 2>&1 || [ -n "${GIT_DIR:-}" ]; then
  git_root=$(git rev-parse --show-toplevel)
  git config --global --add safe.directory "${git_root}" || true
  # Best‑effort: fetch full history and all branches, even if CI did a partial clone
  # Ensure remote tracks all branches
  if git remote get-url origin >/dev/null 2>&1; then
    (git remote set-branches origin "*" 2>/dev/null || true)
  fi
  # Fetch all refs and unshallow if needed
  (git fetch --all --prune --tags --recurse-submodules=no 2>/dev/null || true)
  (git fetch --unshallow 2>/dev/null || git fetch --depth=0 2>/dev/null || true)
fi

# Collect raw logs for the time window
if git rev-parse --is-inside-work-tree >/dev/null 2>&1 || [ -n "${GIT_DIR:-}" ]; then
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

# No implicit fallback; timeframes are strictly one of: yesterday, last_week, last_month

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

# Build rich commit context for AI (includes per-commit numstat and bodies)
ctx_file="tmp/git-context-${label}.md"
{
  echo "# Repo Commit Context"
  echo "- Time window: ${start_ts} .. ${end_ts} (${TZ})"
  echo "- Repo: ${repo_slug}"
  echo "- Branch (trigger): ${CNB_BRANCH:-}"
  echo
  if ${has_commits}; then
    echo "## Author counts"
    awk -F "\t" '{print $2}' "${log_file}" | sed '/^$/d' | sort | uniq -c | sort -nr \
      | awk '{c=$1; $1=""; sub(/^ /, ""); printf("- %s: %d commits\n", $0, c)}'
    echo
    echo "## Commit details (numstat + message)"
    echo "NOTE: This is a summary (numstat) view, not full patches."
    git log --all \
      --since="${start_ts}" --until="${end_ts}" \
      --date=iso-local \
      --numstat \
      --pretty=format:'---%ncommit %H%nauthor %an <%ae>%ndate %ad%ntitle %s%nbody %b' \
      | head -n "${numstat_lines}"
    echo
    if [ "${diff_mode}" != "stats" ]; then
      echo "## Patch samples"
      echo "Budget: chars=${char_budget}, top_files=${top_files}, hunks_per_file=${max_hunks_per_file}, max_commits=${max_commits}"
    fi
  else
    echo "> No commits found in the selected window."
  fi
} > "${ctx_file}"

# Prepare sampled/full patch context if requested
> "${patch_ctx_file}"
if ${has_commits} && [ "${diff_mode}" != "stats" ]; then
  # Collect candidate commits within window (exclude merges for patch to control size)
  mapfile -t commit_list < <(git rev-list --all --since="${start_ts}" --until="${end_ts}" --no-merges)
  tmp_scores="tmp/commit-scores-${label}.txt"
  : > "${tmp_scores}"
  for sha in "${commit_list[@]}"; do
    sum=$(git show --numstat --format="" "$sha" | awk '{a=$1;b=$2; if(a=="-")a=0; if(b=="-")b=0; s+=a+b} END{print s+0}')
    echo "$sum $sha" >> "${tmp_scores}"
  done
  mapfile -t top_commits < <(sort -nr -k1,1 "${tmp_scores}" | awk '{print $2}' | head -n "${max_commits}")

  total_len=0
  for sha in "${top_commits[@]}"; do
    hdr=$(git show -s --format='commit %H%nauthor %an <%ae>%ndate %ad%ntitle %s%n' --date=iso-local "$sha")
    entry_file="tmp/patch-entry-${sha}.txt"
    printf "%s\n" "$hdr" > "$entry_file"
    mapfile -t files < <(git show --numstat --format="" "$sha" | awk '{a=$1;b=$2;f=$3; if(a=="-")a=0; if(b=="-")b=0; print a+b, f}' | sort -nr -k1,1 | awk '{ $1=""; sub(/^ /, ""); print }')
    count=0
    for f in "${files[@]}"; do
      [ -z "$f" ] && continue
      if matches_exclude "$f"; then
        continue
      fi
      file_patch=$(git show --no-color --unified=3 --format="" "$sha" -- "$f" | redact)
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
    if [ "$(( total_len + entry_size ))" -le "$char_budget" ]; then
      cat "$entry_file" >> "$patch_ctx_file"
      total_len=$(( total_len + entry_size ))
    else
      remain=$(( char_budget - total_len ))
      if [ "$remain" -gt 0 ]; then
        head -c "$remain" "$entry_file" >> "$patch_ctx_file" || true
        total_len=$char_budget
      fi
      break
    fi
  done
  if [ -s "$patch_ctx_file" ]; then
    {
      echo
      echo "## Patch samples (truncated by budget if necessary)"
      cat "$patch_ctx_file"
    } >> "$ctx_file"
  fi
fi

# Try opencode to author a polished report from commit context
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
  prompt="你是资深工程经理，需基于提交记录（含 numstat 与提交信息，且在必要时包含采样的补丁片段）生成中文${period_label}。\n"
  prompt+="时间范围：${start_ts} 至 ${end_ts}（${TZ}）。\n"
  prompt+="请输出结构化 Markdown：\n"
  prompt+="1. 概览：本期工作主题、亮点、影响范围。\n"
  prompt+="2. 统计：总提交数/按作者分布/改动规模（可据 numstat 估算）。\n"
  prompt+="3. 关键变更：按模块或特性归纳（列出文件/目录线索与简要说明）。\n"
  prompt+="4. 风险与待办：潜在风险、技术债、下一步行动（Owner/ETA）。\n"
  prompt+="要求：\n- 语言简练，避免赘述；\n- 不复述完整 diff，仅用 numstat 与文件路径提炼要点；\n- 若无提交，明确说明“本期无提交”。\n"

  # Use non-interactive CLI: place message first to avoid it being parsed as a file
  if opencode run "${prompt}" -m "${OPENCODE_MODEL}" -f "${ctx_file}" > "${out_file}.ai" 2>"tmp/opencode.stderr"; then
    return 0
  fi
  # Fallback: attach server if running (unlikely in CI)
  if opencode run --format json "${prompt}" -m "${OPENCODE_MODEL}" -f "${ctx_file}" > "${out_file}.ai" 2>>"tmp/opencode.stderr"; then
    return 0
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
  echo "## AI 汇总（opencode）" >> "${out_file}"
  echo "" >> "${out_file}"
  cat "${out_file}.ai" >> "${out_file}" || true
  echo "" >> "${out_file}"
else
  # Fallback: native grouped summary
  if [ "${require_ai}" = "1" ]; then
    echo "## AI 汇总生成失败" >> "${out_file}"
    echo "本次手动触发配置要求生成 AI 汇总，但调用失败。请检查 OPENCODE_API_KEY 与网络后重试。" >> "${out_file}"
    exit 1
  else
    if ${has_commits}; then
      echo "## 数据化汇总（备用）" >> "${out_file}"
      echo "" >> "${out_file}"
      # List unique authors preserving locale
      mapfile -t authors < <(awk -F "\t" '{print $2}' "${log_file}" | sed '/^$/d' | sort -u)
      for author in "${authors[@]}"; do
        echo "### ${author}" >> "${out_file}"
        echo "" >> "${out_file}"
        git log \
          --all \
          --since="${start_ts}" \
          --until="${end_ts}" \
          --author="${author}" \
          --pretty=format:'- %h %s' \
          >> "${out_file}" || true
        echo "" >> "${out_file}"
      done
    else
      echo "> 本时段暂无提交记录。" >> "${out_file}"
    fi
  fi
fi

# Append raw logs for reference (collapsed in most viewers)
echo "<details><summary>原始提交记录 / 作者分布</summary>" >> "${out_file}"
echo "" >> "${out_file}"
if ${has_commits}; then
  echo "作者分布：" >> "${out_file}"
  awk -F "\t" '{print $2}' "${log_file}" | sed '/^$/d' | sort | uniq -c | sort -nr \
    | awk '{c=$1; $1=""; sub(/^ /, ""); printf("- %s：%d 次提交\n", $0, c)}' >> "${out_file}"
  echo "" >> "${out_file}"
fi
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
