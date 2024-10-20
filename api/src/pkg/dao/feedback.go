package dao

import (
	"context"

	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/types"
)

type FeedbackDao interface {
	SaveFeedback(context.Context, types.SubmitGenuityFeedbackRequest) error
	GetFeedback(context.Context) ([]types.SubmitGenuityFeedbackRequest, error)
}
