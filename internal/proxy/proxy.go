package proxy

import (
	"context"
	"fmt"
	"log"
	"os"

	"go-llm-proxy/internal/backend"
	"go-llm-proxy/internal/config"
	"go-llm-proxy/internal/models"
	"go-llm-proxy/internal/streaming"
	"go-llm-proxy/internal/types"

	"github.com/gin-gonic/gin"
)

// ProxyServerV2 is the refactored proxy server
type ProxyServerV2 struct {
	Config           *config.Config
	ModelRegistry    *models.ModelRegistry
	BackendManager   *backend.BackendManager
	StreamingHandler *streaming.StreamingHandler
}

// NewProxyServerV2 creates a new refactored proxy server
func NewProxyServerV2() *ProxyServerV2 {
	// Load configuration
	cfg := config.LoadConfig()

	// Create backend factory and manager first
	backendFactory := backend.NewBackendFactory(cfg.AnthropicAPIKey, cfg.OpenAIAPIKey)
	backendManager := backendFactory.CreateBackends()

	// Create model registry with dynamic fetching
	// Try to load from config file first, fall back to environment variables
	configPath := os.Getenv("MODEL_CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml" // Default config file
	}
	modelRegistry, err := models.NewModelRegistryWithDynamicFetching(cfg, backendManager, configPath)
	if err != nil {
		// Fail fast if dynamic fetching fails - no fallback
		log.Fatalf("Failed to fetch models dynamically: %v\n", err)
	}

	// Create streaming handler
	streamingHandler := streaming.NewStreamingHandler(backendManager, modelRegistry)

	return &ProxyServerV2{
		Config:           cfg,
		ModelRegistry:    modelRegistry,
		BackendManager:   backendManager,
		StreamingHandler: streamingHandler,
	}
}

