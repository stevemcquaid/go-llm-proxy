package main

import (
	"fmt"
)

// ExampleAddModel demonstrates how to add a new model
func ExampleAddModel() {
	// Create a new model registry
	registry := NewModelRegistry()

	// Add a new model
	newModel := ModelConfig{
		Name:         "gpt-4-turbo",
		DisplayName:  "GPT-4 Turbo",
		Backend:      BackendOpenAI,
		BackendModel: "gpt-4-turbo-preview",
		Family:       "gpt",
		Description:  "Latest GPT-4 Turbo model",
		MaxTokens:    128000,
		Enabled:      true,
	}

	registry.AddModel(newModel)

	// Verify the model was added
	if model, exists := registry.GetModel("gpt-4-turbo"); exists {
		fmt.Printf("Added model: %s (%s)\n", model.DisplayName, model.Description)
	}

	// List all models
	models := registry.GetAllModels()
	fmt.Printf("Total models: %d\n", len(models))
	for _, model := range models {
		fmt.Printf("- %s (%s)\n", model.Name, model.Backend)
	}
}

// ExampleDisableModel demonstrates how to disable a model
func ExampleDisableModel() {
	registry := NewModelRegistry()

	// Disable a model
	registry.DisableModel("gpt-3.5-turbo")

	// Check if it's disabled
	model, exists := registry.GetModel("gpt-3.5-turbo")
	if exists {
		fmt.Printf("Model %s enabled: %v\n", model.Name, model.Enabled)
	}
}

// ExampleAddBackend demonstrates how to add a new backend
func ExampleAddBackend() {
	// This would be implemented in a separate file
	// For example, to add a new backend like Cohere:

	// 1. Create a new backend type
	// const BackendCohere BackendType = "cohere"

	// 2. Implement the BackendHandler interface
	// type CohereBackend struct {
	//     apiKey string
	//     client *http.Client
	// }

	// func (cb *CohereBackend) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
	//     // Implementation
	// }

	// func (cb *CohereBackend) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	//     // Implementation
	// }

	// func (cb *CohereBackend) IsAvailable() bool {
	//     return cb.apiKey != ""
	// }

	// func (cb *CohereBackend) GetName() string {
	//     return "cohere"
	// }

	// 3. Register the backend in the factory
	// func (bf *BackendFactory) CreateBackends() *BackendManager {
	//     manager := NewBackendManager()
	//
	//     if bf.cohereAPIKey != "" {
	//         cohereBackend := NewCohereBackend(bf.cohereAPIKey)
	//         manager.RegisterBackend(BackendCohere, cohereBackend)
	//     }
	//
	//     return manager
	// }

	fmt.Println("To add a new backend, implement the BackendHandler interface")
	fmt.Println("and register it in the BackendFactory.CreateBackends() method")
}

func main() {
	fmt.Println("=== Example: Adding a Model ===")
	ExampleAddModel()

	fmt.Println("\n=== Example: Disabling a Model ===")
	ExampleDisableModel()

	fmt.Println("\n=== Example: Adding a Backend ===")
	ExampleAddBackend()
}
