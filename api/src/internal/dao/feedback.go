package dao

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/dao"
	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/types"
)

type feedbackDao struct {
	*sql.Tx
}

func MustMakeFeedbackDao(tx *sql.Tx) dao.FeedbackDao {
	return feedbackDao{tx}
}

func (f feedbackDao) SaveFeedback(ctx context.Context, req types.SubmitGenuityFeedbackRequest) error {
	result, err := sq.
		Insert("feedback").Columns("review_id", "pros", "cons", "original_is_genuine", "user_is_genuine", "created_by", "created_at").
		Values(req.ReviewID, req.Pros, req.Cons, req.Original.IsGenuine, req.Feedback.IsGenuine, req.UserID, "NOW()").
		RunWith(f).ExecContext(ctx)

	if err != nil {
		return fmt.Errorf("failed to insert feedback: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected != 1 {
		return fmt.Errorf("expected 1 row to be affected, got %d", rowsAffected)
	}

	return nil
}
