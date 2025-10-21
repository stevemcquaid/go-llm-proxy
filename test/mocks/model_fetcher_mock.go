package mocks

import (
	"context"
	"fmt"
	"go-llm-proxy/internal/config"
	"go-llm-proxy/internal/types"
)

// MockModelFetcher is a mock implementation of the model fetcher for testing
type MockModelFetcher struct {
	apiClient *MockAPIClient
	config    *config.Config
}

// NewMockModelFetcher creates a new mock model fetcher
func NewMockModelFetcher(cfg *config.Config) *MockModelFetcher {
	return &MockModelFetcher{
		apiClient: NewMockAPIClient(),
		config:    cfg,
	}
}

// LoadConfigFromFile is a no-op for testing
func (m *MockModelFetcher) LoadConfigFromFile(configPath string) error {
	return nil
}

// FetchAllModels fetches models using the mock API client
func (m *MockModelFetcher) FetchAllModels(ctx context.Context) ([]types.ModelConfig, error) {
	var allModels []types.ModelConfig

	// Fetch Anthropic models if enabled
	if m.config.ModelFilters.Anthropic.Enabled && m.config.AnthropicAPIKey != "" {
		anthropicModels, err := m.fetchAnthropicModels(ctx)
		if err != nil {
			return nil, err
		}
		allModels = append(allModels, anthropicModels...)
	}

	// Fetch OpenAI models if enabled
	if m.config.ModelFilters.OpenAI.Enabled && m.config.OpenAIAPIKey != "" {
		openaiModels, err := m.fetchOpenAIModels(ctx)
		if err != nil {
			return nil, err
		}
		allModels = append(allModels, openaiModels...)
	}

	if len(allModels) == 0 {
		return nil, fmt.Errorf("no models could be fetched from any backend")
	}

	return allModels, nil
}

// fetchAnthropicModels fetches and filters Anthropic models using mock
func (m *MockModelFetcher) fetchAnthropicModels(ctx context.Context) ([]types.ModelConfig, error) {
	apiModels, err := m.apiClient.FetchAnthropicModels(ctx, m.config.AnthropicAPIKey)
	if err != nil {
		return nil, err
	}

	var models []types.ModelConfig
	for _, apiModel := range apiModels {
		// Apply filters (simplified for testing)
		if !m.matchesFilters(apiModel.ID, m.config.ModelFilters.Anthropic) {
			continue
		}

		// Convert to our ModelConfig format
		model := types.ModelConfig{
			Name:         m.generateModelName(apiModel.ID, types.BackendAnthropic),
			DisplayName:  m.generateDisplayName(apiModel.ID, types.BackendAnthropic),
			Backend:      types.BackendAnthropic,
			BackendModel: apiModel.ID,
			Family:       m.extractFamily(apiModel.ID, types.BackendAnthropic),
			Description:  apiModel.Description,
			MaxTokens:    apiModel.ContextSize,
			Enabled:      true,
		}

		models = append(models, model)
	}

	return models, nil
}

// fetchOpenAIModels fetches and filters OpenAI models using mock
func (m *MockModelFetcher) fetchOpenAIModels(ctx context.Context) ([]types.ModelConfig, error) {
	apiModels, err := m.apiClient.FetchOpenAIModels(ctx, m.config.OpenAIAPIKey)
	if err != nil {
		return nil, err
	}

	var models []types.ModelConfig
	for _, apiModel := range apiModels {
		// Apply filters (simplified for testing)
		if !m.matchesFilters(apiModel.ID, m.config.ModelFilters.OpenAI) {
			continue
		}

		// Convert to our ModelConfig format
		model := types.ModelConfig{
			Name:         m.generateModelName(apiModel.ID, types.BackendOpenAI),
			DisplayName:  m.generateDisplayName(apiModel.ID, types.BackendOpenAI),
			Backend:      types.BackendOpenAI,
			BackendModel: apiModel.ID,
			Family:       m.extractFamily(apiModel.ID, types.BackendOpenAI),
			Description:  m.generateDescription(apiModel.ID, types.BackendOpenAI),
			MaxTokens:    m.estimateMaxTokens(apiModel.ID, types.BackendOpenAI),
			Enabled:      true,
		}

		models = append(models, model)
	}

	return models, nil
}

// Helper methods (simplified versions from the real fetcher)
func (m *MockModelFetcher) matchesFilters(modelID string, filter config.ModelFilterConfig) bool {
	// Simplified filter logic for testing
	if !filter.Enabled {
		return false
	}

	// If no include patterns, include all
	if len(filter.IncludePatterns) == 0 {
		return true
	}

	// Simple pattern matching for testing
	for _, pattern := range filter.IncludePatterns {
		if pattern == "*" || modelID == pattern {
			return true
		}
	}

	return false
}

func (m *MockModelFetcher) generateModelName(apiModelID string, backend types.BackendType) string {
	// Simplified name generation for testing
	return apiModelID
}

func (m *MockModelFetcher) generateDisplayName(apiModelID string, backend types.BackendType) string {
	// Simplified display name generation for testing
	return apiModelID
}

func (m *MockModelFetcher) extractFamily(apiModelID string, backend types.BackendType) string {
	// Simplified family extraction for testing
	if backend == types.BackendAnthropic {
		return "claude"
	}
	return "gpt"
}

func (m *MockModelFetcher) generateDescription(apiModelID string, backend types.BackendType) string {
	// Simplified description generation for testing
	return "Test model description"
}

func (m *MockModelFetcher) estimateMaxTokens(apiModelID string, backend types.BackendType) int {
	// Simplified token estimation for testing
	return 4096
}

// GetAPIClient returns the mock API client for configuration
func (m *MockModelFetcher) GetAPIClient() *MockAPIClient {
	return m.apiClient
}
