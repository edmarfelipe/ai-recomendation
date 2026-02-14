package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/edmarfelipe/ai-recomendation/internal"
	"github.com/edmarfelipe/ai-recomendation/internal/db"
	"github.com/edmarfelipe/ai-recomendation/internal/gemini"
)

func main() {
	slog.Info("Starting AI Recommendation Service")
	ctx := context.Background()

	if err := run(ctx); err != nil {
		slog.Error("service failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := internal.LoadEnvConfig()
	if err != nil {
		return err
	}
	gemini, err := gemini.NewService(ctx, cfg.GeminiAPIKey, cfg.GeminiModel, cfg.GeminiRateLimit)
	if err != nil {
		return fmt.Errorf("unable to create Gemini service: %v", err)
	}

	conn, closeFunc, err := db.OpenDB(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("unable to open database: %v", err)
	}
	defer closeFunc()
	slog.Info("Database initialized")

	return internal.NewServer(cfg.ServerAddr, gemini, conn).Start()
}
