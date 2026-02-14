package movie

import (
	"context"
	"fmt"
	"strings"

	"github.com/edmarfelipe/ai-recomendation/internal/gemini"
)

type Service struct {
	repo   Repository
	gemini *gemini.Service
}

func NewMovieService(repo Repository, gemini *gemini.Service) *Service {
	return &Service{
		repo:   repo,
		gemini: gemini,
	}
}

func (s *Service) Create(ctx context.Context, movie *Movie) error {
	embedding, err := s.gemini.Embed(ctx, movie)
	if err != nil {
		return fmt.Errorf("unable to embed documents: %v", err)
	}

	movie.Embedding = embedding

	err = s.repo.Create(ctx, movie)
	if err != nil {
		return fmt.Errorf("unable to create movie: %v", err)
	}
	return nil
}

func (s *Service) Search(ctx context.Context, query string, limit int) ([]Movie, error) {
	embedding, err := s.gemini.Embed(ctx, strings.TrimSpace(query))
	if err != nil {
		return nil, fmt.Errorf("unable to embed query: %v", err)
	}

	movies, err := s.repo.Search(ctx, embedding, limit)
	if err != nil {
		return nil, fmt.Errorf("unable to search movies: %v", err)
	}
	return movies, nil
}
