# Issue 标签通知 API - 快速测试指南

## 功能说明

本 API 端点用于接收来自 CNB `job-get-issues-info` 的 issue 标签变更通知，替代原有的飞书机器人直接通知方案。

**端点：** `POST /api/v1/issues/label-notify`

## 快速启动（最小配置）

### 1. 准备环境变量

创建 `.env` 文件（或设置环境变量）：

```bash
# 最小配置 - 可以启动服务
PORT=8080
DATABASE_URL=  # 可选：如果没有数据库，服务仍可启动（内存去重）
TIMEZONE=Local
ENCRYPTION_KEY=change-me-32-bytes-key-12345

# CNB 认证令牌（必需用于测试 API）
INTEGRATION_TOKEN=test-token-for-development-only

# 其他可选配置（暂时留空）
PLANE_CLIENT_ID=
PLANE_CLIENT_SECRET=
LARK_APP_ID=
LARK_APP_SECRET=
CNB_APP_TOKEN=
```

**重要提示：**
- ❌ **Plane OAuth 配置（`PLANE_CLIENT_ID`/`PLANE_CLIENT_SECRET`）为空不影响服务启动**
- ✅ 只有访问 OAuth 端点（`/plane/oauth/*`）时才需要这些配置
- ✅ Issue 标签通知 API 只需要 `INTEGRATION_TOKEN` 即可工作

### 2. 启动服务

```bash
# 方式 1：直接运行
go run ./cmd/server

# 方式 2：编译后运行
go build -o plane-integration ./cmd/server
./plane-integration
```

启动成功会看到：
```
plane-integration 0.0.1 listening on :8080
```

### 3. 验证服务健康

```bash
curl http://localhost:8080/healthz
```

预期响应：
```json
{"status":"ok","version":"0.0.1"}
```

### 4. 测试 API

#### 使用自动化测试脚本

```bash
export INTEGRATION_TOKEN=test-token-for-development-only
./scripts/test-label-notify.sh
```

#### 手动测试（成功案例）

```bash
curl -X POST "http://localhost:8080/api/v1/issues/label-notify" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token-for-development-only" \
  -H "X-Delivery-ID: test-$(date +%s)" \
  -d '{
    "repo_slug": "1024hub/Demo",
    "issue_number": 74,
    "issue_url": "https://cnb.cool/1024hub/Demo/-/issues/74",
    "title": "实现用户登录功能",
    "state": "open",
    "author": {"username": "zhangsan", "nickname": "张三"},
    "description": "功能描述",
    "labels": ["bug", "feature"],
    "label_trigger": "bug",
    "updated_at": "2025-01-30T06:00:00Z",
    "event_context": {"event_type": "push", "branch": "main"}
  }' | jq .
```

预期响应：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "issue_number": 74,
    "processed_at": "2025-01-30T06:25:10Z"
  }
}
```

#### 测试鉴权失败（401）

```bash
curl -X POST "http://localhost:8080/api/v1/issues/label-notify" \
  -H "Content-Type: application/json" \
  -d '{"repo_slug": "test/repo", "issue_number": 1}' | jq .
```

预期响应：
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

#### 测试参数校验（400）

```bash
curl -X POST "http://localhost:8080/api/v1/issues/label-notify" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token-for-development-only" \
  -d '{"repo_slug": "test/repo"}' | jq .
```

预期响应：
```json
{
  "error": {
    "code": "missing_fields",
    "message": "缺少必填字段：issue_number",
    ...
  }
}
```

## 架构特性

### ✅ 高内聚、低耦合设计

- **Handler 层**（`internal/handlers/issue_label_notify.go`）：
  - 职责：HTTP 请求处理、参数校验、鉴权、响应封装
  - 依赖：`config.Config`、`store.DB`、`Deduper`
  - 独立性：不依赖具体业务逻辑（Plane/Lark 客户端）

- **业务处理层**（`processIssueLabelNotify`）：
  - 异步执行，不阻塞 HTTP 响应
  - 预留扩展点：同步 Plane、发送飞书通知、触发工作流

- **数据层**（`internal/store`）：
  - 统一数据访问接口
  - 支持可选数据库（无 DB 时降级为内存去重）

### ✅ 统一的错误处理

- 使用 `writeError` 统一错误响应格式
- 结构化日志记录（包含 `request_id`、`source`、`endpoint`）
- 错误码机器可读（`invalid_token`、`missing_fields` 等）

### ✅ 幂等性保证

- **内存级**：5 分钟 TTL，基于 `delivery_id` + `payload_sha256`
- **数据库级**：持久化到 `event_deliveries` 表（如有数据库）
- 支持 `X-Delivery-ID` 请求头自定义幂等键

### ✅ 安全性

- Bearer token 认证（与 CNB `.cnb.yml` 回调共享 `INTEGRATION_TOKEN`）
- 请求体大小限制（Echo 框架默认 2MB）
- 结构化日志不输出敏感信息

## 部署配置

### 生产环境建议

```bash
# 数据库（强烈推荐）
DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=require

