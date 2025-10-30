# Issue 标签同步 API 参考手册

## 📋 概述

本 API 用于接收 CNB `job-get-issues-info` 的 issue 标签变更通知，并同步到 Plane 项目。

---

## 🔗 API 端点

### 1. 简化版 API（推荐）

**端点：** `POST /api/v1/issues/label-sync`

**用途：** 快速标签同步，只需 3 个核心字段

**请求头：**
```
Content-Type: application/json
Authorization: Bearer <INTEGRATION_TOKEN>
X-Delivery-ID: <可选，用于幂等性>
```

**请求体：**
```json
{
  "repo_slug": "1024hub/Demo",
  "issue_number": 74,
  "labels": ["🚧 处理中_CNB", "bug_CNB"]
}
```

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "issue_number": 74,
    "processed_at": "2025-10-30T07:53:09Z"
  }
}
```

---

### 2. 完整版 API

**端点：** `POST /api/v1/issues/label-notify`

**用途：** 完整事件记录，包含 11 个字段

**请求头：**
```
Content-Type: application/json
Authorization: Bearer <INTEGRATION_TOKEN>
X-Delivery-ID: <可选，用于幂等性>
```

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

**响应：**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "issue_number": 74,
    "processed_at": "2025-10-30T07:53:09Z"
  }
}
```

---

## 📝 字段说明

### 简化版（3 个字段）

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| repo_slug | string | 是 | 仓库标识（格式：owner/repo） |
| issue_number | integer | 是 | Issue 编号 |
| labels | []string | 是 | 标签列表（建议以 _CNB 结尾） |

### 完整版（11 个字段）

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| repo_slug | string | 是 | 仓库标识 |
| issue_number | integer | 是 | Issue 编号 |
| issue_url | string | 是 | Issue 完整 URL |
| title | string | 是 | Issue 标题 |
| state | string | 是 | 状态（open/closed） |
| author | object | 是 | 作者信息 |
| author.username | string | 是 | 作者用户名 |
| author.nickname | string | 是 | 作者昵称 |
| description | string | 否 | Issue 描述 |
| labels | []string | 是 | 标签列表 |
| label_trigger | string | 是 | 触发标签 |
| updated_at | string | 是 | 更新时间（RFC3339） |
| event_context | object | 否 | 事件上下文 |
| event_context.event_type | string | 否 | 事件类型 |
| event_context.branch | string | 否 | 分支名称 |

---

## ⚠️ 错误响应

### 401 Unauthorized - 鉴权失败

```json
{
  "error": {
    "code": "invalid_token",
    "message": "鉴权失败（Bearer token 不匹配）",
    "details": {},
    "request_id": "..."
  }
}
```

### 400 Bad Request - 缺少必填字段

```json
{
  "error": {
    "code": "missing_fields",
    "message": "缺少必填字段：repo_slug/issue_number/labels",
    "details": {},
    "request_id": "..."
  }
}
```

### 422 Unprocessable Entity - JSON 解析失败

```json
{
  "error": {
    "code": "invalid_json",
    "message": "JSON 解析失败",
    "details": {
      "error": "invalid character 'i' looking for beginning of value"
    },
    "request_id": "..."
  }
}
```

---

## 🧪 测试示例

### 成功案例

```bash
curl -X POST "http://localhost:8080/api/v1/issues/label-sync" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-integration-token" \
  -H "X-Delivery-ID: test-$(date +%s)" \
  -d '{
    "repo_slug": "1024hub/Demo",
    "issue_number": 74,
    "labels": ["🚧 处理中_CNB", "bug_CNB"]
  }'
```

### 测试鉴权失败（401）

```bash
curl -X POST "http://localhost:8080/api/v1/issues/label-sync" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer WRONG_TOKEN" \
  -d '{"repo_slug": "test/repo", "issue_number": 1, "labels": ["test_CNB"]}'
```

### 测试缺少字段（400）

```bash
curl -X POST "http://localhost:8080/api/v1/issues/label-sync" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-integration-token" \
  -d '{"repo_slug": "test/repo"}'
```

### 测试无效 JSON（422）

```bash
curl -X POST "http://localhost:8080/api/v1/issues/label-sync" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-integration-token" \
  -d 'invalid json here'
```

### 测试幂等性（重复请求）

