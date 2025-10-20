#!/bin/bash

# Setup script for go-llm-proxy

echo "Setting up go-llm-proxy..."

# Copy environment file
if [ ! -f .env ]; then
    cp env.example .env
    echo "Created .env file from env.example"
    echo "Please edit .env and add your Anthropic API key"
else
    echo ".env file already exists"
fi

# Install dependencies
echo "Installing dependencies..."
go mod tidy

# Build the project
echo "Building the project..."
go build -o llm-proxy .

if [ $? -eq 0 ]; then
    echo "Build successful! Binary created as 'llm-proxy'"
    echo ""
    echo "To run the proxy:"
    echo "  ./llm-proxy"
    echo ""
    echo "Or with go run:"
    echo "  go run ."
    echo ""
    echo "Don't forget to:"
    echo "1. Edit .env and add your Anthropic API key"
    echo "2. Configure your JetBrains IDE to use localhost:11434"
else
    echo "Build failed. Please check the error messages above."
    exit 1
fi
