package gemini

import (
	"context"
	"fmt"
	"time"

	"github.com/toon-format/toon-go"
	"golang.org/x/time/rate"
	"google.golang.org/genai"
)

type Service struct {
	gemini  *genai.Client
	model   string
	limiter *rate.Limiter
}

func NewService(ctx context.Context, apiKey string, model string, rpm int) (*Service, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}
	return &Service{
		gemini:  client,
		model:   model,
		limiter: rate.NewLimiter(rate.Every(time.Minute/time.Duration(rpm)), 1),
	}, nil
}

func (s *Service) Embed(ctx context.Context, input any) ([]float32, error) {
	encoded, err := toon.Marshal(input, toon.WithLengthMarkers(true))
	if err != nil {
		return nil, fmt.Errorf("unable to marshal movie: %v", err)
	}

	return s.getEmbedding(ctx, encoded)
}

func (s *Service) getEmbedding(ctx context.Context, encoded []byte) ([]float32, error) {
	if err := s.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %v", err)
	}

	contents := []*genai.Content{
		genai.NewContentFromText(string(encoded), genai.RoleUser),
	}

	result, err := s.gemini.Models.EmbedContent(ctx, s.model, contents, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get embeddings from Gemini: %v", err)
	}

	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned from Gemini")
	}

	return result.Embeddings[0].Values, nil
}
