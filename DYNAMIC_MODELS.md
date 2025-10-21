# Dynamic Model Fetching

This document describes the new dynamic model fetching feature that pulls model configurations from APIs on startup.

## Overview

Instead of using hardcoded model configurations, the proxy now dynamically fetches available models from:
- Anthropic API (https://api.anthropic.com/v1/models)
- OpenAI API (https://api.openai.com/v1/models)

**Important**: The proxy will fail to start if dynamic model fetching fails. There is no fallback to hardcoded models.

## Configuration

### Environment Variables

- `MODEL_CONFIG_PATH`: Path to YAML configuration file for model filtering (defaults to `config.yaml`)
- `ANTHROPIC_API_KEY`: Your Anthropic API key
- `OPENAI_API_KEY`: Your OpenAI API key

### YAML Configuration File

Create a YAML file to filter which models to include/exclude:

```yaml
model_filters:
  anthropic:
    enabled: true
    include_patterns:
      - "claude-*"
    exclude_patterns:
      - "claude-3-5-*"  # Exclude older 3.5 models
      - "claude-2-*"    # Exclude 2.x models
  
  openai:
    enabled: true
    include_patterns:
      - "gpt-*"
    exclude_patterns:
      - "gpt-3.5-turbo"  # Exclude specific model
      - "text-*"         # Exclude text completion models
      - "davinci-*"      # Exclude older models
```

### Pattern Matching

- `include_patterns`: Only models matching these patterns will be included
- `exclude_patterns`: Models matching these patterns will be excluded
- Patterns use filepath.Match syntax (supports `*` wildcards)
- If no include_patterns are specified, all models are included (after exclusions)

## Usage

1. Set your API keys as environment variables:
   ```bash
   export ANTHROPIC_API_KEY="your-anthropic-key"
   export OPENAI_API_KEY="your-openai-key"
   ```

2. Optionally, create a config file (defaults to `config.yaml` in the current directory):
   ```bash
   # Uses default config.yaml
   # OR specify a custom path:
   export MODEL_CONFIG_PATH="/path/to/custom-config.yaml"
   ```

3. Start the proxy:
   ```bash
   go run cmd/llm-proxy/main.go
   ```

## Fallback Behavior

If dynamic fetching fails (API errors, network issues, etc.), the proxy will:
1. Print a warning message to the console
2. Fall back to the hardcoded model configurations
3. Continue running normally

## Model Name Mapping

The system automatically maps API model names to cleaner proxy names:

### Anthropic
- `claude-3-5-sonnet-20241022` → `claude-3.5-sonnet`
- `claude-3-5-haiku-20241022` → `claude-3.5-haiku`

### OpenAI
- `gpt-4o` → `gpt-4o` (unchanged)
- `gpt-3.5-turbo` → `gpt-3.5-turbo` (unchanged)

## Error Handling

- API errors are logged as warnings and printed to console
- Individual backend failures don't prevent other backends from loading
- If all backends fail, the system falls back to hardcoded models
- Network timeouts are set to 30 seconds per API call

## Example Output

```
Starting LLM Proxy server v2 on port 11434
Loaded 8 models dynamically from APIs
Available backends: [anthropic openai]
Total models: 8
```

## Troubleshooting

1. **No models loaded**: Check that your API keys are set correctly
2. **Specific models missing**: Check your include/exclude patterns in the config file
3. **API errors**: Check your internet connection and API key validity
4. **Fallback to hardcoded**: Check the console output for error messages
