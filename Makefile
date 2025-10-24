ifneq (,$(wildcard ./.env))
    include .env
    export
endif

## ENV
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):5432/$(POSTGRES_DB)?sslmode=disable
MIGRATIONS_DIR=internal/migration
DOCKER_COMPOSE_FILE=compose.yaml
SEED_FILE=seed.sql


.PHONY: run build clean

run:
	@echo "Starting server"
	@go run cmd/go_social/main.go

build:
	@echo "Building binary"
	@go build -o bin/api ./cmd/go_social

clean:
	@echo "Removing old build"
	@rm -rf bin

gen-data:
	@go run cmd/seed/seed.go

.PHONY: up down logs

up:
	@echo "Starting containers..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) up -d

down:
	@echo "Stopping containers..."
	@docker compose -f $(DOCKER_COMPOSE_FILE) down

logs:
	@docker compose logs -f


.PHONY:migrate-up migrate-down migrate-new migration-reset

migrate-up:
	@echo "Running migrations..."
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-reset:
	@echo "Removing all migrations..."
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" drop -f

migrate-new:
	@if [ -z "$(name)" ]; then echo "name is empty"; exit 1; fi
	@echo "Creating new migration: $(name)"
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)


.PHONY: test swagger

test:
	@echo "Running tests..."
	@go test ./...

swagger:
	@echo "Generation docs..."
	@swag init --parseDependency --parseInternal -g  cmd/go_social/main.go -o docs
