package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all bot configuration
type Config struct {
	// Telegram
	TelegramToken string

	// LLM Provider
	LLMProvider    string // mistral, groq, openrouter
	LLMAPIKey      string
	LLMModel       string
	LLMBaseURL     string
	LLMTemperature float64
	LLMMaxTokens   int

	// Database
	DBPath string

	// MCP
	MCPServersDir string

	// Skills
	SkillsDir string

	// Build
	GitHubToken string
	GitHubOwner string
	GitHubRepo  string

	// Bot behavior
	MaxSteps              int
	ConfirmationLevel     string // low, medium, high
	ContextEnabled        bool
	ContextTurns          int
	ContextCompaction     bool
	ContextCompactionChars int
}

// Load reads configuration from .env file and environment variables
func Load() (*Config, error) {
	// Try to load .env file, ignore error if not found
	_ = godotenv.Load()

	cfg := &Config{
		LLMProvider:            getEnv("LLM_PROVIDER", "mistral"),
		LLMModel:               getEnv("LLM_MODEL", ""),
		LLMTemperature:         getEnvFloat("LLM_TEMPERATURE", 0.7),
		LLMMaxTokens:           getEnvInt("LLM_MAX_TOKENS", 1200),
		DBPath:                 getEnv("DB_PATH", "./data/xli.db"),
		MCPServersDir:          getEnv("MCP_SERVERS_DIR", "./mcp_servers"),
		SkillsDir:              getEnv("SKILLS_DIR", "./skills"),
		GitHubOwner:            getEnv("GITHUB_OWNER", ""),
		GitHubRepo:             getEnv("GITHUB_REPO", "xli-build-runner"),
		MaxSteps:               getEnvInt("MAX_STEPS", 15),
		ConfirmationLevel:      getEnv("CONFIRMATION_LEVEL", "medium"),
		ContextEnabled:         getEnvBool("CONTEXT_ENABLED", true),
		ContextTurns:           getEnvInt("CONTEXT_TURNS", 10),
		ContextCompaction:      getEnvBool("CONTEXT_COMPACTION", true),
		ContextCompactionChars: getEnvInt("CONTEXT_COMPACTION_CHARS", 18000),
	}

	// Required fields
	cfg.TelegramToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	if cfg.TelegramToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	cfg.LLMAPIKey = os.Getenv("LLM_API_KEY")
	if cfg.LLMAPIKey == "" {
		return nil, fmt.Errorf("LLM_API_KEY is required")
	}

	cfg.LLMBaseURL = os.Getenv("LLM_BASE_URL")
	if cfg.LLMBaseURL == "" {
		// Default URLs per provider
		switch cfg.LLMProvider {
		case "mistral":
			cfg.LLMBaseURL = "https://api.mistral.ai/v1"
		case "groq":
			cfg.LLMBaseURL = "https://api.groq.com/openai/v1"
		case "openrouter":
			cfg.LLMBaseURL = "https://openrouter.ai/api/v1"
		default:
			return nil, fmt.Errorf("LLM_BASE_URL is required for provider %s", cfg.LLMProvider)
		}
	}

	cfg.GitHubToken = os.Getenv("GITHUB_TOKEN")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		switch strings.ToLower(value) {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		}
	}
	return defaultValue
}
