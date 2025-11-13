# Cabb

将 Plane 作为"工作项中枢"，对接 CNB 代码托管与飞书的企业级集成服务。

## 🎯 核心功能

- **Plane Webhook 处理**：接收 Plane 事件，自动创建快照用于预览和通知
- **CNB 代码集成**：通过 `.cnb.yml` 回调实现 Issue/PR/分支与 Plane Issue 的双向同步
- **飞书集成**：支持卡片预览、评论同步、项目通知等功能
- **管理控制台**：Web 界面管理映射关系、凭据和配置
- **事件去重**：幂等处理避免重复事件，支持重试机制

## 🏗️ 架构设计

```
┌─────────────┐    Webhook    ┌──────────────────┐    API     ┌─────────────┐
│   Plane     │ ──────────► │ Integration      │ ─────────► │   Feishu    │
│   (工作项)   │             │ Service          │           │   (飞书)     │
└─────────────┘             └──────────────────┘           └─────────────┘
                                      ▲
                                      │ .cnb.yml 回调
                                      │
                             ┌──────────────────┐
                             │       CNB        │
                             │    (代码托管)     │
                             └──────────────────┘
```

## 🚀 快速开始

### 环境要求

- Go 1.24+
- PostgreSQL 16+
- Node.js 18+ (管理控制台)

### 1. 克隆项目

```bash
git clone <repository-url>
cd plane-integration
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，配置必要的环境变量
```

### 3. 初始化数据库

```bash
# 确保 DATABASE_URL 已正确配置
make migrate
```

### 4. 启动服务

```bash
# 开发模式
make run

# 或构建后运行
make build
./bin/server
```

服务启动后：
- API 服务：`http://localhost:8080`
- 健康检查：`http://localhost:8080/healthz`
- 管理控制台：`http://localhost:8080/admin`

## ⚙️ 配置说明

### 必需配置

| 环境变量 | 说明 | 示例 |
|----------|------|------|
| `DATABASE_URL` | PostgreSQL 连接字符串 | `postgres://user:pass@localhost:5432/plane_integration` |
| `PLANE_WEBHOOK_SECRET` | Plane Webhook 签名密钥 | `your-secret-key` |
| `ENCRYPTION_KEY` | 32字节加密密钥 | `change-me-32-bytes-long-key` |

### Plane 集成

| 环境变量 | 说明 | 获取方式 |
|----------|------|----------|
| `PLANE_BASE_URL` | Plane API 地址 | 默认 `https://api.plane.so` |
| `PLANE_SERVICE_TOKEN` | Plane Service Token | Plane 工作区设置 → Service Token |

### 飞书集成

| 环境变量 | 说明 | 获取方式 |
|----------|------|----------|
| `LARK_APP_ID` | 飞书应用 ID | 飞书开发者后台 |
| `LARK_APP_SECRET` | 飞书应用密钥 | 飞书开发者后台 |
| `LARK_ENCRYPT_KEY` | 事件加密密钥 | 飞书开发者后台 |
| `LARK_VERIFICATION_TOKEN` | 验证令牌 | 飞书开发者后台 |

### CNB 集成

| 环境变量 | 说明 | 获取方式 |
|----------|------|----------|
| `CNB_APP_TOKEN` | CNB API 访问令牌 | CNB 个人设置 → 访问令牌 |
| `INTEGRATION_TOKEN` | 回调验证令牌 | 自行生成 (`openssl rand -hex 32`) |
| `CNB_BASE_URL` | CNB API 地址 | 默认 `https://api.cnb.cool` |

### AI/自动分支命名

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `AI_PROVIDER` | 首选供应商：`cerebras` 或 `openai` | `cerebras` |
| `CEREBRAS_API_KEY` | Cerebras API Key | - |
| `CEREBRAS_BASE_URL` | Cerebras API 基地址 | `https://api.cerebras.ai` |
| `CEREBRAS_MODEL` | Cerebras 模型 | `gpt-oss-120b` |
| `OPENAI_API_KEY` | OpenAI API Key（可选：作为备选供应商） | - |
| `OPENAI_BASE_URL` | 可选，OpenAI Base URL（自定义网关/代理） | - |
| `OPENAI_MODEL` | OpenAI 模型 | `gpt-4o-mini` |

### 管理控制台

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `ADMIN_BOOTSTRAP_EMAIL` | 初始管理员邮箱 | `admin@example.com` |
| `ADMIN_BOOTSTRAP_PASSWORD` | 初始管理员密码 | `ChangeMe123!` |
| `ADMIN_SESSION_TTL_HOURS` | 会话有效期（小时） | `12` |

## 📡 API 端点

### Webhook 接收

- `POST /webhooks/plane` - Plane 事件接收（启用 AI 自动分支后，会在 Issue 创建事件中生成并创建分支）
- `POST /webhooks/lark/events` - 飞书事件接收
- `POST /webhooks/lark/interactivity` - 飞书交互回调

### CNB 回调

- `POST /ingest/cnb/issue` - CNB Issue 回调
- `POST /ingest/cnb/pr` - CNB PR 回调
- `POST /ingest/cnb/branch` - CNB 分支回调

### CNB API v1

