package llmproxy_unit_test

import (
	"os"
	"testing"

	"go-llm-proxy/internal/proxy"
	"go-llm-proxy/internal/types"
	"go-llm-proxy/test/helpers"

	"github.com/stretchr/testify/assert"
)

// TestProxyServerV2Creation tests the creation of the new proxy server
func TestProxyServerV2Creation(t *testing.T) {
	// Skip this test if no API keys are available (since we now fail fast)
	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if anthropicKey == "" && openaiKey == "" {
		t.Skip("Skipping TestProxyServerV2Creation: No API keys available (proxy now fails fast without keys)")
	}

	proxy := proxy.NewProxyServerV2()

	assert.NotNil(t, proxy)
	assert.NotNil(t, proxy.Config)
	assert.NotNil(t, proxy.ModelRegistry)
	assert.NotNil(t, proxy.BackendManager)
	assert.NotNil(t, proxy.StreamingHandler)
}

// TestModelRegistryDefaultModels tests that default models are loaded
func TestModelRegistryDefaultModels(t *testing.T) {
	registry := helpers.CreateTestModelRegistry()

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
