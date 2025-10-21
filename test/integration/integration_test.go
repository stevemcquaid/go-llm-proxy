package llmproxy_integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-llm-proxy/internal/backend"
	"go-llm-proxy/internal/config"
	"go-llm-proxy/internal/proxy"
	"go-llm-proxy/internal/streaming"
	"go-llm-proxy/internal/types"
	"go-llm-proxy/test/helpers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRouter creates a test router with the proxy server
func setupTestRouter(proxy *proxy.ProxyServerV2) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Add routes
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Ollama is running in proxy mode.")
	})

	router.POST("/api/generate", proxy.HandleGenerate)
	router.POST("/api/chat", proxy.HandleChat)
	router.GET("/api/tags", proxy.HandleTags)
	router.GET("/api/version", proxy.HandleVersion)
	router.GET("/api/show/:model", proxy.HandleShow)
	router.GET("/status", func(c *gin.Context) {
		status := proxy.GetHealthStatus()
		c.JSON(200, status)
	})
	router.GET("/health", func(c *gin.Context) {
		status := proxy.GetHealthStatus()
		c.JSON(200, status)
	})

	return router
}

// MockBackend is a mock backend for testing
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

// TestProxyIntegration tests the full proxy integration
func TestProxyIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("FullProxyWorkflow", func(t *testing.T) {
		// Create a proxy with mock backends
		proxy := createTestProxy()
		router := setupTestRouter(proxy)

		// Test 1: Check root endpoint
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Ollama is running in proxy mode.", w.Body.String())

		// Test 2: List available models
		req = httptest.NewRequest("GET", "/api/tags", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var tagsResponse types.OllamaTagsResponse
		err := json.Unmarshal(w.Body.Bytes(), &tagsResponse)
		require.NoError(t, err)
		assert.NotEmpty(t, tagsResponse.Models)

		// Test 3: Get model info
		req = httptest.NewRequest("GET", "/api/show/gpt-4o", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var modelResponse types.OllamaModel
		err = json.Unmarshal(w.Body.Bytes(), &modelResponse)
		require.NoError(t, err)
		assert.Equal(t, "gpt-4o", modelResponse.Name)

		// Test 4: Test chat endpoint (will use mock backend)
		chatReq := types.OllamaChatRequest{
			Model: "gpt-4o",
			Messages: []types.OllamaMessage{
				{Role: "user", Content: "Hello"},
			},
			Stream: false,
		}

		reqBody, _ := json.Marshal(chatReq)
		req = httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var chatResponse types.OllamaChatResponse
		err = json.Unmarshal(w.Body.Bytes(), &chatResponse)
		require.NoError(t, err)
		assert.Equal(t, "gpt-4o", chatResponse.Model)
		assert.True(t, chatResponse.Done)
		assert.Equal(t, "assistant", chatResponse.Message.Role)
		assert.NotEmpty(t, chatResponse.Message.Content)

		// Test 5: Test streaming chat
		chatReq.Stream = true
		reqBody, _ = json.Marshal(chatReq)
		req = httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/x-ndjson", w.Header().Get("Content-Type"))

		// Test 6: Test generate endpoint
		generateReq := types.OllamaGenerateRequest{
			Model:  "gpt-4o",
			Prompt: "Hello",
			Stream: false,
		}

		reqBody, _ = json.Marshal(generateReq)
		req = httptest.NewRequest("POST", "/api/generate", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var generateResponse types.OllamaGenerateResponse
		err = json.Unmarshal(w.Body.Bytes(), &generateResponse)
		require.NoError(t, err)
		assert.Equal(t, "gpt-4o", generateResponse.Model)
		assert.True(t, generateResponse.Done)
		assert.NotEmpty(t, generateResponse.Response)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		proxy := createTestProxy()
		router := setupTestRouter(proxy)

		// Test invalid JSON
		req := httptest.NewRequest("POST", "/api/chat", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Test missing model
		chatReq := types.OllamaChatRequest{
			Messages: []types.OllamaMessage{
				{Role: "user", Content: "Hello"},
			},
		}
		reqBody, _ := json.Marshal(chatReq)
		req = httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Test unknown model
		chatReq.Model = "unknown-model"
		reqBody, _ = json.Marshal(chatReq)
		req = httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("CORSHandling", func(t *testing.T) {
		proxy := createTestProxy()
		router := setupTestRouter(proxy)

		// Test OPTIONS request
		req := httptest.NewRequest("OPTIONS", "/api/chat", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	})

	t.Run("HealthEndpoints", func(t *testing.T) {
		proxy := createTestProxy()
		router := setupTestRouter(proxy)

		// Test /status endpoint
		req := httptest.NewRequest("GET", "/status", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var statusResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &statusResponse)
		require.NoError(t, err)
		assert.Contains(t, statusResponse, "status")
		assert.Contains(t, statusResponse, "available_backends")
		assert.Contains(t, statusResponse, "total_models")

		// Test /health endpoint
		req = httptest.NewRequest("GET", "/health", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestStreamingIntegration tests streaming functionality
func TestStreamingIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("StreamingChatResponse", func(t *testing.T) {
		proxy := createTestProxy()
		router := setupTestRouter(proxy)

		chatReq := types.OllamaChatRequest{
			Model: "gpt-4o",
			Messages: []types.OllamaMessage{
				{Role: "user", Content: "Hello"},
			},
			Stream: true,
		}

		reqBody, _ := json.Marshal(chatReq)
		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/x-ndjson", w.Header().Get("Content-Type"))
		assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
		assert.Equal(t, "keep-alive", w.Header().Get("Connection"))

		// Parse streaming response
		body := w.Body.String()
		lines := strings.Split(strings.TrimSpace(body), "\n")
		assert.Greater(t, len(lines), 0, "Should have at least one streaming chunk")

		// Verify each chunk is valid JSON
		for i, line := range lines {
			if line == "" {
				continue
			}
			var chunk types.OllamaChatResponse
			err := json.Unmarshal([]byte(line), &chunk)
			assert.NoError(t, err, "Line %d should be valid JSON: %s", i, line)
			assert.Equal(t, "gpt-4o", chunk.Model)
			assert.NotEmpty(t, chunk.CreatedAt)
			assert.NotNil(t, chunk.Context)
			assert.Equal(t, "assistant", chunk.Message.Role)
		}
	})

	t.Run("StreamingGenerateResponse", func(t *testing.T) {
		proxy := createTestProxy()
		router := setupTestRouter(proxy)

		generateReq := types.OllamaGenerateRequest{
			Model:  "gpt-4o",
			Prompt: "Hello",
			Stream: true,
		}

		reqBody, _ := json.Marshal(generateReq)
		req := httptest.NewRequest("POST", "/api/generate", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Generate endpoint should return error for streaming
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestModelManagementIntegration tests model management functionality
func TestModelManagementIntegration(t *testing.T) {
	t.Run("AddRemoveModels", func(t *testing.T) {
		registry := helpers.CreateTestModelRegistry()
		initialCount := len(registry.GetAllModels())

		// Add a new model
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
		assert.Equal(t, initialCount+1, len(registry.GetAllModels()))
		model, exists := registry.GetModel("test-model")
		assert.True(t, exists)
		assert.Equal(t, newModel, model)

		// Remove the model
		registry.RemoveModel("test-model")
		assert.Equal(t, initialCount, len(registry.GetAllModels()))
		_, exists = registry.GetModel("test-model")
		assert.False(t, exists)
	})

	t.Run("EnableDisableModels", func(t *testing.T) {
		registry := helpers.CreateTestModelRegistry()

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
}

// TestBackendIntegration tests backend integration
func TestBackendIntegration(t *testing.T) {
	t.Run("BackendManager", func(t *testing.T) {
		manager := backend.NewBackendManager()

		// Register mock backends
		openaiBackend := &MockBackend{name: "openai", available: true}
		anthropicBackend := &MockBackend{name: "anthropic", available: false}

		manager.RegisterBackend(types.BackendOpenAI, openaiBackend)
		manager.RegisterBackend(types.BackendAnthropic, anthropicBackend)

		// Test getting backends
		backend, exists := manager.GetBackend(types.BackendOpenAI)
		assert.True(t, exists)
		assert.Equal(t, openaiBackend, backend)

		// Test available backends
		availableBackends := manager.GetAvailableBackends()
		assert.Len(t, availableBackends, 1)
		assert.Contains(t, availableBackends, types.BackendOpenAI)
		assert.NotContains(t, availableBackends, types.BackendAnthropic)
	})

	t.Run("BackendFactory", func(t *testing.T) {
		factory := backend.NewBackendFactory("anthropic-key", "openai-key")
		manager := factory.CreateBackends()

		assert.NotNil(t, manager)

		// Both backends should be available
		availableBackends := manager.GetAvailableBackends()
		assert.Len(t, availableBackends, 2)
		assert.Contains(t, availableBackends, types.BackendOpenAI)
		assert.Contains(t, availableBackends, types.BackendAnthropic)
	})
}

// createTestProxy creates a test proxy with mock backends
func createTestProxy() *proxy.ProxyServerV2 {
	// Create backend manager with mock backends
	backendManager := backend.NewBackendManager()
	mockOpenAI := &MockBackend{name: "openai", available: true}
	mockAnthropic := &MockBackend{name: "anthropic", available: true}
	backendManager.RegisterBackend(types.BackendOpenAI, mockOpenAI)
	backendManager.RegisterBackend(types.BackendAnthropic, mockAnthropic)

	// Create model registry with available backends
	modelRegistry := helpers.CreateTestModelRegistry()

	// Create streaming handler
	streamingHandler := streaming.NewStreamingHandler(backendManager, modelRegistry)

	// Create config
	config := &config.Config{
		Port:               "11434",
		GinMode:            "test",
		AnthropicAPIKey:    "test-key",
		OpenAIAPIKey:       "test-key",
		DefaultMaxTokens:   4096,
		StreamingChunkSize: 3,
		StreamingDelay:     50,
	}

	return &proxy.ProxyServerV2{
		Config:           config,
		ModelRegistry:    modelRegistry,
		BackendManager:   backendManager,
		StreamingHandler: streamingHandler,
	}
}
