#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# AI 汇报流水线脚本（日报/周报/月报）
# ==============================================================================
#
# 功能概览：
#   - 跨分支聚合：统计指定时间窗（昨天/上一周/上个月）内"所有分支"的提交
#   - AI 汇总优先：基于提交统计 + 采样补丁，调用 opencode 生成中文结构化汇报
#   - 幂等发布：按固定命名将汇报推送到目标仓库（同一时间窗覆盖旧文件）
#   - 控制台预览：流水线末尾输出最新汇报的前 80 行
#
# 设计要点：
#   - 安全性：不在日志泄露密钥；补丁上下文会做脱敏处理
#   - 稳健性：CI 若为浅克隆会自动 unshallow；git safe.directory 避免权限告警
#   - 一致性：手动/定时共享同一实现与参数，差异仅在时间窗
#
# 依赖：
#   - git, curl, jq
#   - opencode CLI (可选，用于 AI 汇总)
#
# ==============================================================================

# ------------------------------------------------------------------------------
# 时区设置
# ------------------------------------------------------------------------------
export TZ=${TZ:-Asia/Shanghai}

# ==============================================================================
# 第一部分：环境变量读取与默认值设置
# ==============================================================================

# ------------------------------------------------------------------------------
# 汇报时间范围配置
# ------------------------------------------------------------------------------
# 支持值：yesterday（日报）、last_week（周报）、last_month（月报）
timeframe="${REPORT_TIMEFRAME:-yesterday}"

# ------------------------------------------------------------------------------
# OpenCode AI 配置
# ------------------------------------------------------------------------------
# 是否启用 Zen 模式（非交互式）
zen_mode="${ZEN_MODE:-1}"
# AI 模型标识，格式：opencode/<model-id>
model="${OPENCODE_MODEL:-opencode/grok-code}"
# 是否要求 AI 汇总必须成功（1=必须成功，0=允许回退到基础统计）
require_ai="${REQUIRE_AI_SUMMARY:-0}"

# ------------------------------------------------------------------------------
# 汇报发布配置
# ------------------------------------------------------------------------------
# 是否启用发布到目标仓库
publish_enable="${REPORT_PUBLISH_ENABLE:-1}"
# 目标仓库 URL
publish_repo_url="${REPORT_PUBLISH_REPO_URL:-https://cnb.cool/1024hub/plane-test}"
# 目标分支
publish_branch="${REPORT_PUBLISH_BRANCH:-main}"
# 目标目录（会按 daily/weekly/monthly 子目录组织）
publish_dir_root="${REPORT_PUBLISH_DIR:-reports}"
# 汇报模板文件路径
template_file="scripts/template.md"

# ------------------------------------------------------------------------------
# Git Diff 采样配置
# ------------------------------------------------------------------------------
# Diff 模式：stats=仅统计, sampled_patch=采样补丁, full_patch=完整补丁
diff_mode="${REPORT_DIFF_MODE:-stats}"
# Diff 内容字符预算上限
char_budget="${REPORT_DIFF_CHAR_BUDGET:-200000}"
# 每个 commit 最多采样的文件数
top_files="${REPORT_DIFF_TOP_FILES:-20}"
# 每个文件最多采样的 hunk（代码块）数
max_hunks_per_file="${REPORT_DIFF_HUNKS_PER_FILE:-3}"
# 最多处理的 commit 数量
max_commits="${REPORT_DIFF_MAX_COMMITS:-100}"
# 排除的文件 glob 模式（分号分隔）
exclude_globs="${REPORT_EXCLUDE_GLOBS:-node_modules/**;dist/**;vendor/**;*.min.*;*.png;*.jpg;*.jpeg;*.gif;*.webp;*.svg;*.mp4;*.mov;*.zip;*.tar;*.gz;*.lock;*.class;*.bin;*.exe;*.dylib;*.so;*.dll}"
# numstat 输出最大行数
numstat_lines="${REPORT_NUMSTAT_LINES:-5000}"

# ------------------------------------------------------------------------------
# 输出格式配置
# ------------------------------------------------------------------------------
# 是否输出 JSON 格式（中间格式，供后端读取转发飞书）
output_json="${REPORT_OUTPUT_JSON:-1}"
# 是否输出 Markdown 格式
output_markdown="${REPORT_OUTPUT_MARKDOWN:-1}"

# ------------------------------------------------------------------------------
# 后端 API 配置
# ------------------------------------------------------------------------------
# 是否从后端 API 拉取多仓库配置
report_config_from_api="${REPORT_CONFIG_FROM_API:-1}"
# 手动指定的仓库列表（JSON 数组格式）
report_repo_list="${REPORT_REPO_LIST:-}"
# Cabb 后端 API 地址
cabb_api_base="${CABB_API_BASE:-}"
# 集成 Token（用于 API 认证）
integration_token="${INTEGRATION_TOKEN:-}"

