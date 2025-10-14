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
- Plane：`PLANE_BASE_URL`、`PLANE_CLIENT_ID`、`PLANE_CLIENT_SECRET`、`PLANE_REDIRECT_URI`、`PLANE_WEBHOOK_SECRET`、`PLANE_APP_BASE_URL`
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

## 在 Plane 设置页面配置与安装（全流程）
以下指引面向“从 Plane 的设置页面完成安装”的场景，覆盖跳转链路与服务端对 Plane 的 API 调用。协议细节以 `docs/` 为准。

结论（要做什么）
- 在 Plane 中创建应用，设置 `Setup URL`/`Redirect URI`/`Webhook URL` 指向本服务。
- 用户在 Plane 设置页点击“Install/Authorize”后，浏览器完成 2 次跳转返回本服务的回调；本服务随后向 Plane 发起 1–2 次令牌交换与 1 次安装信息查询。
- 安装完成后，Plane 的事件将投递到 `POST /webhooks/plane`，由本服务处理。

前置准备
- 本服务已可访问（本地或线上），并设置：
  - `PLANE_BASE_URL`（默认 `https://api.plane.so`）
  - `PLANE_CLIENT_ID`、`PLANE_CLIENT_SECRET`
  - `PLANE_REDIRECT_URI`（如 `https://<your-domain>/plane/oauth/callback`）
  - `PLANE_WEBHOOK_SECRET`（用于校验 Plane Webhook 签名）

一、在 Plane 配置应用
1) 进入 Plane 工作区设置 → Apps/Integrations 或 Developer/Apps → 新建应用（例如 `plane-cnb-app` 或 `plane-feishu-app`）。
2) 填写：
   - Setup URL：`https://<your-domain>/plane/oauth/start`
   - Redirect URI：`https://<your-domain>/plane/oauth/callback`
   - Webhook URL：`https://<your-domain>/webhooks/plane`
3) 记录发给你的 `client_id` 与 `client_secret`，填入本服务环境变量。
4) 保存并发布（如有“启用/可见范围/权限”等选项，按业务需要开启）。

二、安装/授权跳转链路（浏览器）
1) 用户在 Plane 的应用详情页点击 Install/Authorize。
2) 浏览器跳转至本服务的 Setup URL：
   - Plane → `GET /plane/oauth/start?state=<opaque>`（由 Plane 传入或你自定义的 `state` 会被回传）。
3) 本服务立即 302 重定向到 Plane 的授权/安装确认页（URL 由 Plane 控制，参数包含 `client_id`、`redirect_uri`、`response_type=code`、`state`）。
4) 用户在 Plane 同意后，浏览器被重定向回本服务回调（该页面会自动返回到 Plane）：
   - Plane → `GET /plane/oauth/callback?app_installation_id=<uuid>&code=<code>&state=<opaque>`
   - 说明：Plane 可能同时携带 `app_installation_id`（用于应用级 Bot 授权）与 `code`（用于用户授权）。本服务会分别处理。
    - 回调页面行为：
      - 默认返回一个 HTML 页面，自动跳回 Plane；调试或 API 场景可用 `?format=json` 强制返回 JSON。
      - 跳转目标优先级：`return_to`（query）→ `https://app.plane.so/{workspace_slug}/settings/integrations/`（若拿到工作区 slug）→ 将 `state` 作为 URL（若安全）→ `PLANE_APP_BASE_URL` → 从 `PLANE_BASE_URL` 推断（`api.*` → `app.*`）→ `PLANE_BASE_URL`。
      - 安全：仅允许跳转到与 `PLANE_APP_BASE_URL`/`PLANE_BASE_URL` 同主机（或 `api.*` → `app.*`）的站点，避免开放重定向。

三、服务端向 Plane 发起的 API 请求
- 1) 交换 Bot Token（应用安装维度，可选但推荐）
  - 触发条件：回调中含有 `app_installation_id`。
  - 请求：
```
POST {PLANE_BASE_URL}/auth/o/token
Content-Type: application/x-www-form-urlencoded
Authorization: Basic base64({PLANE_CLIENT_ID}:{PLANE_CLIENT_SECRET})

grant_type=client_credentials&app_installation_id={UUID}
```
  - 响应（示例）：
