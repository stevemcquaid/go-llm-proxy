package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

// StreamingHandler handles streaming responses
type StreamingHandler struct {
	backendManager *BackendManager
	modelRegistry  *ModelRegistry
}

// NewStreamingHandler creates a new streaming handler
func NewStreamingHandler(backendManager *BackendManager, modelRegistry *ModelRegistry) *StreamingHandler {
	return &StreamingHandler{
		backendManager: backendManager,
		modelRegistry:  modelRegistry,
	}
}

// HandleStreamingChat handles streaming chat requests
func (sh *StreamingHandler) HandleStreamingChat(c *gin.Context, req OllamaChatRequest) {
	// Set headers for streaming
	c.Header("Content-Type", "application/x-ndjson")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Get model configuration
	modelConfig, exists := sh.modelRegistry.GetModel(req.Model)
	if !exists {
		c.JSON(500, gin.H{"error": "model not found"})
		return
	}

	// Create non-streaming request for backend
	chatReq := ConvertOllamaToChatRequest(req, modelConfig.MaxTokens)
	chatReq.Model = modelConfig.BackendModel

	// Get response from backend
	ctx := context.Background()
	resp, err := sh.backendManager.ProcessRequest(ctx, modelConfig, chatReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	chatResp, ok := resp.(*ChatResponse)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid response type"})
		return
	}

	// Convert to Ollama format and stream
	ollamaResp := ConvertChatToOllamaResponse(chatResp, req.Model)
	sh.streamResponse(c, ollamaResp)
}

// HandleStreamingGenerate handles streaming generate requests
func (sh *StreamingHandler) HandleStreamingGenerate(c *gin.Context, req OllamaGenerateRequest) {
	// Set headers for streaming
	c.Header("Content-Type", "application/x-ndjson")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// Get model configuration
	modelConfig, exists := sh.modelRegistry.GetModel(req.Model)
	if !exists {
		c.JSON(500, gin.H{"error": "model not found"})
		return
	}

	// Create non-streaming request for backend
	generateReq := ConvertOllamaToGenerateRequest(req, modelConfig.MaxTokens)
	generateReq.Model = modelConfig.BackendModel

	// Get response from backend
	ctx := context.Background()
	resp, err := sh.backendManager.ProcessRequest(ctx, modelConfig, generateReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	generateResp, ok := resp.(*GenerateResponse)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid response type"})
		return
	}

	// Convert to Ollama format and stream
	ollamaResp := ConvertGenerateToOllamaResponse(generateResp, req.Model)
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
	case OllamaChatResponse:
		content = resp.Message.Content
		model = resp.Model
		createdAt = resp.CreatedAt
	case OllamaGenerateResponse:
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
		case OllamaChatResponse:
			streamResp = OllamaChatResponse{
				Model:     model,
				CreatedAt: createdAt,
				Message: OllamaMessage{
					Role:    "assistant",
					Content: chunk,
				},
				Done:    done,
				Context: []int{},
			}
		case OllamaGenerateResponse:
			streamResp = OllamaGenerateResponse{
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
