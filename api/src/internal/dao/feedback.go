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
		Insert("feedback").Columns("review_id", "rating", "pros", "cons", "original_is_genuine", "user_is_genuine", "created_by", "created_at").
		Values(req.ReviewID, req.OverallRating, req.Pros, req.Cons, req.Original.IsGenuine, req.Feedback.IsGenuine, req.UserID, "DATE()").
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

func (f feedbackDao) GetFeedback(ctx context.Context) ([]types.SubmitGenuityFeedbackRequest, error) {
	rows, err := sq.
		Select("review_id", "rating", "pros", "cons", "original_is_genuine", "user_is_genuine", "created_by").
		From("feedback").
		RunWith(f).
		QueryContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch feedback: %w", err)
	}
	defer rows.Close()

	var feedbacks []types.SubmitGenuityFeedbackRequest
	for rows.Next() {
		feedback := types.SubmitGenuityFeedbackRequest{
			Original: types.OriginalGenuity{},
			Feedback: types.UserFeedback{},
		}

		err := rows.Scan(
			&feedback.ReviewID,
			&feedback.OverallRating,
			&feedback.Pros,
			&feedback.Cons,
			&feedback.Original.IsGenuine,
			&feedback.Feedback.IsGenuine,
			&feedback.UserID,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to read row: %w", err)
		}
		feedbacks = append(feedbacks, feedback)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over results: %w", err)
	}

	return feedbacks, nil
}
