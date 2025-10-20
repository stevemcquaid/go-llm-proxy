package backend

import (
	"context"
	"fmt"
	"go-llm-proxy/internal/types"
	"go-llm-proxy/pkg/anthropic"
	"go-llm-proxy/pkg/openai"
)

// BackendManager manages all available backends
type BackendManager struct {
	backends map[types.BackendType]types.BackendHandler
}

// NewBackendManager creates a new backend manager
func NewBackendManager() *BackendManager {
	return &BackendManager{
		backends: make(map[types.BackendType]types.BackendHandler),
	}
}

// RegisterBackend registers a new backend
func (bm *BackendManager) RegisterBackend(backendType types.BackendType, handler types.BackendHandler) {
	bm.backends[backendType] = handler
}

// GetBackend returns a backend handler by type
func (bm *BackendManager) GetBackend(backendType types.BackendType) (types.BackendHandler, bool) {
	handler, exists := bm.backends[backendType]
	return handler, exists
}

// GetAvailableBackends returns all available backends
func (bm *BackendManager) GetAvailableBackends() []types.BackendType {
	var available []types.BackendType
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
		anthropicBackend := anthropic.NewAnthropicBackend(bf.anthropicAPIKey)
		manager.RegisterBackend(types.BackendAnthropic, anthropicBackend)
	}

	// Create OpenAI backend if API key is available
	if bf.openaiAPIKey != "" {
		openaiBackend := openai.NewOpenAIBackend(bf.openaiAPIKey)
		manager.RegisterBackend(types.BackendOpenAI, openaiBackend)
	}

	return manager
}

// ProcessRequest processes a request using the appropriate backend
func (bm *BackendManager) ProcessRequest(ctx context.Context, modelConfig types.ModelConfig, req interface{}) (interface{}, error) {
	backend, exists := bm.GetBackend(modelConfig.Backend)
	if !exists {
		return nil, fmt.Errorf("backend %s not available", modelConfig.Backend)
	}

	if !backend.IsAvailable() {
		return nil, fmt.Errorf("backend %s is not available", modelConfig.Backend)
	}

	// Route request based on type
	switch r := req.(type) {
	case types.GenerateRequest:
		return backend.Generate(ctx, r)
	case types.ChatRequest:
		return backend.Chat(ctx, r)
	default:
		return nil, fmt.Errorf("unsupported request type")
	}
}
