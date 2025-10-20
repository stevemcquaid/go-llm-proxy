package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the proxy
type Config struct {
	// Server configuration
	Port    string `json:"port"`
	GinMode string `json:"gin_mode"`

	// API Keys
	AnthropicAPIKey string `json:"anthropic_api_key"`
	OpenAIAPIKey    string `json:"openai_api_key"`

	// Model configuration
	DefaultMaxTokens int `json:"default_max_tokens"`

	// Streaming configuration
	StreamingChunkSize int `json:"streaming_chunk_size"`
	StreamingDelay     int `json:"streaming_delay_ms"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	config := &Config{
		// Default values
		Port:               GetEnv("PORT", "11434"),
		GinMode:            GetEnv("GIN_MODE", "release"),
		AnthropicAPIKey:    GetEnv("ANTHROPIC_API_KEY", ""),
		OpenAIAPIKey:       GetEnv("OPENAI_API_KEY", ""),
		DefaultMaxTokens:   GetEnvInt("DEFAULT_MAX_TOKENS", 4096),
		StreamingChunkSize: GetEnvInt("STREAMING_CHUNK_SIZE", 3),
		StreamingDelay:     GetEnvInt("STREAMING_DELAY_MS", 50),
	}

	return config
}

// GetEnv gets an environment variable with a default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt gets an environment variable as an integer with a default value
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// IsValid checks if the configuration is valid
func (c *Config) IsValid() error {
	if c.AnthropicAPIKey == "" && c.OpenAIAPIKey == "" {
		return fmt.Errorf("at least one API key must be provided")
	}

	if c.Port == "" {
		return fmt.Errorf("port must be specified")
	}

	return nil
}

// HasAnthropic returns true if Anthropic API key is configured
func (c *Config) HasAnthropic() bool {
	return c.AnthropicAPIKey != ""
}

// HasOpenAI returns true if OpenAI API key is configured
func (c *Config) HasOpenAI() bool {
	return c.OpenAIAPIKey != ""
}
