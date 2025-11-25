.PHONY: all build run test clean docker-up docker-down swagger migrate seed

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=inventorypulse
MAIN_PATH=./cmd/api

# Build the application
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	$(GORUN) $(MAIN_PATH)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -cover -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Start PostgreSQL container
docker-up:
	docker-compose up -d

# Stop PostgreSQL container
docker-down:
	docker-compose down

# Stop and remove volumes
docker-clean:
	docker-compose down -v

# Generate Swagger documentation
swagger:
	swag init -g cmd/api/main.go -o docs

# Install swag CLI tool
install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Show help
help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make deps         - Download and tidy dependencies"
	@echo "  make docker-up    - Start PostgreSQL container"
	@echo "  make docker-down  - Stop PostgreSQL container"
	@echo "  make docker-clean - Stop container and remove volumes"
	@echo "  make swagger      - Generate Swagger documentation"
	@echo "  make install-swag - Install swag CLI tool"
	@echo "  make dev          - Run with hot reload (requires air)"

