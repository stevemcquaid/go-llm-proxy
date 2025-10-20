package streaming

import (
	"context"
	"encoding/json"
	"fmt"
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

	// Create non-streaming request for backend
	chatReq := types.ConvertOllamaToChatRequest(req, modelConfig.MaxTokens)
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

	// Create non-streaming request for backend
	generateReq := types.ConvertOllamaToGenerateRequest(req, modelConfig.MaxTokens)
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

	var content string
	var model string
	var createdAt string

	// Extract content based on response type
	switch resp := response.(type) {
	case types.OllamaChatResponse:
		content = resp.Message.Content
		model = resp.Model
		createdAt = resp.CreatedAt
	case types.OllamaGenerateResponse:
		content = resp.Response
		model = resp.Model
		createdAt = resp.CreatedAt
	default:
		c.JSON(500, gin.H{"error": "unsupported response type"})
		return
	}

	// Break content into chunks
	chunkSize := 3 // Small chunks for demonstration
	for i := 0; i < len(content); i += chunkSize {
		end := i + chunkSize
		if end > len(content) {
			end = len(content)
		}

		chunk := content[i:end]
		done := end >= len(content)

		// Create streaming response based on type
		var streamResp interface{}
		switch response.(type) {
		case types.OllamaChatResponse:
			streamResp = types.OllamaChatResponse{
				Model:     model,
				CreatedAt: createdAt,
				Message: types.OllamaMessage{
					Role:    "assistant",
					Content: chunk,
				},
				Done:    done,
				Context: []int{},
			}
		case types.OllamaGenerateResponse:
			streamResp = types.OllamaGenerateResponse{
				Model:     model,
				CreatedAt: createdAt,
				Response:  chunk,
				Done:      done,
				Context:   []int{},
			}
		}

		// Write the chunk as JSON
		jsonData, _ := json.Marshal(streamResp)
		c.Writer.Write(jsonData)
		c.Writer.WriteString("\n")
		c.Writer.Flush()

		// Small delay to simulate streaming
		time.Sleep(50 * time.Millisecond)
	}
}
