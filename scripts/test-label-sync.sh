#!/bin/bash
# 测试简化版 API

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"
TOKEN="${INTEGRATION_TOKEN:-test-token-123}"

echo "=== 测试简化版 Issue Label Sync API ==="
echo "URL: ${BASE_URL}/api/v1/issues/label-sync"
echo ""

# 简化请求体（只需 3 个字段）
REQUEST_BODY='{
  "repo_slug": "1024hub/Demo",
  "issue_number": 74,
  "labels": ["🚧 处理中_CNB", "🧑🏻‍💻 进行中：前端_CNB", "bug", "feature"]
}'

echo "发送简化版请求..."
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "${BASE_URL}/api/v1/issues/label-sync" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "X-Delivery-ID: test-simple-$(date +%s)" \
  -d "$REQUEST_BODY")

HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE:" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_CODE:/d')

echo ""
echo "HTTP 状态码: $HTTP_CODE"
echo "响应体:"
echo "$BODY" | jq . 2>/dev/null || echo "$BODY"

if [ "$HTTP_CODE" = "200" ]; then
    echo ""
    echo "✅ 简化版 API 测试通过"
    exit 0
else
    echo ""
    echo "❌ 简化版 API 测试失败"
    exit 1
fi
