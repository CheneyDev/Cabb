# CNB Job 集成指南

## 概述

本文档描述如何将 CNB 的 `job-get-issues-info` 与 Plane Integration 服务集成，实现 issue 标签变更的自动通知与处理。

## 背景

CNB 的 `job-get-issues-info` 负责监控 issue 的变动（特别是标签变更），原本设计为向飞书机器人发送通知。本集成将通知目标改为 Plane Integration 后端服务，由后端统一处理业务逻辑（如同步到 Plane、发送飞书通知等）。

## API 端点

### POST /api/v1/issues/label-notify

接收来自 CNB job 的 issue 标签变更通知。

**请求头：**
- `Content-Type: application/json`
- `Authorization: Bearer <INTEGRATION_TOKEN>`
- `X-Delivery-ID: <可选，用于幂等性>`

**请求体：**
```json
{
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
}
```

**字段说明：**

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| repo_slug | string | 是 | 仓库标识（格式：owner/repo） |
| issue_number | integer | 是 | Issue 编号 |
| issue_url | string | 是 | Issue 完整访问地址 |
| title | string | 是 | Issue 标题 |
| state | string | 是 | Issue 状态（open/closed） |
| author.username | string | 是 | Issue 作者用户名 |
| author.nickname | string | 是 | Issue 作者昵称 |
| description | string | 否 | Issue 描述内容 |
| labels | []string | 是 | Issue 当前所有标签列表 |
| label_trigger | string | 是 | 触发本次通知的标签 |
| updated_at | string | 是 | Issue 最后更新时间（ISO 8601） |
| event_context.event_type | string | 否 | 触发事件类型 |
| event_context.branch | string | 否 | 触发事件的分支名称 |

**成功响应（200 OK）：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "issue_number": 74,
    "processed_at": "2025-10-29T03:25:10Z"
  }
}
```

**错误响应示例：**

```json
{
  "error": {
    "code": "invalid_token",
    "message": "鉴权失败（Bearer token 不匹配）",
    "details": {},
    "request_id": "req-12345"
  }
}
```

**错误码：**

| 状态码 | code | 说明 |
|--------|------|------|
| 400 | missing_fields | 缺少必填字段 |
| 400 | invalid_body | 请求体读取失败 |
| 401 | invalid_token | 鉴权失败 |
| 422 | invalid_json | JSON 解析失败 |
| 500 | - | 服务器内部错误 |

## CNB Job 修改

### 原始实现（飞书通知）

原始代码使用 Docker 容器调用飞书机器人：

```bash
docker run --rm \
  -e WEBHOOK_URL="$WEBHOOK_URL" \
  -e SIGN_SECRET="$SIGN_SECRET" \
  registry.cnb.cool/plugins/feishu-robot:latest \
  send --card "$CARD_JSON"
```

### 新实现（调用后端 API）

替换为 curl 调用 Go 后端：

```bash
# 构建 JSON 请求体
REQUEST_BODY=$(jq -n \
  --arg repo_slug "$CNB_REPO_SLUG" \
  --arg issue_number "$issue_number" \
  --arg issue_url "$EVENT_URL" \
  --arg issue_title "$ISSUE_TITLE" \
  --arg issue_state "$ISSUE_STATE" \
  --arg author_username "$ISSUE_AUTHOR_USERNAME" \
  --arg author_nickname "$ISSUE_AUTHOR" \
  --arg description "$ISSUE_DESCRIPTION" \
  --argjson labels "$(echo \"$ALL_LABELS\" | jq -R 'split(\", \")')" \
  --arg label_trigger "$addLabel" \
  --arg updated_at "$UPDATED_AT" \
  --arg event_type "${CNB_EVENT_TYPE:-unknown}" \
  --arg branch "${CNB_BRANCH:-}" \
  '{
    "repo_slug": $repo_slug,
    "issue_number": ($issue_number | tonumber),
    "issue_url": $issue_url,
    "title": $issue_title,
    "state": $issue_state,
    "author": {
      "username": $author_username,
      "nickname": $author_nickname
    },
    "description": $description,
    "labels": $labels,
    "label_trigger": $label_trigger,
    "updated_at": $updated_at,
    "event_context": {
      "event_type": $event_type,
      "branch": $branch
    }
  }')

