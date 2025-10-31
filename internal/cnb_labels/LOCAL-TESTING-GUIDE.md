# 本地测试指南

## 📋 概述

本指南用于快速启动服务并运行自动化测试。

**前置准备：**
- 确保已配置 `.env`（参考 [ConfigNote.md](./ConfigNote.md)）
- Go 1.24+ 已安装

---

## 🚀 启动服务

### 方法 1：使用启动脚本（推荐）

```bash
./scripts/start-server.sh
```

服务将在前台运行，按 `Ctrl+C` 停止。

### 方法 2：手动启动

```bash
go run ./cmd/server
```

### 方法 3：编译后运行

```bash
go build -o plane-integration ./cmd/server
./plane-integration
```

---

## ✅ 验证服务

```bash
curl http://localhost:8080/healthz
```

**预期响应：**
```json
{
  "db": "not_connected",
  "status": "ok",
  "time": "2025-10-30T15:53:09+08:00",
  "version": "0.1.0-dev"
}
```

---

## 🧪 自动化测试

### 运行所有单元测试

```bash
go test ./internal/handlers/... -v
```

### 运行特定测试场景

```bash
# 测试成功案例
go test ./internal/handlers -run TestIssueLabelNotify_Success -v
go test ./internal/handlers -run TestIssueLabelSync_Success -v

# 测试鉴权失败
go test ./internal/handlers -run TestIssueLabelNotify_Unauthorized -v

# 测试字段校验
go test ./internal/handlers -run TestIssueLabelNotify_MissingFields -v

# 测试 JSON 解析
go test ./internal/handlers -run TestIssueLabelNotify_InvalidJSON -v

# 测试幂等性
go test ./internal/handlers -run TestIssueLabelNotify_Idempotency -v

# 测试标签过滤
go test ./internal/handlers -run TestFilterCNBLabels -v
```

### 查看测试覆盖率

```bash
go test ./internal/handlers/... -cover
```

### 生成覆盖率报告

```bash
go test ./internal/handlers/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 📊 测试覆盖场景

| 测试场景 | 测试函数 | 验证点 |
|---------|---------|--------|
| 成功请求 | `TestIssueLabelNotify_Success` | 200, 正确响应格式 |
| 鉴权失败 | `TestIssueLabelNotify_Unauthorized` | 401, invalid_token |
| 缺少字段 | `TestIssueLabelNotify_MissingFields` | 400, missing_fields |
| 无效 JSON | `TestIssueLabelNotify_InvalidJSON` | 422, invalid_json |
| 幂等性 | `TestIssueLabelNotify_Idempotency` | 200, status=duplicate |
| 简化 API | `TestIssueLabelSync_Success` | 200, 3 字段版本 |
| 标签过滤 | `TestFilterCNBLabels` | 过滤逻辑正确 |



---

## 🔧 故障排查

### 问题 1：端口被占用

```bash
# 查找占用进程
lsof -i:8080

# 停止进程
kill -9 <PID>
```

### 问题 2：服务启动失败

检查日志：
```bash
go run ./cmd/server 2>&1 | tee server.log
```

### 问题 3：鉴权失败

确保 `INTEGRATION_TOKEN` 环境变量正确设置：
```bash
echo $INTEGRATION_TOKEN
```

### 问题 4：数据库连接失败

数据库连接失败不影响服务启动，只是会跳过需要数据库的功能（如持久化去重、标签映射查询等）。

---

## 📊 日志查看

服务会输出结构化 JSON 日志：

```json
{
  "time": "2025-10-30T07:53:09Z",
  "level": "info",
  "request_id": "...",
  "method": "POST",
  "endpoint": "/api/v1/issues/label-sync",
  "status": 200,
  "latency_ms": 0,
  "result": "success"
}
```