- `POST /api/v1/issues/label-notify` - **完整版** Issue 标签通知（11 个字段）
- `POST /api/v1/issues/label-sync` - **简化版** Issue 标签同步（3 个字段）

#### 用途
接收 CNB job-get-issues-info 发送的 Issue 标签变更通知，自动同步标签到 Plane Issue。

#### 请求示例

**完整版**（推荐用于完整事件记录）：
```json
{
  "repo_slug": "1024hub/Demo",
  "issue_number": 74,
  "issue_url": "https://cnb.cool/1024hub/Demo/-/issues/74",
  "title": "实现用户登录功能",
  "state": "open",
  "author": {"username": "zhangsan", "nickname": "张三"},
  "description": "需要实现用户登录功能...",
  "labels": ["🚧 处理中_CNB", "🧑🏻‍💻 进行中：前端_CNB"],
  "label_trigger": "🚧 处理中_CNB",
  "updated_at": "2025-10-29T03:25:06Z",
  "event_context": {"event_type": "push", "branch": "feature/74-user-login"}
}
```

**简化版**（最小字段）：
```json
{
  "repo_slug": "1024hub/Demo",
  "issue_number": 74,
  "labels": ["🚧 处理中_CNB", "🧑🏻‍💻 进行中：前端_CNB"]
}
```

#### 响应示例
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

#### 业务逻辑
- 根据 `repo_slug` 查找映射的 Plane 项目
- 根据 `issue_number` 查找对应的 Plane Issue
- 同步 CNB labels 到 Plane Issue（自动创建不存在的标签）
- 记录同步事件到数据库

### 管理端点

- `GET /admin` - 管理控制台首页
- `POST /api/admin/login` - 管理员登录
- `GET /api/admin/mappings` - 查看映射关系
- `POST /api/admin/mappings/repo-project` - 创建仓库-项目映射
- `GET /admin/links/issues` - 查看 CNB Issue 链接
- `GET /admin/links/lark-threads` - 查看 Lark 线程链接
- `GET /admin/links/branches` - 查看分支链接（支持 plane_issue_id/cnb_repo_id/branch/active/limit 过滤）

### 系统

- `GET /healthz` - 健康检查
- `GET /` - 根路径（重定向到管理控制台）

## 🔧 开发指南

### 项目结构

```
├── cmd/server/          # 服务入口
├── internal/
│   ├── handlers/        # HTTP 处理器
│   ├── store/          # 数据访问层
│   ├── cnb/            # CNB 客户端
│   ├── plane/          # Plane 客户端
│   └── lark/           # 飞书客户端
├── pkg/config/         # 配置管理
├── db/migrations/      # 数据库迁移
├── web/               # 管理控制台前端
└── docs/              # 设计文档
```

### 数据库迁移

```bash
# 运行所有迁移
make migrate

# 手动执行特定迁移
psql "$DATABASE_URL" -f db/migrations/0001_init.sql
```

### 构建和部署

```bash
# 本地构建
make build

# Docker 构建
make docker-build

# Docker 运行
make docker-run

# CI 验证
make ci-verify
```

## 🔒 安全特性

- **签名验证**：Plane Webhook 使用 HMAC-SHA256 签名验证
- **Bearer Token**：CNB 回调使用 Bearer Token 验证
- **加密存储**：敏感信息（如 Service Token）透明加密存储
- **会话管理**：管理控制台基于 Cookie 的会话管理
- **请求去重**：基于 `delivery_id` 和内容哈希的幂等处理

## 📊 监控和日志

### 结构化日志

服务输出结构化 JSON 日志，包含：
- `request_id` - 请求唯一标识
- `source` - 事件来源（plane/cnb/lark）
- `event_type` - 事件类型
- `latency_ms` - 处理延迟
- `result` - 处理结果

### 健康检查

`GET /healthz` 返回服务状态：
```json
{
  "status": "ok",
  "version": "v1.0.0",
  "database": "connected",
  "timestamp": "2025-01-02T15:04:05Z"
}
```

## 🧪 测试

```bash
# 运行测试
go test ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📚 相关文档

- [架构设计](docs/design/) - 详细的架构和设计文档
- [API 规范](docs/design/integration-feature-spec-webhook-only.md) - 完整的 API 功能清单
- [CNB 集成指南](docs/design/cnb-integration-webhook-only.md) - CNB 集成详细说明
- [飞书集成指南](docs/design/feishu-integration-webhook-only.md) - 飞书集成详细说明

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

请遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范编写提交信息。

## 📄 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 🆘 故障排除

### 常见问题

**Q: 数据库连接失败**
```
A: 检查 DATABASE_URL 格式和数据库服务状态
```

**Q: Plane Webhook 验证失败**
```
A: 确认 PLANE_WEBHOOK_SECRET 与 Plane 配置一致
```

**Q: 飞书事件接收失败**
```
A: 检查 LARK_ENCRYPT_KEY 和 LARK_VERIFICATION_TOKEN 配置
```

### 日志级别

通过环境变量 `LOG_LEVEL` 控制日志级别：
- `debug` - 详细调试信息
- `info` - 一般信息（默认）
- `warn` - 警告信息
- `error` - 错误信息

---

如有问题或建议，请提交 Issue 或联系维护团队。
