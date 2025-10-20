#!/bin/bash

# Test script for go-llm-proxy

echo "Testing go-llm-proxy..."

# Check if server is running
if ! curl -s http://localhost:11434/health > /dev/null; then
    echo "Error: Proxy server is not running on localhost:11434"
    echo "Please start the server first with: ./llm-proxy or go run ."
    exit 1
fi

echo "✓ Server is running"

# Test version endpoint
echo "Testing version endpoint..."
VERSION_RESPONSE=$(curl -s http://localhost:11434/api/version)
if echo "$VERSION_RESPONSE" | grep -q "version"; then
    echo "✓ Version endpoint working"
    echo "  Response: $VERSION_RESPONSE"
else
    echo "✗ Version endpoint failed"
    exit 1
fi

# Test tags endpoint
echo "Testing tags endpoint..."
TAGS_RESPONSE=$(curl -s http://localhost:11434/api/tags)
if echo "$TAGS_RESPONSE" | grep -q "models"; then
    echo "✓ Tags endpoint working"
    echo "  Available models:"
    echo "$TAGS_RESPONSE" | jq -r '.models[].name' 2>/dev/null || echo "$TAGS_RESPONSE"
else
    echo "✗ Tags endpoint failed"
    exit 1
fi

# Test generate endpoint (if API key is set)
if [ -n "$ANTHROPIC_API_KEY" ] && [ "$ANTHROPIC_API_KEY" != "your_anthropic_api_key_here" ]; then
    echo "Testing generate endpoint..."
    GENERATE_RESPONSE=$(curl -s -X POST http://localhost:11434/api/generate \
        -H "Content-Type: application/json" \
        -d '{
            "model": "claude-3.5-sonnet",
            "prompt": "Hello, how are you?",
            "stream": false
        }')
    
    if echo "$GENERATE_RESPONSE" | grep -q "response"; then
        echo "✓ Generate endpoint working"
        echo "  Response: $(echo "$GENERATE_RESPONSE" | jq -r '.response' 2>/dev/null || echo 'Check response manually')"
    else
        echo "✗ Generate endpoint failed"
        echo "  Response: $GENERATE_RESPONSE"
    fi
else
    echo "⚠ Skipping generate test - ANTHROPIC_API_KEY not set or using placeholder"
    echo "  Set your API key in .env file to test the generate endpoint"
fi

echo ""
echo "Test completed!"
