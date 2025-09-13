.PHONY: help build run test clean docker-build docker-run swagger

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build the application
build: ## Build the Go application
	go build -o bin/api main.go

# Run the application
run: ## Run the application locally
	go run main.go

# Run tests
test: ## Run tests
	go test -v ./...

# Clean build artifacts
clean: ## Clean build artifacts
	rm -rf bin/
	go clean

# Install dependencies
deps: ## Install dependencies
	go mod download
	go mod tidy

# Generate Swagger documentation
swagger: ## Generate Swagger documentation
	swag init -g main.go -o ./docs

# Build Docker image
docker-build: ## Build Docker image
	docker build -t go-api-template .

# Run with Docker Compose
docker-run: ## Run with Docker Compose
	docker-compose up --build

# Stop Docker Compose
docker-stop: ## Stop Docker Compose
	docker-compose down

# Format code
fmt: ## Format Go code
	go fmt ./...

# Lint code
lint: ## Lint Go code
	golangci-lint run

# Install development tools
install-tools: ## Install development tools
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest