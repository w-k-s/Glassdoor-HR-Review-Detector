package embedding

import (
	"context"

	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/services"
	"github.com/openai/openai-go" // imported as openai
	"github.com/openai/openai-go/option"
)

type openaiEmbeddingService struct {
	client *openai.Client
}

func MustOpenAIEmbeddingService(apiKey string) services.EmbeddingsService {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	
	return &openaiEmbeddingService{
		client: client,
	}
}

func (svc openaiEmbeddingService) GetEmbeddings(ctx context.Context, inputs []string, dimensions int) ([][]float64, error) {
	createEmbeddingResponse, err := svc.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input:          openai.F[openai.EmbeddingNewParamsInputUnion](openai.EmbeddingNewParamsInputArrayOfStrings(inputs)),
		Model:          openai.F(openai.EmbeddingModelTextEmbedding3Small),
		Dimensions:     openai.F(int64(dimensions)),
		EncodingFormat: openai.F(openai.EmbeddingNewParamsEncodingFormatFloat),
	})
	if err != nil {
		return nil, err
	}

	data := createEmbeddingResponse.Data
	embeddings := make([][]float64, len(data))
    for i, d := range data {
        embeddings[i] = d.Embedding
    }
    return embeddings,nil
}
