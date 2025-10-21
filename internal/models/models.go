package models

import (
	"context"
	"fmt"
	"log"

	"go-llm-proxy/internal/backend"
	"go-llm-proxy/internal/config"
	"go-llm-proxy/internal/fetcher"
	"go-llm-proxy/internal/types"
)

// ModelRegistry manages all available models
type ModelRegistry struct {
	models map[string]types.ModelConfig
}

// NewTestModelRegistry creates a new empty model registry for testing
func NewTestModelRegistry() *ModelRegistry {
	return &ModelRegistry{
		models: make(map[string]types.ModelConfig),
	}
}

// NewModelRegistryWithDynamicFetching creates a new model registry with dynamically fetched models
func NewModelRegistryWithDynamicFetching(cfg *config.Config, backendManager *backend.BackendManager, configPath string) (*ModelRegistry, error) {
	registry := &ModelRegistry{
		models: make(map[string]types.ModelConfig),
	}

	// Create model fetcher
	modelFetcher := fetcher.NewModelFetcher(cfg)

	// Load config from file if provided
	if configPath != "" {
		if err := modelFetcher.LoadConfigFromFile(configPath); err != nil {
			log.Printf("Warning: Failed to load config from file %s: %v", configPath, err)
		}
	}

	// Fetch models from APIs
	ctx := context.Background()
	dynamicModels, err := modelFetcher.FetchAllModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch models: %w", err)
	}

	// Filter models to only include those for available backends
	availableBackends := backendManager.GetAvailableBackends()
	backendMap := make(map[types.BackendType]bool)
	for _, backendType := range availableBackends {
		backendMap[backendType] = true
	}

	// Add only models for available backends
	for _, model := range dynamicModels {
		if backendMap[model.Backend] {
			registry.models[model.Name] = model
		}
	}

	log.Printf("Loaded %d models dynamically from APIs", len(registry.models))
	return registry, nil
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
