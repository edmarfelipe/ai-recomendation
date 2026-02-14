package internal

import (
	"fmt"
	"log/slog"
	"os"
)

type EnvConfig struct {
	ServerAddr      string
	GeminiAPIKey    string
	GeminiModel     string
	GeminiRateLimit int
	DatabaseURL     string
}

func LoadEnvConfig() (*EnvConfig, error) {
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":9999"
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is required but not set")
	}

	model := os.Getenv("GEMINI_MODEL")
	if model == "" {
		return nil, fmt.Errorf("GEMINI_MODEL environment variable is required but not set")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required but not set")
	}

	slog.Info("Env variables loaded successfully")

	return &EnvConfig{
		ServerAddr:      serverAddr,
		GeminiAPIKey:    apiKey,
		DatabaseURL:     databaseURL,
		GeminiModel:     model,
		GeminiRateLimit: 90, // 90 requests per minute
	}, nil
}