# 发送 HTTP 请求到 Go 后端
HTTP_CODE=$(curl -s -w "%{http_code}" -o /tmp/response.json \
  -X POST "${GO_BACKEND_URL}/api/v1/issues/label-notify" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${GO_BACKEND_TOKEN}" \
  -H "X-Delivery-ID: cnb-job-${CNB_PIPELINE_ID:-unknown}-$(date +%s)" \
  -d "$REQUEST_BODY")

if [ "$HTTP_CODE" = "200" ]; then
  echo "✅ 通知发送成功"
  cat /tmp/response.json | jq .
else
  echo "❌ 通知发送失败 (HTTP $HTTP_CODE)"
  cat /tmp/response.json
  exit 1
fi
```

### 环境变量配置

需要在 CNB 流水线中配置以下环境变量：

- `GO_BACKEND_URL` - Go 后端服务地址（如 `https://your-domain.com`）
- `GO_BACKEND_TOKEN` - 与 Go 后端 `INTEGRATION_TOKEN` 一致的认证令牌

## 幂等性保证

服务实现了双重幂等性保证：

1. **内存级去重**：基于 `X-Delivery-ID` + `payload_sha256`，5 分钟内的重复请求直接返回成功
2. **数据库级去重**：将事件记录到 `event_deliveries` 表，持久化防重

建议在调用时设置 `X-Delivery-ID` 请求头（格式：`cnb-job-{pipeline_id}-{timestamp}`）。

## 安全配置

### 生成 INTEGRATION_TOKEN

```bash
# macOS/Linux
openssl rand -hex 32

# 或使用 Python
python3 -c 'import secrets; print(secrets.token_hex(32))'
```

### 在 Go 后端配置

```bash
# .env
INTEGRATION_TOKEN=your-generated-token-here
```

### 在 CNB 流水线配置

在 CNB 控制台的"仓库设置 → 流水线密钥"中添加：
- 密钥名称：`GO_BACKEND_TOKEN`
- 密钥值：与 Go 后端 `INTEGRATION_TOKEN` 相同

## 测试

### 使用测试脚本

```bash
cd /path/to/plane-integration
export INTEGRATION_TOKEN=your-token
export BASE_URL=http://localhost:8080

./scripts/test-label-notify.sh
```

### 手动测试

```bash
curl -X POST "http://localhost:8080/api/v1/issues/label-notify" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-token" \
  -H "X-Delivery-ID: test-$(date +%s)" \
  -d '{
    "repo_slug": "test/repo",
    "issue_number": 1,
    "issue_url": "https://cnb.cool/test/repo/-/issues/1",
    "title": "测试 Issue",
    "state": "open",
    "author": {"username": "test", "nickname": "测试用户"},
    "labels": ["bug"],
    "label_trigger": "bug",
    "updated_at": "2025-01-30T00:00:00Z"
  }'
```

## 后续扩展

当前实现为占位处理器（placeholder），后续可扩展以下功能：

1. **同步到 Plane**：根据 repo-project 映射查找对应的 Plane 项目，同步标签变更
2. **发送飞书通知**：根据 channel-project 映射，向绑定的飞书群组发送标签变更卡片
3. **触发自动化**：基于特定标签触发工作流（如"待评审"标签自动分配审核人）
4. **数据分析**：记录标签变更历史，用于项目进度分析

实现位置：`internal/handlers/issue_label_notify.go` 中的 `processIssueLabelNotify` 方法。

## 日志与监控

服务会记录结构化日志，包含以下关键字段：

- `request_id` - 请求唯一标识
- `source` - 固定为 `issue.label.notify`
- `endpoint` - API 路径
- `status` - HTTP 状态码
- `result` - 处理结果（success/error）
- `latency_ms` - 处理耗时

可通过日志聚合工具（如 ELK）监控 API 调用情况。

## 故障排查

### 401 Unauthorized

检查：
1. `Authorization` 请求头格式是否为 `Bearer <token>`
2. Go 后端 `INTEGRATION_TOKEN` 与 CNB `GO_BACKEND_TOKEN` 是否一致
3. 令牌是否包含特殊字符需要转义

### 400 Bad Request

检查：
1. 请求体 JSON 格式是否正确
2. 必填字段是否完整
3. `issue_number` 是否为正整数
4. `labels` 数组是否为空

### 重复处理

如果担心重复处理，确保设置 `X-Delivery-ID` 请求头。相同 delivery_id + payload_sha256 的请求会被自动去重。
