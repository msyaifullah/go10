
COVERAGE_DIR=coverage
COVERAGE_PROFILE=$(COVERAGE_DIR)/coverage.out
COVERAGE_HTML=$(COVERAGE_DIR)/coverage.html

# Variables
APP_NAME=go10
CLI_NAME=cli
MAIN_PATH=cmd/server/main.go
CLI_PATH=cmd/cli/main.go
DOCKER_IMAGE=go10-api
DOCKER_TAG=latest

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Build variables
VERSION?=1.0.0
BUILD_DIR=bin
COMMIT_SHA=$(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Environment variables
export GO111MODULE=on
export GOPATH=$(shell go env GOPATH)
export PATH:=$(GOPATH)/bin:$(PATH)

# Make all
.PHONY: all
all: clean deps setup-dev test build

# Clean build directory
.PHONY: clean
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(COVERAGE_DIR)
	@go clean -testcache

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v -race -cover ./...

# Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

#coverage
.PHONY: coverage
coverage:
	@echo "Generating coverage reports..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -race -coverprofile=$(COVERAGE_PROFILE) -covermode=atomic ./...
	@go tool cover -html=$(COVERAGE_PROFILE) -o $(COVERAGE_HTML)
	@go tool cover -func=$(COVERAGE_PROFILE)


.PHONY: coverage-check
coverage-check:
	@echo "Checking coverage threshold..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -race -coverprofile=$(COVERAGE_PROFILE) -covermode=atomic ./...
	@coverage_status=`go tool cover -func=$(COVERAGE_PROFILE) | grep total | awk '{print $$3}' | sed 's/%//'`; \
	if [ $$(echo "$$coverage_status < 80" | bc) -eq 1 ]; then \
		echo "Code coverage is below 80%. Current coverage: $$coverage_status%"; \
		exit 1; \
	fi

# Build application
.PHONY: build
build: build-server build-cli

# Build server application
.PHONY: build-server
build-server:
	@echo "Building server application..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-X main.Version=$(VERSION) -X main.CommitSHA=$(COMMIT_SHA) -X main.BuildDate=$(BUILD_DATE)" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Build CLI application
.PHONY: build-cli
build-cli:
	@echo "Building CLI application..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-X main.Version=$(VERSION) -X main.CommitSHA=$(COMMIT_SHA) -X main.BuildDate=$(BUILD_DATE)" -o $(BUILD_DIR)/$(CLI_NAME) $(CLI_PATH)
	@chmod +x $(BUILD_DIR)/$(CLI_NAME)

# Run database migrations
.PHONY: migrate
migrate:
	@echo "Running database migrations..."
	@./$(BUILD_DIR)/$(CLI_NAME) migrate

# Run database migrations with debug
.PHONY: migrate-debug
migrate-debug:
	@echo "Running database migrations with debug..."
	@./$(BUILD_DIR)/$(CLI_NAME) migrate --debug

# Rollback database migrations
.PHONY: migrate-down
migrate-down:
	@echo "Rolling back database migrations..."
	@./$(BUILD_DIR)/$(CLI_NAME) migrate-down

# Rollback database migrations with debug
.PHONY: migrate-down-debug
migrate-down-debug:
	@echo "Rolling back database migrations with debug..."
	@./$(BUILD_DIR)/$(CLI_NAME) migrate-down --debug

# Run application
.PHONY: run
run:
	@echo "Running server application..."
	@go run $(MAIN_PATH)

# Run application with Air (live reload)
.PHONY: dev
dev:
	@echo "Starting development server with Air (live reload)..."
	@mkdir -p tmp
	@air

# Install Air for development
.PHONY: install-air
install-air:
	@echo "Installing Air for live reloading..."
	@go install github.com/air-verse/air@latest

# Setup development environment (install Air + dependencies)
.PHONY: setup-dev
setup-dev: install-air deps
	@echo "Setting up configuration files..."
	@if [ ! -f configs/cli.yaml ]; then \
		cp configs/cli.yaml.template configs/cli.yaml; \
		echo "Created configs/cli.yaml from template"; \
	else \
		echo "configs/cli.yaml already exists, skipping..."; \
	fi
	@if [ ! -f configs/local.toml ]; then \
		cp configs/local.toml.template configs/local.toml; \
		echo "Created configs/local.toml from template"; \
	else \
		echo "configs/local.toml already exists, skipping..."; \
	fi
	
	@echo "Development environment setup complete!"
	@echo ""
	@echo "To start development with live reload:"
	@echo "  make dev"
	@echo ""
	@echo "To run without live reload:"
	@echo "  make run"

# Build docker image
.PHONY: docker-build
docker-build:
	@echo "Building docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run docker container
.PHONY: docker-run
docker-run:
	@echo "Running docker container..."
	@docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker compose up
.PHONY: docker-compose-up
docker-compose-up:
	@echo "Starting docker compose..."
	@docker compose up -d

# Docker compose down
.PHONY: docker-compose-down
docker-compose-down:
	@echo "Stopping docker compose..."
	@docker compose down

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo ""
	@echo "General:"
	@echo "  make all                    - Clean, download dependencies, run tests, and build"
	@echo "  make clean                  - Clean build directory"
	@echo "  make deps                   - Download dependencies"
	@echo "  make setup-dev              - Setup development environment (Air + deps)"
	@echo "  make test                   - Run tests"
	@echo "  make bench                  - Run benchmarks"
	@echo "  make coverage               - Run tests coverage"
	@echo "  make coverage-check         - Run coverage check threshold"
	@echo ""
	@echo "Server:"
	@echo "  make build-server           - Build server application"
	@echo "  make run                    - Run server application"
	@echo "  make dev                    - Run server with Air (live reload)"
	@echo "  make install-air            - Install Air for development"
	@echo ""
	@echo "CLI:"
	@echo "  make build-cli              - Build CLI application"
	@echo "  make migrate                - Run database migrations"
	@echo "  make migrate-debug          - Run database migrations with debug"
	@echo "  make migrate-down           - Rollback database migrations"
	@echo "  make migrate-down-debug     - Rollback database migrations with debug"
	@echo ""
	@echo "Build:"
	@echo "  make build                  - Build both server and CLI"
#	@echo "  make docker-build           - Build docker image"
#	@echo "  make docker-run             - Run docker container"
#	@echo "  make docker-compose-up      - Start docker compose"
#	@echo "  make docker-compose-down    - Stop docker compose"