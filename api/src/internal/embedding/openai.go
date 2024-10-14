package embedding

import (
	"context"
	"fmt"
	"log"
	"time"

	"com.github/w-k-s/glassdoor-hr-review-detector/internal"
	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/services"
	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
)

type openaiEmbeddingService struct {
	client           *openai.Client
	cache            internal.Cache
	inflightRequests map[string]chan bool
}

func MustOpenAIEmbeddingService(apiKey string, cache internal.Cache) services.EmbeddingsService {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &openaiEmbeddingService{
		client:           client,
		cache:            cache,
		inflightRequests: make(map[string]chan bool),
	}
}

func (svc openaiEmbeddingService) GetEmbeddings(ctx context.Context, inputs []string, dimensions int) ([][]float64, error) {

	embeddingMap := make(map[string][]float64)
	newInputs := make([]string, 0)

	for i, input := range inputs {
		if len(input) == 0 {
			return nil, fmt.Errorf("input can not contain empty string. %d in %q", i, inputs)
		}

		// To avoid sending multiple concurrent requests for the same inputs.
		if c, ok := svc.inflightRequests[input]; ok {
			waitForInlfightEmbedding(c)
		}

		if untypedEmbedding, ok := svc.cache.Get(input); ok {
			if embedding, ok := untypedEmbedding.([]float64); ok {
				embeddingMap[input] = embedding
			}
		} else {
			newInputs = append(newInputs, input)
		}
	}

	// Fetch embeddings for remaining items
	if len(newInputs) > 0 {
		log.Printf("fetching embeddings for %d new inputs", len(newInputs))

		done := make(chan bool)
		for _, newInput := range newInputs {
			svc.inflightRequests[newInput] = done
		}

		createEmbeddingResponse, err := svc.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
			Input:          openai.F[openai.EmbeddingNewParamsInputUnion](openai.EmbeddingNewParamsInputArrayOfStrings(newInputs)),
			Model:          openai.F(openai.EmbeddingModelTextEmbedding3Small),
			Dimensions:     openai.F(int64(dimensions)),
			EncodingFormat: openai.F(openai.EmbeddingNewParamsEncodingFormatFloat),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get embedding for inputs: %w", err)
		}

		data := createEmbeddingResponse.Data

		if len(data) != len(newInputs) {
			return nil, fmt.Errorf("unexpected: the number of embeddings returned does not match the number of inputs")
		}

		for i, newInput := range newInputs {
			embedding := data[i].Embedding
			svc.cache.Set(newInput, embedding)
			embeddingMap[newInput] = embedding
		}

		close(done)
	}

	embeddings := make([][]float64, len(inputs))
	for i, d := range inputs {
		if embedding, ok := embeddingMap[d]; ok {
			embeddings[i] = embedding
		}
	}

	return embeddings, nil
}

func waitForInlfightEmbedding(done <-chan bool) {
	timeout := time.After(30 * time.Second)

	for {
		select {
		case <-done:
			return
		case <-timeout:
			log.Printf("Timeout")
			return
		}
	}
}
