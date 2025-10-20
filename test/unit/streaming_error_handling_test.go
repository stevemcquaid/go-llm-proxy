package llmproxy_unit_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-llm-proxy/internal/backend"
	"go-llm-proxy/internal/models"
	"go-llm-proxy/internal/streaming"
	"go-llm-proxy/internal/types"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStreamingErrorHandling tests that streaming errors are returned in proper streaming format
func TestStreamingErrorHandling(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create backend manager with mock backends
	backendManager := backend.NewBackendManager()
	mockBackend := &MockBackend{name: "test-backend", available: true}
	backendManager.RegisterBackend(types.BackendOpenAI, mockBackend)

	// Create model registry with available backends
	modelRegistry := models.NewModelRegistryWithBackends(backendManager)

	// Create streaming handler
	streamingHandler := streaming.NewStreamingHandler(backendManager, modelRegistry)

	t.Run("ModelNotFoundError", func(t *testing.T) {
		// Test that model not found errors are returned in streaming format
		req := types.OllamaChatRequest{
			Model: "non-existent-model",
			Messages: []types.OllamaMessage{
				{Role: "user", Content: "Hello"},
			},
			Stream: true,
		}

		// Create test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Call the streaming handler
		streamingHandler.HandleStreamingChat(c, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/x-ndjson", w.Header().Get("Content-Type"))
		assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
		assert.Equal(t, "keep-alive", w.Header().Get("Connection"))

		// Parse the streaming response
		lines := bytes.Split(w.Body.Bytes(), []byte("\n"))
		require.Greater(t, len(lines), 0, "Should have at least one line in streaming response")

		// Concatenate all content chunks to get the full error message
		var fullContent string
		var lastResponse types.OllamaChatResponse
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			var response types.OllamaChatResponse
			err := json.Unmarshal(line, &response)
			if err != nil {
				continue
			}
			fullContent += response.Message.Content
			lastResponse = response
		}

		assert.Equal(t, "non-existent-model", lastResponse.Model)
		assert.Contains(t, fullContent, "Error: model not found")
		assert.True(t, lastResponse.Done, "Error response should be marked as done")
	})

	t.Run("BackendProcessingError", func(t *testing.T) {
		// Test that backend processing errors are returned in streaming format
		// First, add a model to the registry
		modelConfig := types.ModelConfig{
			Name:         "test-model",
			DisplayName:  "Test Model",
			Backend:      types.BackendOpenAI,
			BackendModel: "test-model",
			Family:       "test",
			Description:  "Test model",
			MaxTokens:    1000,
			Enabled:      true,
		}
		modelRegistry.AddModel(modelConfig)

		// Create a mock backend that returns an error
		errorBackend := &MockErrorBackend{name: "error-backend", available: true}
		backendManager.RegisterBackend(types.BackendOpenAI, errorBackend)

		req := types.OllamaChatRequest{
			Model: "test-model",
			Messages: []types.OllamaMessage{
				{Role: "user", Content: "Hello"},
			},
			Stream: true,
		}

		// Create test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Call the streaming handler
		streamingHandler.HandleStreamingChat(c, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/x-ndjson", w.Header().Get("Content-Type"))

		// Parse the streaming response
		lines := bytes.Split(w.Body.Bytes(), []byte("\n"))
		require.Greater(t, len(lines), 0, "Should have at least one line in streaming response")

		// Concatenate all content chunks to get the full error message
		var fullContent string
		var lastResponse types.OllamaChatResponse
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			var response types.OllamaChatResponse
			err := json.Unmarshal(line, &response)
			if err != nil {
				continue
			}
			fullContent += response.Message.Content
			lastResponse = response
		}

		assert.Equal(t, "test-model", lastResponse.Model)
		assert.Contains(t, fullContent, "Error: backend processing failed")
		assert.True(t, lastResponse.Done, "Error response should be marked as done")
	})

	t.Run("InvalidResponseTypeError", func(t *testing.T) {
		// Test that invalid response type errors are returned in streaming format
		// This would require a more complex setup to trigger the type assertion failure
		// For now, we'll test the error handling logic indirectly
		req := types.OllamaChatRequest{
			Model: "non-existent-model",
			Messages: []types.OllamaMessage{
				{Role: "user", Content: "Hello"},
			},
			Stream: true,
		}

		// Create test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Call the streaming handler
		streamingHandler.HandleStreamingChat(c, req)

		// Verify that we get a streaming response, not a regular JSON error
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/x-ndjson", w.Header().Get("Content-Type"))
		assert.NotContains(t, w.Body.String(), `{"error":`)
	})

	t.Run("GenerateRequestErrorHandling", func(t *testing.T) {
		// Test that generate request errors are also handled properly
		req := types.OllamaGenerateRequest{
			Model:  "non-existent-model",
			Prompt: "Hello",
			Stream: true,
		}

		// Create test context
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Call the streaming handler
		streamingHandler.HandleStreamingGenerate(c, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/x-ndjson", w.Header().Get("Content-Type"))

		// Parse the streaming response
		lines := bytes.Split(w.Body.Bytes(), []byte("\n"))
		require.Greater(t, len(lines), 0, "Should have at least one line in streaming response")

		// Concatenate all content chunks to get the full error message
		var fullContent string
		var lastResponse types.OllamaGenerateResponse
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			var response types.OllamaGenerateResponse
			err := json.Unmarshal(line, &response)
			if err != nil {
				continue
			}
			fullContent += response.Response
			lastResponse = response
		}

		assert.Equal(t, "non-existent-model", lastResponse.Model)
		assert.Contains(t, fullContent, "Error: model not found")
		assert.True(t, lastResponse.Done, "Error response should be marked as done")
	})
}

