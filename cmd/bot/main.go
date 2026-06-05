package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/oblachko/xli-bot/internal/agent"
	"github.com/oblachko/xli-bot/internal/config"
	"github.com/oblachko/xli-bot/internal/llm"
	"github.com/oblachko/xli-bot/internal/transport"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("XLI Bot starting...")
	log.Printf("LLM Provider: %s", cfg.LLMProvider)
	log.Printf("LLM Model: %s", cfg.LLMModel)
	log.Printf("DB Path: %s", cfg.DBPath)

	// Create LLM client
	llmClient := llm.Factory(cfg.LLMProvider, cfg.LLMAPIKey, cfg.LLMBaseURL)

	// Create agent
	ag := agent.NewAgent(llmClient, cfg)

	// Create Telegram transport
	transport, err := transport.NewTelegramTransport(cfg, ag)
	if err != nil {
		log.Fatalf("Failed to create transport: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		os.Exit(0)
	}()

	// Start bot
	log.Println("Starting Telegram bot...")
	if err := transport.Start(); err != nil {
		log.Fatalf("Transport error: %v", err)
	}
}
