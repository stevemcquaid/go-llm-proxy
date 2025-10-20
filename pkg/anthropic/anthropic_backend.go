package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-llm-proxy/internal/types"
	"io"
	"net/http"
)

// AnthropicBackend implements the BackendHandler interface for Anthropic
type AnthropicBackend struct {
	apiKey string
	client *http.Client
}

// NewAnthropicBackend creates a new Anthropic backend
func NewAnthropicBackend(apiKey string) *AnthropicBackend {
	return &AnthropicBackend{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// Generate handles text generation requests
func (ab *AnthropicBackend) Generate(ctx context.Context, req types.GenerateRequest) (*types.GenerateResponse, error) {
	anthropicReq := AnthropicRequest{
		Model:     req.Model,
		MaxTokens: req.MaxTokens,
		Messages: []AnthropicMessage{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
	}

	resp, err := ab.makeRequest(ctx, anthropicReq)
	if err != nil {
		return nil, err
	}

	return &types.GenerateResponse{
		Model:     req.Model,
		Content:   resp.Content[0].Text,
		CreatedAt: resp.ID,
	}, nil
}

// Chat handles chat completion requests
func (ab *AnthropicBackend) Chat(ctx context.Context, req types.ChatRequest) (*types.ChatResponse, error) {
	var anthropicMessages []AnthropicMessage
	for _, msg := range req.Messages {
		anthropicMessages = append(anthropicMessages, AnthropicMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	anthropicReq := AnthropicRequest{
		Model:     req.Model,
		MaxTokens: req.MaxTokens,
		Messages:  anthropicMessages,
	}

	resp, err := ab.makeRequest(ctx, anthropicReq)
	if err != nil {
		return nil, err
	}

	return &types.ChatResponse{
		Model: req.Model,
		Message: types.ChatMessage{
			Role:    "assistant",
			Content: resp.Content[0].Text,
		},
		CreatedAt: resp.ID,
	}, nil
}

// IsAvailable checks if the backend is available
func (ab *AnthropicBackend) IsAvailable() bool {
	return ab.apiKey != ""
}

// GetName returns the backend name
func (ab *AnthropicBackend) GetName() string {
	return "anthropic"
}

// makeRequest makes a request to the Anthropic API
func (ab *AnthropicBackend) makeRequest(ctx context.Context, req AnthropicRequest) (*AnthropicResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", ab.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := ab.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("anthropic API error: %s", string(body))
	}

	var anthropicResp AnthropicResponse
	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		return nil, err
	}

	return &anthropicResp, nil
}

// AnthropicRequest represents a request to the Anthropic API
type AnthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []AnthropicMessage `json:"messages"`
}

// AnthropicMessage represents a message in the Anthropic API
type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AnthropicResponse represents a response from the Anthropic API
type AnthropicResponse struct {
	ID      string `json:"id"`
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}
