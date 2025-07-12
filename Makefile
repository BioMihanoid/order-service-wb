.PHONY: goose-install migrate-up migrate-down migrate-status migrate-create migrate-force

goose-install:
	@echo "Installing goose..."
	go install github.com/pressly/goose/v3/cmd/goose@latest

include .env
export

check-env:
	@if [ -z "${POSTGRES_HOST}" ] || [ -z "${POSTGRES_PORT}" ] || [ -z "${POSTGRES_USER}" ] || [ -z "${POSTGRES_PASSWORD}" ] || [ -z "${POSTGRES_DB}" ]; then \
		echo "Error: Database environment variables not set"; \
		exit 1; \
	fi

migrate-up: check-env
	@echo "Applying all available migrations..."
	goose -dir=migrations postgres "host=${POSTGRES_HOST} port=${POSTGRES_PORT} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${order_service_wb_db} sslmode=disable" up

migrate-down: check-env
	@echo "Reverting last migration..."
	goose -dir=migrations postgres "host=${POSTGRES_HOST} port=${POSTGRES_PORT} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${order_service_wb_db} sslmode=disable" down

migrate-down-to: check-env
	@echo "Reverting migrations down to version $(version)..."
	goose -dir=migrations postgres "host=${POSTGRES_HOST} port=${POSTGRES_PORT} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${order_service_wb_db} sslmode=disable" down-to $(version)

migrate-status: check-env
	@echo "Migration status:"
	goose -dir=migrations postgres "host=${POSTGRES_HOST} port=${POSTGRES_PORT} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${order_service_wb_db} sslmode=disable" status

migrate-create:
	@echo "Creating new migration files..."
	goose -dir=migrations create $(name) sql

migrate-force: check-env
	@echo "Forcing version to $(version)..."
	goose -dir=migrations postgres "host=${POSTGRES_HOST} port=${POSTGRES_PORT} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${order_service_wb_db} sslmode=disable" force $(version)

up: migrate-up
down: migrate-down
status: migrate-status
create: migrate-create
force: migrate-force

run:
	go run ./cmd/app

generator:
	go run ./cmd/generator

test:
	go test cover ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down

lint:
	golangci-lint run