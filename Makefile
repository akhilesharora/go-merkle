# Variables
BINARY_DIR=bin
SERVER_BINARY=$(BINARY_DIR)/server
CLIENT_BINARY=$(BINARY_DIR)/client
DOCKER_TAG=latest
APP_NAME=go-merkle-app
DOCKER_IMAGE_NAME=$(APP_NAME):$(DOCKER_TAG)

.PHONY: all build build-server build-client run-server run-client test test-coverage clean docker-compose-up docker-compose-down docker-build docker-build-server docker-build-client docker-build-ui help

# Default target
all: clean build test

# Build Flags
LDFLAGS=-ldflags "-s -w"

# Build Commands
GOBUILD=go build $(LDFLAGS)
GOTEST=go test -v -race

# Build all components
build: build-server build-client

# Build the server
build-server:
	@echo "Building server..."
	@$(GOBUILD) -o $(SERVER_BINARY) ./cmd/server

# Build the client
build-client:
	@echo "Building client..."
	@$(GOBUILD) -o $(CLIENT_BINARY) ./cmd/client

# Run the server
run-server: build-server
	@echo "Running server..."
	@./$(SERVER_BINARY)

# Run the client
run-client: build-client
	@echo "Running client..."
	@./$(CLIENT_BINARY)

# Run tests with race condition checks
test:
	@echo "Running tests..."
	@$(GOTEST) ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@$(GOTEST) -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean up binaries and coverage report
clean:
	@echo "Cleaning up..."
	@rm -rf $(BINARY_DIR) coverage.out coverage.html

# Docker Compose up
docker-compose-up:
	@echo "Starting Docker Compose services..."
	@docker-compose up --build -d

# Docker Compose down
docker-compose-down:
	@echo "Stopping Docker Compose services..."
	@docker-compose down

# Build all Docker images
docker-build: docker-build-server docker-build-client docker-build-ui

# Build server Docker image
docker-build-server:
	@echo "Building server Docker image..."
	@docker build -t $(APP_NAME)-server:$(DOCKER_TAG) -f Dockerfile.server .

# Build client Docker image
docker-build-client:
	@echo "Building client Docker image..."
	@docker build -t $(APP_NAME)-client:$(DOCKER_TAG) -f Dockerfile.client .

# Build UI Docker image
docker-build-ui:
	@echo "Building UI Docker image..."
	@docker build -t $(APP_NAME)-ui:$(DOCKER_TAG) -f Dockerfile.ui .

# Display help information
help:
	@echo "Usage: make [TARGET]"
	@echo ""
	@echo "Targets:"
	@echo "  all                Build all components and run tests"
	@echo "  build              Build server and client"
	@echo "  build-server       Build the server"
	@echo "  build-client       Build the client"
	@echo "  run-server         Build and run the server"
	@echo "  run-client         Build and run the client"
	@echo "  test               Run tests with race condition checks"
	@echo "  test-coverage      Run tests with coverage and generate a report"
	@echo "  clean              Remove binary files and coverage report"
	@echo "  docker-compose-up  Start Docker Compose services"
	@echo "  docker-compose-down Stop Docker Compose services"
	@echo "  docker-build       Build all Docker images"
	@echo "  docker-build-server Build server Docker image"
	@echo "  docker-build-client Build client Docker image"
	@echo "  docker-build-ui    Build UI Docker image"
	@echo "  help               Show this help message"