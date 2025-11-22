# Simple Makefile for a Go project

include .env
export

DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

all: build test

build:
	@echo "Building..."
	@go build -o main.exe cmd/api/main.go

run:
	@go run cmd/api/main.go

docker-run:
	@docker compose up --build

docker-down:
	@docker compose down

test:
	@echo "Testing..."
	@go test ./... -v

itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

clean:
	@echo "Cleaning..."
	@rm -f main

watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
       air; \
       Write-Output 'Watching...'; \
    } else { \
       Write-Output 'Installing air...'; \
       go install github.com/air-verse/air@latest; \
       air; \
       Write-Output 'Watching...'; \
    }"

migrate-create:
	@powershell -Command "$$name = Read-Host 'Enter migration name'; migrate create -ext sql -dir migrations -seq $$name"

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

migrate-rollback:
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-fresh:
	migrate -path migrations -database "$(DB_URL)" drop -f
	migrate -path migrations -database "$(DB_URL)" up

migrate-status:
	migrate -path migrations -database "$(DB_URL)" version

migrate-force:
	@powershell -Command "$$version = Read-Host 'Enter version to force'; migrate -path migrations -database '$(DB_URL)' force $$version"

.PHONY: all build run test clean watch docker-run docker-down itest migrate-create migrate-up migrate-down migrate-rollback migrate-fresh migrate-status migrate-force
