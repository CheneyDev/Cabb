# Plane Integration Service（Plane 集成服务）

将 Plane 作为统一的“工作项中枢”，无缝衔接 CNB 与飞书（Feishu）协作场景，打通需求、开发与沟通的全流程。本仓库提供一个可运行的 Go/Echo 脚手架、数据库迁移，以及后续功能实现的落地入口。

- 设计文档：`docs/design/cnb-integration.md`、`docs/design/feishu-integration.md`、`docs/design/integration-feature-spec.md`
- 架构说明：`docs/ARCHITECTURE.md`

## 重要提醒（必须）
- 在进行 Plane/CNB 相关的 API 调用、Webhook 处理、`.cnb.yml` 回调/示例、字段映射、签名校验等代码编写或调整之前，务必先查阅 `docs/` 目录中已下载的最新版官方文档与本仓库的设计文档。
- 若官方文档与现有代码/本文档存在冲突，协议细节（端点/字段/签名/状态码/限流）以 `docs/` 为准；实现需据此更新，并同步修改本 README、`docs/ARCHITECTURE.md` 与相关设计文档。
- 提交前自检：是否逐项对照 `docs/` 校验了端点、字段、签名、示例与错误码？是否更新了示例与文档？

## 功能与场景（概览）
- CNB × Plane
  - 仓库↔项目映射（单向/双向同步可配）。
  - Issue/评论/标签/指派同步；PR 生命周期驱动 Plane 状态。
  - 分支事件联动与每日提交 AI 摘要评论（里程碑 M4）。
- 飞书 × Plane
  - 在飞书中创建/链接/预览工作项；命令与就地操作（指派、状态、评论）。
  - 线程消息 ↔ Plane 评论双向同步；项目新建推送到频道。
- 通用能力
  - 统一的安全校验（签名/令牌）、幂等与重试（指数退避）、结构化日志与指标。

## 架构概览
- 技术栈：Go 1.24+、Echo Web、Postgres 16。
- 模块分层
  - Connectors：`plane-connector`（OAuth/Webhook/API）、`cnb-connector`（API + `.cnb.yml` 回调）、`lark-connector`（飞书事件/卡片/命令）、`ai-connector`（提交摘要）。
  - Sync Core：字段/状态映射、方向控制、防回环与去重、评论与线程编排。
  - Storage：凭据、映射、链接与事件日志（令牌透明加密）。
  - Jobs/Scheduler：入站重试队列、每日提交摘要 CRON。
- 触发机制
  - CNB：无原生 Webhook，依赖仓库 `.cnb.yml` 在 issue/pr/branch 事件中回调本服务。
  - 飞书与 Plane：标准事件订阅与 Webhook（含验签）。

## 快速开始
### 前置要求
- Go 1.24+
- Postgres 16（或兼容版本）

### 配置环境变量
- 复制 `.env.example` 到 `.env` 并按需修改，或通过环境变量直接注入。

常用变量（完整列表见 `.env.example`）：
- 通用：`PORT`、`DATABASE_URL`、`TIMEZONE`、`ENCRYPTION_KEY`
- Plane：`PLANE_BASE_URL`、`PLANE_CLIENT_ID`、`PLANE_CLIENT_SECRET`、`PLANE_REDIRECT_URI`、`PLANE_WEBHOOK_SECRET`
- 飞书：`LARK_APP_ID`、`LARK_APP_SECRET`、`LARK_ENCRYPT_KEY`、`LARK_VERIFICATION_TOKEN`
- CNB：`CNB_APP_TOKEN`、`INTEGRATION_TOKEN`

### 初始化数据库
- 创建数据库并执行迁移：

```
createdb plane_integration
psql "$DATABASE_URL" -f db/migrations/0001_init.sql
```

### 启动服务
- 开发运行：

```
go run ./cmd/server
```

- 或构建二进制：

```
go build -o bin/plane-integration ./cmd/server
./bin/plane-integration
```

启动后访问健康检查：`GET http://localhost:8080/healthz`

### Makefile（本地便捷命令）
```
# 编译二进制
make build

# 本地运行（可用 PORT=xxxx 指定端口）
make run PORT=8080

# 构建并运行 Docker 镜像
make docker-build
make docker-run PORT=8080

# 执行迁移（需设置 DATABASE_URL）
make migrate

# 本地 CI 验证：启动服务并检查 /healthz
make ci-verify
```

## 部署（Render）
- 方式一：使用 Blueprint（推荐）
  - 本仓库已包含 `render.yaml` 与 `Dockerfile`。
  - 在 Render 仪表盘选择“Blueprint”→ 连接你的 GitHub 仓库 → 部署。
  - Blueprint 会创建：
    - Web 服务：基于 `Dockerfile` 构建，健康检查路径 `/healthz`。
    - Postgres 数据库：`plane-integration-db`，连接串注入为 `DATABASE_URL`。
  - 首次部署会自动执行迁移：容器入口脚本 `scripts/entrypoint.sh` 会在启动时运行 `psql "$DATABASE_URL" -f /app/db/migrations/0001_init.sql`。
  - 在服务的环境变量中补齐：
    - 必填：`PLANE_CLIENT_ID`、`PLANE_CLIENT_SECRET`
    - 可选：`PLANE_WEBHOOK_SECRET`、`INTEGRATION_TOKEN`、`LARK_*`
    - `PLANE_REDIRECT_URI`：部署后改为 `https://<your-service>.onrender.com/plane/oauth/callback`
  - 验证：
    - 健康检查：`https://<your-service>.onrender.com/healthz`
    - OAuth 起始：`https://<your-service>.onrender.com/plane/oauth/start`