# ------------------------------------------------------------------------------
# 配置信息日志输出（用于调试）
# ------------------------------------------------------------------------------
token_present="no"
if [ -n "${integration_token}" ]; then
  token_present="yes"
fi
echo "[info] cfg: report_config_from_api=${report_config_from_api} cabb_api_base=${cabb_api_base:-<empty>} token_present=${token_present}" >&2

if [ -n "${CNB_TOKEN:-}" ]; then
  echo "[info] cfg: CNB_TOKEN present (length ${#CNB_TOKEN})" >&2
else
  echo "[warn] cfg: CNB_TOKEN missing" >&2
fi

echo "[info] cfg: report_config_from_api=${report_config_from_api} cabb_api_base=${cabb_api_base:-<empty>} token_present=$([ -n \"${integration_token}\" ] && echo yes || echo no)" >&2

# ------------------------------------------------------------------------------
# API Key 检测
# 支持多种变量名，按优先级查找，最终导出为 OPENCODE_API_KEY
# ------------------------------------------------------------------------------
api_key="${OPENCODE_API_KEY:-}"
if [ -z "${api_key}" ]; then api_key="${OPENCODE_TOKEN:-}"; fi
if [ -z "${api_key}" ]; then api_key="${OPENCODE_KEY:-}"; fi
if [ -z "${api_key}" ]; then api_key="${OPENCODE_ZEN_API_KEY:-}"; fi
if [ -z "${api_key}" ]; then api_key="${ZEN_API_KEY:-}"; fi
if [ -z "${api_key}" ]; then api_key="${OC_API_KEY:-}"; fi
if [ -z "${api_key}" ]; then api_key="${OPENAI_API_KEY:-}"; fi
if [ -z "${api_key}" ]; then api_key="${ANTHROPIC_API_KEY:-}"; fi
if [ -z "${api_key}" ]; then api_key="${OPENROUTER_API_KEY:-}"; fi

# ==============================================================================
# 第二部分：时间窗口计算
# ==============================================================================

# 根据 timeframe 计算时间范围和标签
report_type="daily"
case "${timeframe}" in
  yesterday)
    # 日报：昨天 00:00:00 ~ 今天 00:00:00
    start_ts=$(date -d "yesterday 00:00:00" +'%Y-%m-%d 00:00:00')
    end_ts=$(date -d "today 00:00:00" +'%Y-%m-%d 00:00:00')
    label=$(date -d "yesterday" +%Y-%m-%d)
    report_type="daily"
    ;;
  last_week)
    # 周报：上一自然周（周一至周日）
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
    # 月报：上一个自然月
    first_of_this_month=$(date +%Y-%m-01)
    start_month=$(date -d "${first_of_this_month} -1 month" +%Y-%m-01)
    start_ts="${start_month} 00:00:00"
    end_ts="${first_of_this_month} 00:00:00"
    label=$(date -d "${first_of_this_month} -1 month" +%Y-%m)
    report_type="monthly"
    ;;
  *)
    # 默认：日报
    start_ts=$(date -d "yesterday 00:00:00" +'%Y-%m-%d 00:00:00')
    end_ts=$(date -d "today 00:00:00" +'%Y-%m-%d 00:00:00')
    label=$(date -d "yesterday" +%Y-%m-%d)
    report_type="daily"
    ;;
esac

# ==============================================================================
# 第三部分：初始化与依赖检查
# ==============================================================================

# 获取当前仓库 slug
repo_slug="${CNB_REPO_SLUG:-}"
if [ -z "${repo_slug}" ]; then
  # 回退：使用本地仓库目录名
  if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    repo_slug=$(basename "$(git rev-parse --show-toplevel)")
  else
    repo_slug="repo"
  fi
fi

# 创建输出目录
mkdir -p reports tmp || true

# 输出文件路径
json_file="reports/report-${label}.json"
md_file="reports/report-${label}.md"
out_file="${md_file}"  # 兼容旧逻辑
# Git 上下文文件路径（用于 AI 输入）
ctx_file="tmp/git-context-${label}.md"

# 检查 jq 依赖（用于 JSON 解析）
if ! command -v jq >/dev/null 2>&1; then
  echo "[error] jq is required but not found" >&2
  exit 1
fi

# ==============================================================================
# 第四部分：工具函数
# ==============================================================================

