APP ?= cabb
IMAGE ?= $(APP):dev
PORT ?= 8080

.PHONY: build run docker-build docker-run migrate ci-verify

build:
	go build -o bin/server ./cmd/server

run:
	@set -a; [ -f .env ] && . ./.env || true; set +a; \
	PORT=$(PORT) go run ./cmd/server

docker-build:
	docker build -t $(IMAGE) .

docker-run: docker-build
	docker run --rm -p $(PORT):$(PORT) -e PORT=$(PORT) $(IMAGE)

migrate:
	@: $${DATABASE_URL?Set DATABASE_URL}; \
	psql "$$DATABASE_URL" -v ON_ERROR_STOP=1 -f db/migrations/0001_init.sql; \
	psql "$$DATABASE_URL" -v ON_ERROR_STOP=1 -f db/migrations/0002_drop_plane_credentials.sql; \
	psql "$$DATABASE_URL" -v ON_ERROR_STOP=1 -f db/migrations/0003_add_workspace_slug_to_mappings.sql; \
	psql "$$DATABASE_URL" -v ON_ERROR_STOP=1 -f db/migrations/0004_cleanup_unused_tables.sql

ci-verify: build
	./scripts/start-server-for-ci.sh 18080 20 2; \
	SERVER_PID=$$(cat server.pid 2>/dev/null || echo ""); \
	if [ -n "$$SERVER_PID" ]; then \
		kill $$SERVER_PID 2>/dev/null || true; \
	fi; \
	rm -f server.pid
