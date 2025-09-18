APP_NAME := [[ app_name ]]
BUILD_DIR := build
MAIN := ./main.go
DOCKER_COMPOSE := deployment/docker/docker-compose.yml
DOCKER_DIR := deployment/docker
VOLUMES_DIR := $(DOCKER_DIR)/volumes
INIT_DIR := $(DOCKER_DIR)/init
MIGRATIONS_DIR := infrastructure/db/migrations

COLOR_RESET=\033[0m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_RED=\033[31m

GOLANGCI_LINT := golangci-lint

.PHONY: all build run test lint lint-fix format docker-build docker-up docker-down migrate-up migrate-down clean

## Default target
all: build

## Build Go binary
build:
	@echo "$(COLOR_YELLOW)🔨 Building app...$(COLOR_RESET)"
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)
	@echo "$(COLOR_GREEN)✅ Build complete.$(COLOR_RESET)"

## Run app locally
run: build
	@echo "$(COLOR_YELLOW)🚀 Running app...$(COLOR_RESET)"
	@./$(BUILD_DIR)/$(APP_NAME)

## Install all required dev tools
setup-dev:
	@echo "$(COLOR_YELLOW)📦 Installing development tools...$(COLOR_RESET)"
	@go install github.com/segmentio/golines@latest
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/gordonklaus/ineffassign@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/client9/misspell/cmd/misspell@latest
	@echo "$(COLOR_GREEN)✅ Dev tools installed successfully.$(COLOR_RESET)"

## Run tests with race detector and coverage
test-c:
	@echo "$(COLOR_YELLOW)🧪 Running tests...$(COLOR_RESET)"
	@go test ./... -race -coverprofile=coverage.out -covermode=atomic
	@echo "$(COLOR_GREEN)✅ Tests passed.$(COLOR_RESET)"

## Run tests with race detector
test:
	@echo "$(COLOR_YELLOW)🧪 Running tests...$(COLOR_RESET)"
	@GO_ENV=test go test ./... -v -race
	@echo "$(COLOR_GREEN)✅ Tests passed.$(COLOR_RESET)"

## Run benchmark
bench:
	@echo "$(COLOR_YELLOW)🧪 Running tests...$(COLOR_RESET)"
	@GO_ENV=test go test ./... -bench=. -benchmem
	@echo "$(COLOR_GREEN)✅ Tests passed.$(COLOR_RESET)"

## Run linter
lint:
	@echo "$(COLOR_YELLOW)🧹 Running golangci-lint...$(COLOR_RESET)"
	@$(GOLANGCI_LINT) run --timeout 2m ./...
	@echo "$(COLOR_GREEN)✅ Lint passed.$(COLOR_RESET)"

## Run linter with autofix
lint-fix:
	@echo "$(COLOR_YELLOW)🔧 Running golangci-lint with --fix...$(COLOR_RESET)"
	@$(GOLANGCI_LINT) run --fix --timeout 2m ./... || true
	@echo "$(COLOR_GREEN)✅ Lint autofix done.$(COLOR_RESET)"

## Format code
format:
	@echo "$(COLOR_YELLOW)🎨 Formatting code...$(COLOR_RESET)"
	@goimports -w .
	@gofmt -s -w .
	@gofumpt -extra -w .
	@golines --max-len=100 -w .
	@go vet ./...
	@ineffassign ./...
	@misspell -w .
	@echo "$(COLOR_GREEN)✅ Code formatted & checked successfully.$(COLOR_RESET)"


## Build Docker image
docker-build:
	@echo "$(COLOR_YELLOW)🐳 Building Docker image...$(COLOR_RESET)"
	@docker build -t $(APP_NAME):latest -f deployment/docker/Dockerfile .
	@echo "$(COLOR_GREEN)✅ Docker image built.$(COLOR_RESET)"

## Start docker-compose
docker-up:
	@echo "$(COLOR_YELLOW)📦 Starting docker-compose...$(COLOR_RESET)"
	@docker-compose -f $(DOCKER_COMPOSE) up -d
	@echo "$(COLOR_GREEN)✅ Docker services up.$(COLOR_RESET)"

## Stop docker-compose
docker-down:
	@echo "$(COLOR_YELLOW)🛑 Stopping docker-compose...$(COLOR_RESET)"
	@docker-compose -f $(DOCKER_COMPOSE) down
	@echo "$(COLOR_GREEN)✅ Docker services down.$(COLOR_RESET)"

## Delete volumes & init folders
del-volumes: docker-down
	@echo "$(COLOR_YELLOW)🗑️ Removing volumes and init folders...$(COLOR_RESET)"
	@rm -rf $(VOLUMES_DIR) $(INIT_DIR)
	@echo "$(COLOR_GREEN)✅ Volumes and init folders removed.$(COLOR_RESET)"

## Reset docker environment (down → delete volumes/init → up)
reset-docker: del-volumes docker-up
	@echo "$(COLOR_GREEN)♻️ Docker environment reset complete.$(COLOR_RESET)"

## Run database migrations
migrate-up:
	@echo "$(COLOR_YELLOW)📂 Running migrations up...$(COLOR_RESET)"
	@migrate -path $(MIGRATIONS_DIR) -database "postgres://admin:admin@localhost:5432/boiler_db?sslmode=disable" up
	@echo "$(COLOR_GREEN)✅ Migrations applied.$(COLOR_RESET)"

migrate-down:
	@echo "$(COLOR_YELLOW)📂 Running migrations down...$(COLOR_RESET)"
	@migrate -path $(MIGRATIONS_DIR) -database "postgres://admin:admin@localhost:5432/boiler_db?sslmode=disable" down
	@echo "$(COLOR_GREEN)✅ Migrations rolled back.$(COLOR_RESET)"

## Clean build artifacts
clean:
	@echo "$(COLOR_YELLOW)🧹 Cleaning build artifacts...$(COLOR_RESET)"
	@rm -rf $(BUILD_DIR) coverage.out
	@echo "$(COLOR_GREEN)✅ Clean complete.$(COLOR_RESET)"

## Delete all local branches except master and dev
clean-branches:
	@echo "$(COLOR_YELLOW)🗑️ Cleaning local branches (except master & dev)...$(COLOR_RESET)"
	@git branch | grep -vE "master|dev" | xargs git branch -D
	@echo "$(COLOR_GREEN)✅ Local branches cleaned.$(COLOR_RESET)"

## Detect missing comments for exported funcs/types
autodoc:
	@echo "$(COLOR_YELLOW)📝 Checking for missing comments...$(COLOR_RESET)"
	@grep -R --include="*.go" -nE "^(type|func) [A-Z]" ./internal ./pkg \
	| grep -vE "^[[:space:]]*//" || true

## Run static analysis with go vet
check-vet:
	@echo "$(COLOR_YELLOW)🔍 Running static analysis...$(COLOR_RESET)"
	@go vet ./...
	@echo "$(COLOR_GREEN)✅ All checks passed.$(COLOR_RESET)"

## Download and tidy Go module dependencies
deps:
	@echo "📦 Downloading and cleaning dependencies..."
	@go mod tidy
	@go mod download
	@echo "✅ Dependencies are up to date."

## Run full project checks (deps → setup-dev → format → vet → tests)
check: deps setup-dev format check-vet test
	@echo "$(COLOR_GREEN)✅ All checks completed successfully.$(COLOR_RESET)"

## Run lightweight checks for pre-commit (lint + vet + format)
check-lite: lint-fix lint check-vet format
	@echo "$(COLOR_GREEN)✅ Lightweight checks completed.$(COLOR_RESET)"

## Init swagger to doc
swag:
	@echo "📝 Init swagger..."
	@ swag init --generalInfo main.go --output ./docs
	@echo "✅ Done."
