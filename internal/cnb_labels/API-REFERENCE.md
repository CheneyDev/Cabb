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

## 🔄 CNB Job 集成示例

### 在 .cnb.yml 中调用 API

```yaml
job-get-issues-info:
  script:
    - |
      curl -X POST "$GO_BACKEND_URL/api/v1/issues/label-sync" \
        -H "Authorization: Bearer $INTEGRATION_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
          "repo_slug": "'"$CI_PROJECT_PATH"'",
          "issue_number": 74,
          "labels": ["🚧 处理中_CNB", "bug_CNB"]
        }'
```

### 所需 CNB 流水线密钥

在 CNB 仓库"设置 → 流水线密钥"中添加：
- `INTEGRATION_TOKEN`：与 Go 后端 `.env` 中的值一致
- `GO_BACKEND_URL`：Go 后端地址（如 `https://api.example.com`）
