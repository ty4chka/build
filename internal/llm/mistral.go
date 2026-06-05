package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// MistralClient implements the Client interface for Mistral AI
type MistralClient struct {
	apiKey string
	baseURL string
	httpClient *http.Client
}

// NewMistralClient creates a new Mistral client
func NewMistralClient(apiKey, baseURL string) *MistralClient {
	return &MistralClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 180 * time.Second,
		},
	}
}

// mistralRequest represents the API request body
type mistralRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// mistralResponse represents the API response
type mistralResponse struct {
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

// Complete implements the Client interface
func (c *MistralClient) Complete(ctx context.Context, messages []Message, opts *CompletionOpts) (*CompletionResult, error) {
	if opts == nil {
		opts = &CompletionOpts{
			Temperature: 0.7,
			MaxTokens:   1200,
		}
	}

	model := opts.Model
	if model == "" {
		model = "mistral-large-latest"
	}

	// Convert messages
	msgs := make([]message, len(messages))
	for i, m := range messages {
		msgs[i] = message{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	reqBody := mistralRequest{
		Model:       model,
		Messages:    msgs,
		Temperature: opts.Temperature,
		MaxTokens:   opts.MaxTokens,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result mistralResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	return &CompletionResult{
		Content:      result.Choices[0].Message.Content,
		InputTokens:  result.Usage.PromptTokens,
		OutputTokens: result.Usage.CompletionTokens,
		TotalTokens:  result.Usage.TotalTokens,
	}, nil
}

// Stream implements the Client interface (simplified — returns full response)
func (c *MistralClient) Stream(ctx context.Context, messages []Message, opts *CompletionOpts) (<-chan string, error) {
	// For now, just use Complete and return the full response
	// TODO: implement SSE streaming
	result, err := c.Complete(ctx, messages, opts)
	if err != nil {
		return nil, err
	}

	ch := make(chan string, 1)
	ch <- result.Content
	close(ch)
	return ch, nil
}

// Factory creates the appropriate LLM client based on provider
func Factory(provider, apiKey, baseURL string) Client {
	switch provider {
	case "mistral":
		return NewMistralClient(apiKey, baseURL)
	case "groq":
		// Groq uses OpenAI-compatible API
		return NewMistralClient(apiKey, baseURL) // Same format
	case "openrouter":
		return NewMistralClient(apiKey, baseURL) // Same format
	default:
		return NewMistralClient(apiKey, baseURL)
	}
}
