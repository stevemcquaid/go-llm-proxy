package main

import (
	"context"

	"github.com/gin-gonic/gin"
)

// ProxyServerV2 is the refactored proxy server
type ProxyServerV2 struct {
	config           *Config
	modelRegistry    *ModelRegistry
	backendManager   *BackendManager
	streamingHandler *StreamingHandler
}

// NewProxyServerV2 creates a new refactored proxy server
func NewProxyServerV2() *ProxyServerV2 {
	// Load configuration
	config := LoadConfig()

	// Create model registry
	modelRegistry := NewModelRegistry()

	// Create backend factory and manager
	backendFactory := NewBackendFactory(config.AnthropicAPIKey, config.OpenAIAPIKey)
	backendManager := backendFactory.CreateBackends()

	// Create streaming handler
	streamingHandler := NewStreamingHandler(backendManager, modelRegistry)

	return &ProxyServerV2{
		config:           config,
		modelRegistry:    modelRegistry,
		backendManager:   backendManager,
		streamingHandler: streamingHandler,
	}
}

// HandleGenerate handles the /api/generate endpoint
func (p *ProxyServerV2) HandleGenerate(c *gin.Context) {
	var req OllamaGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if streaming is requested
	if req.Stream {
		p.streamingHandler.HandleStreamingGenerate(c, req)
		return
	}

	// Get model configuration
	modelConfig, exists := p.modelRegistry.GetModel(req.Model)
	if !exists {
		c.JSON(400, gin.H{"error": "model not found"})
		return
	}

	// Create request for backend
	generateReq := ConvertOllamaToGenerateRequest(req, modelConfig.MaxTokens)
	generateReq.Model = modelConfig.BackendModel

	// Process request
	ctx := context.Background()
	resp, err := p.backendManager.ProcessRequest(ctx, modelConfig, generateReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Convert response to Ollama format
	generateResp, ok := resp.(*GenerateResponse)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid response type"})
		return
	}

	ollamaResp := ConvertGenerateToOllamaResponse(generateResp, req.Model)
	c.JSON(200, ollamaResp)
}

// HandleChat handles the /api/chat endpoint
func (p *ProxyServerV2) HandleChat(c *gin.Context) {
	var req OllamaChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Check if streaming is requested
	if req.Stream {
		p.streamingHandler.HandleStreamingChat(c, req)
		return
	}

	// Get model configuration
	modelConfig, exists := p.modelRegistry.GetModel(req.Model)
	if !exists {
		c.JSON(400, gin.H{"error": "model not found"})
		return
	}

	// Create request for backend
	chatReq := ConvertOllamaToChatRequest(req, modelConfig.MaxTokens)
	chatReq.Model = modelConfig.BackendModel

	// Process request
	ctx := context.Background()
	resp, err := p.backendManager.ProcessRequest(ctx, modelConfig, chatReq)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Convert response to Ollama format
	chatResp, ok := resp.(*ChatResponse)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid response type"})
		return
	}

	ollamaResp := ConvertChatToOllamaResponse(chatResp, req.Model)
	c.JSON(200, ollamaResp)
}

// HandleTags handles the /api/tags endpoint (list available models)
func (p *ProxyServerV2) HandleTags(c *gin.Context) {
	// Get all available models
	models := p.modelRegistry.GetAllModels()

	// Convert to Ollama format
	var ollamaModels []OllamaModel
	for _, model := range models {
		ollamaModels = append(ollamaModels, model.ToOllamaModel())
	}

	response := OllamaTagsResponse{Models: ollamaModels}
	c.JSON(200, response)
}

// HandleVersion handles the /api/version endpoint
func (p *ProxyServerV2) HandleVersion(c *gin.Context) {
	availableBackends := p.backendManager.GetAvailableBackends()
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
	var req struct {
		Model string `json:"model"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Get model configuration
	modelConfig, exists := p.modelRegistry.GetModel(req.Model)
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
	availableBackends := p.backendManager.GetAvailableBackends()
	modelCount := len(p.modelRegistry.GetAllModels())

	return gin.H{
		"status":             "healthy",
		"available_backends": len(availableBackends),
		"total_models":       modelCount,
		"backends":           availableBackends,
	}
}
