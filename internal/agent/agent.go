package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/oblachko/xli-bot/internal/config"
	"github.com/oblachko/xli-bot/internal/llm"
)

// AgentResult holds the result of agent execution
type AgentResult struct {
	Answer        string
	ThinkingNotes []string
	AgentLog      []string
	TokenUsage    llm.TokenUsage
	Elapsed       time.Duration
}

// Agent is the main AI agent
type Agent struct {
	llmClient llm.Client
	config    *config.Config
}

// NewAgent creates a new agent
func NewAgent(client llm.Client, cfg *config.Config) *Agent {
	return &Agent{
		llmClient: client,
		config:    cfg,
	}
}

// Run executes the agent on a task
func (a *Agent) Run(ctx context.Context, task string) (*AgentResult, error) {
	start := time.Now()

	// Simple completion for now
	messages := []llm.Message{
		{Role: "system", Content: "You are XLI Bot, a helpful AI assistant."},
		{Role: "user", Content: task},
	}

	opts := &llm.CompletionOpts{
		Model:       a.config.LLMModel,
		Temperature: a.config.LLMTemperature,
		MaxTokens:   a.config.LLMMaxTokens,
	}

	result, err := a.llmClient.Complete(ctx, messages, opts)
	if err != nil {
		return nil, fmt.Errorf("llm complete: %w", err)
	}

	return &AgentResult{
		Answer: result.Content,
		TokenUsage: llm.TokenUsage{
			InputTokens:  result.InputTokens,
			OutputTokens: result.OutputTokens,
			TotalTokens:  result.TotalTokens,
		},
		Elapsed: time.Since(start),
	}, nil
}
