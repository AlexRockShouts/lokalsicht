package http

import (
	"encoding/json"
	"net/http"

	"lokalsicht/internal/interfaces/middleware"
)

type MeResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AccountID uint   `json:"accountId"`
	Plan      string `json:"plan"`
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	resp := MeResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		AccountID: user.AccountID,
		Plan:      string(user.Account.Plan),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
