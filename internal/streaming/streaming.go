package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-llm-proxy/internal/backend"
	"go-llm-proxy/internal/models"
	"go-llm-proxy/internal/types"

	"github.com/gin-gonic/gin"
)

// StreamingHandler handles streaming responses
type StreamingHandler struct {
	backendManager *backend.BackendManager
	modelRegistry  *models.ModelRegistry
}

// NewStreamingHandler creates a new streaming handler
func NewStreamingHandler(backendManager *backend.BackendManager, modelRegistry *models.ModelRegistry) *StreamingHandler {
	return &StreamingHandler{
		backendManager: backendManager,
		modelRegistry:  modelRegistry,
	}
}

// HandleStreamingChat handles streaming chat requests
func (sh *StreamingHandler) HandleStreamingChat(c *gin.Context, req types.OllamaChatRequest) {
	// Set headers for streaming
	c.Header("Content-Type", "application/x-ndjson")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Get model configuration
	modelConfig, exists := sh.modelRegistry.GetModel(req.Model)
	if !exists {
		// For streaming responses, we need to return an error in streaming format
		errorResp := types.OllamaChatResponse{
			Model:     req.Model,
			CreatedAt: fmt.Sprintf("%d", time.Now().Unix()),
			Message: types.OllamaMessage{
				Role:    "assistant",
				Content: "Error: model not found",
			},
			Done:    true,
			Context: []int{},
		}
		sh.streamResponse(c, errorResp)
		return
	}

	// Convert messages for validation
	var messages []types.ChatMessage
	for _, msg := range req.Messages {
		messages = append(messages, msg.ToChatMessage())
	}

	// Validate token limits before making the request
	if err := types.ValidateTokenLimits(modelConfig, messages); err != nil {
		// For streaming responses, we need to return an error in streaming format
		errorResp := types.OllamaChatResponse{
			Model:     req.Model,
			CreatedAt: fmt.Sprintf("%d", time.Now().Unix()),
			Message: types.OllamaMessage{
				Role:    "assistant",
				Content: fmt.Sprintf("Error: %s", err.Error()),
			},
			Done:    true,
			Context: []int{},
		}
		sh.streamResponse(c, errorResp)
		return
	}

	// Calculate appropriate max_tokens for this specific request
	maxTokensForRequest := types.CalculateMaxTokensForRequest(modelConfig, messages)

	// Create non-streaming request for backend
	chatReq := types.ConvertOllamaToChatRequest(req, maxTokensForRequest)
	chatReq.Model = modelConfig.BackendModel

	// Get response from backend
	ctx := context.Background()
	resp, err := sh.backendManager.ProcessRequest(ctx, modelConfig, chatReq)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("Error processing streaming chat request: %v\n", err)
		// For streaming responses, we need to return an error in streaming format
		errorResp := types.OllamaChatResponse{
			Model:     req.Model,
			CreatedAt: fmt.Sprintf("%d", time.Now().Unix()),
			Message: types.OllamaMessage{
				Role:    "assistant",
				Content: fmt.Sprintf("Error: %s", err.Error()),
			},
			Done:    true,
			Context: []int{},
		}
		sh.streamResponse(c, errorResp)
		return
	}

	chatResp, ok := resp.(*types.ChatResponse)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid response type"})
		return
	}

	// Convert to Ollama format and stream
	ollamaResp := types.ConvertChatToOllamaResponse(chatResp, req.Model)
	sh.streamResponse(c, ollamaResp)
}

