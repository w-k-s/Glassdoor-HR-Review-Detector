package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

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
	// Fetch feedback in the last 24 hours
	feedback, err := feedbackDao.GetTodaysFeedback(ctx)
	if err != nil {
		return fmt.Errorf("Failed to retrieve today's feedback: %w", err)
	}
	log.Printf("Retrieved %d feedback in the last 24 hours", len(feedback))

	// Fetch Training File
	r, err := t.objectStoreService.Get(trainingFileBucket, trainingFileName)
	if err != nil {
		return fmt.Errorf("Failed to fetch training file %q from bucket %q. Reason: %q", trainingFileName, trainingFileBucket, err)
	}

	content, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("Failed to read training file %q. Reason: %q", trainingFileName, err)
	}

	// Append to Training File as CSV without headers
	rows, err := gocsv.MarshalString(&feedback)
	if err != nil {
		return fmt.Errorf("Failed to marshal training data into csv. %w", err)
	}

	// Put Training File
	newContent := strings.Join([]string{string(content), rows}, "\n")
	err = t.objectStoreService.Put(trainingFileBucket, trainingFileName, newContent)
	if err != nil {
		return fmt.Errorf("Failed to upload new rows to training file %q in bucket %q. Reason: %q", trainingFileName, trainingFileBucket, err)
	}

	return nil
}
