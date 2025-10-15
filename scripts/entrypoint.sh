#!/bin/sh
set -euo pipefail

if [ -n "${DATABASE_URL:-}" ]; then
  echo "[entrypoint] applying migrations to DATABASE_URL"
  # ensure meta table exists
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -c "CREATE TABLE IF NOT EXISTS schema_migrations (filename text PRIMARY KEY, applied_at timestamptz NOT NULL DEFAULT now())" || {
    echo "[entrypoint] failed to ensure schema_migrations" >&2; exit 1; }
  for f in $(ls -1 /app/db/migrations/*.sql 2>/dev/null | sort); do
    base=$(basename "$f")
    exists=$(psql "$DATABASE_URL" -Atc "SELECT 1 FROM schema_migrations WHERE filename='${base}' LIMIT 1" || true)
    if [ "$exists" = "1" ]; then
      echo "[entrypoint] skip already applied $base"
      continue
    fi
    echo "[entrypoint] applying $base"
    psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f "$f" || {
      echo "[entrypoint] migration failed at $base" >&2
      exit 1
    }
    psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -c "INSERT INTO schema_migrations (filename) VALUES ('${base}') ON CONFLICT (filename) DO NOTHING" || {
      echo "[entrypoint] failed to record migration $base" >&2
      exit 1
    }
  done
else
  echo "[entrypoint] DATABASE_URL not set, skipping migrations"
fi

echo "[entrypoint] starting server on PORT=${PORT:-8080}"
exec /server
