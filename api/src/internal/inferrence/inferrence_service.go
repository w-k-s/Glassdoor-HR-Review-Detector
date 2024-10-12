package inferrence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"

	"com.github/w-k-s/glassdoor-hr-review-detector/internal"
	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/services"
	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/types"
)

const vectorDimensions = 256
const inputSize = 1 + (2 * vectorDimensions)

type predictionResponse struct {
	Predictions [][]float64 `json:"predictions"`
}

type inferenceService struct {
	embeddingService services.EmbeddingsService
	cache            internal.Cache
}

func MustInferenceService(
	embeddingService services.EmbeddingsService,
	cache internal.Cache,
) services.InferenceService {
	if embeddingService == nil {
		log.Panic("Embedding Service is nil")
	}

	if cache == nil {
		log.Panic("cache is nil")
	}

	return &inferenceService{
		embeddingService: embeddingService,
		cache:            cache,
	}
}

func (svc inferenceService) GetGenuity(ctx context.Context, reviews []types.Review) ([]types.GenuityResult, error) {
	genuityResults := make([]types.GenuityResult, 0)
	newReviews := make([]types.Review, 0)
	for _, review := range reviews {
		if result, ok := svc.cache.Get(review.ID); ok {
			genuityResults = append(genuityResults, result.(types.GenuityResult))
		} else {
			newReviews = append(newReviews, review)
		}
	}

	if len(newReviews) > 0 {
		embeddingRequestInputs := make([]string, len(newReviews)*2)
		for i, d := range newReviews {
			embeddingRequestInputs[i] = d.Pros
			embeddingRequestInputs[i+1] = d.Cons
			i += 2
		}

		embeddings, err := svc.embeddingService.GetEmbeddings(ctx, embeddingRequestInputs, vectorDimensions)
		if err != nil {
			return nil, fmt.Errorf("embedding Service Error: %w", err)
		}

		x := make([][]float64, len(newReviews))
		for i, d := range newReviews {
			rating := d.Rating
			prosEmbedding := embeddings[i*2]
			consEmbedding := embeddings[1+(i*2)]

			consIndex := vectorDimensions + 1

			x[i] = make([]float64, inputSize)
			x[i][0] = rating
			copy(x[i][1:consIndex], prosEmbedding)
			copy(x[i][consIndex:], consEmbedding)
		}

		jsonData, _ := json.Marshal(struct {
			Instances [][]float64 `json:"instances"`
		}{
			Instances: x,
		})

		// TODO: Externalise API Endpoint
		resp, err := http.Post("http://localhost:8501/v1/models/glassdoor_hr_review_detector:predict", "application/json", bytes.NewReader(jsonData))
		if err != nil {
			return nil, fmt.Errorf("inference API Error: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read inference response: %w", err)
		}

		var predResp predictionResponse
		err = json.Unmarshal(body, &predResp)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal inference response: %w", err)
		}

		predictions := make([]bool, len(predResp.Predictions))
		for i, pred := range predResp.Predictions {
			if len(pred) != 1 {
				return nil, fmt.Errorf("unexpected prediction array length at index %d: got %d, want 1", i, len(pred))
			}
			predictions[i] = sigmoid(pred[0]) < 0.5
		}

		for i, newReview := range newReviews {
			result := types.GenuityResult{
				ReviewID:  newReview.ID,
				IsGenuine: predictions[i],
			}
			svc.cache.Set(newReview.ID, result)
			genuityResults = append(genuityResults, result)
		}
	}

	return genuityResults, nil
}

func sigmoid(v float64) float64 {
	return (1 / (1 + math.Exp(-v)))
}