// HandleStreamingGenerate handles streaming generate requests
func (sh *StreamingHandler) HandleStreamingGenerate(c *gin.Context, req types.OllamaGenerateRequest) {
	// Set headers for streaming
	c.Header("Content-Type", "application/x-ndjson")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Get model configuration
	modelConfig, exists := sh.modelRegistry.GetModel(req.Model)
	if !exists {
		// For streaming responses, we need to return an error in streaming format
		errorResp := types.OllamaGenerateResponse{
			Model:     req.Model,
			CreatedAt: fmt.Sprintf("%d", time.Now().Unix()),
			Response:  "Error: model not found",
			Done:      true,
			Context:   []int{},
		}
		sh.streamResponse(c, errorResp)
		return
	}

	// For generate requests, we need to estimate tokens from the prompt
	// and calculate appropriate max_tokens
	var messages []types.ChatMessage
	messages = append(messages, types.ChatMessage{
		Role:    "user",
		Content: req.Prompt,
	})
	maxTokensForRequest := types.CalculateMaxTokensForRequest(modelConfig, messages)

	// Create non-streaming request for backend
	generateReq := types.ConvertOllamaToGenerateRequest(req, maxTokensForRequest)
	generateReq.Model = modelConfig.BackendModel

	// Get response from backend
	ctx := context.Background()
	resp, err := sh.backendManager.ProcessRequest(ctx, modelConfig, generateReq)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("Error processing streaming generate request: %v\n", err)
		// For streaming responses, we need to return an error in streaming format
		errorResp := types.OllamaGenerateResponse{
			Model:     req.Model,
			CreatedAt: fmt.Sprintf("%d", time.Now().Unix()),
			Response:  fmt.Sprintf("Error: %s", err.Error()),
			Done:      true,
			Context:   []int{},
		}
		sh.streamResponse(c, errorResp)
		return
	}

	generateResp, ok := resp.(*types.GenerateResponse)
	if !ok {
		// For streaming responses, we need to return an error in streaming format
		errorResp := types.OllamaGenerateResponse{
			Model:     req.Model,
			CreatedAt: fmt.Sprintf("%d", time.Now().Unix()),
			Response:  "Error: invalid response type",
			Done:      true,
			Context:   []int{},
		}
		sh.streamResponse(c, errorResp)
		return
	}

	// Convert to Ollama format and stream
	ollamaResp := types.ConvertGenerateToOllamaResponse(generateResp, req.Model)
	sh.streamResponse(c, ollamaResp)
}

// streamResponse streams a response by breaking it into chunks
func (sh *StreamingHandler) streamResponse(c *gin.Context, response interface{}) {
	// For now, we'll simulate streaming by breaking the response into chunks
	// In a real implementation, you might want to use actual streaming from the backend

	content, model, createdAt, err := sh.extractResponseData(response)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Check if this is an error response (contains "Error:")
	isError := strings.Contains(content, "Error:")

	if isError {
		sh.streamErrorResponse(c, response, content, model, createdAt)
	} else {
		sh.streamNormalResponse(c, response, content, model, createdAt)
	}
}

// extractResponseData extracts content, model, and createdAt from response
func (sh *StreamingHandler) extractResponseData(response interface{}) (string, string, string, error) {
	switch resp := response.(type) {
	case types.OllamaChatResponse:
		return resp.Message.Content, resp.Model, resp.CreatedAt, nil
	case types.OllamaGenerateResponse:
		return resp.Response, resp.Model, resp.CreatedAt, nil
	default:
		return "", "", "", fmt.Errorf("unsupported response type")
	}
}

// streamErrorResponse streams an error response as a single chunk
func (sh *StreamingHandler) streamErrorResponse(c *gin.Context, response interface{}, content, model, createdAt string) {
	streamResp := sh.createStreamResponse(response, content, model, createdAt, true)
	sh.writeResponse(c, streamResp)
}

// streamNormalResponse streams a normal response by breaking it into chunks
func (sh *StreamingHandler) streamNormalResponse(c *gin.Context, response interface{}, content, model, createdAt string) {
	chunkSize := 3 // Small chunks for demonstration
	for i := 0; i < len(content); i += chunkSize {
		end := i + chunkSize
		if end > len(content) {
			end = len(content)
		}

		chunk := content[i:end]
		done := end >= len(content)

		streamResp := sh.createStreamResponse(response, chunk, model, createdAt, done)
		sh.writeResponse(c, streamResp)

		// Small delay to simulate streaming
		time.Sleep(50 * time.Millisecond)
	}
}

// createStreamResponse creates a streaming response based on the original response type
func (sh *StreamingHandler) createStreamResponse(response interface{}, content, model, createdAt string, done bool) interface{} {
	switch response.(type) {
	case types.OllamaChatResponse:
		return types.OllamaChatResponse{
			Model:     model,
			CreatedAt: createdAt,
			Message: types.OllamaMessage{
				Role:    "assistant",
				Content: content,
			},
			Done:    done,
			Context: []int{},
		}
	case types.OllamaGenerateResponse:
		return types.OllamaGenerateResponse{
			Model:     model,
			CreatedAt: createdAt,
			Response:  content,
			Done:      done,
			Context:   []int{},
		}
	default:
		return nil
	}
}

// writeResponse writes a response to the client
func (sh *StreamingHandler) writeResponse(c *gin.Context, response interface{}) {
	jsonData, _ := json.Marshal(response)
	if _, err := c.Writer.Write(jsonData); err != nil {
		fmt.Printf("Warning: failed to write JSON data: %v\n", err)
	}
	if _, err := c.Writer.WriteString("\n"); err != nil {
		fmt.Printf("Warning: failed to write newline: %v\n", err)
	}
	c.Writer.Flush()
}
