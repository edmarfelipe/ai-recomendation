package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/edmarfelipe/ai-recomendation/internal"
	"github.com/edmarfelipe/ai-recomendation/internal/db"
	"github.com/edmarfelipe/ai-recomendation/internal/gemini"
	"github.com/edmarfelipe/ai-recomendation/internal/importer"
)

var (
	pFilePath = flag.String("file", "tmdb_5000_movies.csv", "Path to the movies dataset CSV file")
	pSize     = flag.Int("size", 100, "Number of movies to import")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	slog.Info("Starting import movies", "filePath", *pFilePath, "size", *pSize)

	if err := run(ctx, *pFilePath, *pSize); err != nil {
		slog.Error("import movies failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, filePath string, size int) error {
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

	imp := importer.NewMovieImporter(conn, gemini)
	if err := imp.Import(ctx, filePath, size); err != nil {
		return fmt.Errorf("unable to import movies: %v", err)
	}

	slog.Info("Movies imported successfully")
	return nil
}
