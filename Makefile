APP_NAME=init-codex
PORT?=8080

.PHONY: run build tidy test format swag docs

run: ## Run the API locally with pretty logs
	PRETTY_LOGS=true go run ./cmd/api

build: ## Build the API binary
	go build -o bin/$(APP_NAME) ./cmd/api

tidy:
	go mod tidy

test:
	go test ./...

format: ## Format all Go code
	go fmt ./...
	gofmt -s -w .

swag: ## Install swag CLI
	go install github.com/swaggo/swag/cmd/swag@latest

docs: swag ## Generate Swagger docs
	swag init -g cmd/api/main.go -o internal/docs --parseDependency --parseInternal

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

