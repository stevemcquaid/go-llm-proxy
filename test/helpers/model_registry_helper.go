package helpers

import (
	"go-llm-proxy/internal/models"
	"go-llm-proxy/test/fixtures"
)

// CreateTestModelRegistry creates a model registry with test data for testing
func CreateTestModelRegistry() *models.ModelRegistry {
	// Create a new registry
	registry := models.NewTestModelRegistry()

	// Add test models
	testModels := fixtures.GetExpectedModelConfigs()
	for _, model := range testModels {
		registry.AddModel(model)
	}

	return registry
}
