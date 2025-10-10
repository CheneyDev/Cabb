```mermaid
flowchart LR
  U[团队成员] -->|使用 Plane 管理需求| PlaneNode
  U -->|提交代码/Issue/PR| CNBNode

  subgraph PlaneNode[Plane]
    PWH[Webhook]
    PAPI[API]
  end

  subgraph CNBNode[CNB 代码托管]
    PIPE[.cnb.yml 内置任务/流水线]
    CAPI[API]
  end

  SVC[集成服务] --> PAPI
  SVC --> CAPI
  PWH --> SVC
  PIPE -->|HTTP回调| SVC

  DB[(Postgres 映射/凭据/日志)]
  SVC --- DB

  U -. 一次性安装与授权 .- SVC

  CFG[["同步策略<br/>- 单向/双向可配置<br/>- Issue/评论<br/>- PR → 状态自动化"]]
  SVC --- CFG
```

**CNB 集成设计（参考 GitHub 集成）**

- 目标：将公司内部 CNB 代码托管平台与 Plane 双向或单向同步，简化需求/开发协作。参考 Plane × GitHub 集成能力：仓库与项目映射、Issue 双向同步、PR 生命周期联动、评论与标签同步、用户身份映射、可配置的同步方向。

**整体架构**
- 技术栈：Go 1.24、Echo Web 框架、Postgres 16。
- 组件划分：
  - `plane-connector`：通过 Plane OAuth 获取 user/bot token，订阅 Plane Webhook，消费事件并回写 CNB。
  - `cnb-connector`：对接 CNB API；通过仓库 `.cnb.yml` 配置的内置任务在 CNB 事件发生时回调集成服务；回写 Plane。
  - `sync-core`：映射与路由层（项目/仓库/用户/状态/标签/PR 状态映射、方向控制、去重与幂等）。
  - `storage`：Postgres 存储凭据、映射关系、事件投递日志、去重指纹。
  - `jobs`：异步任务与重试（本地队列/表驱动，事务内投递，指数退避）。

**部署与配置**
- 进程配置：
  - `PORT`（服务端口），`DATABASE_URL`，`ENCRYPTION_KEY`（KMS/本地），`LOG_LEVEL`。
  - Plane：`PLANE_CLIENT_ID`，`PLANE_CLIENT_SECRET`，`PLANE_BASE_URL`（如 https://api.plane.so），`PLANE_WEBHOOK_SECRET`。
  - CNB：`CNB_BASE_URL`，`CNB_APP_TOKEN`（服务调用 CNB API 的凭据，按需配置）。
  - 集成入站鉴权：`INTEGRATION_SHARED_SECRET`（供 CNB 流水线以 Bearer 传入）。
- 回调与事件端点：
  - Plane OAuth：`/plane/oauth/start`、`/plane/oauth/callback`
  - Plane Webhook：`/webhooks/plane`
  - CNB 事件入站（由 .cnb.yml 流水线触发）：
    - `POST /ingest/cnb/issue`（issue.open/issue.close/issue.reopen）
    - `POST /ingest/cnb/pr`（pull_request/pull_request.update/pull_request.target/...）

**数据模型（Postgres）**
- 下图仅展示核心表及主要关系，帮助理解“凭据/映射/链接/投递日志”的职责分工。
```mermaid
classDiagram
  class workspaces {
    uuid id
    uuid plane_workspace_id
    uuid app_installation_id
    enum token_type
    text access_token*
    text refresh_token
    timestamptz expires_at
  }
  class cnb_installations {
    uuid id
    string cnb_org_id
    string installation_id
    text access_token*
    text refresh_token
    timestamptz expires_at
  }
  class repo_project_mappings {
    uuid id
    uuid plane_project_id
    uuid plane_workspace_id
    string cnb_repo_id
    uuid issue_open_state_id
    uuid issue_closed_state_id
    enum sync_direction
    bool active
  }
  class pr_state_mappings {
    uuid id
    uuid plane_project_id
    string cnb_repo_id
    uuid draft_state_id
    uuid opened_state_id
    uuid review_requested_state_id
    uuid approved_state_id
    uuid merged_state_id
    uuid closed_state_id
  }
  class user_mappings {
    uuid id
    uuid plane_user_id
    string cnb_user_id
    string display_name
    timestamptz connected_at
  }
  class issue_links {
    uuid id
    uuid plane_issue_id
    string cnb_issue_id
    string cnb_repo_id
    timestamptz linked_at
  }
  class event_deliveries {
    uuid id
    string source
    string event_type
    string delivery_id
    string payload_sha256
    string status
    int retries
    timestamptz next_retry_at
    timestamptz created_at
  }

  workspaces --> repo_project_mappings : plane_workspace_id
  repo_project_mappings --> issue_links : project/repo
  cnb_installations --> repo_project_mappings : installation scope
  user_mappings .. issue_links : assignees/comments via mapping
  event_deliveries ..> repo_project_mappings : idempotency & retry
```
- `workspaces`：存 Plane 侧安装与令牌
  - `id` (pk)、`plane_workspace_id`、`app_installation_id`、`token_type`（user/bot）、`access_token`（加密）、`refresh_token`（加密，可空）、`expires_at`。
