# Plane Integration Service – Architecture Scaffold

This repository hosts a Go service that bridges Plane with CNB (internal code hosting) and Feishu (Lark). It follows the design in `docs/design/*` and provides a minimal, runnable scaffold you can iterate on.

## Tech Stack
- Go 1.24+ (Echo + Postgres 16)
- Echo HTTP server (wired, placeholders only)
- Postgres migrations in `db/migrations/`

## Layout
- `cmd/server/main.go` – entrypoint, HTTP server bootstrap
- `internal/handlers/*` – HTTP handlers (Plane OAuth/Webhook, CNB ingest, Lark endpoints, admin, jobs)
- `internal/store` – DB wiring placeholder
- `internal/version` – version string
- `pkg/config` – env configuration loader
- `db/migrations` – SQL schema (tables/enums/indexes per design)

## Endpoints (scaffold)
- `GET /healthz` – liveness with version & timestamp
- `GET /plane/oauth/start` – redirect to Plane consent (authorize-app)
- `GET /plane/oauth/callback` – handle app_installation_id/code, exchange tokens, return summary
- `POST /webhooks/plane` – HMAC verification stub (`X-Plane-Signature`)
- `POST /ingest/cnb/issue|pr|branch` – `.cnb.yml` callbacks with Bearer auth (`INTEGRATION_TOKEN`)
- `POST /webhooks/lark/events|interactivity|commands` – Feishu endpoints (challenge handled)
  - 群聊绑定：在群内 @ 机器人执行 `/bind <Plane Issue 链接>`（或 `绑定 <链接>`）后，记录 `thread_links(lark_thread_id↔plane_issue_id)`，Plane 侧的“更新/评论”通过线程回复推送到该话题（M1 文本，卡片待后续）。
  - 线程评论：在已绑定话题中回复 `/comment <文本>`（或 `评论 <文本>`）将把文本追加为 Plane Issue 评论（需要能解析出 `workspace_slug` 与 `project_id`）。
- `POST /admin/mappings/repo-project|pr-states|users|channel-project` – stubs
- `POST /jobs/issue-summary/daily` – stub (202 Accepted)

## Configure
Copy `.env.example` to `.env` (or export envs) and run the server:

```
go run ./cmd/server
```

Note: external dependencies (Echo) will be fetched at build time.

## Next Steps
- Wire database connection and run migrations on startup
- Implement token encryption and connectors (Plane/CNB/Lark/AI)
- Build Sync Core, idempotency store, retries & scheduler
- Flesh out admin APIs and webhook processing per `docs/design/*`

## Redis (optional)
Redis is not required for the minimal architecture. Postgres can back durable idempotency (`event_deliveries`), mappings and retries. Redis can be introduced later for:
- Hot idempotency cache: short-circuit duplicates before DB, reduce QPS on hot keys.
- Distributed locks: avoid double-processing across instances for the same key.
- Rate limiting/burst control: protect Plane/CNB APIs and internal workers.
- Ephemeral queues: buffering heavy tasks before handing to workers (optional; DB-backed queues also work).

Recommendation: start without Redis for M1. Add behind a feature flag (`REDIS_URL`) when we scale out workers or need stricter rate limiting/locks.