# 安全令牌（至少 32 字节随机）
INTEGRATION_TOKEN=$(openssl rand -hex 32)

# 加密密钥（用于存储凭据）
ENCRYPTION_KEY=$(openssl rand -hex 16)

# Plane OAuth（如需 OAuth 流程）
PLANE_CLIENT_ID=your-plane-client-id
PLANE_CLIENT_SECRET=your-plane-secret
PLANE_BASE_URL=https://api.plane.so
PLANE_WEBHOOK_SECRET=your-webhook-secret

# 飞书配置（如需飞书通知）
LARK_APP_ID=your-lark-app-id
LARK_APP_SECRET=your-lark-secret
LARK_ENCRYPT_KEY=your-encrypt-key

# CNB API（如需出站调用 CNB）
CNB_APP_TOKEN=your-cnb-token
CNB_BASE_URL=https://api.cnb.cool
CNB_OUTBOUND_ENABLED=true

# 管理后台
ADMIN_BOOTSTRAP_EMAIL=admin@example.com
ADMIN_BOOTSTRAP_PASSWORD=ChangeMe123!
```

### CNB Job 配置

在 CNB 流水线中配置：

```yaml
# .cnb.yml
env:
  GO_BACKEND_URL: https://your-domain.com
  GO_BACKEND_TOKEN: ${GO_BACKEND_TOKEN}  # 从密钥库注入
```

在 CNB 控制台"仓库设置 → 流水线密钥"中添加：
- 密钥名称：`GO_BACKEND_TOKEN`
- 密钥值：与后端 `INTEGRATION_TOKEN` 相同

## 文档资源

- **集成文档**：`docs/cnb-job-integration.md` - CNB job 修改详细步骤
- **API 示例**：`docs/api-examples.md` - 更多 curl 测试用例
- **架构文档**：`docs/ARCHITECTURE.md` - 整体架构说明
- **设计文档**：`docs/design/cnb-integration.md` - CNB 集成设计

## 常见问题

### Q: 服务启动失败，提示数据库连接错误

A: 数据库连接是可选的。如果暂时没有数据库：
- 将 `DATABASE_URL` 留空或删除该配置
- 服务会使用内存去重（5 分钟 TTL）
- 日志会提示 "db connect error"，但不影响启动

### Q: API 返回 401 错误

A: 检查：
1. 请求头格式：`Authorization: Bearer <token>`（注意 "Bearer " 前缀和空格）
2. Token 值与 `INTEGRATION_TOKEN` 环境变量一致
3. Token 没有包含换行符或额外空格

### Q: 如何验证 Plane OAuth 配置是否必需？

A: Plane OAuth 配置**不是**服务启动的必需条件：
- ✅ 可以在没有 OAuth 配置时启动服务
- ✅ Issue 标签通知 API 完全不依赖 OAuth
- ❌ 只有访问 `/plane/oauth/start` 或 `/plane/oauth/callback` 时才需要

### Q: 当前实现会做什么业务处理？

A: 当前为占位实现（placeholder）：
- ✅ 接收并验证请求
- ✅ 记录事件到数据库（如有）
- ✅ 返回成功响应
- ⏳ 业务逻辑（同步 Plane、飞书通知）待后续实现

实现位置：`internal/handlers/issue_label_notify.go` 中的 `processIssueLabelNotify` 方法。

## 后续步骤

1. **完善业务逻辑**：在 `processIssueLabelNotify` 中实现：
   - 根据 `repo_slug` 查询 repo-project 映射
   - 同步标签到对应的 Plane Issue
   - 根据 channel-project 映射发送飞书通知

2. **集成测试**：修改 CNB job 代码，替换飞书通知为 API 调用

3. **监控告警**：配置日志聚合与指标收集，监控 API 调用量与错误率

4. **性能优化**：根据实际负载调整超时、并发数、数据库连接池
