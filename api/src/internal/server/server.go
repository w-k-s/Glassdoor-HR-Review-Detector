package server

import (
	"net/http"
	"os"

	"com.github/w-k-s/glassdoor-hr-review-detector/internal"
	"com.github/w-k-s/glassdoor-hr-review-detector/internal/embedding"
	"com.github/w-k-s/glassdoor-hr-review-detector/internal/inferrence"
	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	listenAddress  string
	genuityService services.GenuityService
}

func NewServer(listenAddress string) *Server {
	inMemoryCache := internal.LocalCache()
	embeddingService := embedding.MustOpenAIEmbeddingService(os.Getenv("OPENAI_API_KEY"), inMemoryCache)
	inferenceService := inferrence.MustInferenceService(embeddingService, inMemoryCache)
	return &Server{
		listenAddress:  listenAddress,
		genuityService: services.MustGenuityService(inferenceService),
	}
}

func (s *Server) Start() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/api/reviews/genuity-check", s.checkReviewsGenuity)
	r.Post("/api/reviews/genuity-feedback", s.submitGenuityFeedback)
	return http.ListenAndServe(":3000", r)
}