// HandleGenerate handles the /api/generate endpoint
func (p *ProxyServerV2) HandleGenerate(c *gin.Context) {
	var req types.OllamaGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if streaming is requested
	if req.Stream {
		c.JSON(400, gin.H{"error": "streaming not supported for generate endpoint"})
		return
	}

	// Get model configuration
	modelConfig, exists := p.ModelRegistry.GetModel(req.Model)
	if !exists {
		c.JSON(400, gin.H{"error": "model not found"})
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

	// Create request for backend
	generateReq := types.ConvertOllamaToGenerateRequest(req, maxTokensForRequest)
	generateReq.Model = modelConfig.BackendModel

	// Process request
	ctx := context.Background()
	resp, err := p.BackendManager.ProcessRequest(ctx, modelConfig, generateReq)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("Error processing generate request: %v\n", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Convert response to Ollama format
	generateResp, ok := resp.(*types.GenerateResponse)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid response type"})
		return
	}

	ollamaResp := types.ConvertGenerateToOllamaResponse(generateResp, req.Model)
	c.JSON(200, ollamaResp)
}

// HandleChat handles the /api/chat endpoint
func (p *ProxyServerV2) HandleChat(c *gin.Context) {
	var req types.OllamaChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if streaming is requested
	if req.Stream {
		p.StreamingHandler.HandleStreamingChat(c, req)
		return
	}

	// Get model configuration
	modelConfig, exists := p.ModelRegistry.GetModel(req.Model)
	if !exists {
		c.JSON(400, gin.H{"error": "model not found"})
		return
	}

	// Convert messages for validation
	var messages []types.ChatMessage
	for _, msg := range req.Messages {
		messages = append(messages, types.ChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Validate token limits before making the request
	if err := types.ValidateTokenLimits(modelConfig, messages); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Calculate appropriate max_tokens for this specific request
	maxTokensForRequest := types.CalculateMaxTokensForRequest(modelConfig, messages)

	// Create request for backend
	chatReq := types.ConvertOllamaToChatRequest(req, maxTokensForRequest)
	chatReq.Model = modelConfig.BackendModel

	// Process request
	ctx := context.Background()
	resp, err := p.BackendManager.ProcessRequest(ctx, modelConfig, chatReq)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("Error processing chat request: %v\n", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Convert response to Ollama format
	chatResp, ok := resp.(*types.ChatResponse)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid response type"})
		return
	}

	ollamaResp := types.ConvertChatToOllamaResponse(chatResp, req.Model)
	c.JSON(200, ollamaResp)
}

// HandleTags handles the /api/tags endpoint (list available models)
func (p *ProxyServerV2) HandleTags(c *gin.Context) {
	// Get all available models
	models := p.ModelRegistry.GetAllModels()

	// Convert to Ollama format
	var ollamaModels []types.OllamaModel
	for _, model := range models {
		ollamaModels = append(ollamaModels, model.ToOllamaModel())
	}

	response := types.OllamaTagsResponse{Models: ollamaModels}
	c.JSON(200, response)
}

// HandleVersion handles the /api/version endpoint
func (p *ProxyServerV2) HandleVersion(c *gin.Context) {
	availableBackends := p.BackendManager.GetAvailableBackends()
	backendNames := make([]string, len(availableBackends))
	for i, backend := range availableBackends {
		backendNames[i] = string(backend)
	}

	c.JSON(200, gin.H{
		"version":  "2.0.0",
		"proxy":    "go-llm-proxy-v2",
		"backends": backendNames,
	})
}

// HandleShow handles the /api/show endpoint
func (p *ProxyServerV2) HandleShow(c *gin.Context) {
	// Get model from URL parameter
	modelName := c.Param("model")
	if modelName == "" {
		c.JSON(400, gin.H{"error": "model parameter is required"})
		return
	}

	// Get model configuration
	modelConfig, exists := p.ModelRegistry.GetModel(modelName)
	if !exists {
		c.JSON(400, gin.H{"error": "model not found"})
		return
	}

	// Return model information in Ollama format
	model := modelConfig.ToOllamaModel()
	c.JSON(200, model)
}

// HandlePull handles the /api/pull endpoint (not applicable for cloud backends)
func (p *ProxyServerV2) HandlePull(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success", "message": "Models are managed by backends"})
}

// HandlePush handles the /api/push endpoint (not applicable for cloud backends)
func (p *ProxyServerV2) HandlePush(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success", "message": "Models are managed by backends"})
}

// HandleDelete handles the /api/delete endpoint (not applicable for cloud backends)
func (p *ProxyServerV2) HandleDelete(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success", "message": "Models are managed by backends"})
}

// HandleCreate handles the /api/create endpoint (not applicable for cloud backends)
func (p *ProxyServerV2) HandleCreate(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success", "message": "Models are managed by backends"})
}

// HandleCopy handles the /api/copy endpoint (not applicable for cloud backends)
func (p *ProxyServerV2) HandleCopy(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success", "message": "Models are managed by backends"})
}

// HandleEmbeddings handles the /api/embeddings endpoint (not implemented)
func (p *ProxyServerV2) HandleEmbeddings(c *gin.Context) {
	c.JSON(501, gin.H{"error": "embeddings not implemented"})
}

// HandlePs handles the /api/ps endpoint (not applicable for cloud backends)
func (p *ProxyServerV2) HandlePs(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success", "message": "No local processes"})
}

// HandleStop handles the /api/stop endpoint (not applicable for cloud backends)
func (p *ProxyServerV2) HandleStop(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success", "message": "No local processes to stop"})
}

// GetHealthStatus returns the health status of the proxy
func (p *ProxyServerV2) GetHealthStatus() gin.H {
	availableBackends := p.BackendManager.GetAvailableBackends()
	modelCount := len(p.ModelRegistry.GetAllModels())

	return gin.H{
		"status":             "healthy",
		"available_backends": len(availableBackends),
		"total_models":       modelCount,
		"backends":           availableBackends,
	}
}
