#!/bin/bash
# 测试 issue label notify API

set -e

# 默认配置
BASE_URL="${BASE_URL:-http://localhost:8080}"
TOKEN="${INTEGRATION_TOKEN:-test-token-123}"

# 测试请求体
REQUEST_BODY='{
  "repo_slug": "1024hub/Demo",
  "issue_number": 74,
  "issue_url": "https://cnb.cool/1024hub/Demo/-/issues/74",
  "title": "实现用户登录功能",
  "state": "open",
  "author": {
    "username": "zhangsan",
    "nickname": "张三"
  },
  "description": "需要实现用户登录功能，包括账号密码登录和第三方登录",
  "labels": ["🚧 处理中_CNB", "🧑🏻‍💻 进行中：前端_CNB"],
  "label_trigger": "🚧 处理中_CNB",
  "updated_at": "2025-10-29T03:25:06Z",
  "event_context": {
    "event_type": "push",
    "branch": "feature/74-user-login"
  }
}'

echo "=== 测试 Issue Label Notify API ==="
echo "URL: ${BASE_URL}/api/v1/issues/label-notify"
echo ""

# 发送请求
echo "发送请求..."
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "${BASE_URL}/api/v1/issues/label-notify" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "X-Delivery-ID: test-$(date +%s)" \
  -d "$REQUEST_BODY")

# 提取 HTTP 状态码
HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_CODE:/d')

echo ""
echo "HTTP 状态码: $HTTP_CODE"
echo "响应体:"
echo "$BODY" | jq . 2>/dev/null || echo "$BODY"

# 验证结果
if [ "$HTTP_CODE" = "200" ]; then
    echo ""
    echo "✅ 测试通过"
    exit 0
else
    echo ""
    echo "❌ 测试失败"
    exit 1
fi
