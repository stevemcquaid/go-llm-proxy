package llmproxy_unit_test

import (
	"context"
	"testing"

	"go-llm-proxy/internal/backend"
	"go-llm-proxy/internal/models"
	"go-llm-proxy/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestModelRegistry tests the model registry functionality
func TestModelRegistry(t *testing.T) {
	t.Run("CreateRegistry", func(t *testing.T) {
		registry := models.NewModelRegistry()
		assert.NotNil(t, registry)
		// Note: Cannot access unexported field registry.models from test package
	})

	t.Run("AddModel", func(t *testing.T) {
		registry := models.NewModelRegistry()

		newModel := types.ModelConfig{
			Name:         "test-model",
			DisplayName:  "Test Model",
			Backend:      types.BackendOpenAI,
			BackendModel: "test-model-backend",
			Family:       "test",
			Description:  "A test model",
			MaxTokens:    1000,
			Enabled:      true,
		}

		registry.AddModel(newModel)

		// Verify model was added
		model, exists := registry.GetModel("test-model")
		assert.True(t, exists)
		assert.Equal(t, newModel, model)
	})

	t.Run("GetModel", func(t *testing.T) {
		registry := models.NewModelRegistry()

		// Test existing model
		model, exists := registry.GetModel("gpt-4o")
		assert.True(t, exists)
		assert.Equal(t, "gpt-4o", model.Name)
		assert.Equal(t, types.BackendOpenAI, model.Backend)

		// Test non-existing model
		_, exists = registry.GetModel("non-existing-model")
		assert.False(t, exists)
	})

	t.Run("GetModelsByBackend", func(t *testing.T) {
		registry := models.NewModelRegistry()

		// Get OpenAI models
		openaiModels := registry.GetModelsByBackend(types.BackendOpenAI)
		assert.Greater(t, len(openaiModels), 0)

		for _, model := range openaiModels {
			assert.Equal(t, types.BackendOpenAI, model.Backend)
			assert.True(t, model.Enabled)
		}

		// Get Anthropic models
		anthropicModels := registry.GetModelsByBackend(types.BackendAnthropic)
		assert.Greater(t, len(anthropicModels), 0)

		for _, model := range anthropicModels {
			assert.Equal(t, types.BackendAnthropic, model.Backend)
			assert.True(t, model.Enabled)
		}
	})

	t.Run("GetAllModels", func(t *testing.T) {
		registry := models.NewModelRegistry()

		allModels := registry.GetAllModels()
		assert.Greater(t, len(allModels), 0)

		for _, model := range allModels {
			assert.True(t, model.Enabled)
		}
	})

	t.Run("EnableDisableModel", func(t *testing.T) {
		registry := models.NewModelRegistry()

		// Disable a model
		registry.DisableModel("gpt-3.5-turbo")
		model, exists := registry.GetModel("gpt-3.5-turbo")
		require.True(t, exists)
		assert.False(t, model.Enabled)

		// Re-enable the model
		registry.EnableModel("gpt-3.5-turbo")
		model, exists = registry.GetModel("gpt-3.5-turbo")
		require.True(t, exists)
		assert.True(t, model.Enabled)
	})

	t.Run("RemoveModel", func(t *testing.T) {
		registry := models.NewModelRegistry()

		// Add a test model
		testModel := types.ModelConfig{
			Name:         "temp-model",
			DisplayName:  "Temporary Model",
			Backend:      types.BackendOpenAI,
			BackendModel: "temp-model-backend",
			Family:       "test",
			Description:  "A temporary model",
			MaxTokens:    1000,
			Enabled:      true,
		}
		registry.AddModel(testModel)

		// Verify it exists
		_, exists := registry.GetModel("temp-model")
		assert.True(t, exists)

		// Remove it
		registry.RemoveModel("temp-model")

		// Verify it's gone
		_, exists = registry.GetModel("temp-model")
		assert.False(t, exists)
	})

	t.Run("ToOllamaModel", func(t *testing.T) {
		model := types.ModelConfig{
			Name:         "test-model",
			DisplayName:  "Test Model",
			Backend:      types.BackendOpenAI,
			BackendModel: "test-model-backend",
			Family:       "test",
			Description:  "A test model",
			MaxTokens:    1000,
			Enabled:      true,
		}

		ollamaModel := model.ToOllamaModel()

		assert.Equal(t, model.Name, ollamaModel.Name)
		assert.Equal(t, model.Name, ollamaModel.Model)
		assert.NotEmpty(t, ollamaModel.ModifiedAt)
		assert.Equal(t, int64(1000000000), ollamaModel.Size)
		assert.True(t, len(ollamaModel.Digest) > 0)
		assert.True(t, ollamaModel.Digest[:7] == "sha256:")
	})
}

// TestBackendManager tests the backend manager functionality
func TestBackendManager(t *testing.T) {
	t.Run("CreateBackendManager", func(t *testing.T) {
		manager := backend.NewBackendManager()
		assert.NotNil(t, manager)
		// Note: Cannot access unexported field manager.backends from test package
	})

	t.Run("RegisterBackend", func(t *testing.T) {
		manager := backend.NewBackendManager()

		// Create a mock backend
		mockBackend := &MockBackend{name: "test-backend", available: true}

		manager.RegisterBackend(types.BackendOpenAI, mockBackend)

		// Verify backend was registered
		backend, exists := manager.GetBackend(types.BackendOpenAI)
		assert.True(t, exists)
		assert.Equal(t, mockBackend, backend)
	})

	t.Run("GetAvailableBackends", func(t *testing.T) {
		manager := backend.NewBackendManager()

		// Register available and unavailable backends
		availableBackend := &MockBackend{name: "available", available: true}
		unavailableBackend := &MockBackend{name: "unavailable", available: false}

		manager.RegisterBackend(types.BackendOpenAI, availableBackend)
		manager.RegisterBackend(types.BackendAnthropic, unavailableBackend)

		availableBackends := manager.GetAvailableBackends()
		assert.Len(t, availableBackends, 1)
		assert.Contains(t, availableBackends, types.BackendOpenAI)
		assert.NotContains(t, availableBackends, types.BackendAnthropic)
	})
}

// TestBackendFactory tests the backend factory functionality
func TestBackendFactory(t *testing.T) {
	t.Run("CreateBackendFactory", func(t *testing.T) {
		factory := backend.NewBackendFactory("anthropic-key", "openai-key")
		assert.NotNil(t, factory)
		// Note: Cannot access unexported fields factory.anthropicAPIKey and factory.openaiAPIKey from test package
	})

	t.Run("CreateBackends", func(t *testing.T) {
		factory := backend.NewBackendFactory("anthropic-key", "openai-key")
		manager := factory.CreateBackends()

		assert.NotNil(t, manager)

		// Check that both backends are registered
		anthropicBackend, exists := manager.GetBackend(types.BackendAnthropic)
		assert.True(t, exists)
		assert.True(t, anthropicBackend.IsAvailable())

		openaiBackend, exists := manager.GetBackend(types.BackendOpenAI)
		assert.True(t, exists)
		assert.True(t, openaiBackend.IsAvailable())
	})

	t.Run("CreateBackendsWithMissingKeys", func(t *testing.T) {
		factory := backend.NewBackendFactory("", "")
		manager := factory.CreateBackends()

		assert.NotNil(t, manager)

		// Check that no backends are registered
		availableBackends := manager.GetAvailableBackends()
		assert.Len(t, availableBackends, 0)
	})
}

// TestRequestConversion tests the request/response conversion functions
func TestRequestConversion(t *testing.T) {
	t.Run("llmproxy.ConvertOllamaToGenerateRequest", func(t *testing.T) {
		ollamaReq := types.OllamaGenerateRequest{
			Model:  "gpt-4o",
			Prompt: "Hello world",
		}

		generateReq := types.ConvertOllamaToGenerateRequest(ollamaReq, 1000)

		assert.Equal(t, ollamaReq.Model, generateReq.Model)
		assert.Equal(t, ollamaReq.Prompt, generateReq.Prompt)
		assert.Equal(t, 1000, generateReq.MaxTokens)
	})

	t.Run("llmproxy.ConvertOllamaToChatRequest", func(t *testing.T) {
		ollamaReq := types.OllamaChatRequest{
			Model: "gpt-4o",
			Messages: []types.OllamaMessage{
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi there!"},
			},
		}

		chatReq := types.ConvertOllamaToChatRequest(ollamaReq, 1000)

		assert.Equal(t, ollamaReq.Model, chatReq.Model)
		assert.Equal(t, len(ollamaReq.Messages), len(chatReq.Messages))
		assert.Equal(t, 1000, chatReq.MaxTokens)

		for i, msg := range chatReq.Messages {
			assert.Equal(t, ollamaReq.Messages[i].Role, msg.Role)
			assert.Equal(t, ollamaReq.Messages[i].Content, msg.Content)
		}
	})

	t.Run("llmproxy.ConvertGenerateToOllamaResponse", func(t *testing.T) {
		generateResp := &types.GenerateResponse{
			Model:     "gpt-4o",
			Content:   "Hello world",
			CreatedAt: "1234567890",
		}

		ollamaResp := types.ConvertGenerateToOllamaResponse(generateResp, "gpt-4o")

		assert.Equal(t, "gpt-4o", ollamaResp.Model)
		assert.Equal(t, generateResp.CreatedAt, ollamaResp.CreatedAt)
		assert.Equal(t, generateResp.Content, ollamaResp.Response)
		assert.True(t, ollamaResp.Done)
		assert.NotNil(t, ollamaResp.Context)
	})

	t.Run("llmproxy.ConvertChatToOllamaResponse", func(t *testing.T) {
		chatResp := &types.ChatResponse{
			Model: "gpt-4o",
			Message: types.ChatMessage{
				Role:    "assistant",
				Content: "Hello world",
			},
			CreatedAt: "1234567890",
		}

		ollamaResp := types.ConvertChatToOllamaResponse(chatResp, "gpt-4o")

		assert.Equal(t, "gpt-4o", ollamaResp.Model)
		assert.Equal(t, chatResp.CreatedAt, ollamaResp.CreatedAt)
		assert.Equal(t, chatResp.Message.Role, ollamaResp.Message.Role)
		assert.Equal(t, chatResp.Message.Content, ollamaResp.Message.Content)
		assert.True(t, ollamaResp.Done)
		assert.NotNil(t, ollamaResp.Context)
	})
}

// MockBackend is a mock implementation of BackendHandler for testing
type MockBackend struct {
	name      string
	available bool
}

func (m *MockBackend) Generate(ctx context.Context, req types.GenerateRequest) (*types.GenerateResponse, error) {
	return &types.GenerateResponse{
		Model:     req.Model,
		Content:   "Mock response",
		CreatedAt: "1234567890",
	}, nil
}

func (m *MockBackend) Chat(ctx context.Context, req types.ChatRequest) (*types.ChatResponse, error) {
	return &types.ChatResponse{
		Model: req.Model,
		Message: types.ChatMessage{
			Role:    "assistant",
			Content: "Mock response",
		},
		CreatedAt: "1234567890",
	}, nil
}

func (m *MockBackend) IsAvailable() bool {
	return m.available
}

func (m *MockBackend) GetName() string {
	return m.name
}
