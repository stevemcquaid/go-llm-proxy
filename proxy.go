package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

// BackendType represents the type of backend API
type BackendType string

const (
	BackendAnthropic BackendType = "anthropic"
	BackendOpenAI    BackendType = "openai"
)

// ProxyServer handles the translation between Ollama and backend APIs (Anthropic/OpenAI)
type ProxyServer struct {
	anthropicAPIKey  string
	anthropicBaseURL string
	openaiAPIKey     string
	openaiClient     *openai.Client
	httpClient       *http.Client
}

// NewProxyServer creates a new proxy server instance
func NewProxyServer() *ProxyServer {
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	var openaiClient *openai.Client
	if openaiAPIKey != "" {
		openaiClient = openai.NewClient(openaiAPIKey)
	}

	return &ProxyServer{
		anthropicAPIKey:  os.Getenv("ANTHROPIC_API_KEY"),
		anthropicBaseURL: "https://api.anthropic.com/v1",
		openaiAPIKey:     openaiAPIKey,
		openaiClient:     openaiClient,
		httpClient:       &http.Client{},
	}
}

// Ollama API Structures
type OllamaGenerateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type OllamaGenerateResponse struct {
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	Context            []int  `json:"context,omitempty"`
	TotalDuration      int64  `json:"total_duration,omitempty"`
	LoadDuration       int64  `json:"load_duration,omitempty"`
	PromptEvalCount    int    `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64  `json:"prompt_eval_duration,omitempty"`
	EvalCount          int    `json:"eval_count,omitempty"`
	EvalDuration       int64  `json:"eval_duration,omitempty"`
}

type OllamaChatRequest struct {
	Model    string                 `json:"model"`
	Messages []OllamaMessage        `json:"messages"`
	Stream   bool                   `json:"stream"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OllamaChatResponse struct {
	Model              string        `json:"model"`
	CreatedAt          string        `json:"created_at"`
	Message            OllamaMessage `json:"message"`
	Done               bool          `json:"done"`
	TotalDuration      int64         `json:"total_duration,omitempty"`
	LoadDuration       int64         `json:"load_duration,omitempty"`
	PromptEvalCount    int           `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64         `json:"prompt_eval_duration,omitempty"`
	EvalCount          int           `json:"eval_count,omitempty"`
	EvalDuration       int64         `json:"eval_duration,omitempty"`
}

type OllamaModel struct {
	Name       string `json:"name"`
	Model      string `json:"model"`
	ModifiedAt string `json:"modified_at"`
	Size       int64  `json:"size"`
	Digest     string `json:"digest"`
}

type OllamaTagsResponse struct {
	Models []OllamaModel `json:"models"`
}

// Anthropic API Structures
type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []AnthropicMessage `json:"messages"`
	Stream    bool               `json:"stream,omitempty"`
}

type AnthropicResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence,omitempty"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// OpenAI API Structures
type OpenAIRequest struct {
	Model     string          `json:"model"`
	Messages  []OpenAIMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens,omitempty"`
	Stream    bool            `json:"stream,omitempty"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// ModelInfo contains information about a model and its backend
type ModelInfo struct {
	Backend BackendType
	Model   string
}

// getModelInfo determines which backend to use and maps the model name
func (p *ProxyServer) getModelInfo(ollamaModel string) ModelInfo {
	// Anthropic model mappings
	anthropicModels := map[string]string{
		"claude":            "claude-3-5-sonnet-20241022",
		"claude-3":          "claude-3-5-sonnet-20241022",
		"claude-3-sonnet":   "claude-3-5-sonnet-20241022",
		"claude-3-haiku":    "claude-3-5-haiku-20241022",
		"claude-3-opus":     "claude-3-5-opus-20241022",
		"claude-3.5-sonnet": "claude-3-5-sonnet-20241022",
		"claude-3.5-haiku":  "claude-3-5-haiku-20241022",
		"claude-3.5-opus":   "claude-3-5-opus-20241022",
	}

	// OpenAI model mappings
	openaiModels := map[string]string{
		"gpt-4":             "gpt-4",
		"gpt-4-turbo":       "gpt-4-turbo-preview",
		"gpt-4o":            "gpt-4o",
		"gpt-4o-mini":       "gpt-4o-mini",
		"gpt-3.5-turbo":     "gpt-3.5-turbo",
		"gpt-3.5-turbo-16k": "gpt-3.5-turbo-16k",
		"o1-preview":        "o1-preview",
		"o1-mini":           "o1-mini",
	}

	// Check if it's an Anthropic model
	if mapped, exists := anthropicModels[ollamaModel]; exists {
		return ModelInfo{Backend: BackendAnthropic, Model: mapped}
	}

	// Check if it's an OpenAI model
	if mapped, exists := openaiModels[ollamaModel]; exists {
		return ModelInfo{Backend: BackendOpenAI, Model: mapped}
	}

	// Check if the model name starts with known prefixes
	if strings.HasPrefix(ollamaModel, "claude") {
		return ModelInfo{Backend: BackendAnthropic, Model: "claude-3-5-sonnet-20241022"}
	}
	if strings.HasPrefix(ollamaModel, "gpt") || strings.HasPrefix(ollamaModel, "o1") {
		return ModelInfo{Backend: BackendOpenAI, Model: "gpt-4"}
	}

	// Default to Anthropic (fallback)
	return ModelInfo{Backend: BackendAnthropic, Model: "claude-3-5-sonnet-20241022"}
}

// HandleGenerate handles the /api/generate endpoint
func (p *ProxyServer) HandleGenerate(c *gin.Context) {
	var req OllamaGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Get model information and determine backend
	modelInfo := p.getModelInfo(req.Model)

	// Convert to appropriate backend format and make request
	var response OllamaGenerateResponse
	var err error

	switch modelInfo.Backend {
	case BackendAnthropic:
		response, err = p.handleAnthropicGenerate(req, modelInfo.Model)
	case BackendOpenAI:
		response, err = p.handleOpenAIGenerate(req, modelInfo.Model)
	default:
		c.JSON(500, gin.H{"error": "unsupported backend"})
		return
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}

// HandleChat handles the /api/chat endpoint
func (p *ProxyServer) HandleChat(c *gin.Context) {
	var req OllamaChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Get model information and determine backend
	modelInfo := p.getModelInfo(req.Model)

	// Convert to appropriate backend format and make request
	var response OllamaChatResponse
	var err error

	switch modelInfo.Backend {
	case BackendAnthropic:
		response, err = p.handleAnthropicChat(req, modelInfo.Model)
	case BackendOpenAI:
		response, err = p.handleOpenAIChat(req, modelInfo.Model)
	default:
		c.JSON(500, gin.H{"error": "unsupported backend"})
		return
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}

// HandleTags handles the /api/tags endpoint (list available models)
func (p *ProxyServer) HandleTags(c *gin.Context) {
	var models []OllamaModel

	// Add Anthropic models if API key is available
	if p.anthropicAPIKey != "" {
		anthropicModels := []OllamaModel{
			{
				Name:       "claude-3.5-sonnet",
				Model:      "claude-3.5-sonnet",
				ModifiedAt: "2025-10-20T17:35:00.195Z",
				Size:       1000000000,
				Digest:     "sha256:claude35sonnet",
			},
			{
				Name:       "claude-3.5-haiku",
				Model:      "claude-3.5-haiku",
				ModifiedAt: "2025-10-20T17:35:00.195Z",
				Size:       1000000000,
				Digest:     "sha256:claude35haiku",
			},
			{
				Name:       "claude-3.5-opus",
				Model:      "claude-3.5-opus",
				ModifiedAt: "2025-10-20T17:35:00.195Z",
				Size:       1000000000,
				Digest:     "sha256:claude35opus",
			},
		}
		models = append(models, anthropicModels...)
	}

	// Add OpenAI models if API key is available
	if p.openaiAPIKey != "" {
		openaiModels := []OllamaModel{
			{
				Name:       "gpt-4",
				Model:      "gpt-4",
				ModifiedAt: "2025-10-20T17:35:00.195Z",
				Size:       1000000000,
				Digest:     "sha256:gpt4",
			},
			{
				Name:       "gpt-4o",
				Model:      "gpt-4o",
				ModifiedAt: "2025-10-20T17:35:00.195Z",
				Size:       1000000000,
				Digest:     "sha256:gpt4o",
			},
			{
				Name:       "gpt-4o-mini",
				Model:      "gpt-4o-mini",
				ModifiedAt: "2025-10-20T17:35:00.195Z",
				Size:       1000000000,
				Digest:     "sha256:gpt4omini",
			},
			{
				Name:       "gpt-3.5-turbo",
				Model:      "gpt-3.5-turbo",
				ModifiedAt: "2025-10-20T17:35:00.195Z",
				Size:       1000000000,
				Digest:     "sha256:gpt35turbo",
			},
		}
		models = append(models, openaiModels...)
	}

	response := OllamaTagsResponse{Models: models}
	c.JSON(200, response)
}

// HandleVersion handles the /api/version endpoint
func (p *ProxyServer) HandleVersion(c *gin.Context) {
	c.JSON(200, gin.H{
		"version": "1.0.0",
		"proxy":   "go-llm-proxy",
		"backend": "anthropic",
	})
}

// HandlePull handles the /api/pull endpoint (not applicable for Anthropic)
func (p *ProxyServer) HandlePull(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}

// HandlePush handles the /api/push endpoint (not applicable for Anthropic)
func (p *ProxyServer) HandlePush(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}

// HandleDelete handles the /api/delete endpoint (not applicable for Anthropic)
func (p *ProxyServer) HandleDelete(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}

// HandleCreate handles the /api/create endpoint (not applicable for Anthropic)
func (p *ProxyServer) HandleCreate(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}

// HandleCopy handles the /api/copy endpoint (not applicable for Anthropic)
func (p *ProxyServer) HandleCopy(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}

// HandleEmbeddings handles the /api/embeddings endpoint (not supported by Anthropic)
func (p *ProxyServer) HandleEmbeddings(c *gin.Context) {
	c.JSON(501, gin.H{"error": "embeddings not supported by Anthropic API"})
}

// HandleShow handles the /api/show endpoint
func (p *ProxyServer) HandleShow(c *gin.Context) {
	var req struct {
		Model string `json:"model"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Return model information in Ollama format
	model := OllamaModel{
		Name:       req.Model,
		Model:      req.Model,
		ModifiedAt: "2025-10-20T17:35:00.195Z",
		Size:       1000000000,
		Digest:     "sha256:" + req.Model,
	}

	c.JSON(200, model)
}

// HandlePs handles the /api/ps endpoint (not applicable for Anthropic)
func (p *ProxyServer) HandlePs(c *gin.Context) {
	c.JSON(200, gin.H{"processes": []interface{}{}})
}

// HandleStop handles the /api/stop endpoint (not applicable for Anthropic)
func (p *ProxyServer) HandleStop(c *gin.Context) {
	c.JSON(200, gin.H{"status": "success"})
}

// handleAnthropicGenerate handles generate requests for Anthropic
func (p *ProxyServer) handleAnthropicGenerate(req OllamaGenerateRequest, model string) (OllamaGenerateResponse, error) {
	anthropicReq := AnthropicRequest{
		Model:     model,
		MaxTokens: 4096,
		Messages: []AnthropicMessage{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
		Stream: req.Stream,
	}

	anthropicResp, err := p.makeAnthropicRequest(anthropicReq)
	if err != nil {
		return OllamaGenerateResponse{}, err
	}

	return OllamaGenerateResponse{
		Model:     req.Model,
		CreatedAt: anthropicResp.ID,
		Response:  anthropicResp.Content[0].Text,
		Done:      true,
	}, nil
}

// handleOpenAIGenerate handles generate requests for OpenAI
func (p *ProxyServer) handleOpenAIGenerate(req OllamaGenerateRequest, model string) (OllamaGenerateResponse, error) {
	openaiReq := OpenAIRequest{
		Model:     model,
		MaxTokens: 4096,
		Messages: []OpenAIMessage{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
		Stream: req.Stream,
	}

	openaiResp, err := p.makeOpenAIRequest(openaiReq)
	if err != nil {
		return OllamaGenerateResponse{}, err
	}

	return OllamaGenerateResponse{
		Model:     req.Model,
		CreatedAt: fmt.Sprintf("%d", openaiResp.Created),
		Response:  openaiResp.Choices[0].Message.Content,
		Done:      true,
	}, nil
}

// handleAnthropicChat handles chat requests for Anthropic
func (p *ProxyServer) handleAnthropicChat(req OllamaChatRequest, model string) (OllamaChatResponse, error) {
	var anthropicMessages []AnthropicMessage
	for _, msg := range req.Messages {
		anthropicMessages = append(anthropicMessages, AnthropicMessage(msg))
	}

	anthropicReq := AnthropicRequest{
		Model:     model,
		MaxTokens: 4096,
		Messages:  anthropicMessages,
		Stream:    req.Stream,
	}

	anthropicResp, err := p.makeAnthropicRequest(anthropicReq)
	if err != nil {
		return OllamaChatResponse{}, err
	}

	return OllamaChatResponse{
		Model:     req.Model,
		CreatedAt: anthropicResp.ID,
		Message: OllamaMessage{
			Role:    "assistant",
			Content: anthropicResp.Content[0].Text,
		},
		Done: true,
	}, nil
}

// handleOpenAIChat handles chat requests for OpenAI
func (p *ProxyServer) handleOpenAIChat(req OllamaChatRequest, model string) (OllamaChatResponse, error) {
	var openaiMessages []OpenAIMessage
	for _, msg := range req.Messages {
		openaiMessages = append(openaiMessages, OpenAIMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	openaiReq := OpenAIRequest{
		Model:     model,
		MaxTokens: 4096,
		Messages:  openaiMessages,
		Stream:    req.Stream,
	}

	openaiResp, err := p.makeOpenAIRequest(openaiReq)
	if err != nil {
		return OllamaChatResponse{}, err
	}

	return OllamaChatResponse{
		Model:     req.Model,
		CreatedAt: fmt.Sprintf("%d", openaiResp.Created),
		Message: OllamaMessage{
			Role:    "assistant",
			Content: openaiResp.Choices[0].Message.Content,
		},
		Done: true,
	}, nil
}

// makeAnthropicRequest makes a request to the Anthropic API
func (p *ProxyServer) makeAnthropicRequest(req AnthropicRequest) (*AnthropicResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", p.anthropicBaseURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.anthropicAPIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("anthropic API error: %s", string(body))
	}

	var anthropicResp AnthropicResponse
	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		return nil, err
	}

	return &anthropicResp, nil
}

// makeOpenAIRequest makes a request to the OpenAI API
func (p *ProxyServer) makeOpenAIRequest(req OpenAIRequest) (*OpenAIResponse, error) {
	// Use the OpenAI SDK
	openaiReq := openai.ChatCompletionRequest{
		Model:     req.Model,
		Messages:  convertToOpenAIMessages(req.Messages),
		MaxTokens: req.MaxTokens,
		Stream:    req.Stream,
	}

	resp, err := p.openaiClient.CreateChatCompletion(context.Background(), openaiReq)
	if err != nil {
		return nil, err
	}

	// Convert to our response format
	openaiResp := &OpenAIResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: make([]struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		}, len(resp.Choices)),
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}

	for i, choice := range resp.Choices {
		openaiResp.Choices[i].Index = choice.Index
		openaiResp.Choices[i].Message.Role = choice.Message.Role
		openaiResp.Choices[i].Message.Content = choice.Message.Content
		openaiResp.Choices[i].FinishReason = string(choice.FinishReason)
	}

	return openaiResp, nil
}

// convertToOpenAIMessages converts our message format to OpenAI's format
func convertToOpenAIMessages(messages []OpenAIMessage) []openai.ChatCompletionMessage {
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return openaiMessages
}
