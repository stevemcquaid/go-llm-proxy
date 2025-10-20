# Go LLM Proxy

A maintainable Go-based proxy for LLMs that allows JetBrains IDEs to communicate with various AI providers using the Ollama API format.

## 🏗️ Project Structure

```
go-llm-proxy/
├── bin/                         # Compiled binaries
│   ├── llm-proxy
│   ├── llm-proxy-linux
│   ├── llm-proxy-macos
│   └── llm-proxy-windows.exe
├── cmd/llm-proxy/               # Main application entry point
│   └── main.go
├── internal/                    # Internal packages
│   ├── config/                  # Configuration management
│   │   └── config.go
│   ├── models/                  # Model registry and management
│   │   └── models.go
│   ├── backend/                 # Backend abstraction layer
│   │   └── backend.go
│   ├── proxy/                   # Main proxy server logic
│   │   └── proxy.go
│   ├── streaming/               # Streaming response handling
│   │   └── streaming.go
│   └── types/                   # Shared type definitions
│       └── types.go
├── pkg/                         # Public packages
│   ├── anthropic/               # Anthropic Claude integration
│   │   └── anthropic_backend.go
│   └── openai/                  # OpenAI GPT integration
│       └── openai_backend.go
├── test/                        # Test files
│   ├── unit/                    # Unit tests
│   │   ├── config_test.go
│   │   ├── model_management_test.go
│   │   ├── ollama_api_test.go
│   │   └── proxy_test.go
│   └── integration/             # Integration tests
│       └── integration_test.go
├── scripts/                     # Build and utility scripts
│   ├── setup.sh
│   └── test_proxy.sh
├── examples/                    # Usage examples
│   └── main.go
├── docs/                        # Documentation
│   ├── README.md
│   └── TEST_SUMMARY.md
├── .env.example                 # Environment variables template
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
└── Makefile                     # Build automation
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or later
- API keys for Anthropic and/or OpenAI

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd go-llm-proxy
   ```

2. **Install dependencies:**
   ```bash
   make deps
   ```

3. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your API keys
   ```

4. **Build the proxy:**
   ```bash
   make build
   ```

5. **Run the proxy:**
   ```bash
   ./bin/llm-proxy
   ```

## 🔧 Usage

### Building the Proxy

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Clean build artifacts
make clean
```

### Running the Proxy

```bash
# Run the built binary
./bin/llm-proxy

# Or run directly without building
make run
```

### Testing

```bash
# Run all tests
make test

# Run with verbose output
make test-verbose

# Run with coverage
make test-coverage

# Run specific test suites
make test-api          # API compatibility tests
make test-integration  # Integration tests
make test-models       # Model management tests
make test-config       # Configuration tests
```

## 📁 Source Files

The main source files are organized as follows:

### Core Application
- **`cmd/llm-proxy/main.go`** - Main application entry point

### Internal Packages
- **`internal/types/types.go`** - Shared type definitions and interfaces
- **`internal/config/config.go`** - Configuration management
- **`internal/models/models.go`** - Model registry and management
- **`internal/backend/backend.go`** - Backend abstraction layer
- **`internal/proxy/proxy.go`** - Main proxy server logic
- **`internal/streaming/streaming.go`** - Streaming response handling

### Public Packages
- **`pkg/anthropic/anthropic_backend.go`** - Anthropic Claude integration
- **`pkg/openai/openai_backend.go`** - OpenAI GPT integration

### Test Files

All test files are organized by type:

#### Unit Tests (`test/unit/`)
- **`config_test.go`** - Configuration tests
- **`model_management_test.go`** - Model management tests
- **`ollama_api_test.go`** - Ollama API compatibility tests
- **`proxy_test.go`** - Basic proxy tests

#### Integration Tests (`test/integration/`)
- **`integration_test.go`** - End-to-end integration tests

## 🧪 Testing

The project includes comprehensive tests:

- **Unit Tests** - Individual component testing
- **Integration Tests** - End-to-end workflow testing
- **API Compatibility Tests** - Ollama API compliance verification
- **Model Management Tests** - Model registry functionality
- **Configuration Tests** - Environment variable handling

See `TEST_SUMMARY.md` for detailed test documentation.

## 🔧 Configuration

The proxy can be configured via environment variables:

```bash
# Server configuration
PORT=11434
GIN_MODE=release

# API Keys
ANTHROPIC_API_KEY=your_anthropic_key_here
OPENAI_API_KEY=your_openai_key_here

# Model configuration
DEFAULT_MAX_TOKENS=4096
STREAMING_CHUNK_SIZE=3
STREAMING_DELAY_MS=50
```

## 🎯 Features

- **Ollama API Compatibility** - Full compatibility with Ollama API format
- **JetBrains IDE Support** - Works seamlessly with GoLand AI Assistant
- **Multi-Backend Support** - Anthropic and OpenAI backends
- **Streaming Support** - Both streaming and non-streaming responses
- **Model Management** - Dynamic model addition/removal
- **CORS Support** - Cross-origin request handling
- **Comprehensive Testing** - Full test coverage

## 📚 Documentation

- **`README.md`** - This file (main documentation)
- **`TEST_SUMMARY.md`** - Comprehensive test documentation
- **`examples/main.go`** - Usage examples
- **`scripts/`** - Build and utility scripts

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite: `make test`
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

For issues and questions:
1. Check the documentation
2. Run the test suite to verify setup
3. Check the logs for error messages
4. Open an issue on GitHub