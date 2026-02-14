package internal

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/edmarfelipe/ai-recomendation/internal/gemini"
	"github.com/edmarfelipe/ai-recomendation/internal/movie"
	"github.com/jackc/pgx/v5/pgxpool"
)

type server struct {
	addr         string
	gemini       *gemini.Service
	movieService *movie.Service
}

func NewServer(addr string, gemini *gemini.Service, conn *pgxpool.Pool) *server {
	return &server{
		addr:         addr,
		gemini:       gemini,
		movieService: movie.NewMovieService(movie.NewRepository(conn), gemini),
	}
}

func (s *server) Start() error {
	slog.Info("Starting server", "addr", s.addr)

	http.HandleFunc("/", s.handlerHome)
	http.HandleFunc("GET /v1/movies", HandleWithError(s.handlerSearch))

	err := http.ListenAndServe(s.addr, nil)
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %v", err)
	}
	return nil
}
