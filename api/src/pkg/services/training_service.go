package services

import (
	"context"
	"fmt"
	"log"

	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/dao"
	"github.com/gocarina/gocsv"
)

type TrainingService interface {
	UploadFeedback(context.Context, dao.FeedbackDao, string, string) error
}

type trainingService struct {
	objectStoreService ObjectStoreService
}

func MustTrainingService(os ObjectStoreService) TrainingService {
	if os == nil {
		log.Panic("objectStoreService is nil")
	}

	return &trainingService{
		objectStoreService: os,
	}
}

func (t trainingService) UploadFeedback(
	ctx context.Context,
	feedbackDao dao.FeedbackDao,
	trainingFileBucket string,
	trainingFileName string,
) error {
	feedback, err := feedbackDao.GetFeedback(ctx)
	if err != nil {
		return fmt.Errorf("Failed to retrieve feedback: %w", err)
	}
	log.Printf("Retrieved %d feedback", len(feedback))

	if len(feedback) > 0 {
		csv, err := gocsv.MarshalString(&feedback)
		if err != nil {
			return fmt.Errorf("Failed to marshal training data into csv. %w", err)
		}

		err = t.objectStoreService.Put(trainingFileBucket, trainingFileName, csv)
		if err != nil {
			return fmt.Errorf("Failed to upload new rows to training file %q in bucket %q. Reason: %q", trainingFileName, trainingFileBucket, err)
		}
	}

	return nil
}
