# Go LLM Proxy

A Go-based proxy server that translates Ollama API requests to multiple backend APIs (Anthropic and OpenAI), allowing JetBrains IDEs to use various LLM models through the familiar Ollama interface.

## Features

- **Ollama API Compatibility**: Implements the standard Ollama API endpoints
- **Multi-Backend Support**: Supports both Anthropic Claude and OpenAI GPT models
- **Smart Model Routing**: Automatically routes requests to the appropriate backend based on model name
- **Model Mapping**: Maps Ollama model names to appropriate backend models
- **JetBrains IDE Support**: Works seamlessly with JetBrains IDEs that support Ollama

## Supported Models

The proxy supports both Anthropic Claude and OpenAI GPT models. Model names are automatically routed to the appropriate backend:

### Anthropic Claude Models
- `claude` → `claude-3-5-sonnet-20241022`
- `claude-3` → `claude-3-5-sonnet-20241022`
- `claude-3-sonnet` → `claude-3-5-sonnet-20241022`
- `claude-3-haiku` → `claude-3-5-haiku-20241022`
- `claude-3-opus` → `claude-3-5-opus-20241022`
- `claude-3.5-sonnet` → `claude-3-5-sonnet-20241022`
- `claude-3.5-haiku` → `claude-3-5-haiku-20241022`
- `claude-3.5-opus` → `claude-3-5-opus-20241022`

### OpenAI GPT Models
- `gpt-4` → `gpt-4`
- `gpt-4-turbo` → `gpt-4-turbo-preview`
- `gpt-4o` → `gpt-4o`
- `gpt-4o-mini` → `gpt-4o-mini`
- `gpt-3.5-turbo` → `gpt-3.5-turbo`
- `gpt-3.5-turbo-16k` → `gpt-3.5-turbo-16k`
- `o1-preview` → `o1-preview`
- `o1-mini` → `o1-mini`

### Smart Routing
- Models starting with `claude` → Anthropic
- Models starting with `gpt` or `o1` → OpenAI
- Unknown models → Anthropic (fallback)

### How Model Routing Works

The proxy automatically determines which backend to use based on the model name in your request:

1. **Exact matches**: If the model name exactly matches a predefined mapping, it uses that backend
2. **Prefix matching**: If the model starts with known prefixes (`claude`, `gpt`, `o1`), it routes accordingly
3. **Fallback**: Unknown models default to Anthropic

This means you can use natural model names like `claude-3.5-sonnet` or `gpt-4` and the proxy will automatically route them to the correct API.

## Project Structure

```
go-llm-proxy/
├── main.go              # Main server entry point
├── proxy.go             # Core proxy logic and API handlers
├── proxy_test.go        # Unit tests
├── go.mod               # Go module dependencies
├── Makefile             # Build and run commands
├── setup.sh             # Setup script
├── test_proxy.sh        # Test script
├── env.example          # Environment configuration template
└── README.md            # This file
```

## Quick Start

### Option 1: Automated Setup (Recommended)

1. **Run the setup script**:
   ```bash
   ./setup.sh
   ```

2. **Add your API keys**:
   ```bash
   # Edit .env file and add your API keys
   ANTHROPIC_API_KEY=your_anthropic_api_key_here
   OPENAI_API_KEY=your_openai_api_key_here
   ```
   
   **Note**: You can configure one or both API keys. The proxy will only show models for which you have API keys configured.

3. **Run the proxy**:
   ```bash
   ./llm-proxy
   # or
   go run .
   ```

4. **Test it works**:
   ```bash
   ./test_proxy.sh
   ```

### Option 2: Manual Setup

1. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

2. **Configure Environment**:
   ```bash
   cp env.example .env
   ```
   
   Edit `.env` and add your API keys:
   ```
   ANTHROPIC_API_KEY=your_anthropic_api_key_here
   OPENAI_API_KEY=your_openai_api_key_here
   PORT=11434
   ```

3. **Run the Server**:
   ```bash
   go run .
   ```

   The server will start on port 11434 (default Ollama port).

## Usage

### With JetBrains IDEs

1. Start the proxy server
2. In your JetBrains IDE, configure the Ollama integration to point to `localhost:11434`
3. Use any of the supported model names (e.g., `claude-3.5-sonnet`, `gpt-4`, `gpt-4o`)

### Direct API Usage

The proxy implements the standard Ollama API endpoints:

- `POST /api/generate` - Generate text completions
- `POST /api/chat` - Chat completions
- `GET /api/tags` - List available models
- `GET /api/version` - Get server version
- `POST /api/show` - Show model information

### Example Usage

