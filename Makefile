.PHONY: build run clean test deps

# Build the application
build:
	go build -o llm-proxy .

# Run the application
run:
	go run .

# Clean build artifacts
clean:
	rm -f llm-proxy

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests
test:
	go test ./...

# Build for different platforms
build-linux:
	GOOS=linux GOARCH=amd64 go build -o llm-proxy-linux .

build-windows:
	GOOS=windows GOARCH=amd64 go build -o llm-proxy-windows.exe .

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o llm-proxy-macos .

# Build all platforms
build-all: build-linux build-windows build-macos

# Development mode (with hot reload)
dev:
	air

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  test         - Run tests"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-windows - Build for Windows"
	@echo "  build-macos  - Build for macOS"
	@echo "  build-all    - Build for all platforms"
	@echo "  dev          - Run in development mode (requires air)"
	@echo "  help         - Show this help message"
