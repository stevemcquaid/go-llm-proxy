package main

import (
	"context"
	"fmt"
)

// BackendHandler defines the interface that all backends must implement
type BackendHandler interface {
	// Generate handles text generation requests
	Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)

	// Chat handles chat completion requests
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)

	// IsAvailable checks if the backend is available (has API key, etc.)
	IsAvailable() bool

	// GetName returns the backend name
	GetName() string
}

// GenerateRequest represents a text generation request
type GenerateRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens,omitempty"`
}

// GenerateResponse represents a text generation response
type GenerateResponse struct {
	Model     string `json:"model"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model     string        `json:"model"`
	Messages  []ChatMessage `json:"messages"`
	MaxTokens int           `json:"max_tokens,omitempty"`
}

// ChatMessage represents a single message in a chat
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Model     string      `json:"model"`
	Message   ChatMessage `json:"message"`
	CreatedAt string      `json:"created_at"`
}

// BackendManager manages all available backends
type BackendManager struct {
	backends map[BackendType]BackendHandler
}

// NewBackendManager creates a new backend manager
func NewBackendManager() *BackendManager {
	return &BackendManager{
		backends: make(map[BackendType]BackendHandler),
	}
}

// RegisterBackend registers a new backend
func (bm *BackendManager) RegisterBackend(backendType BackendType, handler BackendHandler) {
	bm.backends[backendType] = handler
}

// GetBackend returns a backend handler by type
func (bm *BackendManager) GetBackend(backendType BackendType) (BackendHandler, bool) {
	handler, exists := bm.backends[backendType]
	return handler, exists
}

// GetAvailableBackends returns all available backends
func (bm *BackendManager) GetAvailableBackends() []BackendType {
	var available []BackendType
	for backendType, handler := range bm.backends {
		if handler.IsAvailable() {
			available = append(available, backendType)
		}
	}
	return available
}

// BackendFactory creates backend handlers
type BackendFactory struct {
	anthropicAPIKey string
	openaiAPIKey    string
}

// NewBackendFactory creates a new backend factory
func NewBackendFactory(anthropicAPIKey, openaiAPIKey string) *BackendFactory {
	return &BackendFactory{
		anthropicAPIKey: anthropicAPIKey,
		openaiAPIKey:    openaiAPIKey,
	}
}

// CreateBackends creates all available backends
func (bf *BackendFactory) CreateBackends() *BackendManager {
	manager := NewBackendManager()

	// Create Anthropic backend if API key is available
	if bf.anthropicAPIKey != "" {
		anthropicBackend := NewAnthropicBackend(bf.anthropicAPIKey)
		manager.RegisterBackend(BackendAnthropic, anthropicBackend)
	}

	// Create OpenAI backend if API key is available
	if bf.openaiAPIKey != "" {
		openaiBackend := NewOpenAIBackend(bf.openaiAPIKey)
		manager.RegisterBackend(BackendOpenAI, openaiBackend)
	}

	return manager
}

// ConvertOllamaToGenerateRequest converts an Ollama generate request to our format
func ConvertOllamaToGenerateRequest(req OllamaGenerateRequest, maxTokens int) GenerateRequest {
	return GenerateRequest{
		Model:     req.Model,
		Prompt:    req.Prompt,
		MaxTokens: maxTokens,
	}
}

// ConvertOllamaToChatRequest converts an Ollama chat request to our format
func ConvertOllamaToChatRequest(req OllamaChatRequest, maxTokens int) ChatRequest {
	var messages []ChatMessage
	for _, msg := range req.Messages {
		messages = append(messages, ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return ChatRequest{
		Model:     req.Model,
		Messages:  messages,
		MaxTokens: maxTokens,
	}
}

// ConvertGenerateToOllamaResponse converts our generate response to Ollama format
func ConvertGenerateToOllamaResponse(resp *GenerateResponse, model string) OllamaGenerateResponse {
	return OllamaGenerateResponse{
		Model:     model,
		CreatedAt: resp.CreatedAt,
		Response:  resp.Content,
		Done:      true,
		Context:   []int{},
	}
}

// ConvertChatToOllamaResponse converts our chat response to Ollama format
func ConvertChatToOllamaResponse(resp *ChatResponse, model string) OllamaChatResponse {
	return OllamaChatResponse{
		Model:     model,
		CreatedAt: resp.CreatedAt,
		Message: OllamaMessage{
			Role:    resp.Message.Role,
			Content: resp.Message.Content,
		},
		Done:    true,
		Context: []int{},
	}
}

// ProcessRequest processes a request using the appropriate backend
func (bm *BackendManager) ProcessRequest(ctx context.Context, modelConfig ModelConfig, req interface{}) (interface{}, error) {
	backend, exists := bm.GetBackend(modelConfig.Backend)
	if !exists {
		return nil, fmt.Errorf("backend %s not available", modelConfig.Backend)
	}

	if !backend.IsAvailable() {
		return nil, fmt.Errorf("backend %s is not available", modelConfig.Backend)
	}

	// Route request based on type
	switch r := req.(type) {
	case GenerateRequest:
		return backend.Generate(ctx, r)
	case ChatRequest:
		return backend.Chat(ctx, r)
	default:
		return nil, fmt.Errorf("unsupported request type")
	}
}
