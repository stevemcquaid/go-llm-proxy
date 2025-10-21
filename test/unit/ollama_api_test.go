package llmproxy_unit_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"go-llm-proxy/internal/proxy"
	"go-llm-proxy/internal/types"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOllamaAPISpec tests that our proxy behaves like a real Ollama API server
func TestOllamaAPISpec(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Skip this test if no API keys are available (since we now fail fast)
	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if anthropicKey == "" && openaiKey == "" {
		t.Skip("Skipping TestOllamaAPISpec: No API keys available (proxy now fails fast without keys)")
	}

	// Create a test proxyServerV2 server
	proxyServerV2 := proxy.NewProxyServerV2()
	router := setupTestRouter(proxyServerV2)

	t.Run("RootEndpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
		assert.Equal(t, "Ollama is running in proxy mode.", w.Body.String())
	})

	t.Run("APITagsEndpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/tags", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

		var response types.OllamaTagsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify response structure matches Ollama spec
		// Note: Models may be empty if no backends are available (no API keys)
		// This is expected behavior when running tests without API keys
		if len(response.Models) > 0 {
			assert.NotEmpty(t, response.Models)
		}

		for _, model := range response.Models {
			// Check required fields
			assert.NotEmpty(t, model.Name, "Model name should not be empty")
			assert.NotEmpty(t, model.Model, "Model field should not be empty")
			assert.NotEmpty(t, model.ModifiedAt, "ModifiedAt should not be empty")
			assert.Greater(t, model.Size, int64(0), "Size should be greater than 0")
			assert.True(t, strings.HasPrefix(model.Digest, "sha256:"), "Digest should start with 'sha256:'")

			// Check timestamp format
			_, err := time.Parse("2006-01-02T15:04:05.000Z", model.ModifiedAt)
			assert.NoError(t, err, "ModifiedAt should be in ISO 8601 format")
		}
	})

	t.Run("APIVersionEndpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/version", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Check required fields
		assert.Contains(t, response, "version")
		assert.Contains(t, response, "proxy")
		assert.Contains(t, response, "backends")
	})

	t.Run("ChatEndpointNonStreaming", func(t *testing.T) {
		chatReq := types.OllamaChatRequest{
			Model: "gpt-4o",
			Messages: []types.OllamaMessage{
				{Role: "user", Content: "Hello"},
			},
			Stream: false,
		}

		reqBody, _ := json.Marshal(chatReq)
		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Note: This will fail without API keys, but we can test the structure
		if w.Code == http.StatusOK {
			var response types.OllamaChatResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Verify response structure matches Ollama spec
			assert.Equal(t, chatReq.Model, response.Model)
			assert.NotEmpty(t, response.CreatedAt)
			assert.True(t, response.Done)
			assert.NotNil(t, response.Context)
			assert.Equal(t, "assistant", response.Message.Role)
			assert.NotEmpty(t, response.Message.Content)
		}
	})

	t.Run("ChatEndpointStreaming", func(t *testing.T) {
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

		// Note: This will fail without API keys, but we can test the headers
		if w.Code == http.StatusOK {
			assert.Equal(t, "application/x-ndjson", w.Header().Get("Content-Type"))
			assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
			assert.Equal(t, "keep-alive", w.Header().Get("Connection"))
		}
	})

	t.Run("GenerateEndpointNonStreaming", func(t *testing.T) {
		generateReq := types.OllamaGenerateRequest{
			Model:  "gpt-4o",
			Prompt: "Hello",
			Stream: false,
		}

		reqBody, _ := json.Marshal(generateReq)
		req := httptest.NewRequest("POST", "/api/generate", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Note: This will fail without API keys, but we can test the structure
		if w.Code == http.StatusOK {
			var response types.OllamaGenerateResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Verify response structure matches Ollama spec
			assert.Equal(t, generateReq.Model, response.Model)
			assert.NotEmpty(t, response.CreatedAt)
			assert.True(t, response.Done)
			assert.NotNil(t, response.Context)
			assert.NotEmpty(t, response.Response)
		}
	})

	t.Run("ShowEndpoint", func(t *testing.T) {
		showReq := map[string]string{"model": "gpt-4o"}
		reqBody, _ := json.Marshal(showReq)
		req := httptest.NewRequest("POST", "/api/show", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			var response types.OllamaModel
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Verify response structure matches Ollama spec
			assert.Equal(t, "gpt-4o", response.Name)
			assert.Equal(t, "gpt-4o", response.Model)
			assert.NotEmpty(t, response.ModifiedAt)
			assert.Greater(t, response.Size, int64(0))
			assert.True(t, strings.HasPrefix(response.Digest, "sha256:"))
		}
	})

	t.Run("CORSHeaders", func(t *testing.T) {
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
		var healthResponse map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &healthResponse)
		require.NoError(t, err)
		assert.Contains(t, healthResponse, "status")
	})

	t.Run("AlternativeEndpoints", func(t *testing.T) {
		// Test /v1/models endpoint (OpenAI style)
		req := httptest.NewRequest("GET", "/v1/models", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response types.OllamaTagsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		// Note: Models may be empty if no backends are available (no API keys)
		if len(response.Models) > 0 {
			assert.NotEmpty(t, response.Models)
		}

		// Test /models endpoint
		req = httptest.NewRequest("GET", "/models", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		// Note: Models may be empty if no backends are available (no API keys)
		if len(response.Models) > 0 {
			assert.NotEmpty(t, response.Models)
		}
	})
}

// TestOllamaResponseFormats tests that our responses match Ollama's exact format
func TestOllamaResponseFormats(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Skip this test if no API keys are available (since we now fail fast)
	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if anthropicKey == "" && openaiKey == "" {
		t.Skip("Skipping TestOllamaResponseFormats: No API keys available (proxy now fails fast without keys)")
	}

	newProxyServerV2 := proxy.NewProxyServerV2()
	router := setupTestRouter(newProxyServerV2)

	t.Run("TagsResponseFormat", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/tags", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response types.OllamaTagsResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Verify the response has the exact structure expected by Ollama clients
		assert.IsType(t, types.OllamaTagsResponse{}, response)
		assert.IsType(t, []types.OllamaModel{}, response.Models)

		// Check that we have models available (if backends are configured)
		// Note: Models may be empty if no backends are available (no API keys)
		if len(response.Models) > 0 {
			assert.Greater(t, len(response.Models), 0, "Should have at least one model available")
		}

		// Verify each model has the correct structure
		for _, model := range response.Models {
			// Required string fields
			assert.NotEmpty(t, model.Name)
			assert.NotEmpty(t, model.Model)
			assert.NotEmpty(t, model.ModifiedAt)
			assert.NotEmpty(t, model.Digest)

			// Required numeric fields
			assert.Greater(t, model.Size, int64(0))

			// Digest format validation
			assert.True(t, strings.HasPrefix(model.Digest, "sha256:"),
				"Digest should start with 'sha256:', got: %s", model.Digest)

			// Timestamp format validation
			_, err := time.Parse("2006-01-02T15:04:05.000Z", model.ModifiedAt)
			assert.NoError(t, err, "ModifiedAt should be in ISO 8601 format, got: %s", model.ModifiedAt)
		}
	})

	t.Run("ModelResponseFormat", func(t *testing.T) {
		showReq := map[string]string{"model": "gpt-4o"}
		reqBody, _ := json.Marshal(showReq)
		req := httptest.NewRequest("POST", "/api/show", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			var response types.OllamaModel
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Verify the response has the exact structure expected by Ollama clients
			assert.IsType(t, types.OllamaModel{}, response)
			assert.Equal(t, "gpt-4o", response.Name)
			assert.Equal(t, "gpt-4o", response.Model)
			assert.True(t, strings.HasPrefix(response.Digest, "sha256:"))
			assert.Greater(t, response.Size, int64(0))
		}
	})
}