# ------------------------------------------------------------------------------
# matches_exclude - 检查文件是否匹配排除模式
# 参数：$1 - 文件路径
# 返回：0=匹配（应排除），1=不匹配
# ------------------------------------------------------------------------------
matches_exclude() {
  local f="$1"
  local IFS=';'
  for pat in ${exclude_globs}; do
    [ -z "${pat}" ] && continue
    case "$f" in
      $pat) return 0 ;;
    esac
  done
  return 1
}

# ------------------------------------------------------------------------------
# redact - 脱敏处理，移除补丁中的敏感信息
# 处理：Authorization 头、Bearer Token、sk- 开头的密钥、CNB_TOKEN
# ------------------------------------------------------------------------------
redact() {
  sed -E \
    -e 's/(Authorization:)[[:space:]]*[^\r\n]*/\1 ***REDACTED***/Ig' \
    -e 's/(Bearer)[[:space:]]+[A-Za-z0-9._-]+/\1 ***REDACTED***/Ig' \
    -e 's/(sk-[A-Za-z0-9_-]{10,})/sk-***REDACTED***/g' \
    -e 's/(CNB_TOKEN=)[^\r\n& ]+/\1***REDACTED***/g'
}

# ------------------------------------------------------------------------------
# sanitize_slug - 清理仓库 slug，移除不安全字符
# 参数：$1 - 原始 slug
# 输出：清理后的 slug
# ------------------------------------------------------------------------------
sanitize_slug() {
  local s="$1"
  s="${s// /-}"
  s="${s//\//-}"
  s="${s//[^A-Za-z0-9._-]/-}"
  echo "$s"
}

# ==============================================================================
# 第五部分：API 配置加载
# ==============================================================================

# ------------------------------------------------------------------------------
# load_report_config - 从后端 API 加载汇报配置
# 功能：获取多仓库列表和发布目标配置
# ------------------------------------------------------------------------------
load_report_config() {
  local resp
  
  # 如果已手动指定仓库列表或未启用 API 配置，则跳过
  if [ -z "${report_repo_list}" ] && [ "${report_config_from_api}" = "1" ]; then
    # 检查必要的配置
    if [ -z "${cabb_api_base}" ] || [ -z "${integration_token}" ]; then
      echo "[warn] skip fetching report config: cabb_api_base='${cabb_api_base:-<empty>}' token_present=${token_present}" >&2
      return
    fi
    
    echo "[info] fetching report config from ${cabb_api_base%/}/jobs/report/config" >&2
    
    # 调用 API 获取配置
    resp=$(curl -w "\n%{http_code}" -sSL -H "Authorization: Bearer ${integration_token}" "${cabb_api_base%/}/jobs/report/config" || true)
    local http_status body
    http_status=$(printf "%s" "${resp}" | tail -n1)
    body=$(printf "%s" "${resp}" | sed '$d')
    
    if [ -n "${http_status}" ] && [ "${http_status}" != "200" ]; then
      echo "[warn] report config api status=${http_status}, body=${body}" >&2
    fi
    
    resp="${body}"
    if [ -n "${resp}" ]; then
      # 解析仓库列表
      report_repo_list=$(echo "${resp}" | jq -c '.report_repos // []' 2>/dev/null)
      
      # 解析发布配置（仅在未手动指定时使用）
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
  
  # 回退：使用当前仓库
  if [ -z "${report_repo_list}" ]; then
    report_repo_list=$(printf '[{"slug":"%s","repo_url":"","branch":""}]' "${repo_slug}")
  fi
  
  # 验证 JSON 格式
  if ! echo "${report_repo_list}" | jq empty >/dev/null 2>&1; then
    echo "[warn] invalid REPORT_REPO_LIST JSON; falling back to current repo only" >&2
    report_repo_list=$(printf '[{"slug":"%s","repo_url":"","branch":""}]' "${repo_slug}")
  fi
}

# ==============================================================================
# 第六部分：仓库准备
# ==============================================================================