- `cnb_installations`：存 CNB 安装与令牌
  - `id`、`cnb_org_id`/`owner_id`、`installation_id`、`access_token`（加密）、`refresh_token`（加密，可空）、`expires_at`。
- `repo_project_mappings`：仓库与项目映射与同步方向
  - `id`、`plane_project_id`、`plane_workspace_id`、`cnb_repo_id`、`issue_open_state_id`、`issue_closed_state_id`、`sync_direction`（`cnb_to_plane`/`bidirectional`）、`active`。
- `pr_state_mappings`：PR 生命周期到 Plane 状态映射
  - `id`、`plane_project_id`、`cnb_repo_id`、`draft_state_id`、`opened_state_id`、`review_requested_state_id`、`approved_state_id`、`merged_state_id`、`closed_state_id`。
- `user_mappings`：用户身份映射
  - `id`、`plane_user_id`、`cnb_user_id`、`display_name`、`connected_at`。
- `issue_links`：跨系统工单关联
  - `id`、`plane_issue_id`、`cnb_issue_id`、`cnb_repo_id`、`linked_at`。
- `event_deliveries`：事件去重与可观测
  - `id`、`source`（plane/cnb）、`event_type`、`delivery_id`/`signature`、`payload_sha256`、`status`、`retries`、`next_retry_at`、`created_at`。

**同步模型与规则**
- 典型的双向/单向同步时序如下（CNB 侧通过 .cnb.yml 触发回调，非 Webhook）：
```mermaid
sequenceDiagram
  participant CNB as CNB
  participant SVC as 集成服务
  participant Plane as Plane
  Note over CNB,Plane: CNB→Plane（Issue 事件触发，服务侧判定含 Plane 标签）
  CNB->>SVC: issue.open/close/reopen（.cnb.yml 回调）
  SVC->>SVC: 查映射/幂等
  SVC->>Plane: 创建/更新 Work Item
  Plane-->>SVC: 201/200 + 链接
  SVC->>CNB: 回写评论附 Plane 链接

  Note over CNB,Plane: Plane→CNB（启用双向 + CNB 标签）
  Plane->>SVC: issue.updated(labels+=CNB)
  SVC->>SVC: 查映射/幂等
  SVC->>CNB: 创建/更新 Issue
  CNB-->>SVC: 201/200 + 链接
  SVC-->>Plane: 回写评论附 CNB 链接
```
- 同步方向（默认单向 CNB→Plane，可切换双向）：参考 Plane×GitHub 的“Sync issues / What gets synced?”
  - 标题、描述、指派人、标签、评论：默认双向；遵循冲突最后写入或来源优先（由 `sync_direction` 决定）。
  - 状态：默认 CNB→Plane（issue open/closed 映射 Plane State）。
- 触发机制（CNB 无 Webhook 的适配）：
  - CNB→Plane：在 `issue.open/issue.close/issue.reopen` 等事件中通过 `.cnb.yml` 以内置任务回调集成服务；服务侧拉取 Issue 详情并判断是否含 `Plane` 标签，满足则创建/链接 Plane Work Item（标签作为“门槛”而非事件源）。
  - Plane→CNB：在 Plane Work Item 打上 `CNB` 标签 → 由 Plane Webhook 触发服务回写 CNB（创建/链接 Issue）。若需“补同步”，可在 CNB 仓库提供 `web_trigger` 按钮手动触发一次回调。
- 评论同步：
  - 双向同步，若未映射用户则以 Bot 身份留言（Plane 侧可绑定“个人账户”以便归属评论，参考 GitHub 文档“Connect personal GitHub account”）。
- 用户映射：
  - 后台支持手动/半自动（邮件/用户名）映射；缺省不阻断同步但可能丢失“指派人”。

**PR 生命周期自动化（参考 GitHub PR mapping）**
- PR 状态到 Plane 状态的自动化映射：
```mermaid
flowchart LR
  A[CNB PR 状态变化] -->|draft| D[Plane 状态: 草稿]
  A -->|opened| O[Plane 状态: 待办]
  A -->|review requested| R[Plane 状态: 评审中]
  A -->|approved| AP[Plane 状态: 待合并]
  A -->|merged| M[Plane 状态: 已完成]
  A -->|closed| C[Plane 状态: 已关闭]
  N["说明：标题或描述含 (KEY-123) 可触发自动化<br/>仅写 KEY-123 为引用"]
  A -.-> N
```
- 支持将 CNB PR 状态映射到 Plane Work Item 状态：draft/opened/review_requested/approved/merged/closed。
- 参考格式（同 GitHub）
  - 标题/描述包含带方括号的 Issue 标识如 `[WEB-344]` → 链接并触发状态自动化。
  - 不带括号 `WEB-344` → 仅建立引用关系，不触发状态变更。

