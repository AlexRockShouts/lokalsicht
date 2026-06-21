package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"lokalsicht/internal/domain/location"
	"lokalsicht/internal/interfaces/middleware"
)

type LocationHandler struct {
	db *gorm.DB
}

func NewLocationHandler(db *gorm.DB) *LocationHandler {
	return &LocationHandler{db: db}
}

type LocationResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	GoogleID string `json:"googleId,omitempty"`
}

func (h *LocationHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var locations []location.Location
	h.db.Where("account_id = ?", user.AccountID).Preload("GoogleProfile").Find(&locations)

	var resp []LocationResponse
	for _, loc := range locations {
		lr := LocationResponse{ID: loc.ID, Name: loc.Name}
		if loc.GoogleProfile != nil {
			lr.GoogleID = loc.GoogleProfile.GoogleID
		}
		resp = append(resp, lr)
	}
	if resp == nil {
		resp = []LocationResponse{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type UpdateLocationRequest struct {
	Name         *string `json:"name,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	Website      *string `json:"website,omitempty"`
	Description  *string `json:"description,omitempty"`
	OpeningHours *string `json:"opening_hours,omitempty"`
}

func (h *LocationHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var loc location.Location
	if err := h.db.Where("id = ? AND account_id = ?", id, user.AccountID).First(&loc).Error; err != nil {
		writeError(w, http.StatusNotFound, "location not found")
		return
	}

	var req UpdateLocationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Website != nil {
		updates["website"] = *req.Website
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.OpeningHours != nil {
		updates["opening_hours"] = *req.OpeningHours
	}

	if len(updates) > 0 {
		h.db.Model(&loc).Updates(updates)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}
