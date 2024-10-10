package services

import (
	"context"

	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/types"
)



type InferenceService interface {
	GetGenuity(ctx context.Context,reviews []types.Review) ([]types.GenuityResult, error)
}

