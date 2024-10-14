package services

import (
	"context"
	"fmt"
	"log"

	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/dao"
	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/types"
)

type GenuityService interface {
	CheckReviewGenuity(ctx context.Context, req types.CheckReviewsGenuityRequest) (types.CheckReviewsGenuityResponse, error)
	SubmitGenuityFeedback(ctx context.Context, dao dao.FeedbackDao, req types.SubmitGenuityFeedbackRequest) error
}

type genuityService struct {
	inferenceService InferenceService
}

func MustGenuityService(
	inferenceService InferenceService,
) GenuityService {
	if inferenceService == nil {
		log.Panic("Embedding Service is nil")
	}

	return &genuityService{
		inferenceService: inferenceService,
	}
}

func (svc genuityService) CheckReviewGenuity(ctx context.Context, req types.CheckReviewsGenuityRequest) (types.CheckReviewsGenuityResponse, error) {

	genuityResult, err := svc.inferenceService.GetGenuity(ctx, req.Reviews)
	if err != nil {
		return types.CheckReviewsGenuityResponse{}, fmt.Errorf("inference Request Faield: %w", err)
	}

	return types.CheckReviewsGenuityResponse{Results: genuityResult}, nil
}

func (svc genuityService) SubmitGenuityFeedback(ctx context.Context, dao dao.FeedbackDao, req types.SubmitGenuityFeedbackRequest) error {
	return dao.SaveFeedback(ctx, req)
}
