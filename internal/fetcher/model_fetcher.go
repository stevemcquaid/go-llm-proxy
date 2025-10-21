package fetcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"go-llm-proxy/internal/config"
	"go-llm-proxy/internal/types"

	"gopkg.in/yaml.v3"
)

// ModelFetcher handles fetching and filtering models from APIs
type ModelFetcher struct {
	apiClient *APIClient
	config    *config.Config
}

// NewModelFetcher creates a new model fetcher
func NewModelFetcher(cfg *config.Config) *ModelFetcher {
	return &ModelFetcher{
		apiClient: NewAPIClient(),
		config:    cfg,
	}
}

// titleCase converts a string to title case
func titleCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// LoadConfigFromFile loads model filter configuration from a YAML file
func (f *ModelFetcher) LoadConfigFromFile(configPath string) error {
	if configPath == "" {
		// Use default config if no path provided
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var configData struct {
		ModelFilters config.ModelFilters `yaml:"model_filters"`
	}

	if err := yaml.Unmarshal(data, &configData); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	f.config.ModelFilters = configData.ModelFilters
	return nil
}

// FetchAllModels fetches models from all enabled backends and applies filters
func (f *ModelFetcher) FetchAllModels(ctx context.Context) ([]types.ModelConfig, error) {
	var allModels []types.ModelConfig

	// Fetch Anthropic models if enabled
	if f.config.ModelFilters.Anthropic.Enabled && f.config.AnthropicAPIKey != "" {
		anthropicModels, err := f.fetchAnthropicModels(ctx)
		if err != nil {
			log.Printf("Warning: Failed to fetch Anthropic models: %v", err)
		} else {
			allModels = append(allModels, anthropicModels...)
		}
	}

	// Fetch OpenAI models if enabled
	if f.config.ModelFilters.OpenAI.Enabled && f.config.OpenAIAPIKey != "" {
		openaiModels, err := f.fetchOpenAIModels(ctx)
		if err != nil {
			log.Printf("Warning: Failed to fetch OpenAI models: %v", err)
		} else {
			allModels = append(allModels, openaiModels...)
		}
	}

	if len(allModels) == 0 {
		return nil, fmt.Errorf("no models could be fetched from any backend")
	}

	return allModels, nil
}

// fetchAnthropicModels fetches and filters Anthropic models
func (f *ModelFetcher) fetchAnthropicModels(ctx context.Context) ([]types.ModelConfig, error) {
	apiModels, err := f.apiClient.FetchAnthropicModels(ctx, f.config.AnthropicAPIKey)
	if err != nil {
		return nil, err
	}

	var models []types.ModelConfig
	for _, apiModel := range apiModels {
		// Apply filters
		if !f.matchesFilters(apiModel.ID, f.config.ModelFilters.Anthropic) {
			continue
		}

		// Convert to our ModelConfig format
		model := types.ModelConfig{
			Name:         f.generateModelName(apiModel.ID, types.BackendAnthropic),
			DisplayName:  f.generateDisplayName(apiModel.ID, types.BackendAnthropic),
			Backend:      types.BackendAnthropic,
			BackendModel: apiModel.ID,
			Family:       f.extractFamily(apiModel.ID, types.BackendAnthropic),
			Description:  apiModel.Description,
			MaxTokens:    apiModel.ContextSize,
			Enabled:      true,
		}

		models = append(models, model)
	}

	return models, nil
}

// fetchOpenAIModels fetches and filters OpenAI models
func (f *ModelFetcher) fetchOpenAIModels(ctx context.Context) ([]types.ModelConfig, error) {
	apiModels, err := f.apiClient.FetchOpenAIModels(ctx, f.config.OpenAIAPIKey)
	if err != nil {
		return nil, err
	}

	var models []types.ModelConfig
	for _, apiModel := range apiModels {
		// Apply filters
		if !f.matchesFilters(apiModel.ID, f.config.ModelFilters.OpenAI) {
			continue
		}

		// Convert to our ModelConfig format
		model := types.ModelConfig{
			Name:         f.generateModelName(apiModel.ID, types.BackendOpenAI),
			DisplayName:  f.generateDisplayName(apiModel.ID, types.BackendOpenAI),
			Backend:      types.BackendOpenAI,
			BackendModel: apiModel.ID,
			Family:       f.extractFamily(apiModel.ID, types.BackendOpenAI),
			Description:  f.generateDescription(apiModel.ID, types.BackendOpenAI),
			MaxTokens:    f.estimateMaxTokens(apiModel.ID, types.BackendOpenAI),
			Enabled:      true,
		}

		models = append(models, model)
	}

	return models, nil
}

// matchesFilters checks if a model ID matches the include/exclude patterns
func (f *ModelFetcher) matchesFilters(modelID string, filter config.ModelFilterConfig) bool {
	// Check exclude patterns first
	for _, pattern := range filter.ExcludePatterns {
		if matched, _ := filepath.Match(pattern, modelID); matched {
			return false
		}
	}

	// If no include patterns, include all (after exclusions)
	if len(filter.IncludePatterns) == 0 {
		return true
	}

	// Check include patterns
	for _, pattern := range filter.IncludePatterns {
		if matched, _ := filepath.Match(pattern, modelID); matched {
			return true
		}
	}

	return false
}

// generateModelName creates a clean model name for our proxy
func (f *ModelFetcher) generateModelName(apiModelID string, backend types.BackendType) string {
	switch backend {
	case types.BackendAnthropic:
		// Convert claude-3-5-sonnet-20241022 to claude-3.5-sonnet
		parts := strings.Split(apiModelID, "-")
		if len(parts) >= 3 {
			// Join first 3 parts and replace last dash with dot
			if len(parts) > 3 {
				parts[2] = strings.Join(parts[2:len(parts)-1], ".")
			}
			return strings.Join(parts[:3], "-")
		}
		return apiModelID
	case types.BackendOpenAI:
		// Use OpenAI model ID as-is for cleaner names
		return apiModelID
	default:
		return apiModelID
	}
}

// generateDisplayName creates a human-readable display name
func (f *ModelFetcher) generateDisplayName(apiModelID string, backend types.BackendType) string {
	switch backend {
	case types.BackendAnthropic:
		// Convert claude-3-5-sonnet-20241022 to Claude 3.5 Sonnet
		parts := strings.Split(apiModelID, "-")
		if len(parts) >= 3 {
			// Capitalize and join
			display := "Claude"
			if len(parts) > 1 {
				display += " " + titleCase(parts[1])
			}
			if len(parts) > 2 {
				display += " " + titleCase(parts[2])
			}
			return display
		}
		return titleCase(apiModelID)
	case types.BackendOpenAI:
		// Convert gpt-4o to GPT-4o
		return strings.ToUpper(apiModelID)
	default:
		return titleCase(apiModelID)
	}
}

// extractFamily extracts the model family from the API model ID
func (f *ModelFetcher) extractFamily(apiModelID string, backend types.BackendType) string {
	switch backend {
	case types.BackendAnthropic:
		// Extract claude from claude-3-5-sonnet-20241022
		parts := strings.Split(apiModelID, "-")
		if len(parts) > 0 {
			return parts[0]
		}
		return "claude"
	case types.BackendOpenAI:
		// Extract gpt from gpt-4o
		parts := strings.Split(apiModelID, "-")
		if len(parts) > 0 {
			return parts[0]
		}
		return "gpt"
	default:
		return "unknown"
	}
}

// generateDescription creates a description for the model
func (f *ModelFetcher) generateDescription(apiModelID string, backend types.BackendType) string {
	switch backend {
	case types.BackendAnthropic:
		return fmt.Sprintf("Anthropic %s model", f.generateDisplayName(apiModelID, backend))
	case types.BackendOpenAI:
		return fmt.Sprintf("OpenAI %s model", f.generateDisplayName(apiModelID, backend))
	default:
		return fmt.Sprintf("%s model", f.generateDisplayName(apiModelID, backend))
	}
}

// estimateMaxTokens estimates max tokens for models where not provided by API
func (f *ModelFetcher) estimateMaxTokens(apiModelID string, backend types.BackendType) int {
	// Common token limits for different model families
	switch backend {
	case types.BackendOpenAI:
		// OpenAI models have known context limits
		if strings.Contains(apiModelID, "gpt-4o") {
			return 128000
		}
		if strings.Contains(apiModelID, "gpt-4") {
			return 8192
		}
		if strings.Contains(apiModelID, "gpt-3.5") {
			return 4096
		}
		return 4096 // Default fallback
	case types.BackendAnthropic:
		// Anthropic models typically have large context windows
		if strings.Contains(apiModelID, "claude-3-5") {
			return 200000
		}
		if strings.Contains(apiModelID, "claude-3") {
			return 200000
		}
		return 100000 // Default fallback
	default:
		return 4096 // Default fallback
	}
}
