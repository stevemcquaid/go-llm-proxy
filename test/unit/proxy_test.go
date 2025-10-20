package llmproxy_unit_test

import (
	"testing"

	"go-llm-proxy/internal/models"
	"go-llm-proxy/internal/proxy"
	"go-llm-proxy/internal/types"

	"github.com/stretchr/testify/assert"
)

// TestProxyServerV2Creation tests the creation of the new proxy server
func TestProxyServerV2Creation(t *testing.T) {
	proxy := proxy.NewProxyServerV2()

	assert.NotNil(t, proxy)
	assert.NotNil(t, proxy.Config)
	assert.NotNil(t, proxy.ModelRegistry)
	assert.NotNil(t, proxy.BackendManager)
	assert.NotNil(t, proxy.StreamingHandler)
}

// TestModelRegistryDefaultModels tests that default models are loaded
func TestModelRegistryDefaultModels(t *testing.T) {
	registry := models.NewModelRegistry()

	// Check that we have models from both backends
	openaiModels := registry.GetModelsByBackend(types.BackendOpenAI)
	anthropicModels := registry.GetModelsByBackend(types.BackendAnthropic)

	assert.Greater(t, len(openaiModels), 0, "Should have OpenAI models")
	assert.Greater(t, len(anthropicModels), 0, "Should have Anthropic models")

	// Check specific models exist
	_, exists := registry.GetModel("gpt-4o")
	assert.True(t, exists, "Should have gpt-4o model")

	_, exists = registry.GetModel("claude-4.5-sonnet")
	assert.True(t, exists, "Should have claude-4.5-sonnet model")
}

// Note: Comprehensive tests are in separate files:
// - ollama_api_test.go: Tests Ollama API compatibility
// - model_management_test.go: Tests model registry and backend management
// - config_test.go: Tests configuration management
// - integration_test.go: Tests full proxy integration