# ------------------------------------------------------------------------------
# prepare_repo - 准备仓库（必要时克隆）并确保完整历史
# 参数：
#   $1 - slug（仓库标识）
#   $2 - repo_url（仓库 URL，空则使用当前仓库）
#   $3 - branch（分支名）
# 输出：仓库本地路径
# ------------------------------------------------------------------------------
prepare_repo() {
  local slug="$1"
  local repo_url="$2"
  local branch="$3"
  local repo_path=""
  
  if [ -z "${repo_url}" ]; then
    # 使用当前工作目录的仓库
    if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
      repo_path=$(git rev-parse --show-toplevel)
    fi
  else
    # 克隆远程仓库到临时目录
    repo_path="tmp/repos/${slug}"
    if [ ! -d "${repo_path}/.git" ]; then
      mkdir -p "tmp/repos"
      echo "[info] cloning ${repo_url} into ${repo_path}..." >&2
      
      # 优先使用带认证的 URL
      if [ -n "${CNB_TOKEN:-}" ]; then
        auth_url="${repo_url/https:\/\//https://cnb:${CNB_TOKEN}@}"
        git clone --no-tags "${auth_url}" "${repo_path}" >/dev/null 2>&1 || true
      fi
      
      # 回退：无认证克隆
      if [ ! -d "${repo_path}/.git" ]; then
        git clone --no-tags "${repo_url}" "${repo_path}" >/dev/null 2>&1 || true
      fi
    fi
  fi
  
  # 检查仓库是否可用
  if [ -z "${repo_path}" ] || [ ! -d "${repo_path}/.git" ]; then
    echo "[warn] repo ${slug} unavailable" >&2
    echo ""
    return
  fi
  
  # 配置 safe.directory 避免权限告警
  git -C "${repo_path}" config --global --add safe.directory "$(cd "${repo_path}" && pwd)" >/dev/null 2>&1 || true
  
  # 获取所有远程分支
  if git -C "${repo_path}" remote get-url origin >/dev/null 2>&1; then
    (git -C "${repo_path}" remote set-branches origin "*" >/dev/null 2>&1 || true)
  fi
  
  # 拉取最新数据
  (git -C "${repo_path}" fetch --all --prune --tags --recurse-submodules=no >/dev/null 2>&1 || true)
  
  # 如果是浅克隆，获取完整历史
  (git -C "${repo_path}" fetch --unshallow >/dev/null 2>&1 || git -C "${repo_path}" fetch --depth=0 >/dev/null 2>&1 || true)
  
  # 获取指定分支
  if [ -n "${branch}" ]; then
    (git -C "${repo_path}" fetch origin "${branch}" --depth=0 >/dev/null 2>&1 || true)
  fi
  
  echo "${repo_path}"
}

# ==============================================================================
# 第七部分：提交数据收集
# ==============================================================================

# 标记是否有提交
has_commits=false

