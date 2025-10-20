package models

import (
	"go-llm-proxy/internal/backend"
	"go-llm-proxy/internal/types"
)

// ModelRegistry manages all available models
type ModelRegistry struct {
	models map[string]types.ModelConfig
}

// NewModelRegistry creates a new model registry with default models
func NewModelRegistry() *ModelRegistry {
	registry := &ModelRegistry{
		models: make(map[string]types.ModelConfig),
	}

	// Add default models
	registry.addDefaultModels()
	return registry
}

// NewModelRegistryWithBackends creates a new model registry with only models for available backends
func NewModelRegistryWithBackends(backendManager *backend.BackendManager) *ModelRegistry {
	registry := &ModelRegistry{
		models: make(map[string]types.ModelConfig),
	}

	// Add models only for available backends
	registry.addModelsForAvailableBackends(backendManager)
	return registry
}

// addDefaultModels adds all the default models to the registry
func (r *ModelRegistry) addDefaultModels() {
	// Anthropic models
	anthropicModels := []types.ModelConfig{
		{
			Name:         "claude-4.5-sonnet",
			DisplayName:  "Claude 4.5 Sonnet",
			Backend:      types.BackendAnthropic,
			BackendModel: "claude-3-5-sonnet-20241022",
			Family:       "claude",
			Description:  "Enhanced coding and reasoning capabilities with improved performance",
			MaxTokens:    200000,
			Enabled:      true,
		},
		{
			Name:         "claude-4.5-haiku",
			DisplayName:  "Claude 4.5 Haiku",
			Backend:      types.BackendAnthropic,
			BackendModel: "claude-3-5-haiku-20241022",
			Family:       "claude",
			Description:  "Lightweight high-speed model optimized for real-time tasks and coding",
			MaxTokens:    200000,
			Enabled:      true,
		},
		{
			Name:         "claude-4.1-opus",
			DisplayName:  "Claude 4.1 Opus",
			Backend:      types.BackendAnthropic,
			BackendModel: "claude-3-5-opus-20241022",
			Family:       "claude",
			Description:  "Most powerful model for complex reasoning and advanced AI solutions",
			MaxTokens:    200000,
			Enabled:      true,
		},
		{
			Name:         "claude-3.7-sonnet",
			DisplayName:  "Claude 3.7 Sonnet",
			Backend:      types.BackendAnthropic,
			BackendModel: "claude-3-5-sonnet-20241022",
			Family:       "claude",
			Description:  "Hybrid reasoning model with rapid and detailed reasoning options",
			MaxTokens:    200000,
			Enabled:      true,
		},
	}

	// OpenAI models
	openaiModels := []types.ModelConfig{
		{
			Name:         "gpt-5",
			DisplayName:  "GPT-5",
			Backend:      types.BackendOpenAI,
			BackendModel: "gpt-5",
			Family:       "gpt",
			Description:  "Latest multimodal model with reasoning capabilities and unified interface",
			MaxTokens:    4096,
			Enabled:      true,
		},
		{
			Name:         "gpt-4.1",
			DisplayName:  "GPT-4.1",
			Backend:      types.BackendOpenAI,
			BackendModel: "gpt-4.1",
			Family:       "gpt",
			Description:  "Enhanced model with improved coding and long-context comprehension",
			MaxTokens:    4096,
			Enabled:      true,
		},
		{
			Name:         "gpt-4o",
			DisplayName:  "GPT-4o",
			Backend:      types.BackendOpenAI,
			BackendModel: "gpt-4o",
			Family:       "gpt",
			Description:  "Most capable GPT-4 model with multimodal capabilities",
			MaxTokens:    4096,
			Enabled:      true,
		},
		{
			Name:         "gpt-4o-mini",
			DisplayName:  "GPT-4o Mini",
			Backend:      types.BackendOpenAI,
			BackendModel: "gpt-4o-mini",
			Family:       "gpt",
			Description:  "Faster, cheaper GPT-4 model with multimodal support",
			MaxTokens:    4096,
			Enabled:      true,
		},
		{
			Name:         "gpt-4",
			DisplayName:  "GPT-4",
			Backend:      types.BackendOpenAI,
			BackendModel: "gpt-4",
			Family:       "gpt",
			Description:  "Classic GPT-4 model for reliable performance",
			MaxTokens:    4096,
			Enabled:      true,
		},
		{
			Name:         "gpt-3.5-turbo",
			DisplayName:  "GPT-3.5 Turbo",
			Backend:      types.BackendOpenAI,
			BackendModel: "gpt-3.5-turbo",
			Family:       "gpt",
			Description:  "Fast and efficient model for general tasks",
			MaxTokens:    4096,
			Enabled:      true,
		},
	}

	// Add all models to registry
	for _, model := range anthropicModels {
		r.models[model.Name] = model
	}
	for _, model := range openaiModels {
		r.models[model.Name] = model
	}
}