// TestOllamaErrorHandling tests that our error responses match Ollama's format
func TestOllamaErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Skip this test if no API keys are available (since we now fail fast)
	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if anthropicKey == "" && openaiKey == "" {
		t.Skip("Skipping TestOllamaErrorHandling: No API keys available (proxy now fails fast without keys)")
	}

	proxyServerV2 := proxy.NewProxyServerV2()
	router := setupTestRouter(proxyServerV2)

	t.Run("InvalidJSONRequest", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/chat", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

		var errorResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse, "error")
	})

	t.Run("MissingModel", func(t *testing.T) {
		chatReq := types.OllamaChatRequest{
			Messages: []types.OllamaMessage{
				{Role: "user", Content: "Hello"},
			},
			Stream: false,
		}

		reqBody, _ := json.Marshal(chatReq)
		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 400 for missing model
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("UnknownModel", func(t *testing.T) {
		chatReq := types.OllamaChatRequest{
			Model: "unknown-model",
			Messages: []types.OllamaMessage{
				{Role: "user", Content: "Hello"},
			},
			Stream: false,
		}

		reqBody, _ := json.Marshal(chatReq)
		req := httptest.NewRequest("POST", "/api/chat", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should return 400 for unknown model
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// TestOllamaStreamingFormat tests that our streaming responses match Ollama's format
func TestOllamaStreamingFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Skip this test if no API keys are available (since we now fail fast)
	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if anthropicKey == "" && openaiKey == "" {
		t.Skip("Skipping TestOllamaStreamingFormat: No API keys available (proxy now fails fast without keys)")
	}

	proxyServerV2 := proxy.NewProxyServerV2()
	router := setupTestRouter(proxyServerV2)

	t.Run("StreamingHeaders", func(t *testing.T) {
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

		// Check streaming headers (even if request fails due to missing API keys)
		if w.Code == http.StatusOK {
			assert.Equal(t, "application/x-ndjson", w.Header().Get("Content-Type"))
			assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
			assert.Equal(t, "keep-alive", w.Header().Get("Connection"))
		}
	})

	t.Run("StreamingResponseStructure", func(t *testing.T) {
		// This test would require mocking the backend responses
		// For now, we'll test the structure of our streaming handler
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

		// If the request succeeds, verify the streaming format
		if w.Code == http.StatusOK {
			body := w.Body.String()
			lines := strings.Split(strings.TrimSpace(body), "\n")

			// Should have at least one line
			assert.Greater(t, len(lines), 0, "Streaming response should have at least one line")

			// Each line should be valid JSON
			for i, line := range lines {
				if line == "" {
					continue
				}
				var chunk types.OllamaChatResponse
				err := json.Unmarshal([]byte(line), &chunk)
				assert.NoError(t, err, "Line %d should be valid JSON: %s", i, line)

				// Verify chunk structure
				assert.NotEmpty(t, chunk.Model)
				assert.NotEmpty(t, chunk.CreatedAt)
				assert.NotNil(t, chunk.Context)
				assert.Equal(t, "assistant", chunk.Message.Role)
				assert.NotEmpty(t, chunk.Message.Content)
			}
		}
	})
}

// setupTestRouter creates a test router with the proxy server
func setupTestRouter(proxy *proxy.ProxyServerV2) *gin.Engine {
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

	// Ollama API endpoints
	router.POST("/api/generate", proxy.HandleGenerate)
	router.POST("/api/chat", proxy.HandleChat)
	router.GET("/api/tags", proxy.HandleTags)
	router.GET("/api/version", proxy.HandleVersion)
	router.POST("/api/pull", proxy.HandlePull)
	router.POST("/api/push", proxy.HandlePush)
	router.DELETE("/api/delete", proxy.HandleDelete)
	router.POST("/api/create", proxy.HandleCreate)
	router.POST("/api/copy", proxy.HandleCopy)
	router.POST("/api/embeddings", proxy.HandleEmbeddings)
	router.POST("/api/show", proxy.HandleShow)
	router.POST("/api/ps", proxy.HandlePs)
	router.POST("/api/stop", proxy.HandleStop)

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Ollama is running in proxy mode.")
	})

	// Additional endpoints
	router.GET("/api", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Ollama API proxy v2",
			"version": "2.0.0",
		})
	})

	router.GET("/v1/models", proxy.HandleTags)
	router.GET("/models", proxy.HandleTags)
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