// MockErrorBackend is a mock backend that always returns an error
type MockErrorBackend struct {
	name      string
	available bool
}

func (m *MockErrorBackend) Generate(ctx context.Context, req types.GenerateRequest) (*types.GenerateResponse, error) {
	return nil, fmt.Errorf("backend processing failed")
}

func (m *MockErrorBackend) Chat(ctx context.Context, req types.ChatRequest) (*types.ChatResponse, error) {
	return nil, fmt.Errorf("backend processing failed")
}

func (m *MockErrorBackend) IsAvailable() bool {
	return m.available
}

func (m *MockErrorBackend) GetName() string {
	return m.name
}

// TestStreamingErrorFormat tests the specific error format that was causing issues
func TestStreamingErrorFormat(t *testing.T) {
	// This test specifically verifies that errors are returned in streaming format
	// instead of regular JSON format, which was the root cause of the 500 error

	gin.SetMode(gin.TestMode)

	// Create a minimal setup
	backendManager := backend.NewBackendManager()
	modelRegistry := models.NewModelRegistryWithBackends(backendManager)
	streamingHandler := streaming.NewStreamingHandler(backendManager, modelRegistry)

	t.Run("ErrorResponseFormat", func(t *testing.T) {
		req := types.OllamaChatRequest{
			Model: "invalid-model",
			Messages: []types.OllamaMessage{
				{Role: "user", Content: "Test"},
			},
			Stream: true,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		streamingHandler.HandleStreamingChat(c, req)

		// Verify that the response is in streaming format, not regular JSON
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/x-ndjson", w.Header().Get("Content-Type"))

		// The response should be a valid streaming response, not a JSON error
		responseBody := w.Body.String()
		assert.NotContains(t, responseBody, `{"error":`)
		assert.Contains(t, responseBody, `"model":"invalid-model"`)
		assert.Contains(t, responseBody, `"done":true`)

		// Concatenate all content chunks to check for the full error message
		lines := bytes.Split(w.Body.Bytes(), []byte("\n"))
		var fullContent string
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			var response types.OllamaChatResponse
			err := json.Unmarshal(line, &response)
			if err != nil {
				continue
			}
			fullContent += response.Message.Content
		}
		assert.Contains(t, fullContent, "Error: model not found")
	})

	t.Run("StreamingHeaders", func(t *testing.T) {
		// Test that proper streaming headers are set
		req := types.OllamaChatRequest{
			Model: "invalid-model",
			Messages: []types.OllamaMessage{
				{Role: "user", Content: "Test"},
			},
			Stream: true,
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		streamingHandler.HandleStreamingChat(c, req)

		// Verify streaming headers
		assert.Equal(t, "application/x-ndjson", w.Header().Get("Content-Type"))
		assert.Equal(t, "no-cache", w.Header().Get("Cache-Control"))
		assert.Equal(t, "keep-alive", w.Header().Get("Connection"))
	})
}
