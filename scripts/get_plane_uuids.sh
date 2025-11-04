#!/bin/bash
# 从 Plane 实例获取 Project 和 Workspace UUID
# 用于配置 repo_project_mappings 映射

set -e

echo "=== 从 Plane 获取配置信息 ==="
echo ""

# 从 .env 加载配置
if [ -f .env ]; then
    export $(grep -v '^#' .env | grep -v '^$' | xargs)
fi

PLANE_BASE_URL="${PLANE_BASE_URL:-https://work.1024hub.org:4430/api}"
echo "Plane API: $PLANE_BASE_URL"
echo ""

# 检查是否配置了 Service Token
if [ -z "$PLANE_SERVICE_TOKEN" ]; then
    echo "❌ 错误：未配置 PLANE_SERVICE_TOKEN"
    echo ""
    echo "获取 Service Token 的步骤："
    echo "1. 访问 Plane 实例：https://work.1024hub.org:4430"
    echo "2. 进入个人设置 → API Tokens"
    echo "3. 创建新 Token（权限：至少需要 project:read）"
    echo "4. 添加到 .env 文件："
    echo "   PLANE_SERVICE_TOKEN=plane_api_xxxxxxxxxxxxx"
    echo ""
    echo "或者，从浏览器手动查找："
    echo "1. 访问项目页面"
    echo "2. URL 格式：https://work.1024hub.org:4430/<workspace_slug>/projects/<project_id>/..."
    echo "3. 打开开发者工具（F12）→ Network → 查看 API 请求"
    echo "4. 找到包含 'projects' 的请求，查看响应中的 'id' 字段"
    exit 1
fi

echo "✓ 已配置 PLANE_SERVICE_TOKEN"
echo ""

# 1. 获取所有 Workspaces
echo "=== 1. 获取 Workspaces ==="
WORKSPACE_RESPONSE=$(curl -s -X GET "$PLANE_BASE_URL/workspaces/" \
    -H "x-api-key: $PLANE_SERVICE_TOKEN" 2>&1)

if echo "$WORKSPACE_RESPONSE" | grep -q "error\|Error\|<!DOCTYPE"; then
    echo "❌ 请求失败，响应："
    echo "$WORKSPACE_RESPONSE" | head -20
    exit 1
fi

echo "$WORKSPACE_RESPONSE" | jq -r '.[] | "  • \(.slug) (ID: \(.id), Name: \(.name))"' 2>/dev/null || {
    echo "响应格式："
    echo "$WORKSPACE_RESPONSE" | head -20
}
echo ""

# 2. 获取指定 Workspace 的项目（默认 my-test）
WORKSPACE_SLUG="${1:-my-test}"
echo "=== 2. 获取 Workspace '$WORKSPACE_SLUG' 的项目 ==="

PROJECT_RESPONSE=$(curl -s -X GET "$PLANE_BASE_URL/workspaces/$WORKSPACE_SLUG/projects/" \
    -H "x-api-key: $PLANE_SERVICE_TOKEN" 2>&1)

if echo "$PROJECT_RESPONSE" | grep -q "error\|Error\|<!DOCTYPE"; then
    echo "❌ 请求失败，可能的原因："
    echo "  - workspace_slug '$WORKSPACE_SLUG' 不存在"
    echo "  - Service Token 权限不足"
    echo ""
    echo "响应："
    echo "$PROJECT_RESPONSE" | head -20
    exit 1
fi

echo "$PROJECT_RESPONSE" | jq -r '.results[]? // .[] | "  • \(.identifier) - \(.name)\n    ID: \(.id)\n    Workspace ID: \(.workspace)"' 2>/dev/null || {
    echo "响应格式："
    echo "$PROJECT_RESPONSE" | head -20
}
echo ""

# 3. 生成 SQL 模板
echo "=== 3. 生成 SQL 配置 ==="
echo ""
echo "复制以下 SQL 到 scripts/fix_be_test_issue_mapping.sql 并替换占位符："
echo ""

# 尝试自动提取 test-notify 项目
PROJECT_ID=$(echo "$PROJECT_RESPONSE" | jq -r '.results[]? // .[] | select(.name=="test-notify" or .identifier=="test-notify") | .id' 2>/dev/null | head -1)
WORKSPACE_ID=$(echo "$WORKSPACE_RESPONSE" | jq -r ".[] | select(.slug==\"$WORKSPACE_SLUG\") | .id" 2>/dev/null | head -1)

if [ -n "$PROJECT_ID" ] && [ "$PROJECT_ID" != "null" ] && [ -n "$WORKSPACE_ID" ] && [ "$WORKSPACE_ID" != "null" ]; then
    echo "-- ✓ 已自动填充 UUID"
    cat << EOF
INSERT INTO repo_project_mappings (
  plane_project_id,
  plane_workspace_id,
  cnb_repo_id,
  workspace_slug,
  active,
  sync_direction,
  created_at,
  updated_at
) VALUES (
  '$PROJECT_ID',
  '$WORKSPACE_ID',
  '1024hub/Demo/BE-test-issue',
  '$WORKSPACE_SLUG',
  true,
  'cnb_to_plane',
  now(),
  now()
)
ON CONFLICT (plane_project_id, cnb_repo_id) DO UPDATE
SET active = true, updated_at = now();
EOF
else
    echo "-- ⚠ 需要手动替换 <占位符>"
    cat << EOF
INSERT INTO repo_project_mappings (
  plane_project_id,
  plane_workspace_id,
  cnb_repo_id,
  workspace_slug,
  active,
  sync_direction,
  created_at,
  updated_at
) VALUES (
  '<从上面的项目列表中复制 ID>',
  '<从上面的 Workspace 列表中复制 ID>',
  '1024hub/Demo/BE-test-issue',
  '$WORKSPACE_SLUG',
  true,
  'cnb_to_plane',
  now(),
  now()
)
ON CONFLICT (plane_project_id, cnb_repo_id) DO UPDATE
SET active = true, updated_at = now();
EOF
fi

echo ""
echo "=== 4. 获取标签 UUID ==="
if [ -n "$PROJECT_ID" ] && [ "$PROJECT_ID" != "null" ]; then
    echo "查询项目标签..."
    LABEL_RESPONSE=$(curl -s -X GET "$PLANE_BASE_URL/workspaces/$WORKSPACE_SLUG/projects/$PROJECT_ID/labels/" \
        -H "x-api-key: $PLANE_SERVICE_TOKEN" 2>&1)
    
    echo "$LABEL_RESPONSE" | jq -r '.[] | "  • \(.name) (ID: \(.id))"' 2>/dev/null || {
        echo "需要手动查询标签，执行："
        echo "  curl -H 'x-api-key: \$PLANE_SERVICE_TOKEN' \\"
        echo "    '$PLANE_BASE_URL/workspaces/$WORKSPACE_SLUG/projects/<project_id>/labels/'"
    }
else
    echo "需要先获取 Project ID"
fi

echo ""
echo "=== ✓ 完成 ==="
echo "下一步："
echo "1. 将上面的 SQL 保存到文件"
echo "2. 执行：psql \"\$DATABASE_URL\" -c \"<SQL>\""
echo "3. 创建标签映射（label_mappings 表）"
