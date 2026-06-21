package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"lokalsicht/internal/domain/review"
	"lokalsicht/internal/interfaces/middleware"
)

type ReviewHandler struct {
	reviewRepo  review.ReviewRepository
	aiGenerator review.AIReplyGenerator
}

func NewReviewHandler(reviewRepo review.ReviewRepository, aiGenerator review.AIReplyGenerator) *ReviewHandler {
	return &ReviewHandler{reviewRepo: reviewRepo, aiGenerator: aiGenerator}
}

type ReviewResponse struct {
	ID        uint    `json:"id"`
	GoogleID  string  `json:"googleId"`
	Author    string  `json:"author"`
	Rating    int     `json:"rating"`
	Text      string  `json:"text"`
	Language  string  `json:"language"`
	ReplyText *string `json:"replyText,omitempty"`
	RepliedAt *string `json:"repliedAt,omitempty"`
	CreatedAt string  `json:"createdAt"`
}

// List returns reviews for a location.
func (h *ReviewHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	locationID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid location id"})
		return
	}

	reviews, err := h.reviewRepo.FindByLocation(r.Context(), uint(locationID), nil)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := make([]ReviewResponse, 0, len(reviews))
	for _, rev := range reviews {
		resp = append(resp, toReviewResponse(rev))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type GenerateReplyRequest struct {
	BusinessContext string `json:"businessContext"`
}

type GenerateReplyResponse struct {
	Variants []string `json:"variants"`
	Language string   `json:"language"`
}

// GenerateReply generates AI reply variants for a review.
func (h *ReviewHandler) GenerateReply(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	reviewID, err := strconv.ParseUint(chi.URLParam(r, "rid"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid review id"})
		return
	}

	rev, err := h.reviewRepo.FindByID(r.Context(), uint(reviewID))
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "review not found"})
		return
	}

	var req GenerateReplyRequest
	if r.Header.Get("Content-Type") == "application/json" {
		json.NewDecoder(r.Body).Decode(&req)
	}

	variants, err := h.aiGenerator.GenerateReply(r.Context(), rev.Text, rev.Language, req.BusinessContext)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	resp := GenerateReplyResponse{Variants: variants, Language: rev.Language}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type ReplyRequest struct {
	Text string `json:"text"`
}

// Reply sends a reply to a review.
func (h *ReviewHandler) Reply(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	reviewID, err := strconv.ParseUint(chi.URLParam(r, "rid"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid review id"})
		return
	}

	var req ReplyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	reply := &review.Reply{
		Text:        req.Text,
		GeneratedBy: "manual",
		CreatedAt:   time.Now(),
	}

	if err := h.reviewRepo.SaveReply(r.Context(), uint(reviewID), reply); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "replied"})
}

func toReviewResponse(rev review.Review) ReviewResponse {
	r := ReviewResponse{
		ID:        rev.ID,
		GoogleID:  rev.GoogleID,
		Author:    rev.Author,
		Rating:    rev.Rating,
		Text:      rev.Text,
		Language:  rev.Language,
		CreatedAt: rev.CreatedAt.Format(time.RFC3339),
	}
	if rev.ReplyText != nil {
		r.ReplyText = rev.ReplyText
		rt := rev.RepliedAt.Format(time.RFC3339)
		r.RepliedAt = &rt
	}
	return r
}
