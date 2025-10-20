package main

import (
	"testing"
)

func TestModelMapping(t *testing.T) {
	proxy := NewProxyServer()

	tests := []struct {
		input           string
		expectedBackend BackendType
		expectedModel   string
	}{
		{"claude", BackendAnthropic, "claude-3-5-sonnet-20241022"},
		{"claude-3", BackendAnthropic, "claude-3-5-sonnet-20241022"},
		{"claude-3-sonnet", BackendAnthropic, "claude-3-5-sonnet-20241022"},
		{"claude-3-haiku", BackendAnthropic, "claude-3-5-haiku-20241022"},
		{"claude-3-opus", BackendAnthropic, "claude-3-5-opus-20241022"},
		{"claude-3.5-sonnet", BackendAnthropic, "claude-3-5-sonnet-20241022"},
		{"claude-3.5-haiku", BackendAnthropic, "claude-3-5-haiku-20241022"},
		{"claude-3.5-opus", BackendAnthropic, "claude-3-5-opus-20241022"},
		{"gpt-4", BackendOpenAI, "gpt-4"},
		{"gpt-4o", BackendOpenAI, "gpt-4o"},
		{"gpt-3.5-turbo", BackendOpenAI, "gpt-3.5-turbo"},
		{"unknown-model", BackendAnthropic, "claude-3-5-sonnet-20241022"}, // Default fallback (when no API keys)
	}

	for _, test := range tests {
		result := proxy.getModelInfo(test.input)
		if result.Backend != test.expectedBackend {
			t.Errorf("getModelInfo(%s).Backend = %s, expected %s", test.input, result.Backend, test.expectedBackend)
		}
		if result.Model != test.expectedModel {
			t.Errorf("getModelInfo(%s).Model = %s, expected %s", test.input, result.Model, test.expectedModel)
		}
	}
}

func TestProxyServerCreation(t *testing.T) {
	proxy := NewProxyServer()

	if proxy == nil {
		t.Error("NewProxyServer() returned nil")
	}

	if proxy.anthropicBaseURL != "https://api.anthropic.com/v1" {
		t.Errorf("Expected base URL to be 'https://api.anthropic.com/v1', got '%s'", proxy.anthropicBaseURL)
	}

	if proxy.httpClient == nil {
		t.Error("HTTP client should not be nil")
	}
}
