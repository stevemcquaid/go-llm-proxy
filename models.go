package main

import "time"

// ModelConfig represents configuration for a model
type ModelConfig struct {
	Name         string      `json:"name"`
	DisplayName  string      `json:"display_name"`
	Backend      BackendType `json:"backend"`
	BackendModel string      `json:"backend_model"`
	Family       string      `json:"family"`
	Description  string      `json:"description"`
	MaxTokens    int         `json:"max_tokens"`
	Enabled      bool        `json:"enabled"`
}

// ModelRegistry manages all available models
type ModelRegistry struct {
	models map[string]ModelConfig
}

// NewModelRegistry creates a new model registry with default models
func NewModelRegistry() *ModelRegistry {
	registry := &ModelRegistry{
		models: make(map[string]ModelConfig),
	}

	// Add default models
	registry.addDefaultModels()
	return registry
}

// addDefaultModels adds all the default models to the registry
func (r *ModelRegistry) addDefaultModels() {
	// Anthropic models
	anthropicModels := []ModelConfig{
		{
			Name:         "claude-3.5-sonnet",
			DisplayName:  "Claude 3.5 Sonnet",
			Backend:      BackendAnthropic,
			BackendModel: "claude-3-5-sonnet-20241022",
			Family:       "claude",
			Description:  "Most capable model for complex tasks",
			MaxTokens:    200000,
			Enabled:      true,
		},
		{
			Name:         "claude-3.5-haiku",
			DisplayName:  "Claude 3.5 Haiku",
			Backend:      BackendAnthropic,
			BackendModel: "claude-3-5-haiku-20241022",
			Family:       "claude",
			Description:  "Fast and efficient model",
			MaxTokens:    200000,
			Enabled:      true,
		},
		{
			Name:         "claude-3.5-opus",
			DisplayName:  "Claude 3.5 Opus",
			Backend:      BackendAnthropic,
			BackendModel: "claude-3-5-opus-20241022",
			Family:       "claude",
			Description:  "Most powerful model for complex reasoning",
			MaxTokens:    200000,
			Enabled:      true,
		},
	}

	// OpenAI models
	openaiModels := []ModelConfig{
		{
			Name:         "gpt-4o",
			DisplayName:  "GPT-4o",
			Backend:      BackendOpenAI,
			BackendModel: "gpt-4o",
			Family:       "gpt",
			Description:  "Most capable GPT-4 model",
			MaxTokens:    16384,
			Enabled:      true,
		},
		{
			Name:         "gpt-4o-mini",
			DisplayName:  "GPT-4o Mini",
			Backend:      BackendOpenAI,
			BackendModel: "gpt-4o-mini",
			Family:       "gpt",
			Description:  "Faster, cheaper GPT-4 model",
			MaxTokens:    16384,
			Enabled:      true,
		},
		{
			Name:         "gpt-4",
			DisplayName:  "GPT-4",
			Backend:      BackendOpenAI,
			BackendModel: "gpt-4",
			Family:       "gpt",
			Description:  "Classic GPT-4 model",
			MaxTokens:    8192,
			Enabled:      true,
		},
		{
			Name:         "gpt-3.5-turbo",
			DisplayName:  "GPT-3.5 Turbo",
			Backend:      BackendOpenAI,
			BackendModel: "gpt-3.5-turbo",
			Family:       "gpt",
			Description:  "Fast and efficient model",
			MaxTokens:    16384,
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

// GetModel returns a model configuration by name
func (r *ModelRegistry) GetModel(name string) (ModelConfig, bool) {
	model, exists := r.models[name]
	return model, exists
}

// GetModelsByBackend returns all models for a specific backend
func (r *ModelRegistry) GetModelsByBackend(backend BackendType) []ModelConfig {
	var models []ModelConfig
	for _, model := range r.models {
		if model.Backend == backend && model.Enabled {
			models = append(models, model)
		}
	}
	return models
}

// GetAllModels returns all enabled models
func (r *ModelRegistry) GetAllModels() []ModelConfig {
	var models []ModelConfig
	for _, model := range r.models {
		if model.Enabled {
			models = append(models, model)
		}
	}
	return models
}

// AddModel adds a new model to the registry
func (r *ModelRegistry) AddModel(model ModelConfig) {
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

// ToOllamaModel converts a ModelConfig to OllamaModel format
func (m ModelConfig) ToOllamaModel() OllamaModel {
	return OllamaModel{
		Name:       m.Name,
		Model:      m.Name,
		ModifiedAt: time.Now().Format("2006-01-02T15:04:05.000Z"),
		Size:       1000000000, // 1GB placeholder
		Digest:     "sha256:" + m.Name,
	}
}
