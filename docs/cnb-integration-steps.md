# CNB 集成操作步骤（从零到打通）

先给结论（要做什么 / 验收点）
- 安装并授权 Plane 应用，使本服务获得工作区上下文（Bot Token、workspace_slug）。
- 建立“仓库 ↔ 项目”映射，完成 Plane Project 与 CNB Repo 的关联（必需）。
- 在 CNB 的 `.cnb.yml` 配置回调（issue/pr/branch）至本服务 `/ingest/cnb/*`（入站单向即可最小闭环）。
- 可选：开启 Plane → CNB 出站以实现双向同步与 fan‑out（按标签路由到多个仓库）。
- 验收：
  - CNB 侧新开 Issue（或模拟回调）后，Plane 中自动创建并建立链接；关闭/评论可同步到 Plane。
  - 开启出站后，Plane 中创建/更新/关闭会写回 CNB；评论会双向同步（HTML）。

——

## 前置准备
- 运行服务与数据库
  - 执行迁移：`psql "$DATABASE_URL" -f db/migrations/0001_init.sql`
  - 启动：`go run ./cmd/server`，健康检查：`GET http://localhost:8080/healthz`
- 环境变量（关键）
  - `INTEGRATION_TOKEN`：CNB 回调鉴权（Bearer），在 `.cnb.yml` 中注入同名变量。
  - `PLANE_*`：Plane OAuth 与 Webhook 校验（参考 `.env.example`）。
  - `CNB_OUTBOUND_ENABLED`、`CNB_BASE_URL`、`CNB_APP_TOKEN`：开启 Plane → CNB 双向写回时需要。
- 端点一览（以代码为准）：
  - Plane OAuth/Webhook：`/plane/oauth/start`、`/plane/oauth/callback`、`/webhooks/plane`（见 `internal/handlers/router.go:21`、`internal/handlers/router.go:23`）。
  - CNB 入站：`/ingest/cnb/issue|pr|branch`（见 `internal/handlers/router.go:26`–`internal/handlers/router.go:28`）。
  - 管理映射：`/admin/mappings/repo-project|pr-states|labels|users`（见 `internal/handlers/router.go:36`–`internal/handlers/router.go:39`）。

## 步骤一：安装 Plane 应用（获取工作区上下文）
- 浏览器打开：`GET /plane/oauth/start` → 跳转 Plane 同意页 → 回跳 `GET /plane/oauth/callback`。
- 回调处理：优先交换 Bot Token（`app_installation_id`），并查询安装信息（工作区 ID 与 slug）。
- 参考实现：`internal/handlers/plane.go`，签名验签：`POST /webhooks/plane`。
- 成功标志：回调返回安装摘要 JSON（不回显敏感令牌）；数据库表 `workspaces` 有对应记录。

最小自测（本地）：
```
open "http://localhost:8080/plane/oauth/start?state=dev"
# Plane 安装后会回到 /plane/oauth/callback?app_installation_id=...&code=...
```

## 步骤二：关联 Plane Project 与 CNB Repo（核心映射）
通过管理端点创建“仓库 ↔ 项目”映射，确定同步归属与路由策略。

- 端点：`POST /admin/mappings/repo-project`（见 `internal/handlers/admin.go:11`）。
- 请求体（必填字段：`cnb_repo_id`、`plane_workspace_id`、`plane_project_id`）：
```
curl -X POST "$INTEGRATION_URL/admin/mappings/repo-project" \
  -H "Content-Type: application/json" \
  -d '{
    "cnb_repo_id": "org/repo-a",
    "plane_workspace_id": "<plane_workspace_uuid>",
    "plane_project_id": "<plane_project_uuid>",
    "issue_open_state_id": "<可选_open_state_uuid>",
    "issue_closed_state_id": "<可选_closed_state_uuid>",
    "sync_direction": "cnb_to_plane",   
    "label_selector": "*",               
    "active": true
  }'
```
- 字段说明：
  - `cnb_repo_id`：仓库标识（示例：`org/repo`）。
  - `plane_workspace_id`：工作区 UUID（来自回调查询或 Plane API）。
  - `plane_project_id`：项目 UUID（Plane UI 或 API 获取）。
  - `issue_open_state_id` / `issue_closed_state_id`：CNB 事件驱动 Plane 状态切换（可选）。
  - `sync_direction`：`cnb_to_plane`（默认单向）或 `bidirectional`（双向，需步骤四开启出站）。
  - `label_selector`：标签路由（多仓库 fan‑out，用于出站；逗号/空格/分号/竖线分隔，大小写不敏感；`*`/`all` 表示任意标签）。

