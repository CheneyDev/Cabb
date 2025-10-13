#!/bin/sh
set -euo pipefail

if [ -n "${DATABASE_URL:-}" ]; then
  echo "[entrypoint] applying migrations to DATABASE_URL"
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f /app/db/migrations/0001_init.sql || {
    echo "[entrypoint] migration failed" >&2
    exit 1
  }
else
  echo "[entrypoint] DATABASE_URL not set, skipping migrations"
fi

echo "[entrypoint] starting server on PORT=${PORT:-8080}"
exec /server

