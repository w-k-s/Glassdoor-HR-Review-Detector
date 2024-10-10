package server

import (
	"encoding/json"
	"net/http"

	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/types"
)

func (s *Server) checkReviewsGenuity(w http.ResponseWriter, req *http.Request){
	var request types.CheckReviewsGenuityRequest

    if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    response, err := s.genuityService.CheckReviewGenuity(req.Context(), request)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")

    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
        return
    }
}

func (s *Server) submitGenuityFeedback(w http.ResponseWriter, req *http.Request){
	w.Write([]byte("welcome"))
}