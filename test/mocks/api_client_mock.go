package mocks

import (
	"context"
	"go-llm-proxy/internal/fetcher"
)

// MockAPIClient is a mock implementation of the API client for testing
type MockAPIClient struct {
	AnthropicModels []fetcher.AnthropicModel
	OpenAIModels    []fetcher.OpenAIModel
	AnthropicError  error
	OpenAIError     error
}

// NewMockAPIClient creates a new mock API client
func NewMockAPIClient() *MockAPIClient {
	return &MockAPIClient{
		AnthropicModels: []fetcher.AnthropicModel{},
		OpenAIModels:    []fetcher.OpenAIModel{},
	}
}

// FetchAnthropicModels returns the mock Anthropic models
func (m *MockAPIClient) FetchAnthropicModels(ctx context.Context, apiKey string) ([]fetcher.AnthropicModel, error) {
	if m.AnthropicError != nil {
		return nil, m.AnthropicError
	}
	return m.AnthropicModels, nil
}

// FetchOpenAIModels returns the mock OpenAI models
func (m *MockAPIClient) FetchOpenAIModels(ctx context.Context, apiKey string) ([]fetcher.OpenAIModel, error) {
	if m.OpenAIError != nil {
		return nil, m.OpenAIError
	}
	return m.OpenAIModels, nil
}

// SetAnthropicModels sets the mock Anthropic models
func (m *MockAPIClient) SetAnthropicModels(models []fetcher.AnthropicModel) {
	m.AnthropicModels = models
}

// SetOpenAIModels sets the mock OpenAI models
func (m *MockAPIClient) SetOpenAIModels(models []fetcher.OpenAIModel) {
	m.OpenAIModels = models
}

// SetAnthropicError sets an error to return for Anthropic API calls
func (m *MockAPIClient) SetAnthropicError(err error) {
	m.AnthropicError = err
}

// SetOpenAIError sets an error to return for OpenAI API calls
func (m *MockAPIClient) SetOpenAIError(err error) {
	m.OpenAIError = err
}
