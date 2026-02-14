package movie

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

type Repository interface {
	Create(ctx context.Context, movie *Movie) error
	Search(ctx context.Context, embedding []float32, limit int) ([]Movie, error)
}

type repository struct {
	conn *pgxpool.Pool
}

func NewRepository(conn *pgxpool.Pool) *repository {
	return &repository{conn: conn}
}

func (r *repository) Create(ctx context.Context, movie *Movie) error {
	query := `
		INSERT INTO movies (title, overview, genres, language, popularity, release_date, embedding, embedded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
	`

	movie.Genres = []string{} // Ensure Genres is not nil to avoid issues with NULL arrays

	_, err := r.conn.Exec(
		ctx,
		query,
		movie.Title,
		movie.Overview,
		movie.Genres,
		movie.Language,
		movie.Popularity,
		movie.ReleaseDate,
		pgvector.NewVector(movie.Embedding),
	)
	if err != nil {
		return fmt.Errorf("failed to insert movie: %v", err)
	}
	return nil
}

func (r *repository) Search(ctx context.Context, embedding []float32, limit int) ([]Movie, error) {
	query := `
		SELECT id, title, overview, genres, language, popularity, release_date, embedded_at
		FROM movies
		ORDER BY embedding <=> $1
		LIMIT $2
	`

	rows, err := r.conn.Query(ctx, query, pgvector.NewVector(embedding), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %v", err)
	}
	defer rows.Close()

	movies, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Movie])
	if err != nil {
		return nil, fmt.Errorf("failed to collect movies: %v", err)
	}
	return movies, nil
}