- 方式二：原生 Go（非 Docker）
  - 在 Render 创建 Web Service（Go），设置 `Build Command: go build -o server ./cmd/server`、`Start Command: ./server`。
  - 配置 `GO_VERSION` 至 1.23.x（若 1.24 暂不支持）；或临时将 `go.mod` 中 `go 1.24` 降为 `1.23` 用于测试（不建议长期）。
  - 迁移：在本地使用 `psql "$DATABASE_URL" -f db/migrations/0001_init.sql` 连接 Render Postgres 执行，或在“Post‑deploy Command”中配置执行。

注意
- 免费实例可能冷启动，首次请求有延迟；Webhook 超时敏感场景建议升级或保持活跃。
- 迁移脚本使用 `pgcrypto` 扩展，Render 托管 Postgres 支持该扩展。

## 使用 Docker 本地运行
```
docker build -t plane-integration:dev .
docker run --rm -p 8080:8080 \
  -e DATABASE_URL="postgres://root:123456@host.docker.internal:15432/plane_intergration?sslmode=disable" \
  -e PLANE_REDIRECT_URI="http://localhost:8080/plane/oauth/callback" \
  plane-integration:dev
```

## API 与端点（脚手架）
- 健康检查
  - `GET /healthz`
- Plane（OAuth/Webhook）
  - `GET /plane/oauth/start`（重定向到 Plane 授权页）
  - `GET /plane/oauth/callback`（处理 app_installation_id 或 code，换取 Token 并回传安装信息摘要）
  - `POST /webhooks/plane`（支持 `X-Plane-Signature` HMAC-SHA256 验签）
- CNB（来自 `.cnb.yml` 的回调）
  - `POST /ingest/cnb/issue`
  - `POST /ingest/cnb/pr`
  - `POST /ingest/cnb/branch`
  - 安全：`Authorization: Bearer $INTEGRATION_TOKEN`
- 飞书（Feishu/Lark）
  - `POST /webhooks/lark/events`（支持 challenge 握手）
  - `POST /webhooks/lark/interactivity`
  - `POST /webhooks/lark/commands`
- 管理映射
  - `POST /admin/mappings/repo-project`
  - `POST /admin/mappings/pr-states`
  - `POST /admin/mappings/users`
  - `POST /admin/mappings/channel-project`
- 任务
  - `POST /jobs/issue-summary/daily`

### 示例：CNB Issue 回调
```
curl -X POST "$INTEGRATION_URL/ingest/cnb/issue" \
  -H "Authorization: Bearer $INTEGRATION_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"event":"issue.open","repo":"group/repo","issue_iid":"42"}'
```

### 示例：Plane OAuth
- 启动安装/同意（浏览器访问，支持可选 `state`）：
```
open "http://localhost:8080/plane/oauth/start?state=dev"
```

- 回调（Plane 完成安装后将携带 `app_installation_id` 与 `code` 调用）：
```
# 示例：本地手动验证（模拟 Plane 回调）
curl "http://localhost:8080/plane/oauth/callback?app_installation_id=<uuid>&code=<code>&state=dev"
```

注意：服务端会调用 Plane 的 `/auth/o/token/` 与 `/auth/o/app-installation/` 完成 Token 交换与安装信息查询；当前仅返回摘要 JSON，不回显敏感 Token。令牌持久化与加密存储将在接入数据库后启用。

## 安全与鉴权
- Plane Webhook：校验 `X-Plane-Signature`（HMAC-SHA256(secret, raw_body)）。
- CNB 回调：校验 `Authorization: Bearer $INTEGRATION_TOKEN`。
- 飞书事件：支持 challenge；正式环境需接入签名/时间戳校验（预留）。
- 令牌安全：后续实现对敏感令牌（access/refresh/tenant）进行透明加密存储。

## 目录结构
```
cmd/server/                # 服务入口与 HTTP 启动
internal/handlers/         # 路由与各端点处理（Plane/CNB/Lark/Admin/Jobs）
internal/store/            # 数据层占位（DB 连接/仓储）
internal/version/          # 版本信息
pkg/config/                # 环境变量加载
db/migrations/             # Postgres 迁移脚本
docs/design/               # 详细设计文档
```

## 里程碑（简要）
- 飞书
  - M1：创建/链接/预览、项目新建推送
  - M2：线程双向同步、卡片就地操作
  - M3：用户映射完善、可观测与配额/权限优化
- CNB
  - M1：最小可用（CNB→Plane 单向：Issue/评论/映射）
  - M2：双向同步与用户映射
  - M3：PR 生命周期自动化
  - M4：分支联动与每日提交 AI 摘要

## 开发下一步（建议）
- 接入数据库连接与启动时迁移。
- 实现令牌加密与各 Connector（Plane/CNB/Lark/AI）。
- 补齐 Sync Core、幂等存储、重试与调度器。
- 按 `docs/design/*` 逐步补全各路由的业务逻辑与管理接口。

---
如需我继续：
- 接入 DB 连接与自动迁移
- 定义各 Connector 接口与最小实现
- 实装 Plane/CNB/飞书的安全校验与事件解包