说明
- 可为“同一 Plane 项目”配置多条映射（不同 `cnb_repo_id`）。出站开启后按 `label_selector` fan‑out 到多个仓库。
- 映射存储：表 `repo_project_mappings`（迁移 `db/migrations/0001_init.sql`、`0006_repo_project_label_selector.sql`）。

## 步骤三：在 CNB 配置 `.cnb.yml` 回调（入站）
本服务接受来自 `.cnb.yml` 的 HTTP 回调，事件入口如下：
- Issue：`POST /ingest/cnb/issue`
- PR：`POST /ingest/cnb/pr`
- Branch：`POST /ingest/cnb/branch`
- 鉴权：`Authorization: Bearer $INTEGRATION_TOKEN`
- 去重：建议设置 `X-CNB-Delivery`（任意唯一字符串），重复投递将按指纹返回 `200 duplicate`。
- 接收字段与事件（以代码为准）：
  - Issue（`internal/handlers/cnb_ingest.go:18` 起）：
    - 事件：`issue.open`、`issue.close`、`issue.reopen`、`issue.update`、`issue.comment`/`comment.create`。
    - 必填：`event`、`repo`、`issue_iid`；可选：`title`、`description`、`labels[]`、`assignees[]`、`comment_html|comment`。
  - PR（`internal/handlers/cnb_ingest.go:30` 起）：
    - 事件（规范化后取 `action`）：`opened`/`ready_for_review`、`review_requested`、`approved`、`merged`、`closed`。
    - 必填：`event`、`repo`、`pr_iid`；可选：`issue_iid`（若提供可建立/补充关联）。
  - Branch（`internal/handlers/cnb_ingest.go:38` 起）：
    - 事件：`create|created`、`delete|deleted`；必填：`event`、`repo`、`branch`；可选：`issue_iid`（用于反查 Plane Issue）。

`.cnb.yml` 示例（语法与内置变量名以 CNB 官方文档为准，待确认）
```
# .cnb.yml （片段，供思路参考，CI 语法/变量名需按 CNB 实际调整）
on:
  issue_opened:
    steps:
      - name: notify-plane-issue-open
        run: |
          curl -sS -X POST "$INTEGRATION_URL/ingest/cnb/issue" \
            -H "Authorization: Bearer $INTEGRATION_TOKEN" \
            -H "Content-Type: application/json" \
            -H "X-CNB-Delivery: issue-open-${CI_ISSUE_IID}-${CI_PIPELINE_ID}" \
            -d "{\"event\":\"issue.open\",\"repo\":\"$CI_REPO\",\"issue_iid\":\"$CI_ISSUE_IID\",\"title\":\"$CI_ISSUE_TITLE\"}"
  issue_commented:
    steps:
      - name: notify-plane-issue-comment
        run: |
          curl -sS -X POST "$INTEGRATION_URL/ingest/cnb/issue" \
            -H "Authorization: Bearer $INTEGRATION_TOKEN" \
            -H "Content-Type: application/json" \
            -H "X-CNB-Delivery: issue-comment-${CI_ISSUE_IID}-${CI_EVENT_ID}" \
            -d "{\"event\":\"issue.comment\",\"repo\":\"$CI_REPO\",\"issue_iid\":\"$CI_ISSUE_IID\",\"comment_html\":\"$CI_COMMENT_HTML\"}"
```

手动联调（不依赖 CI）
```
# 新开 Issue → 在 Plane 创建并建立链接
curl -X POST "$INTEGRATION_URL/ingest/cnb/issue" \
  -H "Authorization: Bearer $INTEGRATION_TOKEN" \
  -H "Content-Type: application/json" \
  -H "X-CNB-Delivery: demo-issue-open-42" \
  -d '{"event":"issue.open","repo":"org/repo-a","issue_iid":"42","title":"Demo from CNB"}'

# 关闭 Issue → 切换到映射的 closed 状态（若配置）
curl -X POST "$INTEGRATION_URL/ingest/cnb/issue" \
  -H "Authorization: Bearer $INTEGRATION_TOKEN" \
  -H "Content-Type: application/json" \
  -H "X-CNB-Delivery: demo-issue-close-42" \
  -d '{"event":"issue.close","repo":"org/repo-a","issue_iid":"42"}'
```

## 步骤四（可选）：开启 Plane → CNB 出站（双向同步）
- 打开环境变量（见 `pkg/config/config.go:57`–`pkg/config/config.go:60`）：
  - `CNB_OUTBOUND_ENABLED=true`
  - `CNB_BASE_URL=https://api.cnb.cool`
  - `CNB_APP_TOKEN=<CNB 平台颁发的访问令牌>`
