.PHONY: help build test clean fmt lint install push pull validate examples all

# Default target
all: fmt lint build test

# Help target
help:
	@echo "Available targets:"
	@echo "  make build      - Build all packages"
	@echo "  make test       - Run all tests"
	@echo "  make fmt        - Format Go code"
	@echo "  make lint       - Run Go linter"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make install    - Install the package"
	@echo "  make push       - Build and run push example"
	@echo "  make pull       - Build and run pull example"
	@echo "  make examples   - Build example binaries"
	@echo "  make all        - Run fmt, lint, build, and test"

# Build all packages
build:
	@echo "Building all packages..."
	@go build ./...
	@echo "✓ Build complete"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...
	@echo "✓ Tests complete"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@gofmt -w .
	@echo "✓ Format complete"

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, running go vet instead"; \
		go vet ./...; \
	fi
	@echo "✓ Lint complete"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@go clean
	@rm -f bin/push bin/pull
	@rm -rf dist/
	@echo "✓ Clean complete"

# Install the package
install:
	@echo "Installing package..."
	@go install ./...
	@echo "✓ Install complete"

# Build example binaries
examples: bin/push bin/pull
	@echo "✓ Examples built"

bin/push:
	@mkdir -p bin
	@echo "Building push example..."
	@go build -o bin/push examples/push/main.go

bin/pull:
	@mkdir -p bin
	@echo "Building pull example..."
	@go build -o bin/pull examples/pull/main.go

# Run push example
push: bin/push
	@echo "Running push example..."
	@echo "Usage: bin/push -spec <file> -registry <url> [-tag <tag>] [-description <desc>] [-source <url>]"
	@echo ""
	@echo "Example:"
	@echo "  bin/push -spec test-spec.yaml -registry ghcr.io/myorg/myartifact -tag latest"

# Run pull example  
pull: bin/pull
	@echo "Running pull example..."
	@echo "Usage: bin/pull -ref <reference> [-output <file>] [-plain-http]"
	@echo "   OR: bin/pull -registry <url> -digest <sha256> [-output <file>] [-plain-http]"
	@echo ""
	@echo "Example:"
	@echo "  bin/pull -ref ghcr.io/myorg/myartifact:latest -output spec.yaml"

# Development targets
.PHONY: dev watch

# Run in development mode with file watching (requires entr)
watch:
	@if command -v entr > /dev/null; then \
		find . -name '*.go' | entr -c make build; \
	else \
		echo "entr not installed. Install with: brew install entr (macOS) or apt-get install entr (Linux)"; \
	fi

# Quick development build
dev: fmt build
	@echo "✓ Development build complete"

# Testing targets
.PHONY: test-unit test-integration test-coverage

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test -v -short ./...

# Run integration tests (if any)
test-integration:
	@echo "Running integration tests..."
	@go test -v -run Integration ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# Docker targets (if needed in future)
.PHONY: docker-build docker-push

# Build Docker image
docker-build:
	@echo "Docker support not yet implemented"
	@echo "Add a Dockerfile to enable this target"

# Dependency management
.PHONY: deps deps-update deps-tidy

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@echo "✓ Dependencies downloaded"

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✓ Dependencies updated"

# Tidy dependencies
deps-tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "✓ Dependencies tidied"

# CI/CD targets
.PHONY: ci cd

# Run CI checks
ci: deps fmt lint build test
	@echo "✓ CI checks passed"

# Run CD deployment (placeholder)
cd:
	@echo "CD deployment not configured"
	@echo "Configure your deployment pipeline here"

# Version management
.PHONY: version

# Show version information
version:
	@echo "eigenruntime-go"
	@echo "Version: $$(git describe --tags --always --dirty 2>/dev/null || echo 'unknown')"
	@echo "Commit:  $$(git rev-parse HEAD 2>/dev/null || echo 'unknown')"
	@echo "Go:      $$(go version)"

# Generate code (if needed)
.PHONY: generate

generate:
	@echo "Running go generate..."
	@go generate ./...
	@echo "✓ Code generation complete"