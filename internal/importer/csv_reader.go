package importer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/edmarfelipe/ai-recomendation/internal/movie"
)

type CSVReader struct {
	filePath string
}

func NewCSVReader(filePath string) *CSVReader {
	return &CSVReader{
		filePath: filePath,
	}
}

func (r *CSVReader) Read(limit int) (<-chan *movie.Movie, <-chan error) {
	moviesChan := make(chan *movie.Movie)
	errorsChan := make(chan error, 1)

	go func() {
		defer close(moviesChan)
		defer close(errorsChan)

		file, err := os.Open(r.filePath)
		if err != nil {
			errorsChan <- fmt.Errorf("unable to open file: %v", err)
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		count := 0

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				errorsChan <- fmt.Errorf("unable to read record: %v", err)
				return
			}

			// Skip header
			if record[3] == "id" {
				continue
			}

			if count >= limit {
				break
			}

			m, err := parseMovieFromRecord(record)
			if err != nil {
				// Skip invalid records
				continue
			}

			moviesChan <- m
			count++
		}
	}()

	return moviesChan, errorsChan
}

func parseMovieFromRecord(record []string) (*movie.Movie, error) {
	id, err := strconv.Atoi(record[3])
	if err != nil {
		return nil, fmt.Errorf("unable to parse id: %v", err)
	}

	popularity, err := strconv.ParseFloat(record[8], 64)
	if err != nil {
		return nil, fmt.Errorf("unable to parse popularity: %v", err)
	}

	genresList := []struct {
		Name string `json:"name"`
	}{}
	if err := json.Unmarshal([]byte(record[1]), &genresList); err != nil {
		return nil, fmt.Errorf("unable to parse categories: %v", err)
	}
	var genres []string
	for _, c := range genresList {
		genres = append(genres, c.Name)
	}

	releaseDate, err := time.Parse("2006-01-02", record[11])
	if err != nil {
		return nil, fmt.Errorf("unable to parse release date: %v", err)
	}

	return &movie.Movie{
		ID:          int64(id),
		Title:       record[17],
		Overview:    record[7],
		Genres:      genres,
		Language:    record[5],
		Popularity:  popularity,
		ReleaseDate: releaseDate,
	}, nil
}