```bash
DELIVERY_ID="test-idempotent-$(date +%s)"

# 第一次请求
curl -X POST "http://localhost:8080/api/v1/issues/label-sync" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-integration-token" \
  -H "X-Delivery-ID: $DELIVERY_ID" \
  -d '{"repo_slug": "test/repo", "issue_number": 99, "labels": ["test_CNB"]}'

# 第二次请求（应该返回 status: "duplicate"）
curl -X POST "http://localhost:8080/api/v1/issues/label-sync" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-integration-token" \
  -H "X-Delivery-ID: $DELIVERY_ID" \
  -d '{"repo_slug": "test/repo", "issue_number": 99, "labels": ["test_CNB"]}'
```

### 使用 jq 格式化响应

```bash
curl -s -X POST "http://localhost:8080/api/v1/issues/label-sync" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your-integration-token" \
  -H "X-Delivery-ID: test-$(date +%s)" \
  -d '{
    "repo_slug": "test/repo",
    "issue_number": 1,
    "labels": ["bug_CNB"]
  }' | jq .
```

---

## 🔄 CNB Job 集成

### 修改 CNB Job 代码

将原有的飞书通知替换为 API 调用：

**原代码（131-144 行）：**
```bash
docker run --rm -e WEBHOOK_URL="$WEBHOOK_URL" ...
```

**新代码：**
```bash
# 构建 JSON 请求体
REQUEST_BODY=$(jq -n \
  --arg repo_slug "$CNB_REPO_SLUG" \
  --arg issue_number "$issue_number" \
  --argjson labels "$(echo $ALL_LABELS | jq -R 'split(\", \")')" \
  '{
    "repo_slug": $repo_slug,
    "issue_number": ($issue_number | tonumber),
    "labels": $labels
  }')

# 发送 HTTP 请求到 Go 后端
curl -X POST "${GO_BACKEND_URL}/api/v1/issues/label-sync" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${GO_BACKEND_TOKEN}" \
  -d "$REQUEST_BODY"
```

### 环境变量配置

**Go 后端（.env）：**
```bash
INTEGRATION_TOKEN=<生成的随机令牌>
```

**CNB 流水线密钥：**
- 密钥名称：`GO_BACKEND_TOKEN`
- 密钥值：与 `INTEGRATION_TOKEN` 相同
- 密钥名称：`GO_BACKEND_URL`
- 密钥值：Go 后端地址（如 `https://api.example.com`）

**生成安全令牌：**
```bash
openssl rand -hex 32
```

---

## 🔐 安全性

### Bearer Token 认证

- 使用 `Authorization: Bearer <token>` 请求头
- Token 与 `.env` 中的 `INTEGRATION_TOKEN` 一致
- 建议使用至少 32 字节的随机字符串

### 幂等性保证

- 基于 `X-Delivery-ID` + `payload_sha256` 去重
- 内存级：5 分钟 TTL
- 数据库级：持久化到 `event_deliveries` 表（如有数据库）
- 重复请求返回 200 并标记 `status: "duplicate"`

### 请求体大小限制

- 默认限制：2MB（Echo 框架）
- 建议：保持请求体简洁，避免传输大量数据

---

## 📊 业务逻辑

### 处理流程

```
接收请求 → Bearer 鉴权 → 解析 JSON → 字段校验 
    ↓
内存去重 → 数据库去重 → 记录事件
    ↓
异步处理（立即返回 200 OK）
    ↓
过滤 _CNB 标签 → 查询映射关系 → 增量更新 Plane → 飞书通知
```

### 标签更新策略（增量更新）

1. 从 Plane API 读取 Issue 当前所有标签
2. 从 `label_mappings` 表查询哪些标签是 CNB 管理的
3. 过滤出非 CNB 管理的标签（需要保留）
4. 合并：保留的标签 + 新的 CNB 标签
5. 去重后更新到 Plane

**示例：**
- Plane 当前标签：`["priority:high", "🚧 处理中_CNB", "bug"]`
- CNB 管理的标签：`["🚧 处理中_CNB"]`
- CNB 新标签：`["✅ 已完成_CNB"]`
- **最终结果：** `["priority:high", "bug", "✅ 已完成_CNB"]`

---

## 📚 相关资源

- **开发者指南**：`.vscode/DEVELOPER-GUIDE.md`
- **测试脚本**：`scripts/test-label-sync.sh`、`scripts/test-label-notify.sh`
- **架构文档**：`docs/ARCHITECTURE.md`
- **本地测试**：`docs/LOCAL-TESTING-GUIDE.md`
