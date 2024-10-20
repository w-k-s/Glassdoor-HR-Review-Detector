package server

import (
	"database/sql"
	"net/http"
	"os"

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
}

func NewServer(listenAddress, migrationsDirectory, awsRegion string) *Server {
	db := dao.Must(sql.Open("sqlite3", ":memory:"))
	pkg.Must(migrations.Exec(migrationsDirectory, db))

	inMemoryCache := internal.LocalCache()
	embeddingService := embedding.MustOpenAIEmbeddingService(os.Getenv("OPENAI_API_KEY"), inMemoryCache)
	inferenceService := inferrence.MustInferenceService(embeddingService, inMemoryCache)
	s3 := objectstore.MustS3(awsRegion)

	return &Server{
		listenAddress:   listenAddress,
		db:              db,
		genuityService:  services.MustGenuityService(inferenceService),
		trainingService: services.MustTrainingService(s3),
	}
}

func (s *Server) Start() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/api/reviews/genuity-check", s.checkReviewsGenuity)
	r.Post("/api/reviews/genuity-feedback", s.submitGenuityFeedback)

	gocron.Every(10).Seconds().Do(s.uploadFeedback)
	gocron.Start()
	return http.ListenAndServe(":3000", r)
}