**接口设计（Echo）**
- `GET /healthz`：健康检查。
- `GET /plane/oauth/start`：重定向至 Plane 授权；`GET /plane/oauth/callback`：换取 token（支持 user/bot token，含 `app_installation_id`）。
- `POST /webhooks/plane`：Plane webhook（校验 `X-Plane-Delivery`/`X-Plane-Event`/`X-Plane-Signature`，HMAC-SHA256）。
- `POST /ingest/cnb/issue`：接收 `.cnb.yml` 流水线在 Issue 事件中回调的最小负载（repo、issue_iid、event）。
- `POST /ingest/cnb/pr`：接收 PR 相关事件回调（iid、action、状态等）。
  请求示例：
  - `/ingest/cnb/issue`
    ```json
    {
      "event": "issue.open", // issue.open | issue.close | issue.reopen
      "repo": "group_slug/repo_name",
      "issue_iid": "42"
    }
    ```
  - `/ingest/cnb/pr`
    ```json
    {
      "event": "pull_request", // pull_request | pull_request.update | ...
      "action": "created|code_update|status_update",
      "repo": "group_slug/repo_name",
      "pr_iid": "13"
    }
    ```
- `POST /admin/mappings/repo-project`：创建/更新仓库-项目映射（含状态映射与方向）。
- `POST /admin/mappings/pr-states`：配置 PR 生命周期映射。
- `POST /admin/mappings/users`：批量用户映射。

**事件与处理（概要）**
- 入站回调处理（校验、去重与重试）：
```mermaid
flowchart TD
  In[接收 Webhook] --> Verify{签名有效?}
  Verify -- 否 --> Rej[返回 403/丢弃]
  Verify -- 是 --> Dedup{重复投递?}
  Dedup -- 是 --> Ack[200 确认, 忽略]
  Dedup -- 否 --> Enq[事务写入日志/入队]
  Enq --> Handle[调用对端 API]
  Handle -->|成功| Done[完成]
  Handle -->|失败| Retry[记录重试/退避]
```
- CNB → 集成服务 → Plane：
  - 由 `.cnb.yml` 流水线在以下事件中回调：Issue `issue.open` / `issue.close` / `issue.reopen`；PR `pull_request` / `pull_request.update` / `pull_request.target` / `pull_request.mergeable` / `pull_request.merged` / `pull_request.approved` / `pull_request.changes_requested` / `pull_request.comment`。
  - 集成服务使用 `CNB_APP_TOKEN` 拉取详情，创建/更新/评论 Work Item，按 `pr_state_mappings` 推动状态同步。
- Plane → 集成服务 → CNB（参考 Plane webhook 文档）：
  - 事件：`issue`、`issue_comment` 为主。
  - 行为：标题/描述/标签/评论变更反写 CNB；Plane State→CNB issue open/close（当 `sync_direction=bidirectional`）。

**安全与幂等**
- 入站鉴权：
  - Plane：校验 `X-Plane-Signature`（HMAC-SHA256(secret, raw_body)）。
  - 来自 CNB 流水线的回调：`Authorization: Bearer {INTEGRATION_SHARED_SECRET}`。
- Token 安全：对 `access_token`/`refresh_token` 透明加密存储；最小权限；定期轮换。
- 幂等处理：Plane 侧以 `delivery_id`+`payload_sha256`；CNB 回调侧按（event+repo+iid+时间窗口）生成指纹；业务侧以 `issue_links`/`PR 引用` 防重复创建。

**错误与重试**
- 同步失败写入 `event_deliveries` 并进入 `jobs` 队列，指数退避（10m/30m/...）。
- 明确不可重试错误（4xx 语义错误）与可重试错误（网络/429/5xx）。

**Go/Echo 参考实现结构**
- 目录：
  - `cmd/server/main.go`（引导、路由、依赖注入）
  - `internal/handlers/plane.go`（Plane Webhook 回调）
  - `internal/handlers/cnb_ingest.go`（接收 `.cnb.yml` 回调入口）
  - `internal/services/sync.go`（核心映射/路由）
  - `internal/connectors/plane.go`、`internal/connectors/cnb.go`（各平台 API 客户端）
  - `internal/store/...`（Repo/Project/User/PR 映射、凭据）
  - `internal/jobs/...`（异步任务/重试）

