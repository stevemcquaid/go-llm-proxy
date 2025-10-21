package fixtures

import (
	"go-llm-proxy/internal/fetcher"
	"go-llm-proxy/internal/types"
)

// GetTestAnthropicModels returns test Anthropic models
func GetTestAnthropicModels() []fetcher.AnthropicModel {
	return []fetcher.AnthropicModel{
		{
			ID:          "claude-3-5-sonnet-20241022",
			Name:        "claude-3-5-sonnet-20241022",
			Description: "Enhanced coding and reasoning capabilities with improved performance",
			ContextSize: 200000,
		},
		{
			ID:          "claude-3-5-haiku-20241022",
			Name:        "claude-3-5-haiku-20241022",
			Description: "Lightweight high-speed model optimized for real-time tasks and coding",
			ContextSize: 200000,
		},
		{
			ID:          "claude-3-5-opus-20241022",
			Name:        "claude-3-5-opus-20241022",
			Description: "Most powerful model for complex reasoning and advanced AI solutions",
			ContextSize: 200000,
		},
		{
			ID:          "claude-3-7-sonnet-20250219",
			Name:        "claude-3-7-sonnet-20250219",
			Description: "Hybrid reasoning model with rapid and detailed reasoning options",
			ContextSize: 8192,
		},
	}
}

// GetTestOpenAIModels returns test OpenAI models
func GetTestOpenAIModels() []fetcher.OpenAIModel {
	return []fetcher.OpenAIModel{
		{
			ID:      "gpt-4o",
			Object:  "model",
			Created: 1700000000,
			OwnedBy: "openai",
		},
		{
			ID:      "gpt-4o-mini",
			Object:  "model",
			Created: 1700000000,
			OwnedBy: "openai",
		},
		{
			ID:      "gpt-4",
			Object:  "model",
			Created: 1700000000,
			OwnedBy: "openai",
		},
		{
			ID:      "gpt-3.5-turbo",
			Object:  "model",
			Created: 1700000000,
			OwnedBy: "openai",
		},
	}
}

// GetExpectedModelConfigs returns the expected ModelConfig objects for testing
func GetExpectedModelConfigs() []types.ModelConfig {
	return []types.ModelConfig{
		// Anthropic models
		{
			Name:         "claude-3.5-sonnet",
			DisplayName:  "Claude 3.5 Sonnet",
			Backend:      types.BackendAnthropic,
			BackendModel: "claude-3-5-sonnet-20241022",
			Family:       "claude",
			Description:  "Anthropic Claude 3.5 Sonnet model",
			MaxTokens:    200000,
			Enabled:      true,
		},
		{
			Name:         "claude-3.5-haiku",
			DisplayName:  "Claude 3.5 Haiku",
			Backend:      types.BackendAnthropic,
			BackendModel: "claude-3-5-haiku-20241022",
			Family:       "claude",
			Description:  "Anthropic Claude 3.5 Haiku model",
			MaxTokens:    200000,
			Enabled:      true,
		},
		{
			Name:         "claude-3.5-opus",
			DisplayName:  "Claude 3.5 Opus",
			Backend:      types.BackendAnthropic,
			BackendModel: "claude-3-5-opus-20241022",
			Family:       "claude",
			Description:  "Anthropic Claude 3.5 Opus model",
			MaxTokens:    200000,
			Enabled:      true,
		},
		{
			Name:         "claude-3.7-sonnet",
			DisplayName:  "Claude 3.7 Sonnet",
			Backend:      types.BackendAnthropic,
			BackendModel: "claude-3-7-sonnet-20250219",
			Family:       "claude",
			Description:  "Anthropic Claude 3.7 Sonnet model",
			MaxTokens:    8192,
			Enabled:      true,
		},
		{
			Name:         "claude-4.5-sonnet",
			DisplayName:  "Claude 4.5 Sonnet",
			Backend:      types.BackendAnthropic,
			BackendModel: "claude-3-5-sonnet-20241022",
			Family:       "claude",
			Description:  "Anthropic Claude 4.5 Sonnet model",
			MaxTokens:    200000,
			Enabled:      true,
		},
		// OpenAI models
		{
			Name:         "gpt-4o",
			DisplayName:  "GPT-4o",
			Backend:      types.BackendOpenAI,
			BackendModel: "gpt-4o",
			Family:       "gpt",
			Description:  "OpenAI GPT-4o model",
			MaxTokens:    16384,
			Enabled:      true,
		},
		{
			Name:         "gpt-4o-mini",
			DisplayName:  "GPT-4o Mini",
			Backend:      types.BackendOpenAI,
			BackendModel: "gpt-4o-mini",
			Family:       "gpt",
			Description:  "OpenAI GPT-4o Mini model",
			MaxTokens:    16384,
			Enabled:      true,
		},
		{
			Name:         "gpt-4",
			DisplayName:  "GPT-4",
			Backend:      types.BackendOpenAI,
			BackendModel: "gpt-4",
			Family:       "gpt",
			Description:  "OpenAI GPT-4 model",
			MaxTokens:    8192,
			Enabled:      true,
		},
		{
			Name:         "gpt-3.5-turbo",
			DisplayName:  "GPT-3.5 Turbo",
			Backend:      types.BackendOpenAI,
			BackendModel: "gpt-3.5-turbo",
			Family:       "gpt",
			Description:  "OpenAI GPT-3.5 Turbo model",
			MaxTokens:    4096,
			Enabled:      true,
		},
	}
}
