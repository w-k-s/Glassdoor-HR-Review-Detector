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

	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/services"
	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/types"
)

const vectorDimensions = 256
const inputSize = 1 + (2*vectorDimensions)

type inferenceService struct {
	embeddingService services.EmbeddingsService
}

type predictionResponse struct {
    Predictions [][]float64 `json:"predictions"`
}

func MustInferenceService(embeddingService services.EmbeddingsService) services.InferenceService {
	if embeddingService == nil {
		log.Panic("Embedding Service is nil")
	}
	
	return &inferenceService{
		embeddingService: embeddingService,
	}
}

func (svc inferenceService) GetGenuity(ctx context.Context,reviews []types.Review) ([]types.GenuityResult, error){
	
	embeddingRequestInputs := make([]string, len(reviews)*2)
    for i, d := range reviews {
        embeddingRequestInputs[i] = d.Pros
		embeddingRequestInputs[i+1] = d.Cons
		i += 2
    }

	embeddings,err := svc.embeddingService.GetEmbeddings(ctx, embeddingRequestInputs, vectorDimensions)
	if err != nil{
		return nil, err
	}

	x := make([][]float64, len(reviews))
	for i, d := range reviews {
		rating := d.Rating
		prosEmbedding := embeddings[i*2]
		consEmbedding := embeddings[1+(i*2)]
		
		consIndex := vectorDimensions + 1

		x[i] = make([]float64, inputSize)
		x[i][0] = rating
		copy(x[i][1:consIndex], prosEmbedding)
		copy(x[i][consIndex:], consEmbedding)
    }

	jsonData,_ := json.Marshal(struct {
		Instances [][]float64 `json:"instances"`
	}{
		Instances: x,
	})

	resp, err := http.Post("http://localhost:8501/v1/models/glassdoor_hr_review_detector:predict", "application/json", bytes.NewReader(jsonData))
    if err != nil {
        return nil, fmt.Errorf("inference API failed: %w", err)
    }
    defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }

	log.Print(string(body))
    var predResp predictionResponse
    err = json.Unmarshal(body, &predResp)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }

    predictions := make([]bool, len(predResp.Predictions))
    for i, pred := range predResp.Predictions {
        if len(pred) != 1 {
            return nil, fmt.Errorf("unexpected prediction array length at index %d: got %d, want 1", i, len(pred))
        }
        predictions[i] = sigmoid(pred[0]) < 0.5
    }

	genuityResults := make([]types.GenuityResult, len(reviews))
	for i, review := range reviews {
		genuityResults[i].ReviewID = review.ID
		genuityResults[i].IsGenuine = predictions[i]
	}

	return genuityResults,nil
}

func sigmoid(v float64) float64 {
	return (1 / (1 + math.Exp(-v)))
}