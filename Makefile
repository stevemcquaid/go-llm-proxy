.PHONY: build run clean test deps

# Build the application
build:
	go build -o bin/llm-proxy ./cmd/llm-proxy

# Run the application
run:
	go run ./cmd/llm-proxy

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests
test:
	go test ./test/...

# Run tests with verbose output
test-verbose:
	go test -v ./test/...

# Run tests with coverage
test-coverage:
	go test -v -cover ./test/...

# Run specific test suite
test-api:
	go test -v -run TestOllamaAPISpec ./test/...

test-integration:
	go test -v -run TestProxyIntegration ./test/...

test-models:
	go test -v -run TestModelRegistry ./test/...

test-config:
	go test -v -run TestConfig ./test/...

# Build for different platforms
build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/llm-proxy-linux ./cmd/llm-proxy

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/llm-proxy-windows.exe ./cmd/llm-proxy

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/llm-proxy-macos ./cmd/llm-proxy

# Build all platforms
build-all: build-linux build-windows build-macos

# Development mode (with hot reload)
dev:
	air

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  test          - Run tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  test-api      - Run API compatibility tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-models   - Run model management tests"
	@echo "  test-config   - Run configuration tests"
	@echo "  build-linux   - Build for Linux"
	@echo "  build-windows - Build for Windows"
	@echo "  build-macos   - Build for macOS"
	@echo "  build-all     - Build for all platforms"
	@echo "  dev           - Run in development mode (requires air)"
	@echo "  help          - Show this help message"
