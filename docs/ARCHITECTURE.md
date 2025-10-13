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
- `GET /plane/oauth/start` – placeholder
- `GET /plane/oauth/callback` – placeholder
- `POST /webhooks/plane` – HMAC verification stub (`X-Plane-Signature`)
- `POST /ingest/cnb/issue|pr|branch` – `.cnb.yml` callbacks with Bearer auth (`INTEGRATION_TOKEN`)
- `POST /webhooks/lark/events|interactivity|commands` – Feishu endpoints (challenge handled)
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
