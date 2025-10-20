package llmproxy_unit_test

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

		cfg := config.LoadConfig()

		assert.Equal(t, "11434", cfg.Port)
		assert.Equal(t, "release", cfg.GinMode)
		assert.Equal(t, "", cfg.AnthropicAPIKey)
		assert.Equal(t, "", cfg.OpenAIAPIKey)
		assert.Equal(t, 4096, cfg.DefaultMaxTokens)
		assert.Equal(t, 3, cfg.StreamingChunkSize)
		assert.Equal(t, 50, cfg.StreamingDelay)
	})

	t.Run("LoadConfigWithEnvironment", func(t *testing.T) {
		// Set environment variables
		_ = os.Setenv("PORT", "8080")
		_ = os.Setenv("GIN_MODE", "debug")
		_ = os.Setenv("ANTHROPIC_API_KEY", "test-anthropic-key")
		_ = os.Setenv("OPENAI_API_KEY", "test-openai-key")
		_ = os.Setenv("DEFAULT_MAX_TOKENS", "8192")
		_ = os.Setenv("STREAMING_CHUNK_SIZE", "5")
		_ = os.Setenv("STREAMING_DELAY_MS", "100")

		cfg := config.LoadConfig()

		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, "debug", cfg.GinMode)
		assert.Equal(t, "test-anthropic-key", cfg.AnthropicAPIKey)
		assert.Equal(t, "test-openai-key", cfg.OpenAIAPIKey)
		assert.Equal(t, 8192, cfg.DefaultMaxTokens)
		assert.Equal(t, 5, cfg.StreamingChunkSize)
		assert.Equal(t, 100, cfg.StreamingDelay)

		// Clean up
		os.Clearenv()
	})

	t.Run("IsValid", func(t *testing.T) {
		// Test with no API keys
		cfg := &config.Config{
			Port:            "11434",
			AnthropicAPIKey: "",
			OpenAIAPIKey:    "",
		}
		err := cfg.IsValid()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one API key must be provided")

		// Test with Anthropic API key
		cfg.AnthropicAPIKey = "test-key"
		err = cfg.IsValid()
		assert.NoError(t, err)

		// Test with OpenAI API key
		cfg.AnthropicAPIKey = ""
		cfg.OpenAIAPIKey = "test-key"
		err = cfg.IsValid()
		assert.NoError(t, err)

		// Test with both API keys
		cfg.AnthropicAPIKey = "test-key"
		cfg.OpenAIAPIKey = "test-key"
		err = cfg.IsValid()
		assert.NoError(t, err)

		// Test with empty port
		cfg.Port = ""
		err = cfg.IsValid()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "port must be specified")
	})

	t.Run("HasAnthropic", func(t *testing.T) {
		cfg := &config.Config{AnthropicAPIKey: ""}
		assert.False(t, cfg.HasAnthropic())

		cfg.AnthropicAPIKey = "test-key"
		assert.True(t, cfg.HasAnthropic())
	})

	t.Run("HasOpenAI", func(t *testing.T) {
		cfg := &config.Config{OpenAIAPIKey: ""}
		assert.False(t, cfg.HasOpenAI())

		cfg.OpenAIAPIKey = "test-key"
		assert.True(t, cfg.HasOpenAI())
	})
}

// TestGetEnv tests the config.GetEnv helper function
func TestGetEnv(t *testing.T) {
	t.Run("GetEnvWithValue", func(t *testing.T) {
		_ = os.Setenv("TEST_VAR", "test-value")
		value := config.GetEnv("TEST_VAR", "default")
		assert.Equal(t, "test-value", value)
		_ = os.Unsetenv("TEST_VAR")
	})

	t.Run("GetEnvWithDefault", func(t *testing.T) {
		value := config.GetEnv("NON_EXISTENT_VAR", "default-value")
		assert.Equal(t, "default-value", value)
	})
}

// TestGetEnvInt tests the config.GetEnvInt helper function
func TestGetEnvInt(t *testing.T) {
	t.Run("GetEnvIntWithValue", func(t *testing.T) {
		_ = os.Setenv("TEST_INT", "42")
		value := config.GetEnvInt("TEST_INT", 0)
		assert.Equal(t, 42, value)
		_ = os.Unsetenv("TEST_INT")
	})

	t.Run("GetEnvIntWithInvalidValue", func(t *testing.T) {
		_ = os.Setenv("TEST_INT", "not-a-number")
		value := config.GetEnvInt("TEST_INT", 100)
		assert.Equal(t, 100, value)
		_ = os.Unsetenv("TEST_INT")
	})

	t.Run("GetEnvIntWithDefault", func(t *testing.T) {
		value := config.GetEnvInt("NON_EXISTENT_INT", 200)
		assert.Equal(t, 200, value)
	})
}
