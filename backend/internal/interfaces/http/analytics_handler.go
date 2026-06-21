package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"lokalsicht/internal/application"
)

type AnalyticsHandler struct {
	insightService *application.InsightService
}

func NewAnalyticsHandler(insightService *application.InsightService) *AnalyticsHandler {
	return &AnalyticsHandler{insightService: insightService}
}

// Analytics returns performance data for a location.
func (h *AnalyticsHandler) Analytics(w http.ResponseWriter, r *http.Request) {
	locationID, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	days := 30
	if d := r.URL.Query().Get("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 && parsed <= 365 {
			days = parsed
		}
	}

	analytics, err := h.insightService.GetAnalytics(r.Context(), uint(locationID), days)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, analytics)
}

type CronHandler struct {
	insightService *application.InsightService
	notification   interface {
		CheckNewReviews() (map[string]interface{}, error)
	}
}

func NewCronHandler(insightService *application.InsightService) *CronHandler {
	return &CronHandler{insightService: insightService}
}

func (h *CronHandler) SyncInsights(w http.ResponseWriter, r *http.Request) {
	result, err := h.insightService.SyncAll(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":            "synced",
		"locations_checked": result.LocationsChecked,
		"synced":            result.Synced,
		"errors":            result.Errors,
		"checked_at":        result.CheckedAt.Format(time.RFC3339),
	})
}
