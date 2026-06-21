package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"lokalsicht/internal/domain/account"
	"lokalsicht/internal/infrastructure/stripe"
	"lokalsicht/internal/interfaces/middleware"
)

type BillingHandler struct {
	db           *gorm.DB
	stripeClient *stripe.Client
}

func NewBillingHandler(db *gorm.DB, stripeClient *stripe.Client) *BillingHandler {
	return &BillingHandler{db: db, stripeClient: stripeClient}
}

type CreateCheckoutRequest struct {
	Plan string `json:"plan"` // "standard" | "pro"
}

type CreateCheckoutResponse struct {
	URL string `json:"url"`
}

// CreateCheckout creates a Stripe checkout session for upgrading the plan.
func (h *BillingHandler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	var req CreateCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	url, err := h.stripeClient.CreateCheckoutSession(r.Context(), stripe.CreateCheckoutParams{
		AccountID: user.AccountID,
		Plan:      req.Plan,
		Email:     user.Email,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, CreateCheckoutResponse{URL: url})
}

// Webhook handles Stripe webhook events.
func (h *BillingHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "cannot read body"})
		return
	}

	event, err := h.stripeClient.HandleWebhook(payload, r.Header.Get("Stripe-Signature"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid signature"})
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		meta := event.Data.Object["metadata"].(map[string]interface{})
		plan := meta["plan"].(string)
		accountIDStr := meta["account_id"].(string)
		accountID, _ := strconv.ParseUint(accountIDStr, 10, 64)

		customerID, _ := event.Data.Object["customer"].(string)
		subscriptionID, _ := event.Data.Object["subscription"].(string)

		var acc account.Account
		if h.db.First(&acc, accountID).Error == nil {
			acc.Plan = account.Plan(plan)
			acc.StripeCustomerID = &customerID
			acc.StripeSubscriptionID = &subscriptionID
			h.db.Save(&acc)
		}

	case "customer.subscription.deleted":
		subscriptionID, _ := event.Data.Object["id"].(string)
		var acc account.Account
		if h.db.Where("stripe_subscription_id = ?", subscriptionID).First(&acc).Error == nil {
			acc.Plan = account.PlanBasic
			acc.StripeSubscriptionID = nil
			h.db.Save(&acc)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
