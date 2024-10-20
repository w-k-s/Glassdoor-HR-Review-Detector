package server

import (
	"database/sql"
	"net/http"

	"com.github/w-k-s/glassdoor-hr-review-detector/internal"
	"com.github/w-k-s/glassdoor-hr-review-detector/internal/dao"
	"com.github/w-k-s/glassdoor-hr-review-detector/internal/embedding"
	"com.github/w-k-s/glassdoor-hr-review-detector/internal/inferrence"
	"com.github/w-k-s/glassdoor-hr-review-detector/internal/migrations"
	"com.github/w-k-s/glassdoor-hr-review-detector/internal/objectstore"
	"com.github/w-k-s/glassdoor-hr-review-detector/pkg"
	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jasonlvhit/gocron"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	listenAddress   string
	db              *sql.DB
	genuityService  services.GenuityService
	trainingService services.TrainingService
	config          *Config
}

func NewServer(config *Config) *Server {
	db := dao.Must(sql.Open("sqlite3", ":memory:"))
	pkg.Must(migrations.Exec(config.Migrations.Directory, db))

	inMemoryCache := internal.LocalCache()
	embeddingService := embedding.MustOpenAIEmbeddingService(config.OpenAi.ApiKey, inMemoryCache)
	inferenceService := inferrence.MustInferenceService(config.Inference.Api.Endpoint, embeddingService, inMemoryCache)
	s3 := objectstore.MustS3(config.S3.Region)

	return &Server{
		listenAddress:   config.Server.ListenAddress,
		db:              db,
		genuityService:  services.MustGenuityService(inferenceService),
		trainingService: services.MustTrainingService(s3),
		config:          config,
	}
}

func (s *Server) Start() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/api/reviews/genuity-check", RESTEndpoint(s.checkReviewsGenuity))
	r.Post("/api/reviews/genuity-feedback", RESTEndpoint(s.submitGenuityFeedback))

	gocron.Every(uint64(s.config.Feedback.Upload.Frequency.Hours)).Hours().Do(s.uploadFeedback)
	gocron.Start()
	return http.ListenAndServe(s.listenAddress, r)
}
