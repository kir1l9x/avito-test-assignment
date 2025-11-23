APP_NAME=avito-test-assignment
MIGRATIONS_DIR=./migrations
DB_DSN=postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

include .env
export $(shell sed 's/=.*//' .env)

.PHONY: run stop stopv logs dev build fmt migrate-up migrate-down lint test test-handlers test-repositories test-env-up test-env-down

run:
	docker compose up --build -d

stop:
	docker compose down

stopv:
	docker compose down -v --remove-orphans

logs:
	docker compose logs -f app

dev:
	go run ./cmd/app

build:
	go build -o bin/app ./cmd/app

fmt:
	go fmt ./...

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down

lint:
	golangci-lint run ./...

test:
	go test ./... -cover

test-handlers:
	go test ./api/handlers -cover

test-repositories:
	go test ./internal/infrastructure/repositories -cover

test-env-up:
	docker compose -f docker-compose.test.yml up -d

test-env-down:
	docker compose -f docker-compose.test.yml down -v --remove-orphans
