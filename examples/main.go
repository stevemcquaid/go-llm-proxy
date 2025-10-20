package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from parent directory
	envPath := filepath.Join("..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	fmt.Println("=== Go LLM Proxy Examples ===")
	fmt.Println()

	// Example 1: Adding a new model
	fmt.Println("1. Adding a new model:")
	exampleAddModel()
	fmt.Println()

	// Example 2: Disabling a model
	fmt.Println("2. Disabling a model:")
	exampleDisableModel()
	fmt.Println()

	// Example 3: Listing models by backend
	fmt.Println("3. Listing models by backend:")
	exampleListModelsByBackend()
	fmt.Println()

	// Example 4: Configuration example
	fmt.Println("4. Configuration example:")
	exampleConfiguration()
	fmt.Println()

	fmt.Println("=== Examples completed ===")
}

// exampleAddModel demonstrates how to add a new model
func exampleAddModel() {
	// This would be implemented in the main package
	fmt.Println("To add a new model, use the ModelRegistry:")
	fmt.Println("  registry := NewModelRegistry()")
	fmt.Println("  newModel := ModelConfig{...}")
	fmt.Println("  registry.AddModel(newModel)")
}

// exampleDisableModel demonstrates how to disable a model
func exampleDisableModel() {
	fmt.Println("To disable a model:")
	fmt.Println("  registry.DisableModel(\"model-name\")")
	fmt.Println("To re-enable a model:")
	fmt.Println("  registry.EnableModel(\"model-name\")")
}

// exampleListModelsByBackend demonstrates how to list models by backend
func exampleListModelsByBackend() {
	fmt.Println("To list models by backend:")
	fmt.Println("  openaiModels := registry.GetModelsByBackend(BackendOpenAI)")
	fmt.Println("  anthropicModels := registry.GetModelsByBackend(BackendAnthropic)")
}

// exampleConfiguration demonstrates configuration usage
func exampleConfiguration() {
	fmt.Println("Configuration is loaded from environment variables:")
	fmt.Println("  ANTHROPIC_API_KEY=your_key_here")
	fmt.Println("  OPENAI_API_KEY=your_key_here")
	fmt.Println("  PORT=11434")
	fmt.Println("  GIN_MODE=release")
}