```
{
  "access_token": "<bot-token>",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

- 2) 交换用户 Token（用户维度，可选，需在 Plane 显示用户同意）
  - 触发条件：回调中含有 `code`。
  - 请求：
```
POST {PLANE_BASE_URL}/auth/o/token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&code={CODE}&redirect_uri={PLANE_REDIRECT_URI}&client_id={PLANE_CLIENT_ID}&client_secret={PLANE_CLIENT_SECRET}
```
  - 响应（示例）：
```
{
  "access_token": "<user-access-token>",
  "refresh_token": "<user-refresh-token>",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

- 3) 获取安装信息（用于落库工作区/安装上下文）
  - 触发条件：通常在拿到 Bot Token 后立即调用。
  - 请求：
```
GET {PLANE_BASE_URL}/auth/o/app-installation?id={app_installation_id}
Authorization: Bearer <bot-token>
```
  - 响应要点：`workspace_id`、`workspace_slug`、`app_bot`。
  - 注意：回跳 URL 使用 `workspace_slug` 作为路径段（不是 `workspace_id`）。

四、Webhook 投递与验签（安装后）
- Plane 将把工作区内的事件投递到 `POST /webhooks/plane`。
- 安全校验：`X-Plane-Signature: sha256=...`；本服务使用 `HMAC-SHA256(PLANE_WEBHOOK_SECRET, raw_body)` 做常量时间比较。
- 幂等：以 `delivery_id + payload_sha256` 去重，重复投递返回 `200` 并可标注 `status=duplicate`。

五、时序图（端到端）
```mermaid
sequenceDiagram
  participant User as 用户
  participant Plane as Plane(设置/授权)
  participant Svc as 集成服务

  User->>Plane: 点击 Install/Authorize
  Plane-->>Svc: GET /plane/oauth/start?state=...
  Svc-->>Plane: 302 到授权确认页(带 client_id/redirect_uri/state)
  Plane-->>Svc: GET /plane/oauth/callback?app_installation_id&code&state
  Svc->>Plane: POST /auth/o/token (client_credentials)
  Plane-->>Svc: 200 access_token(bot)
  Svc->>Plane: GET /auth/o/app-installation?id=...
  Plane-->>Svc: 200 workspace/install 信息
  Note over Plane,Svc: (可选) 用户授权
  Svc->>Plane: POST /auth/o/token (authorization_code)
  Plane-->>Svc: 200 access_token(user), refresh_token
  Plane->>Svc: POST /webhooks/plane (后续事件)
```

六、常见错误与排查
- 回调未带参数：`/plane/oauth/callback` 缺少 `app_installation_id` 与 `code` → 返回 `400 invalid_request`。
- 令牌交换失败：Plane 返回 `400 invalid_client/invalid_grant` → 检查 `client_id/secret` 与 `redirect_uri` 是否一致。
- Webhook 验签失败：检查 `PLANE_WEBHOOK_SECRET` 是否与 Plane 配置一致；确认使用原始请求体计算 HMAC。

七、验收要点（DoD）
- 在 Plane 设置页完成安装后，`/plane/oauth/callback` 返回安装摘要 JSON（当前实现不回显敏感令牌）。
- 本服务日志打印：安装成功（含 `workspace_slug`/`app_installation_id`），并在 Bot Token 过期前可自动续期。
- 触发一次 Plane 事件（如 issue 更新），`/webhooks/plane` 能验签并记录投递。

八、本地调试最小示例
- 启动授权（浏览器）：
```
open "http://localhost:8080/plane/oauth/start?state=dev"
```
- 模拟 Plane 回调：
```
curl "http://localhost:8080/plane/oauth/callback?app_installation_id=<uuid>&code=<code>&state=dev"
```

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

可选参数
- `return_to`：指定回跳 URL（需与 Plane 主机一致，否则忽略）。
- `format=json`：强制以 JSON 返回（默认浏览器为 HTML 自动跳转）。

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