# ------------------------------------------------------------------------------
# collect_repo_context - 收集单个仓库的提交上下文和补丁采样
# 参数：
#   $1 - slug（仓库标识）
#   $2 - repo_path（本地路径）
#   $3 - branch（分支名，空则统计所有分支）
#   $4 - display_name（显示名称）
#   $5 - repo_url（仓库 URL）
#   $6 - budget（字符预算）
# ------------------------------------------------------------------------------
collect_repo_context() {
  local slug="$1"
  local repo_path="$2"
  local branch="$3"
  local display_name="$4"
  local repo_url="$5"
  local budget="$6"
  local title="${display_name:-$slug}"
  local revset="--all"

  # 写入仓库基本信息
  echo "## Repo: ${title}" >> "${ctx_file}"
  echo "- slug: ${slug}" >> "${ctx_file}"
  if [ -n "${repo_url}" ]; then
    echo "- repo_url: ${repo_url}" >> "${ctx_file}"
  fi
  
  # 确定分支范围
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

  # 获取时间窗口内的提交列表
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

  # 统计提交数量
  local commit_count=0
  if [ -s "${log_file}" ]; then
    commit_count=$(wc -l < "${log_file}" | tr -d ' ')
  fi
  echo "[info] repo ${slug}: ${commit_count} commits found in time window" >&2

  # 无提交则跳过
  if [ ! -s "${log_file}" ]; then
    echo "> No commits found in this repo for the selected window." >> "${ctx_file}"
    echo >> "${ctx_file}"
    return
  fi

  has_commits=true

  # 写入提交统计
  echo "**本仓库在此时间范围内共有 ${commit_count} 条提交记录。**" >> "${ctx_file}"
  echo "" >> "${ctx_file}"

  # 按作者统计提交数
  echo "### Author counts" >> "${ctx_file}"
  awk -F "\t" '{print $2}' "${log_file}" | sed '/^$/d' | sort | uniq -c | sort -nr \
    | awk '{c=$1; $1=""; sub(/^ /, ""); printf("- %s: %d commits\n", $0, c)}' >> "${ctx_file}"
  echo >> "${ctx_file}"

  # 输出提交详情（包含 numstat）
  # 注意：使用 || true 忽略 SIGPIPE（当 head 提前退出时 git 会收到此信号）
  echo "### Commit details (numstat + message)" >> "${ctx_file}"
  if [ "${revset}" = "--all" ]; then
    git -C "${repo_path}" log --all \
      --since="${start_ts}" --until="${end_ts}" \
      --date=iso-local \
      --numstat \
      --pretty=format:'---%ncommit %H%nauthor %an <%ae>%ndate %ad%ntitle %s%nbody %b' \
      | head -n "${numstat_lines}" >> "${ctx_file}" || true
  else
    git -C "${repo_path}" log "${revset}" \
      --since="${start_ts}" --until="${end_ts}" \
      --date=iso-local \
      --numstat \
      --pretty=format:'---%ncommit %H%nauthor %an <%ae>%ndate %ad%ntitle %s%nbody %b' \
      | head -n "${numstat_lines}" >> "${ctx_file}" || true
  fi
  echo >> "${ctx_file}"

  # 如果仅需要统计模式，则到此结束
  if [ "${diff_mode}" = "stats" ]; then
    echo >> "${ctx_file}"
    return
  fi

  # --------------------------------------------------------------------------
  # 补丁采样逻辑
  # --------------------------------------------------------------------------
  
  # 获取非合并提交列表
  local commit_list=()
  if [ "${revset}" = "--all" ]; then
    mapfile -t commit_list < <(git -C "${repo_path}" rev-list --all --since="${start_ts}" --until="${end_ts}" --no-merges)
  else
    mapfile -t commit_list < <(git -C "${repo_path}" rev-list "${revset}" --since="${start_ts}" --until="${end_ts}" --no-merges)
  fi
  
  # 按变更行数对提交排序（优先采样大变更）
  local tmp_scores="tmp/commit-scores-${slug}-${label}.txt"
  : > "${tmp_scores}"
  for sha in "${commit_list[@]}"; do
    sum=$(git -C "${repo_path}" show --numstat --format="" "$sha" | awk '{a=$1;b=$2; if(a=="-")a=0; if(b=="-")b=0; s+=a+b} END{print s+0}')
    echo "$sum $sha" >> "${tmp_scores}"
  done
  mapfile -t top_commits < <(sort -nr -k1,1 "${tmp_scores}" | awk '{print $2}' | head -n "${max_commits}")

  # 采样补丁
  local total_len=0
  local patch_ctx_file="tmp/git-patch-${slug}-${label}.txt"
  : > "${patch_ctx_file}"

  for sha in "${top_commits[@]}"; do
    # 提交头信息
    hdr=$(git -C "${repo_path}" show -s --format='commit %H%nauthor %an <%ae>%ndate %ad%ntitle %s%n' --date=iso-local "$sha")
    entry_file="tmp/patch-entry-${slug}-${sha}.txt"
    printf "%s\n" "$hdr" > "$entry_file"
    
    # 获取文件列表并按变更大小排序
    mapfile -t files < <(git -C "${repo_path}" show --numstat --format="" "$sha" | awk '{a=$1;b=$2;f=$3; if(a=="-")a=0; if(b=="-")b=0; print a+b, f}' | sort -nr -k1,1 | awk '{ $1=""; sub(/^ /, ""); print }')
    
    count=0
    for f in "${files[@]}"; do
      [ -z "$f" ] && continue
      
      # 跳过排除的文件
      if matches_exclude "$f"; then
        continue
      fi
      
      # 获取文件补丁并脱敏
      file_patch=$(git -C "${repo_path}" show --no-color --unified=3 --format="" "$sha" -- "$f" | redact)
      
      if [ -n "$file_patch" ]; then
        # 采样模式下限制 hunk 数量
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
    
    # 检查字符预算
    entry_size=$(wc -c < "$entry_file" | tr -d ' ')
    if [ "$(( total_len + entry_size ))" -le "$budget" ]; then
      cat "$entry_file" >> "$patch_ctx_file"
      total_len=$(( total_len + entry_size ))
    else
      # 超出预算，截断并退出
      remain=$(( budget - total_len ))
      if [ "$remain" -gt 0 ]; then
        head -c "$remain" "$entry_file" >> "$patch_ctx_file" || true
        total_len=$budget
      fi
      break
    fi
  done
  
  # 写入补丁采样结果
  if [ -s "$patch_ctx_file" ]; then
    {
      echo "### Patch samples (truncated by budget if necessary)"
      echo "Budget: chars=${budget}, top_files=${top_files}, hunks_per_file=${max_hunks_per_file}, max_commits=${max_commits}"
      cat "$patch_ctx_file"
    } >> "$ctx_file"
    echo >> "$ctx_file"
  fi
}

# ==============================================================================
# 第八部分：主流程 - 数据收集
# ==============================================================================

# 加载 API 配置
load_report_config

# 解析仓库列表
mapfile -t repo_entries < <(echo "${report_repo_list}" | jq -c '.[]')
if [ "${#repo_entries[@]}" -eq 0 ]; then
  repo_entries=()
  repo_entries+=("$(printf '{"slug":"%s","repo_url":"","branch":""}' "${repo_slug}")")