```bash
# List available models
curl http://localhost:11434/api/tags

# Generate text with Claude
curl -X POST http://localhost:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3.5-sonnet",
    "prompt": "Hello, how are you?",
    "stream": false
  }'

# Generate text with GPT-4
curl -X POST http://localhost:11434/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4",
    "prompt": "Hello, how are you?",
    "stream": false
  }'

# Chat completion with Claude
curl -X POST http://localhost:11434/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "claude-3.5-sonnet",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ],
    "stream": false
  }'

# Chat completion with GPT-4
curl -X POST http://localhost:11434/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4o",
    "messages": [
      {"role": "user", "content": "Hello!"}
    ],
    "stream": false
  }'
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/generate` | POST | Generate text completions |
| `/api/chat` | POST | Chat completions |
| `/api/tags` | GET | List available models |
| `/api/version` | GET | Get server version |
| `/api/show` | POST | Show model information |
| `/api/pull` | POST | Pull model (no-op for Anthropic) |
| `/api/push` | POST | Push model (no-op for Anthropic) |
| `/api/delete` | DELETE | Delete model (no-op for Anthropic) |
| `/api/create` | POST | Create model (no-op for Anthropic) |
| `/api/copy` | POST | Copy model (no-op for Anthropic) |
| `/api/embeddings` | POST | Get embeddings (not supported) |
| `/api/ps` | POST | List processes (no-op for backends) |
| `/api/stop` | POST | Stop generation (no-op for backends) |
| `/` | GET | Root endpoint (JetBrains compatibility) |
| `/api` | GET | API info endpoint |
| `/health` | GET | Health check |

## Configuration

Environment variables:

- `ANTHROPIC_API_KEY` (optional): Your Anthropic API key for Claude models
- `OPENAI_API_KEY` (optional): Your OpenAI API key for GPT models
- `PORT` (optional): Server port (default: 11434)
- `GIN_MODE` (optional): Gin mode (release/debug, default: release)

**Note**: At least one API key must be configured. The proxy will only show models for which you have the corresponding API key.

## Available Commands

### Make Commands

```bash
make build        # Build the application
make run          # Run the application
make clean        # Clean build artifacts
make deps         # Install dependencies
make test         # Run tests
make build-linux  # Build for Linux
make build-windows # Build for Windows
make build-macos  # Build for macOS
make build-all    # Build for all platforms
make help         # Show all available commands
```

### Direct Go Commands

```bash
go run .          # Run the server
go build -o llm-proxy .  # Build binary
go test ./...     # Run tests
go mod tidy       # Install dependencies
```

### Scripts

```bash
./setup.sh        # Initial setup (creates .env, installs deps, builds)
./test_proxy.sh   # Test the proxy endpoints
```

## Building

### Build Binary

```bash
# Build for current platform
go build -o llm-proxy .

# Or use make
make build
```

### Cross-Platform Builds

```bash
# Build for different platforms
make build-linux    # Linux binary
make build-windows  # Windows binary
make build-macos    # macOS binary
make build-all      # All platforms
```

## Troubleshooting

### Common Issues

1. **"Proxy server is not running" error**:
   - Make sure the proxy is running: `./llm-proxy` or `go run .`
   - Check if port 11434 is available: `lsof -i :11434`

2. **"API error"**:
   - Verify your API keys are correct in `.env`
   - Check your API keys have sufficient credits
   - Ensure you have access to the requested models
   - For Anthropic: Check your Claude model access
   - For OpenAI: Check your GPT model access

3. **JetBrains IDE not connecting**:
   - Verify the proxy is running on `localhost:11434`
   - Check your IDE's Ollama configuration points to the correct URL
   - Try restarting your IDE after starting the proxy

4. **Build errors**:
   - Run `go mod tidy` to ensure dependencies are installed
   - Check Go version (requires Go 1.21+)

### Testing the Proxy

```bash
# Test if server is running
curl http://localhost:11434/

# Test API endpoint
curl http://localhost:11434/api

# Test version endpoint
curl http://localhost:11434/api/version

# Test available models
curl http://localhost:11434/api/tags

# Run comprehensive tests
./test_proxy.sh
```

## Limitations

- **Streaming responses** are not yet implemented (planned for future release)
- **Model management** features (pull, push, delete) are no-ops since backends manage models
- **Embeddings** are not supported (neither Anthropic nor OpenAI provide embeddings through chat APIs)
- **Model size information** is not available from backend APIs
- **Custom model parameters** may not be fully supported (uses backend defaults)
- **Mixed backends**: Each request uses only one backend; you cannot mix Anthropic and OpenAI in a single conversation

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestModelMapping
```

### Adding New Model Mappings

Edit the `getModelInfo` function in `proxy.go`:

```go
// Add to the appropriate model map
anthropicModels := map[string]string{
    "your-ollama-name": "anthropic-model-name",
    // ... existing mappings
}

openaiModels := map[string]string{
    "your-ollama-name": "openai-model-name",
    // ... existing mappings
}
```

### Adding New Backends

To add support for a new backend (e.g., Google Gemini):

1. Add a new `BackendType` constant
2. Add the backend to the `ProxyServer` struct
3. Implement the backend-specific handler functions
4. Update the `getModelInfo` function to include the new backend
5. Add the new backend to the switch statements in handlers

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run `go test ./...` to ensure tests pass
6. Submit a pull request

## License

MIT License - see LICENSE file for details.