// addModelsForAvailableBackends adds models only for available backends
func (r *ModelRegistry) addModelsForAvailableBackends(backendManager *backend.BackendManager) {
	availableBackends := backendManager.GetAvailableBackends()

	// Check if Anthropic is available
	anthropicAvailable := false
	for _, backendType := range availableBackends {
		if backendType == types.BackendAnthropic {
			anthropicAvailable = true
			break
		}
	}

	// Check if OpenAI is available
	openaiAvailable := false
	for _, backendType := range availableBackends {
		if backendType == types.BackendOpenAI {
			openaiAvailable = true
			break
		}
	}

	// Add Anthropic models only if backend is available
	if anthropicAvailable {
		anthropicModels := []types.ModelConfig{
			{
				Name:         "claude-4.5-sonnet",
				DisplayName:  "Claude 4.5 Sonnet",
				Backend:      types.BackendAnthropic,
				BackendModel: "claude-3-5-sonnet-20241022",
				Family:       "claude",
				Description:  "Enhanced coding and reasoning capabilities with improved performance",
				MaxTokens:    200000,
				Enabled:      true,
			},
			{
				Name:         "claude-4.5-haiku",
				DisplayName:  "Claude 4.5 Haiku",
				Backend:      types.BackendAnthropic,
				BackendModel: "claude-3-5-haiku-20241022",
				Family:       "claude",
				Description:  "Lightweight high-speed model optimized for real-time tasks and coding",
				MaxTokens:    200000,
				Enabled:      true,
			},
			{
				Name:         "claude-4.1-opus",
				DisplayName:  "Claude 4.1 Opus",
				Backend:      types.BackendAnthropic,
				BackendModel: "claude-3-5-opus-20241022",
				Family:       "claude",
				Description:  "Most powerful model for complex reasoning and advanced AI solutions",
				MaxTokens:    200000,
				Enabled:      true,
			},
			{
				Name:         "claude-3.7-sonnet",
				DisplayName:  "Claude 3.7 Sonnet",
				Backend:      types.BackendAnthropic,
				BackendModel: "claude-3-5-sonnet-20241022",
				Family:       "claude",
				Description:  "Hybrid reasoning model with rapid and detailed reasoning options",
				MaxTokens:    200000,
				Enabled:      true,
			},
		}

		for _, model := range anthropicModels {
			r.models[model.Name] = model
		}
	}

	// Add OpenAI models only if backend is available
	if openaiAvailable {
		openaiModels := []types.ModelConfig{
			{
				Name:         "gpt-5",
				DisplayName:  "GPT-5",
				Backend:      types.BackendOpenAI,
				BackendModel: "gpt-5",
				Family:       "gpt",
				Description:  "Latest multimodal model with reasoning capabilities and unified interface",
				MaxTokens:    4096,
				Enabled:      true,
			},
			{
				Name:         "gpt-4.1",
				DisplayName:  "GPT-4.1",
				Backend:      types.BackendOpenAI,
				BackendModel: "gpt-4.1",
				Family:       "gpt",
				Description:  "Enhanced model with improved coding and long-context comprehension",
				MaxTokens:    4096,
				Enabled:      true,
			},
			{
				Name:         "gpt-4o",
				DisplayName:  "GPT-4o",
				Backend:      types.BackendOpenAI,
				BackendModel: "gpt-4o",
				Family:       "gpt",
				Description:  "Most capable GPT-4 model with multimodal capabilities",
				MaxTokens:    4096,
				Enabled:      true,
			},
			{
				Name:         "gpt-4o-mini",
				DisplayName:  "GPT-4o Mini",
				Backend:      types.BackendOpenAI,
				BackendModel: "gpt-4o-mini",
				Family:       "gpt",
				Description:  "Faster, cheaper GPT-4 model with multimodal support",
				MaxTokens:    4096,
				Enabled:      true,
			},
			{
				Name:         "gpt-4",
				DisplayName:  "GPT-4",
				Backend:      types.BackendOpenAI,
				BackendModel: "gpt-4",
				Family:       "gpt",
				Description:  "Classic GPT-4 model for reliable performance",
				MaxTokens:    4096,
				Enabled:      true,
			},
			{
				Name:         "gpt-3.5-turbo",
				DisplayName:  "GPT-3.5 Turbo",
				Backend:      types.BackendOpenAI,
				BackendModel: "gpt-3.5-turbo",
				Family:       "gpt",
				Description:  "Fast and efficient model for general tasks",
				MaxTokens:    4096,
				Enabled:      true,
			},
		}

		for _, model := range openaiModels {
			r.models[model.Name] = model
		}
	}
}

// GetModel returns a model configuration by name
func (r *ModelRegistry) GetModel(name string) (types.ModelConfig, bool) {
	model, exists := r.models[name]
	return model, exists
}

// GetModelsByBackend returns all models for a specific backend
func (r *ModelRegistry) GetModelsByBackend(backend types.BackendType) []types.ModelConfig {
	var models []types.ModelConfig
	for _, model := range r.models {
		if model.Backend == backend && model.Enabled {
			models = append(models, model)
		}
	}
	return models
}

// GetAllModels returns all enabled models
func (r *ModelRegistry) GetAllModels() []types.ModelConfig {
	var models []types.ModelConfig
	for _, model := range r.models {
		if model.Enabled {
			models = append(models, model)
		}
	}
	return models
}

// AddModel adds a new model to the registry
func (r *ModelRegistry) AddModel(model types.ModelConfig) {
	r.models[model.Name] = model
}

// RemoveModel removes a model from the registry
func (r *ModelRegistry) RemoveModel(name string) {
	delete(r.models, name)
}

// EnableModel enables a model
func (r *ModelRegistry) EnableModel(name string) {
	if model, exists := r.models[name]; exists {
		model.Enabled = true
		r.models[name] = model
	}
}

// DisableModel disables a model
func (r *ModelRegistry) DisableModel(name string) {
	if model, exists := r.models[name]; exists {
		model.Enabled = false
		r.models[name] = model
	}
}