fi

# 输出配置信息
repo_count=${#repo_entries[@]}
echo "[info] report repo entries (raw JSON): ${report_repo_list}" >&2
echo "[info] report repo count: ${repo_count}" >&2

for entry in "${repo_entries[@]}"; do
  slug_dbg=$(echo "${entry}" | jq -r '.slug // ""')
  url_dbg=$(echo "${entry}" | jq -r '.repo_url // ""')
  branch_dbg=$(echo "${entry}" | jq -r '.branch // ""')
  echo "[info] repo entry -> slug=${slug_dbg:-<empty>} url=${url_dbg:-<empty>} branch=${branch_dbg:-<all>}" >&2
done

# 计算每个仓库的字符预算
budget_per_repo="${char_budget}"
if [ "${repo_count}" -gt 0 ] && [ "${char_budget}" -gt 0 ]; then
  budget_per_repo=$(( char_budget / repo_count ))
  # 确保最小预算
  if [ "${budget_per_repo}" -lt 50000 ] && [ "${char_budget}" -ge 50000 ]; then
    budget_per_repo=50000
  fi
  if [ "${budget_per_repo}" -le 0 ]; then
    budget_per_repo="${char_budget}"
  fi
fi

# 初始化上下文文件
{
  echo "# Multi-repo Commit Context"
  echo "- Time window: ${start_ts} .. ${end_ts} (${TZ})"
  echo "- Branch (trigger): ${CNB_BRANCH:-}"
  echo "- Repos configured: ${repo_count}"
  echo "- Repo entries raw: ${report_repo_list}"
  echo
} > "${ctx_file}"

# 遍历仓库收集数据
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
  
  # 准备仓库
  repo_path=$(prepare_repo "${slug}" "${repo_url}" "${repo_branch}")
  
  if [ -z "${repo_path}" ]; then
    echo "## Repo: ${slug}" >> "${ctx_file}"
    echo "> repository unavailable; skipped" >> "${ctx_file}"
    echo >> "${ctx_file}"
    continue
  fi
  
  # 收集提交数据
  collect_repo_context "${slug}" "${repo_path}" "${repo_branch}" "${repo_display}" "${repo_url}" "${budget_per_repo}"
done

# ==============================================================================
# 第九部分：AI 汇总生成
# ==============================================================================

# ------------------------------------------------------------------------------
# try_opencode - 尝试使用 OpenCode AI 生成结构化 JSON 汇报
# 返回：0=成功，1=失败
# ------------------------------------------------------------------------------
try_opencode() {
  # 检查 opencode 是否可用
  command -v opencode >/dev/null 2>&1 || return 1

  # 设置环境变量（避免在日志中泄露密钥）
  export OPENCODE_MODEL="${model}"
  if [ -n "${api_key}" ]; then
    export OPENCODE_API_KEY="${api_key}"
  fi

  # 根据时间范围确定报告类型标签
  case "${timeframe}" in
    yesterday) period_label="日报" ;;
    last_week) period_label="周报" ;;
    last_month) period_label="月报" ;;
    *) period_label="汇报" ;;
  esac
  
  # 从上下文文件提取提交统计摘要
  commit_summary=$(grep '本仓库在此时间范围内共有.*条提交记录' "${ctx_file}" 2>/dev/null | sed 's/\*\*//g' || echo "")
  
  # JSON 模板文件路径
  local template_json="scripts/report-template.json"
  
  # 将模板文件追加到上下文
  if [ -f "${template_json}" ]; then
    echo -e "\n\n# JSON 输出模板（请严格按此结构输出）\n" >> "${ctx_file}"
    cat "${template_json}" >> "${ctx_file}"
  fi
  
  # 构建 AI Prompt（输出 JSON 格式）
  prompt="你是资深工程经理，需基于多个仓库的提交记录（含 numstat/补丁采样）生成中文 ${period_label}。\n"
  prompt+="【重要】以下是各仓库的提交统计，请务必为每个有提交的仓库生成详细总结：\n${commit_summary}\n\n"
  prompt+="时间范围：${start_ts} 至 ${end_ts}（${TZ}）。\n\n"
  prompt+="【输出格式】请严格按照上下文中提供的 JSON 模板结构输出，直接输出 JSON（不要输出 markdown 代码块标记）。\n"
  prompt+="模板展示的是多仓库场景示例，请根据实际提交数据调整：\n"
  prompt+="- 仓库数量：根据实际涉及的仓库数量调整 repos 数组长度\n"
  prompt+="- 成员数量：根据每个仓库的实际贡献者数量调整该仓库的 members 数组\n"
  prompt+="- 如果只有 1 个仓库，则 repos 数组只保留 1 个元素\n\n"
  prompt+="字段填充说明：\n"
  prompt+="- meta.type: \"${report_type}\"\n"
  prompt+="- meta.label: \"${label}\"\n"
  prompt+="- meta.time_range: 使用 \"${start_ts}\" 至 \"${end_ts}\"，timezone 为 \"${TZ}\"\n"
  prompt+="- summary: 概括所有仓库的整体进展\n"
  prompt+="- highlights: 提炼 3-5 个重要亮点，用 [display_name] 前缀标注来源（如 [Cabb 前端]）\n"
  prompt+="- repos: 每个仓库包含 slug、display_name、commit_count 和 members 数组\n"
  prompt+="- repos[].members: 该仓库下的贡献者列表\n\n"
  prompt+="要求：\n"
  prompt+="- 结论先行，语言简练\n"
  prompt+="- 每个仓库的提交数据都在对应的 Repo 部分，请逐一总结，不要遗漏\n"
  prompt+="- 确保输出是合法的 JSON，可以被 jq 解析\n"

  # 调试信息
  echo "[debug] ctx_file path: $(pwd)/${ctx_file}" >&2
  echo "[debug] ctx_file repos found: $(grep -c '^## Repo:' "${ctx_file}" 2>/dev/null || echo 0)" >&2
  echo "[debug] ctx_file author sections: $(grep -c '^### Author counts' "${ctx_file}" 2>/dev/null || echo 0)" >&2
  grep '^## Repo:' "${ctx_file}" 2>/dev/null | while read -r line; do echo "[debug] $line" >&2; done

  # 调用 OpenCode AI
  if opencode run "${prompt}" -m "${OPENCODE_MODEL}" -f "${ctx_file}" > "${json_file}.tmp" 2>"tmp/opencode.stderr"; then
    # 验证 JSON 格式
    if jq empty "${json_file}.tmp" 2>/dev/null; then
      mv "${json_file}.tmp" "${json_file}"
      echo "[info] JSON report generated: ${json_file}" >&2
      return 0
    else
      echo "[warn] AI output is not valid JSON, attempting to extract JSON..." >&2
      # 尝试从输出中提取 JSON（AI 可能输出了额外内容）
      if grep -o '{.*}' "${json_file}.tmp" | jq empty 2>/dev/null; then
        grep -o '{.*}' "${json_file}.tmp" > "${json_file}"
        echo "[info] Extracted valid JSON from AI output" >&2
        return 0
      fi
      echo "[error] Failed to extract valid JSON from AI output" >&2
      cat "${json_file}.tmp" >&2
      return 1
    fi
  fi
  
  # 输出错误信息（不泄露密钥）
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

