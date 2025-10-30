# API 测试示例

## 健康检查

```bash
curl http://localhost:8080/healthz
```

## Issue 标签通知（成功案例）

```bash
curl -X POST "http://localhost:8080/api/v1/issues/label-notify" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-integration-token" \
  -H "X-Delivery-ID: test-$(date +%s)" \
  -d '{
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
```

## 缺少认证（401 错误）

```bash
curl -X POST "http://localhost:8080/api/v1/issues/label-notify" \
  -H "Content-Type: application/json" \
  -d '{
    "repo_slug": "test/repo",
    "issue_number": 1,
    "issue_url": "https://cnb.cool/test/repo/-/issues/1",
    "title": "测试",
    "state": "open",
    "author": {"username": "test", "nickname": "测试"},
    "labels": ["bug"],
    "label_trigger": "bug",
    "updated_at": "2025-01-30T00:00:00Z"
  }'
```

## 缺少必填字段（400 错误）

```bash
curl -X POST "http://localhost:8080/api/v1/issues/label-notify" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-integration-token" \
  -d '{
    "repo_slug": "test/repo",
    "issue_number": 1
  }'
```

## JSON 格式错误（422 错误）

```bash
curl -X POST "http://localhost:8080/api/v1/issues/label-notify" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-integration-token" \
  -d 'invalid json'
```

## 幂等性测试（重复请求）

```bash
DELIVERY_ID="test-idempotent-123"

# 第一次请求
curl -X POST "http://localhost:8080/api/v1/issues/label-notify" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-integration-token" \
  -H "X-Delivery-ID: $DELIVERY_ID" \
  -d '{
    "repo_slug": "test/repo",
    "issue_number": 99,
    "issue_url": "https://cnb.cool/test/repo/-/issues/99",
    "title": "幂等性测试",
    "state": "open",
    "author": {"username": "test", "nickname": "测试"},
    "labels": ["test"],
    "label_trigger": "test",
    "updated_at": "2025-01-30T00:00:00Z"
  }'

# 第二次请求（应该返回 duplicate 状态）
curl -X POST "http://localhost:8080/api/v1/issues/label-notify" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-integration-token" \
  -H "X-Delivery-ID: $DELIVERY_ID" \
  -d '{
    "repo_slug": "test/repo",
    "issue_number": 99,
    "issue_url": "https://cnb.cool/test/repo/-/issues/99",
    "title": "幂等性测试",
    "state": "open",
    "author": {"username": "test", "nickname": "测试"},
    "labels": ["test"],
    "label_trigger": "test",
    "updated_at": "2025-01-30T00:00:00Z"
  }'
```

## 使用 jq 格式化响应

```bash
curl -s -X POST "http://localhost:8080/api/v1/issues/label-notify" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-integration-token" \
  -H "X-Delivery-ID: test-$(date +%s)" \
  -d '{
    "repo_slug": "test/repo",
    "issue_number": 1,
    "issue_url": "https://cnb.cool/test/repo/-/issues/1",
    "title": "测试",
    "state": "open",
    "author": {"username": "test", "nickname": "测试"},
    "labels": ["bug"],
    "label_trigger": "bug",
    "updated_at": "2025-01-30T00:00:00Z"
  }' | jq .
```
