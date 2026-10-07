.PHONY: build-all test vet migrate-up migrate-down docker-up docker-down dev e2e content-validate content-sync test-docker

DATABASE_URL ?= postgres://prepio:prepio@localhost:5432/prepio?sslmode=disable

build-all:
	go build ./...

test:
	go test ./... -count=1 -timeout 20m

test-docker:
	./scripts/test-docker.sh

test-short:
	go test ./services/gateway/... ./shared/... -count=1 -timeout 2m

vet:
	go vet ./...

migrate-up:
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "$(DATABASE_URL)" up; \
	else \
		docker compose --profile services run --rm migrate; \
	fi

migrate-down:
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path migrations -database "$(DATABASE_URL)" down 1; \
	else \
		echo "install golang-migrate CLI or use docker compose run migrate manually"; \
		exit 1; \
	fi

content-validate:
	go run ./services/question/cmd/content-sync --validate

content-sync:
	go run ./services/question/cmd/content-sync

docker-up:
	docker compose up -d postgres redis kafka

docker-down:
	docker compose down

docker-up-all:
	docker compose --profile services up -d

dev:
	./scripts/start-dev.sh

dev-all:
	./scripts/start-all.sh

e2e:
	node scripts/lesson-smoke.mjs
