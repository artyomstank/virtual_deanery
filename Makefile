# Makefile

.PHONY: help build run test clean fmt lint docker-build docker-up docker-down migrate

BINARY_NAME=myapp
GO=go
DOCKER_IMAGE=myapp:latest

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: ## Build binary
	$(GO) build -o bin/$(BINARY_NAME) cmd/api/main.go

run: ## Run application
	$(GO) run cmd/api/main.go

test: ## Run tests
	$(GO) test -v -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	$(GO) tool cover -html=coverage.out

clean: ## Clean build artifacts
	rm -rf bin/
	$(GO) clean

fmt: ## Format code
	$(GO) fmt ./...
	gofmt -s -w .

lint: ## Run linter
	golangci-lint run

vet: ## Run go vet
	$(GO) vet ./...

mod-download: ## Download dependencies
	$(GO) mod download

mod-tidy: ## Tidy dependencies
	$(GO) mod tidy

docker-build: ## Build Docker image
	docker build -f build/package/Dockerfile -t $(DOCKER_IMAGE) .

docker-up: ## Start services with docker-compose
	docker-compose -f docker-compose.yml up -d

docker-down: ## Stop services with docker-compose
	docker-compose -f docker-compose.yml down

migrate-up: ## Run migrations up
	# TODO: Use migrate tool or custom script
	# migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down: ## Run migrations down
	# TODO: Use migrate tool
	# migrate -path migrations -database "$(DATABASE_URL)" down

dev: ## Run in development mode with hot reload
	# TODO: Use air or similar tool
	# air

.DEFAULT_GOAL := help