- 出站行为（见 `internal/handlers/plane.go`）：
  - Issue create：按项目下的 repo‑project 映射与 `label_selector` fan‑out 到匹配的仓库（新建 Issue 并建立链接）。
  - Issue update/close：回写已链接的 CNB Issue（按默认 swagger 路径）。
  - Issue comment：将 Plane 评论 HTML 追加到所有已链接的 CNB Issue。
- CNB API 路径（默认，见 `internal/cnb/client.go:16`–`internal/cnb/client.go:28`）：
  - Create: `/{repo}/-/issues`
  - Update: `/{repo}/-/issues/{number}`
  - Comment: `/{repo}/-/issues/{number}/comments`
- 说明：如网关路径与 swagger 不一致，可在代码层自定义 `IssueCreatePath/IssueUpdatePath/IssueCommentPath`（高级用法，暂未暴露到环境变量）。

## 步骤五（可选）：PR 状态映射、标签与用户映射
- PR → Plane 状态映射（`POST /admin/mappings/pr-states`，见 `internal/handlers/admin.go:48` 起）：
```
curl -X POST "$INTEGRATION_URL/admin/mappings/pr-states" \
  -H "Content-Type: application/json" \
  -d '{
    "cnb_repo_id": "org/repo-a",
    "plane_project_id": "<plane_project_uuid>",
    "opened_state_id": "<uuid>",
    "review_requested_state_id": "<uuid>",
    "approved_state_id": "<uuid>",
    "merged_state_id": "<uuid>",
    "closed_state_id": "<uuid>"
  }'
```
- 标签映射（CNB 标签名 → Plane `label_id`，见 `internal/handlers/admin.go:103` 起）：
```
curl -X POST "$INTEGRATION_URL/admin/mappings/labels" \
  -H "Content-Type: application/json" \
  -d '{
    "cnb_repo_id": "org/repo-a",
    "plane_project_id": "<plane_project_uuid>",
    "items": [
      {"cnb_label": "bug", "plane_label_id": "<plane_label_uuid>"}
    ]
  }'
```
- 用户映射（CNB 用户 ID → Plane 用户 UUID，见 `internal/handlers/admin.go:80` 起）：
```
curl -X POST "$INTEGRATION_URL/admin/mappings/users" \
  -H "Content-Type: application/json" \
  -d '{
    "mappings": [
      {"cnb_user_id": "u123", "plane_user_id": "<uuid>", "display_name": "Alice"}
    ]
  }'
```

## 验证与排查
- 常见响应
  - `202 Accepted`：已入队异步处理。
  - `200 OK status=duplicate`：命中幂等；检查 `X-CNB-Delivery` 与 payload 是否一致。
  - `401 invalid_token`：`INTEGRATION_TOKEN` 不匹配（见 `internal/handlers/cnb_ingest.go:204` 起）。
  - `400/422`：缺少字段或 JSON 解析失败。
- 端到端最小 DoD
  - 完成 OAuth 安装回调；`/plane/oauth/callback` 返回安装摘要，`workspaces` 落库。
  - 创建一条 repo‑project 映射。
  - `issue.open` 回调创建 Plane Issue 并在 `issue_links` 写入关联；`issue.close`/`comment` 生效。
  - （可选出站）Plane 新建/更新/关闭/评论回写 CNB 成功。
- 实用提示
  - 观察日志的 `request_id/source/event_type` 与 `event_deliveries` 表记录（去重/重试）。
  - 多仓库 fan‑out：确保 Plane Issue 打上能命中 `label_selector` 的标签；未命中的仓库不会创建链接。

## 安全与契约
- 入站鉴权：`Authorization: Bearer $INTEGRATION_TOKEN`。
- Plane Webhook 验签：`X-Plane-Signature: sha256=...`，HMAC‑SHA256(secret, raw_body)。
- 时间与格式：`application/json; charset=utf-8`，时间使用 RFC3339 UTC。
- 幂等：`delivery_id + payload_sha256` 去重；重复返回 `200` 且 `status=duplicate`（可选）。

## 限制与待确认
- `.cnb.yml` 触发器名称与内置变量名以 CNB 官方文档为准；本文示例为思路参考（待确认）。
- 出站路径模板如需自定义，目前需改代码注入客户端的可选字段（未暴露为环境变量）。
- 令牌与机密的透明加密存储与续期策略将按里程碑逐步补齐。

