APP ?= plane-integration
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
	psql "$$DATABASE_URL" -v ON_ERROR_STOP=1 -f db/migrations/0002_unique_workspaces.sql; \
	psql "$$DATABASE_URL" -v ON_ERROR_STOP=1 -f db/migrations/0003_pr_links.sql; \
	psql "$$DATABASE_URL" -v ON_ERROR_STOP=1 -f db/migrations/0004_label_mappings.sql

ci-verify: build
	PORT=18080 ./bin/server & echo $$! > .server.pid; \
	sleep 2; curl -fsS http://localhost:18080/healthz; \
	kill `cat .server.pid`; rm .server.pid