# ------------------------------------------------------------------------------
# render_markdown - 将 JSON 转换为 Markdown 格式
# 参数：$1 - JSON 文件路径，$2 - 输出 Markdown 文件路径
# ------------------------------------------------------------------------------
render_markdown() {
  local json_in="$1"
  local md_out="$2"
  
  if [ ! -f "${json_in}" ]; then
    echo "[error] JSON file not found: ${json_in}" >&2
    return 1
  fi
  
  # 使用 jq 生成完整的 Markdown
  jq -r '
    "# 项目工作汇报（" + .meta.label + "）\n",
    "- 时间范围：" + .meta.time_range.start + " 至 " + .meta.time_range.end + "（" + .meta.time_range.timezone + "）",
    "- 仓库：" + ([.repos[].slug] | join("、")) + "\n",
    "## 概览\n",
    .summary + "\n",
    (if (.highlights | length) > 0 then
      "### 亮点\n\n" + ([.highlights[] | "- " + .] | join("\n")) + "\n"
    else "" end),
    "## 仓库汇总\n",
    (.repos[] | 
      "### [" + .slug + "] " + .display_name + "\n",
      (if .brief then "**简述**：" + .brief + "\n" else "" end),
      "**主要贡献者**：\n",
      (.members[] |
        "#### " + .name + "（" + (.commits | tostring) + " 次提交）\n",
        (.achievements[] | "- " + .),
        ""
      ),
      (if .impact then "**影响与价值**：" + .impact + "\n" else "" end)
    )
  ' "${json_in}" > "${md_out}"
  
  echo "[info] Markdown report generated: ${md_out}" >&2
  return 0
}

# ==============================================================================
# 第十部分：主流程 - AI 生成与多格式输出
# ==============================================================================

# 尝试 AI 汇总生成 JSON
ai_success=false
if try_opencode; then
  echo "[info] AI JSON summary generated." >&2
  ai_success=true
else
  echo "[warn] AI summary generation failed." >&2
