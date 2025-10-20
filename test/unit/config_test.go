package llmproxy_test

import (
	"os"
	"testing"

	"go-llm-proxy/internal/config"

	"github.com/stretchr/testify/assert"
)

// TestConfig tests the configuration management functionality
func TestConfig(t *testing.T) {
	t.Run("LoadConfigWithDefaults", func(t *testing.T) {
		// Clear environment variables
		os.Clearenv()

		config := config.LoadConfig()

		assert.Equal(t, "11434", config.Port)
		assert.Equal(t, "release", config.GinMode)
		assert.Equal(t, "", config.AnthropicAPIKey)
		assert.Equal(t, "", config.OpenAIAPIKey)
		assert.Equal(t, 4096, config.DefaultMaxTokens)
		assert.Equal(t, 3, config.StreamingChunkSize)
		assert.Equal(t, 50, config.StreamingDelay)
	})

	t.Run("LoadConfigWithEnvironment", func(t *testing.T) {
		// Set environment variables
		os.Setenv("PORT", "8080")
		os.Setenv("GIN_MODE", "debug")
		os.Setenv("ANTHROPIC_API_KEY", "test-anthropic-key")
		os.Setenv("OPENAI_API_KEY", "test-openai-key")
		os.Setenv("DEFAULT_MAX_TOKENS", "8192")
		os.Setenv("STREAMING_CHUNK_SIZE", "5")
		os.Setenv("STREAMING_DELAY_MS", "100")

		config := config.LoadConfig()

		assert.Equal(t, "8080", config.Port)
		assert.Equal(t, "debug", config.GinMode)
		assert.Equal(t, "test-anthropic-key", config.AnthropicAPIKey)
		assert.Equal(t, "test-openai-key", config.OpenAIAPIKey)
		assert.Equal(t, 8192, config.DefaultMaxTokens)
		assert.Equal(t, 5, config.StreamingChunkSize)
		assert.Equal(t, 100, config.StreamingDelay)

		// Clean up
		os.Clearenv()
	})

	t.Run("IsValid", func(t *testing.T) {
		// Test with no API keys
		config := &config.Config{
			Port:            "11434",
			AnthropicAPIKey: "",
			OpenAIAPIKey:    "",
		}
		err := config.IsValid()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one API key must be provided")

		// Test with Anthropic API key
		config.AnthropicAPIKey = "test-key"
		err = config.IsValid()
		assert.NoError(t, err)

		// Test with OpenAI API key
		config.AnthropicAPIKey = ""
		config.OpenAIAPIKey = "test-key"
		err = config.IsValid()
		assert.NoError(t, err)

		// Test with both API keys
		config.AnthropicAPIKey = "test-key"
		config.OpenAIAPIKey = "test-key"
		err = config.IsValid()
		assert.NoError(t, err)

		// Test with empty port
		config.Port = ""
		err = config.IsValid()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "port must be specified")
	})

	t.Run("HasAnthropic", func(t *testing.T) {
		config := &config.Config{AnthropicAPIKey: ""}
		assert.False(t, config.HasAnthropic())

		config.AnthropicAPIKey = "test-key"
		assert.True(t, config.HasAnthropic())
	})

	t.Run("HasOpenAI", func(t *testing.T) {
		config := &config.Config{OpenAIAPIKey: ""}
		assert.False(t, config.HasOpenAI())

		config.OpenAIAPIKey = "test-key"
		assert.True(t, config.HasOpenAI())
	})
}

// TestGetEnv tests the config.GetEnv helper function
func TestGetEnv(t *testing.T) {
	t.Run("GetEnvWithValue", func(t *testing.T) {
		os.Setenv("TEST_VAR", "test-value")
		value := config.GetEnv("TEST_VAR", "default")
		assert.Equal(t, "test-value", value)
		os.Unsetenv("TEST_VAR")
	})

	t.Run("GetEnvWithDefault", func(t *testing.T) {
		value := config.GetEnv("NON_EXISTENT_VAR", "default-value")
		assert.Equal(t, "default-value", value)
	})
}

// TestGetEnvInt tests the config.GetEnvInt helper function
func TestGetEnvInt(t *testing.T) {
	t.Run("GetEnvIntWithValue", func(t *testing.T) {
		os.Setenv("TEST_INT", "42")
		value := config.GetEnvInt("TEST_INT", 0)
		assert.Equal(t, 42, value)
		os.Unsetenv("TEST_INT")
	})

	t.Run("GetEnvIntWithInvalidValue", func(t *testing.T) {
		os.Setenv("TEST_INT", "not-a-number")
		value := config.GetEnvInt("TEST_INT", 100)
		assert.Equal(t, 100, value)
		os.Unsetenv("TEST_INT")
	})

	t.Run("GetEnvIntWithDefault", func(t *testing.T) {
		value := config.GetEnvInt("NON_EXISTENT_INT", 200)
		assert.Equal(t, 200, value)
	})
}
