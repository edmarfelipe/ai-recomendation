package internal

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type SearchResponseItem struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Overview    string    `json:"overview"`
	Genres      []string  `json:"genres"`
	Language    string    `json:"language"`
	Popularity  float64   `json:"popularity"`
	ReleaseDate time.Time `json:"release_date"`
}

func (s *server) handlerSearch(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query().Get("query")
	if query == "" {
		return &APIError{"missing query parameter", http.StatusBadRequest}
	}

	if len(query) < 3 {
		return &APIError{"query parameter too short", http.StatusBadRequest}
	}

	slog.Info("Handling search", "query", query)
	movies, err := s.movieService.Search(r.Context(), query, 5)
	if err != nil {
		return err
	}

	var resp []SearchResponseItem
	for _, movie := range movies {
		item := SearchResponseItem{
			ID:          movie.ID,
			Title:       movie.Title,
			Overview:    movie.Overview,
			Genres:      movie.Genres,
			Language:    movie.Language,
			Popularity:  movie.Popularity,
			ReleaseDate: movie.ReleaseDate,
		}
		resp = append(resp, item)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}
