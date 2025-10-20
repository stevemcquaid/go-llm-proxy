package helpers_test

import (
	"context"

	"go-llm-proxy/internal/types"
)

// MockBackend is a mock implementation of BackendHandler for testing
type MockBackend struct {
	name      string
	available bool
}

func (m *MockBackend) Generate(_ context.Context, req types.GenerateRequest) (*types.GenerateResponse, error) {
	return &types.GenerateResponse{
		Model:     req.Model,
		Content:   "Mock response",
		CreatedAt: "2025-10-20T17:00:00Z",
	}, nil
}

func (m *MockBackend) Chat(_ context.Context, req types.ChatRequest) (*types.ChatResponse, error) {
	return &types.ChatResponse{
		Model: req.Model,
		Message: types.ChatMessage{
			Role:    "assistant",
			Content: "Mock response",
		},
		CreatedAt: "2025-10-20T17:00:00Z",
	}, nil
}

func (m *MockBackend) IsAvailable() bool {
	return m.available
}

func (m *MockBackend) GetName() string {
	return m.name
}

// isNewerModel checks if the model is a newer model that doesn't support MaxTokens
// This is a copy of the function from the openai package for testing purposes
func IsNewerModel(model string) bool {
	// Models that require MaxCompletionTokens instead of MaxTokens
	newerModels := []string{
		"gpt-4o",
		"gpt-4o-mini",
		"gpt-5",
		"gpt-4.1",
		"gpt-4.5",
	}

	for _, newerModel := range newerModels {
		if model == newerModel {
			return true
		}
	}
	return false
}
