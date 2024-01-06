# Variables
BINARY_DIR=bin
APP_BINARY=$(BINARY_DIR)/app
DOCKER_TAG=latest
APP_NAME=go-merkle-app
DOCKER_IMAGE_NAME=$(APP_NAME):$(DOCKER_TAG)

.PHONY: all build run test test-coverage clean docker docker-build docker-run help

# Default target
all: clean build test

# Build Flags
LDFLAGS=-ldflags "-s -w"

# Build Commands
GOBUILD=go build $(LDFLAGS)
GOTEST=go test -v -race

# Build the application
build:
	@echo "Building application..."
	@$(GOBUILD) -o $(APP_BINARY) ./cmd

# Run the application
run: build
	@echo "Running application..."
	@./$(APP_BINARY)

# Run tests with race condition checks
test:
	@echo "Running tests..."
	@$(GOTEST) ./pkg/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@$(GOTEST) -coverprofile=coverage.out ./pkg/...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean up binaries and coverage report
clean:
	@echo "Cleaning up..."
	@rm -rf $(BINARY_DIR) coverage.out coverage.html

# Build dockerfile
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE_NAME) .

# Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run $(DOCKER_IMAGE_NAME)

# Display help information
help:
	@echo "Usage: make [TARGET]"
	@echo ""
	@echo "Targets:"
	@echo "  all            Build the application and run tests"
	@echo "  build          Build the application"
	@echo "  run            Build and run the application"
	@echo "  test           Run tests with race condition checks"
	@echo "  test-coverage  Run tests with coverage and generate a report"
	@echo "  clean          Remove binary files and coverage report"
	@echo "  docker         Build Docker image"
	@echo "  help           Show this help message"




