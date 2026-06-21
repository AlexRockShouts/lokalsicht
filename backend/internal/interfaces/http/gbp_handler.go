package http

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"

	"lokalsicht/internal/domain/location"
	"lokalsicht/internal/infrastructure/gbp"
	"lokalsicht/internal/interfaces/middleware"
)

type GBPHandler struct {
	client      *gbp.Client
	db          *gorm.DB
	frontendURL string
}

type oauthState struct {
	State     string `gorm:"primaryKey"`
	UserID    uint   `gorm:"index"`
	ExpiresAt time.Time
}

func NewGBPHandler(client *gbp.Client, db *gorm.DB, frontendURL string) *GBPHandler {
	db.AutoMigrate(&oauthState{})
	return &GBPHandler{client: client, db: db, frontendURL: frontendURL}
}

type ConnectResponse struct {
	AuthURL string `json:"authUrl"`
}

// Connect initiates GBP OAuth and returns the Google authorization URL.
func (h *GBPHandler) Connect(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	b := make([]byte, 16)
	rand.Read(b)
	state := hex.EncodeToString(b)

	h.db.Create(&oauthState{
		State:     state,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	})

	resp := ConnectResponse{AuthURL: h.client.AuthCodeURL(state)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type CallbackResponse struct {
	Status    string `json:"status"`
	Locations int    `json:"locations"`
}

// Callback handles the Google OAuth redirect, stores tokens, and syncs locations.
func (h *GBPHandler) Callback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	code := r.URL.Query().Get("code")
	if state == "" || code == "" {
		http.Redirect(w, r, h.frontendURL+"/onboarding?error=missing_params", http.StatusFound)
		return
	}

	// Verify state
	var st oauthState
	if err := h.db.Where("state = ? AND expires_at > ?", state, time.Now()).First(&st).Error; err != nil {
		http.Redirect(w, r, h.frontendURL+"/onboarding?error=invalid_state", http.StatusFound)
		return
	}
	h.db.Delete(&st)

	// Exchange code for token
	token, err := h.client.Exchange(r.Context(), code)
	if err != nil {
		http.Redirect(w, r, h.frontendURL+"/onboarding?error=token_exchange_failed", http.StatusFound)
		return
	}

	// Encrypt refresh token
	encrypted, err := h.client.EncryptToken(token.RefreshToken)
	if err != nil {
		http.Redirect(w, r, h.frontendURL+"/onboarding?error=encryption_failed", http.StatusFound)
		return
	}

	// List GBP accounts
	accounts, err := h.client.ListAccounts(r.Context(), token)
	if err != nil || len(accounts) == 0 {
		http.Redirect(w, r, h.frontendURL+"/onboarding?error=no_gbp_account", http.StatusFound)
		return
	}

	// Sync locations for the first account
	synced := 0
	for _, acc := range accounts {
		locations, err := h.client.ListLocations(r.Context(), token, acc.Name)
		if err != nil {
			continue
		}
		for _, loc := range locations {
			// Upsert location + google profile
			var existing location.Location
			// Try to find by GoogleProfile.GoogleID
			result := h.db.Joins("JOIN google_profiles ON google_profiles.location_id = locations.id").
				Where("google_profiles.google_id = ?", loc.Name).
				First(&existing)
			if result.Error != nil {
				// Create new
				newLoc := location.Location{
					AccountID: st.UserID,
					Name:      loc.Title,
					Address:   loc.Address,
					Phone:     loc.Phone,
					Website:   loc.Website,
					GoogleProfile: &location.GoogleProfile{
						GoogleID:     loc.Name,
						RefreshToken: encrypted,
						Scope:        "https://www.googleapis.com/auth/business.manage",
						TokenType:    token.TokenType,
						Expiry:       token.Expiry,
					},
				}
				h.db.Create(&newLoc)
			} else {
				// Update existing
				existing.Name = loc.Title
				if loc.Address != "" {
					existing.Address = loc.Address
				}
				h.db.Save(&existing)
				// Update google profile token
				h.db.Model(&location.GoogleProfile{}).
					Where("location_id = ?", existing.ID).
					Updates(map[string]interface{}{
						"refresh_token": encrypted,
						"expiry":        token.Expiry,
					})
			}
			synced++
		}
	}

	http.Redirect(w, r, fmt.Sprintf("%s/onboarding?success=true&locations=%d", h.frontendURL, synced), http.StatusFound)
}

func writeError(w http.ResponseWriter, status int, detail string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error":"%s"}`, detail)
}
