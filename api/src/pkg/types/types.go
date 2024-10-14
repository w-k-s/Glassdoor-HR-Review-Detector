package types

type CheckReviewsGenuityRequest struct {
	Reviews []Review `json:"reviews"`
}

type Review struct {
	ID     string  `json:"id"`
	Rating float64 `json:"rating"`
	Pros   string  `json:"pros"`
	Cons   string  `json:"cons"`
}

type CheckReviewsGenuityResponse struct {
	Results []GenuityResult `json:"results"`
}

type GenuityResult struct {
	ReviewID  string `json:"reviewId"`
	IsGenuine bool   `json:"isGenuine"`
}

type SubmitGenuityFeedbackRequest struct {
	UserID   string          `json:"userId"`
	ReviewID string          `json:"reviewId"`
	Pros     string          `json:"pros"`
	Cons     string          `json:"cons"`
	Original OriginalGenuity `json:"original"`
	Feedback UserFeedback    `json:"feedback"`
}

type OriginalGenuity struct {
	IsGenuine bool `json:"isGenuine"`
}

type UserFeedback struct {
	IsGenuine bool `json:"isGenuine"`
}
