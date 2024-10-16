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
	UserID        string          `json:"userId" csv:"-"`
	ReviewID      string          `json:"reviewId" csv:"Review ID"`
	OverallRating float32         `json:"rating" csv:"Rating"`
	Pros          string          `json:"pros" csv:"Pros"`
	Cons          string          `json:"cons" csv:"Cons"`
	Original      OriginalGenuity `json:"original" csv:"Prediction"`
	Feedback      UserFeedback    `json:"feedback" csv:"Feedback"`
}

type OriginalGenuity struct {
	IsGenuine bool `json:"isGenuine"`
}

type UserFeedback struct {
	IsGenuine bool `json:"isGenuine"`
}
