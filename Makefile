# Variables
SERVICE ?= posts
BINARY = $(SERVICE)
DOCKER_IMAGE = hornet/$(SERVICE):latest

# Build the binary
.PHONY: build
build:
	@echo "Building $(SERVICE) app..."
	@go build -o ./bin/$(BINARY) ./cmd/$(SERVICE)/main.go

# Build the Docker container
.PHONY: build-container
build-container:
	@echo "Building Docker container for $(SERVICE)..."
	@docker build --build-arg SERVICE=$(SERVICE) -t $(DOCKER_IMAGE) .

# Run the binary directly
.PHONY: run
run:
	@echo "Running $(SERVICE) binary..."
	@./bin/$(BINARY)

# Run lint checks
.PHONY: lint
lint:
	@echo "Linting code..."
	@golangci-lint run ./...

# Format the Go code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Cleanup build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up binaries..."
	@rm -rf ./bin/$(BINARY)

# Help
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make build			   Build the binary for the specified service"
	@echo "  make build-container  Build the Docker container for the specified service"
	@echo "  make run			   Run the service binary directly"
	@echo "  make lint             Run lint checks"
	@echo "  make fmt              Format the Go code"
	@echo "  make clean            Clean up build artifacts"
	@echo "  make help             Display this help message"
	@echo
	@echo "Variables:"
	@echo "  SERVICE               Specify the service to build (e.g., 'posts', 'reactions', 'connections'). Default Value is 'posts'"
	@echo
