.PHONY: help run build test clean docker-build docker-run docker-stop lint fmt deps dev

help:
	@echo "Available commands:"
	@echo "  make run          - Run the application locally"
	@echo "  make build        - Build the application"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run application in Docker"
	@echo "  make docker-stop  - Stop Docker containers"
	@echo "  make lint         - Run linter"
	@echo "  make fmt          - Format code"
	@echo "  make deps         - Download dependencies"
	@echo "  make dev          - Run in development mode with hot reload"

run:
	go run cmd/api/main.go

build:
	go build -ldflags="-w -s" -o bin/api cmd/api/main.go

test:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf bin/ coverage.out coverage.html

docker-build:
	docker build -t go-api-template:latest .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

lint:
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run

fmt:
	go fmt ./...
	gofmt -s -w .

deps:
	go mod download
	go mod tidy

dev:
	@if ! command -v air &> /dev/null; then \
		echo "Installing air for hot reload..."; \
		go install github.com/air-verse/air@latest; \
	fi
	air