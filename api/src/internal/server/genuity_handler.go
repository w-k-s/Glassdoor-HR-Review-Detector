package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"com.github/w-k-s/glassdoor-hr-review-detector/internal/dao"
	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/types"
)

func (s *Server) checkReviewsGenuity(req *http.Request) (any, error, int) {
	var request types.CheckReviewsGenuityRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, fmt.Errorf("Invalid request body"), http.StatusBadRequest
	}

	response, err := s.genuityService.CheckReviewGenuity(req.Context(), request)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return response, nil, http.StatusOK
}

func (s *Server) submitGenuityFeedback(req *http.Request) (any, error, int) {
	var request types.SubmitGenuityFeedbackRequest

	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		return nil, fmt.Errorf("Invalid request body"), http.StatusBadRequest
	}

	tx, err := s.db.BeginTx(req.Context(), nil)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	defer tx.Rollback()

	feedbackDao := dao.MustMakeFeedbackDao(tx)
	err = s.genuityService.SubmitGenuityFeedback(req.Context(), feedbackDao, request)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	err = tx.Commit()
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return nil, nil, http.StatusCreated
}

func (s *Server) uploadFeedback() {
	log.Printf("Running Upload Feedback Job")

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		log.Printf("failed to begin transaction. Reason: %q", err)
	}

	err = s.trainingService.UploadFeedback(context.Background(), dao.MustMakeFeedbackDao(tx), s.config.S3.Bucket, "training-data/feedback.csv")
	if err != nil {
		log.Printf("failed to upload feedback. Reason: %q", err)
	}

}
