# Go LLM Proxy v2 - Refactored

A refactored, maintainable Go-based proxy for LLMs that allows JetBrains IDEs to communicate with various AI providers using the Ollama API format.

## 🏗️ Architecture

The refactored version is organized into clear, maintainable components:

### Core Components

- **`config.go`** - Configuration management
- **`models.go`** - Model registry and configuration
- **`backend.go`** - Backend abstraction layer
- **`streaming.go`** - Streaming response handling
- **`proxy_v2.go`** - Main proxy server logic

### Backend Implementations

- **`backends/anthropic.go`** - Anthropic Claude integration
- **`backends/openai.go`** - OpenAI GPT integration

### Examples

- **`examples/add_model.go`** - How to add new models
- **`main_v2.go`** - Refactored main application

## 🚀 Key Improvements

### 1. **Easy Model Management**
```go
// Add a new model
newModel := ModelConfig{
    Name:        "gpt-4-turbo",
    DisplayName: "GPT-4 Turbo",
    Backend:     BackendOpenAI,
    BackendModel: "gpt-4-turbo-preview",
    Family:      "gpt",
    Description: "Latest GPT-4 Turbo model",
    MaxTokens:   128000,
    Enabled:     true,
}
registry.AddModel(newModel)
```

### 2. **Easy Backend Addition**
To add a new backend (e.g., Cohere):

1. **Create backend type**:
```go
const BackendCohere BackendType = "cohere"
```

2. **Implement BackendHandler interface**:
```go
type CohereBackend struct {
    apiKey string
    client *http.Client
}

func (cb *CohereBackend) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
    // Implementation
}

func (cb *CohereBackend) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
    // Implementation
}

func (cb *CohereBackend) IsAvailable() bool {
    return cb.apiKey != ""
}

func (cb *CohereBackend) GetName() string {
    return "cohere"
}
```

3. **Register in factory**:
```go
func (bf *BackendFactory) CreateBackends() *BackendManager {
    manager := NewBackendManager()
    
    if bf.cohereAPIKey != "" {
        cohereBackend := NewCohereBackend(bf.cohereAPIKey)
        manager.RegisterBackend(BackendCohere, cohereBackend)
    }
    
    return manager
}
```

### 3. **Configuration Management**
```go
// Environment variables with defaults
PORT=11434
GIN_MODE=release
ANTHROPIC_API_KEY=your_key_here
OPENAI_API_KEY=your_key_here
DEFAULT_MAX_TOKENS=4096
STREAMING_CHUNK_SIZE=3
STREAMING_DELAY_MS=50
```

### 4. **Clean Separation of Concerns**
- **Models**: Centralized model configuration
- **Backends**: Pluggable backend implementations
- **Streaming**: Dedicated streaming handler
- **Config**: Environment-based configuration

## 📁 File Structure

```
go-llm-proxy/
├── main_v2.go              # Refactored main application
├── proxy_v2.go             # Main proxy server logic
├── config.go               # Configuration management
├── models.go               # Model registry
├── backend.go              # Backend abstraction
├── streaming.go            # Streaming handler
├── backends/
│   ├── anthropic.go        # Anthropic backend
│   └── openai.go           # OpenAI backend
├── examples/
│   └── add_model.go        # Usage examples
└── README_v2.md            # This file
```

## 🔧 Usage

### Running the Refactored Version

```bash
# Build and run the refactored version
go build -o llm-proxy-v2 main_v2.go proxy_v2.go config.go models.go backend.go streaming.go backends/anthropic.go backends/openai.go
./llm-proxy-v2
```

### Adding New Models

```go
// In your application code
registry := NewModelRegistry()

// Add a new model
newModel := ModelConfig{
    Name:        "claude-3-opus",
    DisplayName: "Claude 3 Opus",
    Backend:     BackendAnthropic,
    BackendModel: "claude-3-opus-20240229",
    Family:      "claude",
    Description: "Most powerful Claude model",
    MaxTokens:   200000,
    Enabled:     true,
}
registry.AddModel(newModel)
```

### Disabling Models

```go
// Disable a model
registry.DisableModel("gpt-3.5-turbo")

// Re-enable a model
registry.EnableModel("gpt-3.5-turbo")
```

## 🎯 Benefits of Refactoring

### 1. **Maintainability**
- Clear separation of concerns
- Easy to understand and modify
- Consistent patterns across components

### 2. **Extensibility**
- Easy to add new backends
- Easy to add new models
- Pluggable architecture

### 3. **Testability**
- Each component can be tested independently
- Clear interfaces for mocking
- Isolated functionality

### 4. **Configuration**
- Environment-based configuration
- Easy to customize behavior
- Clear defaults

### 5. **Readability**
- Self-documenting code
- Clear naming conventions
- Logical file organization

## 🔄 Migration from v1

The refactored version maintains the same API endpoints and behavior as v1, so it's a drop-in replacement. The main differences are:

1. **Better organization** - Code is split into logical files
2. **Easier maintenance** - Clear separation of concerns
3. **Better extensibility** - Easy to add new backends/models
4. **Better configuration** - Environment-based config management

## 🚀 Future Enhancements

With the refactored architecture, it's now easy to add:

- **New backends** (Cohere, Google, etc.)
- **New models** (just add to registry)
- **New features** (caching, rate limiting, etc.)
- **New endpoints** (embeddings, image generation, etc.)
- **Better monitoring** (metrics, logging, etc.)

The refactored codebase provides a solid foundation for future development while maintaining the simplicity and reliability of the original proxy.
