package services

import (
	"context"
)

type EmbeddingsService interface {
	GetEmbeddings(ctx context.Context, text []string, dimensions int) ([][]float64, error)
}
