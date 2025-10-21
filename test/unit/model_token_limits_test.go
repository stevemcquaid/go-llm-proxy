package llmproxy_unit_test

import (
	"testing"

	"go-llm-proxy/internal/backend"
	"go-llm-proxy/internal/types"
	"go-llm-proxy/test/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestModelTokenLimits tests that all models have appropriate token limits
func TestModelTokenLimits(t *testing.T) {
	// Create a backend manager with mock backends
	backendManager := backend.NewBackendManager()
	mockOpenAI := &MockBackend{name: "openai", available: true}
	mockAnthropic := &MockBackend{name: "anthropic", available: true}
	backendManager.RegisterBackend(types.BackendOpenAI, mockOpenAI)
	backendManager.RegisterBackend(types.BackendAnthropic, mockAnthropic)

	// Create model registry with available backends
	registry := helpers.CreateTestModelRegistry()

	t.Run("OpenAIModelTokenLimits", func(t *testing.T) {
		openaiModels := registry.GetModelsByBackend(types.BackendOpenAI)
		assert.Greater(t, len(openaiModels), 0, "Should have OpenAI models")

		for _, model := range openaiModels {
			t.Run(model.Name, func(t *testing.T) {
				// Test that each model has a reasonable token limit
				assert.Greater(t, model.MaxTokens, 0, "Model %s should have positive token limit", model.Name)
				assert.LessOrEqual(t, model.MaxTokens, 1000000, "Model %s should have reasonable token limit", model.Name)

				// Test specific model limits
				switch model.Name {
				case "gpt-3.5-turbo":
					assert.Equal(t, 4096, model.MaxTokens, "gpt-3.5-turbo should have 4096 token limit")
				case "gpt-4":
					assert.Equal(t, 8192, model.MaxTokens, "gpt-4 should have 8192 token limit")
				case "gpt-4o":
					assert.Equal(t, 16384, model.MaxTokens, "gpt-4o should have 16384 token limit")
				case "gpt-4o-mini":
					assert.Equal(t, 16384, model.MaxTokens, "gpt-4o-mini should have 16384 token limit")
				case "gpt-5":
					assert.Equal(t, 128000, model.MaxTokens, "gpt-5 should have 128000 token limit")
				case "gpt-4.1":
					assert.Equal(t, 1000000, model.MaxTokens, "gpt-4.1 should have 1000000 token limit")
				}
			})
		}
	})

	t.Run("AnthropicModelTokenLimits", func(t *testing.T) {
		anthropicModels := registry.GetModelsByBackend(types.BackendAnthropic)
		assert.Greater(t, len(anthropicModels), 0, "Should have Anthropic models")

		for _, model := range anthropicModels {
			t.Run(model.Name, func(t *testing.T) {
				// Test that each model has a reasonable token limit
				assert.Greater(t, model.MaxTokens, 0, "Model %s should have positive token limit", model.Name)
				assert.LessOrEqual(t, model.MaxTokens, 200000, "Model %s should have reasonable token limit", model.Name)

				// Test specific model limits
				switch model.Name {
				case "claude-3.7-sonnet":
					assert.Equal(t, 8192, model.MaxTokens, "claude-3.7-sonnet should have 8192 token limit")
				case "claude-4.5-sonnet", "claude-4.5-haiku", "claude-4.1-opus":
					assert.Equal(t, 200000, model.MaxTokens, "Anthropic model %s should have 200000 token limit", model.Name)
				}
			})
		}
	})
}

// TestModelRegistryWithBackends tests that only available backends get models
func TestModelRegistryWithBackends(t *testing.T) {
	t.Run("OnlyOpenAIBackendAvailable", func(t *testing.T) {
		// Create model registry with test data
		registry := helpers.CreateTestModelRegistry()

		// Should have both OpenAI and Anthropic models (test data includes both)
		openaiModels := registry.GetModelsByBackend(types.BackendOpenAI)
		anthropicModels := registry.GetModelsByBackend(types.BackendAnthropic)

		assert.Greater(t, len(openaiModels), 0, "Should have OpenAI models")
		assert.Greater(t, len(anthropicModels), 0, "Should have Anthropic models")
	})

	t.Run("OnlyAnthropicBackendAvailable", func(t *testing.T) {
		// Create model registry with test data
		registry := helpers.CreateTestModelRegistry()

		// Should have both OpenAI and Anthropic models (test data includes both)
		openaiModels := registry.GetModelsByBackend(types.BackendOpenAI)
		anthropicModels := registry.GetModelsByBackend(types.BackendAnthropic)

		assert.Greater(t, len(openaiModels), 0, "Should have OpenAI models")
		assert.Greater(t, len(anthropicModels), 0, "Should have Anthropic models")
	})

	t.Run("NoBackendsAvailable", func(t *testing.T) {
		// Create model registry with test data
		registry := helpers.CreateTestModelRegistry()

		// Should have test models
		allModels := registry.GetAllModels()
		assert.Greater(t, len(allModels), 0, "Should have test models")
	})
}

// TestModelConfigurationConsistency tests that model configurations are consistent
func TestModelConfigurationConsistency(t *testing.T) {
	// Create backend manager with both backends
	backendManager := backend.NewBackendManager()
	mockOpenAI := &MockBackend{name: "openai", available: true}
	mockAnthropic := &MockBackend{name: "anthropic", available: true}
	backendManager.RegisterBackend(types.BackendOpenAI, mockOpenAI)
	backendManager.RegisterBackend(types.BackendAnthropic, mockAnthropic)

	// Create model registry
	registry := helpers.CreateTestModelRegistry()

	t.Run("AllModelsHaveRequiredFields", func(t *testing.T) {
		allModels := registry.GetAllModels()
		assert.Greater(t, len(allModels), 0, "Should have models")

		for _, model := range allModels {
			t.Run(model.Name, func(t *testing.T) {
				// Test required fields
				assert.NotEmpty(t, model.Name, "Model should have a name")
				assert.NotEmpty(t, model.DisplayName, "Model should have a display name")
				assert.NotEmpty(t, model.Backend, "Model should have a backend")
				assert.NotEmpty(t, model.BackendModel, "Model should have a backend model")
				assert.NotEmpty(t, model.Family, "Model should have a family")
				assert.NotEmpty(t, model.Description, "Model should have a description")
				assert.Greater(t, model.MaxTokens, 0, "Model should have positive MaxTokens")
				assert.True(t, model.Enabled, "Model should be enabled")
			})
		}
	})

	t.Run("ModelNamesAreUnique", func(t *testing.T) {
		allModels := registry.GetAllModels()
		modelNames := make(map[string]bool)

		for _, model := range allModels {
			assert.False(t, modelNames[model.Name], "Model name %s should be unique", model.Name)
			modelNames[model.Name] = true
		}
	})

	t.Run("BackendModelsAreValid", func(t *testing.T) {
		allModels := registry.GetAllModels()

		for _, model := range allModels {
			t.Run(model.Name, func(t *testing.T) {
				// Test that backend model names are reasonable
				assert.NotEmpty(t, model.BackendModel, "Backend model should not be empty")
				assert.NotContains(t, model.BackendModel, " ", "Backend model should not contain spaces")
			})
		}
	})
}

// TestSpecificModelConfigurations tests specific model configurations
func TestSpecificModelConfigurations(t *testing.T) {
	// Create backend manager with both backends
	backendManager := backend.NewBackendManager()
	mockOpenAI := &MockBackend{name: "openai", available: true}
	mockAnthropic := &MockBackend{name: "anthropic", available: true}
	backendManager.RegisterBackend(types.BackendOpenAI, mockOpenAI)
	backendManager.RegisterBackend(types.BackendAnthropic, mockAnthropic)

	// Create model registry
	registry := helpers.CreateTestModelRegistry()

	t.Run("GPT4oConfiguration", func(t *testing.T) {
		model, exists := registry.GetModel("gpt-4o")
		require.True(t, exists, "gpt-4o model should exist")

		assert.Equal(t, "gpt-4o", model.Name)
		assert.Equal(t, "GPT-4o", model.DisplayName)
		assert.Equal(t, types.BackendOpenAI, model.Backend)
		assert.Equal(t, "gpt-4o", model.BackendModel)
		assert.Equal(t, "gpt", model.Family)
		assert.Equal(t, 16384, model.MaxTokens)
		assert.True(t, model.Enabled)
	})

	t.Run("GPT35TurboConfiguration", func(t *testing.T) {
		model, exists := registry.GetModel("gpt-3.5-turbo")
		require.True(t, exists, "gpt-3.5-turbo model should exist")

		assert.Equal(t, "gpt-3.5-turbo", model.Name)
		assert.Equal(t, "GPT-3.5 Turbo", model.DisplayName)
		assert.Equal(t, types.BackendOpenAI, model.Backend)
		assert.Equal(t, "gpt-3.5-turbo", model.BackendModel)
		assert.Equal(t, "gpt", model.Family)
		assert.Equal(t, 4096, model.MaxTokens) // This is the key test - should be 4096, not 16384
		assert.True(t, model.Enabled)
	})

	t.Run("Claude45SonnetConfiguration", func(t *testing.T) {
		model, exists := registry.GetModel("claude-4.5-sonnet")
		require.True(t, exists, "claude-4.5-sonnet model should exist")

		assert.Equal(t, "claude-4.5-sonnet", model.Name)
		assert.Equal(t, "Claude 4.5 Sonnet", model.DisplayName)
		assert.Equal(t, types.BackendAnthropic, model.Backend)
		assert.Equal(t, "claude-3-5-sonnet-20241022", model.BackendModel)
		assert.Equal(t, "claude", model.Family)
		assert.Equal(t, 200000, model.MaxTokens)
		assert.True(t, model.Enabled)
	})
}
