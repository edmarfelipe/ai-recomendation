package importer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/edmarfelipe/ai-recomendation/internal/gemini"
	"github.com/edmarfelipe/ai-recomendation/internal/movie"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Importer struct {
	serv *movie.Service
}

func NewMovieImporter(conn *pgxpool.Pool, gemini *gemini.Service) *Importer {
	return &Importer{
		serv: movie.NewMovieService(movie.NewRepository(conn), gemini),
	}
}

func (i *Importer) Import(ctx context.Context, filePath string, size int) error {
	reader := NewCSVReader(filePath)
	moviesChan, errorsChan := reader.Read(size)

	totalRecords := 0

	for moviesChan != nil || errorsChan != nil {
		select {
		case err, ok := <-errorsChan:
			if !ok {
				errorsChan = nil
				continue
			}
			if err != nil {
				return err
			}
		case m, ok := <-moviesChan:
			if !ok {
				moviesChan = nil
				continue
			}

			if err := i.createMovie(ctx, m); err != nil {
				slog.WarnContext(ctx, "Failed to import movie", "id", m.ID, "title", m.Title, "error", err)
				continue
			}

			totalRecords++
		}
	}

	slog.InfoContext(ctx, "Import completed", "total", totalRecords)
	return nil
}

func (i *Importer) createMovie(ctx context.Context, m *movie.Movie) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	slog.InfoContext(ctx, "Importing movie", "id", m.ID, "title", m.Title)
	if err := i.serv.Create(ctx, m); err != nil {
		return fmt.Errorf("unable to create movie: %v", err)
	}
	return nil
}