fi

# 如果 AI 成功，进行格式转换
if [ "${ai_success}" = "true" ] && [ -f "${json_file}" ]; then
  
  # 转换为 Markdown
  if [ "${output_markdown}" = "1" ]; then
    if render_markdown "${json_file}" "${md_file}"; then
      out_file="${md_file}"
    fi
  fi
  
else
  # AI 失败的回退处理
  if [ "${require_ai}" = "1" ]; then
    {
      echo "# 项目工作汇报（${label}）"
      echo ""
      echo "AI 汇总生成失败。请检查 OPENCODE_API_KEY/网络后重试。"
    } > "${md_file}"
    out_file="${md_file}"
    exit 1
  else
    {
      echo "# 项目工作汇报（${label}）"
      echo ""
      if ${has_commits}; then
        echo "AI 汇总生成失败，但检测到提交记录。"
      else
        echo "无提交记录。"
      fi
    } > "${md_file}"
    out_file="${md_file}"
  fi
fi

echo "[info] Report generated: ${out_file}" >&2

# ==============================================================================
# 第十部分：汇报发布
# ==============================================================================

# ------------------------------------------------------------------------------
# publish_report - 发布汇报到目标仓库
# 目录结构：{publish_dir_root}/{daily|weekly|monthly}/report-{label}.md
# ------------------------------------------------------------------------------
publish_report() {
  # 检查必要条件
  if [ -z "${publish_repo_url}" ] || [ -z "${CNB_TOKEN:-}" ]; then
    return 0
  fi

  # 根据报告类型计算目标路径和文件名
  target_subdir="${publish_dir_root}/${report_type}"
  case "${report_type}" in
    daily)
      target_filename="report-${label}.md"
      range_text="${label}"
      ;;
    weekly)
      if [ -n "${week_start_date:-}" ] && [ -n "${week_end_date:-}" ]; then
        target_filename="report-${week_start_date}_to_${week_end_date}.md"
        range_text="${week_start_date}~${week_end_date}"
      else
        target_filename="report-${label}.md"
        range_text="${label}"
      fi
      ;;
    monthly)
      target_filename="report-${label}.md"
      range_text="${label}"
      ;;
    *)
      target_filename="report-${label}.md"
      range_text="${label}"
      ;;
  esac

  # 克隆目标仓库
  workdir="tmp/publish"
  rm -rf "${workdir}" && mkdir -p "${workdir}"
  auth_url="${publish_repo_url/https:\/\//https://cnb:${CNB_TOKEN}@}"
  git clone "${auth_url}" "${workdir}" >/dev/null 2>&1 || return 0
  
  cd "${workdir}"
  git config --global --add safe.directory "$(pwd)" || true
  
  # 切换到目标分支（不存在则创建）
  if git show-ref --verify --quiet "refs/heads/${publish_branch}"; then
    git checkout "${publish_branch}" >/dev/null 2>&1 || true
  else
    git checkout -b "${publish_branch}" >/dev/null 2>&1 || true
  fi
  
  # 拉取最新变更（避免 non-fast-forward）
  git pull --rebase origin "${publish_branch}" >/dev/null 2>&1 || true
  
  # 复制汇报文件
  mkdir -p "${target_subdir}"
  cp -f "${OLDPWD}/${out_file}" "${target_subdir}/${target_filename}"
  
  # 同时发布 JSON 文件（后端飞书推送需要）
  json_target_filename="${target_filename%.md}.json"
  if [ -f "${OLDPWD}/${json_file}" ]; then
    cp -f "${OLDPWD}/${json_file}" "${target_subdir}/${json_target_filename}"
    echo "[info] JSON file published: ${target_subdir}/${json_target_filename}" >&2
  fi
  
  # 配置 Git 用户
  git config user.name "cabb-report-bot"
  git config user.email "bot@cabb.local"
  
  # 提交变更
  git add -A "${target_subdir}" || true
  if git diff --cached --quiet >/dev/null 2>&1; then
    # 无变更
    cd - >/dev/null 2>&1 || true
    return 0
  fi
  
  git commit -m "chore(report): ${report_type} ${range_text}" >/dev/null 2>&1 || true
  
  # 推送（失败则重试一次）
  if ! git push origin "${publish_branch}" >/dev/null 2>&1; then
    git pull --rebase origin "${publish_branch}" >/dev/null 2>&1 || true
    git push origin "${publish_branch}" >/dev/null 2>&1 || true
  fi
  
  cd - >/dev/null 2>&1 || true
}

# 执行发布
if [ "${publish_enable}" = "1" ] || [ "${publish_enable}" = "true" ]; then
  publish_report || true
fi
