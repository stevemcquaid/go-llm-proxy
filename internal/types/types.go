package types

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// BackendType represents the type of backend API
type BackendType string

const (
	BackendAnthropic BackendType = "anthropic"
	BackendOpenAI    BackendType = "openai"
)

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
	Context            []int  `json:"context"`
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

// ToChatMessage converts an OllamaMessage to ChatMessage
func (om OllamaMessage) ToChatMessage() ChatMessage {
	return ChatMessage(om)
}

type OllamaChatResponse struct {
	Model              string        `json:"model"`
	CreatedAt          string        `json:"created_at"`
	Message            OllamaMessage `json:"message"`
	Done               bool          `json:"done"`
	Context            []int         `json:"context"`
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
type AnthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []AnthropicMessage `json:"messages"`
}

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string      `json:"model"`
	StopReason   string      `json:"stop_reason"`
	StopSequence interface{} `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// ConvertOllamaToGenerateRequest converts an Ollama generate request to our format
func ConvertOllamaToGenerateRequest(req OllamaGenerateRequest, maxTokens int) GenerateRequest {
	return GenerateRequest{
		Model:     req.Model,
		Prompt:    req.Prompt,
		MaxTokens: maxTokens,
	}
}

// ConvertOllamaToChatRequest converts an Ollama chat request to our format
func ConvertOllamaToChatRequest(req OllamaChatRequest, maxTokens int) ChatRequest {
	var messages []ChatMessage
	for _, msg := range req.Messages {
		messages = append(messages, msg.ToChatMessage())
	}

	return ChatRequest{
		Model:     req.Model,
		Messages:  messages,
		MaxTokens: maxTokens,
	}
}

// ConvertGenerateToOllamaResponse converts our generate response to Ollama format
func ConvertGenerateToOllamaResponse(resp *GenerateResponse, model string) OllamaGenerateResponse {
	return OllamaGenerateResponse{
		Model:     model,
		CreatedAt: resp.CreatedAt,
		Response:  resp.Content,
		Done:      true,
		Context:   []int{},
	}
}

// BackendHandler defines the interface that all backends must implement
type BackendHandler interface {
	// Generate handles text generation requests
	Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error)

	// Chat handles chat completion requests
	Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)

	// IsAvailable checks if the backend is available (has API key, etc.)
	IsAvailable() bool

	// GetName returns the backend name
	GetName() string
}

// GenerateRequest represents a text generation request
type GenerateRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens,omitempty"`
}

// GenerateResponse represents a text generation response
type GenerateResponse struct {
	Model     string `json:"model"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model     string        `json:"model"`
	Messages  []ChatMessage `json:"messages"`
	MaxTokens int           `json:"max_tokens,omitempty"`
}

// ChatMessage represents a single message in a chat
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	Model     string      `json:"model"`
	Message   ChatMessage `json:"message"`
	CreatedAt string      `json:"created_at"`
}

// ModelConfig represents configuration for a model
type ModelConfig struct {
	Name         string      `json:"name"`
	DisplayName  string      `json:"display_name"`
	Backend      BackendType `json:"backend"`
	BackendModel string      `json:"backend_model"`
	Family       string      `json:"family"`
	Description  string      `json:"description"`
	MaxTokens    int         `json:"max_tokens"`
	Enabled      bool        `json:"enabled"`
}

// ToOllamaModel converts a ModelConfig to OllamaModel format
func (m ModelConfig) ToOllamaModel() OllamaModel {
	return OllamaModel{
		Name:       m.Name,
		Model:      m.Name,
		ModifiedAt: time.Now().Format("2006-01-02T15:04:05.000Z"),
		Size:       1000000000, // 1GB placeholder
		Digest:     "sha256:" + m.Name,
	}
}

// ConvertChatToOllamaResponse converts our chat response to Ollama format
func ConvertChatToOllamaResponse(resp *ChatResponse, model string) OllamaChatResponse {
	return OllamaChatResponse{
		Model:     model,
		CreatedAt: resp.CreatedAt,
		Message: OllamaMessage{
			Role:    resp.Message.Role,
			Content: resp.Message.Content,
		},
		Done:    true,
		Context: []int{},
	}
}

// EstimateTokens provides a rough estimation of token count for text
// This is a simple approximation: ~4 characters per token for English text
func EstimateTokens(text string) int {
	if text == "" {
		return 0
	}
	// Simple approximation: ~4 characters per token
	// This is conservative and may overestimate, but that's safer for validation
	return (len(strings.TrimSpace(text)) + 3) / 4
}

// EstimateChatTokens estimates the total token count for a chat request
func EstimateChatTokens(messages []ChatMessage) int {
	total := 0
	for _, msg := range messages {
		// Add tokens for role and content
		total += EstimateTokens(msg.Role) + EstimateTokens(msg.Content)
		// Add some overhead for message formatting
		total += 4
	}
	// Add some overhead for the overall request structure
	return total + 10
}

// CalculateMaxTokensForRequest calculates the appropriate max_tokens for an API request
// based on the model's context limit and input token count
func CalculateMaxTokensForRequest(modelConfig ModelConfig, messages []ChatMessage) int {
	estimatedInputTokens := EstimateChatTokens(messages)

	// Reserve some buffer for the input tokens and calculate remaining for output
	// Use a conservative approach: total context - input tokens - buffer
	buffer := 100 // Small buffer for safety
	availableForOutput := modelConfig.MaxTokens - estimatedInputTokens - buffer

	// Ensure we don't go negative and have a reasonable minimum
	if availableForOutput < 100 {
		availableForOutput = 100 // Minimum 100 tokens for output
	}

	// Cap at a reasonable maximum to avoid very long responses
	maxOutputTokens := 4000
	if availableForOutput > maxOutputTokens {
		availableForOutput = maxOutputTokens
	}

	return availableForOutput
}

// ValidateTokenLimits checks if a request would exceed the model's token limits
func ValidateTokenLimits(modelConfig ModelConfig, messages []ChatMessage) error {
	estimatedTokens := EstimateChatTokens(messages)

	// For models with small context windows, be more conservative
	// Reserve at least 25% of context for output tokens
	var maxInputTokens int
	if modelConfig.MaxTokens <= 8192 {
		// For small context models, reserve 25% for output
		maxInputTokens = int(float64(modelConfig.MaxTokens) * 0.75)
	} else {
		// For larger context models, reserve 50% for output
		maxInputTokens = int(float64(modelConfig.MaxTokens) * 0.5)
	}

	if estimatedTokens > maxInputTokens {
		return fmt.Errorf("request too long: estimated %d tokens exceeds model limit of %d tokens (max input: %d tokens). Please reduce the length of your messages",
			estimatedTokens, modelConfig.MaxTokens, maxInputTokens)
	}

	return nil
}
