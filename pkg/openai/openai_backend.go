package openai

import (
	"context"
	"fmt"

	"go-llm-proxy/internal/types"

	"github.com/sashabaranov/go-openai"
)

// OpenAIBackend implements the BackendHandler interface for OpenAI
type OpenAIBackend struct {
	apiKey string
	client *openai.Client
}

// NewOpenAIBackend creates a new OpenAI backend
func NewOpenAIBackend(apiKey string) *OpenAIBackend {
	client := openai.NewClient(apiKey)
	return &OpenAIBackend{
		apiKey: apiKey,
		client: client,
	}
}

// Generate handles text generation requests
func (ob *OpenAIBackend) Generate(ctx context.Context, req types.GenerateRequest) (*types.GenerateResponse, error) {
	openaiReq := openai.ChatCompletionRequest{
		Model: req.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: req.Prompt,
			},
		},
	}

	// Only set MaxTokens for models that support it
	// Newer models like GPT-4o use MaxCompletionTokens instead
	if req.MaxTokens > 0 && !isNewerModel(req.Model) {
		openaiReq.MaxTokens = req.MaxTokens
	}

	resp, err := ob.client.CreateChatCompletion(ctx, openaiReq)
	if err != nil {
		return nil, err
	}

	return &types.GenerateResponse{
		Model:     req.Model,
		Content:   resp.Choices[0].Message.Content,
		CreatedAt: fmt.Sprintf("%d", resp.Created),
	}, nil
}

// Chat handles chat completion requests
func (ob *OpenAIBackend) Chat(ctx context.Context, req types.ChatRequest) (*types.ChatResponse, error) {
	var messages []openai.ChatCompletionMessage
	for _, msg := range req.Messages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	openaiReq := openai.ChatCompletionRequest{
		Model:    req.Model,
		Messages: messages,
	}

	// Only set MaxTokens for models that support it
	// Newer models like GPT-4o use MaxCompletionTokens instead
	if req.MaxTokens > 0 && !isNewerModel(req.Model) {
		openaiReq.MaxTokens = req.MaxTokens
	}

	resp, err := ob.client.CreateChatCompletion(ctx, openaiReq)
	if err != nil {
		return nil, err
	}

	return &types.ChatResponse{
		Model: req.Model,
		Message: types.ChatMessage{
			Role:    resp.Choices[0].Message.Role,
			Content: resp.Choices[0].Message.Content,
		},
		CreatedAt: fmt.Sprintf("%d", resp.Created),
	}, nil
}

// IsAvailable checks if the backend is available
func (ob *OpenAIBackend) IsAvailable() bool {
	return ob.apiKey != ""
}

// GetName returns the backend name
func (ob *OpenAIBackend) GetName() string {
	return "openai"
}

// isNewerModel checks if the model is a newer model that doesn't support MaxTokens
func isNewerModel(model string) bool {
	// Models that require MaxCompletionTokens instead of MaxTokens
	newerModels := []string{
		"gpt-4o",
		"gpt-4o-mini",
		"gpt-5",
		"gpt-4.1",
		"gpt-4.5",
	}

	for _, newerModel := range newerModels {
		if model == newerModel {
			return true
		}
	}
	return false
}
