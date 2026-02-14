package movie

import (
	"time"
)

type Movie struct {
	ID          int64
	Title       string
	Overview    string
	Genres      []string
	Language    string
	Popularity  float64
	ReleaseDate time.Time
	Embedding   []float32  `json:"-" db:"-" toon:"-"`
	EmbeddedAt  *time.Time `json:"-" toon:"-"`
}