**关键流程（文字时序）**
- CNB Issue 事件（open/close/reopen）→ `.cnb.yml` 回调 → 服务查映射（repo→project）→ 拉取 Issue 详情（含标签）→ 若含 Plane 标签则在 Plane 创建/链接 Work Item → 记录 `issue_links` → 回写 CNB 评论附 Plane 链接。
- Plane Work Item 打上 `CNB` 标签 → Plane Webhook → 查映射 → 在 CNB 创建 Issue → 记录 `issue_links` → 回写 Plane 评论附 CNB 链接。
- PR 标题含 `[ABC-123]` → `.cnb.yml` 在 PR 事件中回调 → 查 `pr_state_mappings` → 推动 Plane Work Item 状态流转；反向在 PR 留言引用。

**.cnb.yml 配置示例（关键片段）**
- 在 Issue 和 PR 相关事件中，以内置任务调用集成服务，实现 CNB→Plane 同步。
```yaml
$:
  issue.open:
    - stages:
        - name: notify plane sync (issue open)
          image: curlimages/curl:8.7.1
          script: >
            curl -sS -X POST "$INTEGRATION_URL/ingest/cnb/issue" \
              -H "Authorization: Bearer $INTEGRATION_TOKEN" \
              -H "Content-Type: application/json" \
              -d '{"event":"issue.open","repo":"'"$CNB_REPO_SLUG"'","issue_iid":"'"$CNB_ISSUE_IID"'"}'

  issue.close:
    - stages:
        - name: notify plane sync (issue close)
          image: curlimages/curl:8.7.1
          script: >
            curl -sS -X POST "$INTEGRATION_URL/ingest/cnb/issue" \
              -H "Authorization: Bearer $INTEGRATION_TOKEN" \
              -H "Content-Type: application/json" \
              -d '{"event":"issue.close","repo":"'"$CNB_REPO_SLUG"'","issue_iid":"'"$CNB_ISSUE_IID"'"}'

  issue.reopen:
    - stages:
        - name: notify plane sync (issue reopen)
          image: curlimages/curl:8.7.1
          script: >
            curl -sS -X POST "$INTEGRATION_URL/ingest/cnb/issue" \
              -H "Authorization: Bearer $INTEGRATION_TOKEN" \
              -H "Content-Type: application/json" \
              -d '{"event":"issue.reopen","repo":"'"$CNB_REPO_SLUG"'","issue_iid":"'"$CNB_ISSUE_IID"'"}'

  pull_request:
    - stages:
        - name: notify plane sync (pr)
          image: curlimages/curl:8.7.1
          script: >
            curl -sS -X POST "$INTEGRATION_URL/ingest/cnb/pr" \
              -H "Authorization: Bearer $INTEGRATION_TOKEN" \
              -H "Content-Type: application/json" \
              -d '{"event":"pull_request","action":"'"$CNB_PULL_REQUEST_ACTION"'","repo":"'"$CNB_REPO_SLUG"'","pr_iid":"'"$CNB_PULL_REQUEST_IID"'"}'
```

提示：如何注入 `$INTEGRATION_URL` / `$INTEGRATION_TOKEN`
- 使用 `.cnb.yml` 的 `env` 或 `imports` 将密钥仓库中的变量注入为环境变量，避免明文暴露。
  - 参考：docs/cnb-docs/docs/build/env.md:9、docs/cnb-docs/docs/repo/secret.md:54
  - 环境变量命名限制参考：docs/cnb-docs/docs/build/env.md:61

**证据（CNB 无 Webhook，使用内置任务/事件）**
- 事件列表：`docs/cnb-docs/docs/build/trigger-rule.md:282`
  - 列出了 Issue 事件：`issue.open` / `issue.close` / `issue.reopen`；PR 事件：`pull_request` 系列。
- 内置任务：`docs/cnb-docs/docs/build/internal-steps/README.md:135`
  - 展示流水线 `script/commands` 能力，可用 curl 发送 HTTP 回调。
- 环境变量：`docs/cnb-docs/docs/build/build-in-env.md:286`, `docs/cnb-docs/docs/build/build-in-env.md:352`
  - 提供 `CNB_ISSUE_*` 与 `CNB_PULL_REQUEST_*` 变量，可拼装回调负载。

**测试与验收**
- 单元：映射、签名校验、幂等、方向策略。
- 集成：模拟 Webhook 负载（Plane 示例见 intro-webhooks.mdx），本地回调。
- 试点：单项目、双向同步关闭（仅 CNB→Plane）；稳定后开放双向与 PR 自动化。

**里程碑**
- M1：最小可用（CNB→Plane 单向：Issue 创建/更新/评论、基础映射 UI/接口）。
- M2：Plane→CNB 双向、用户映射、评论归属。
- M3：PR 生命周期自动化、引用行为、性能与可观测性完善。
