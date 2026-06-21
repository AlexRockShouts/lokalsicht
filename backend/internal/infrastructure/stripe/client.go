package stripe

import (
	"context"
	"fmt"

	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/checkout/session"
	"github.com/stripe/stripe-go/v83/webhook"
)

type Client struct {
	secretKey       string
	webhookSecret   string
	standardPriceID string
	proPriceID      string
	frontendURL     string
}

func NewClient(secretKey, webhookSecret, standardPriceID, proPriceID, frontendURL string) *Client {
	stripe.Key = secretKey
	return &Client{
		secretKey:       secretKey,
		webhookSecret:   webhookSecret,
		standardPriceID: standardPriceID,
		proPriceID:      proPriceID,
		frontendURL:     frontendURL,
	}
}

type CreateCheckoutParams struct {
	AccountID uint
	Plan      string // "standard" | "pro"
	Email     string
}

func (c *Client) CreateCheckoutSession(ctx context.Context, p CreateCheckoutParams) (string, error) {
	priceID := c.standardPriceID
	if p.Plan == "pro" {
		priceID = c.proPriceID
	}

	params := &stripe.CheckoutSessionParams{
		Mode:              stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		CustomerEmail:     stripe.String(p.Email),
		ClientReferenceID: stripe.String(fmt.Sprintf("account_%d", p.AccountID)),
		Metadata: map[string]string{
			"plan":       p.Plan,
			"account_id": fmt.Sprintf("%d", p.AccountID),
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(c.frontendURL + "/settings?checkout=success"),
		CancelURL:  stripe.String(c.frontendURL + "/settings?checkout=cancelled"),
	}

	s, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("create checkout: %w", err)
	}
	return s.URL, nil
}

func (c *Client) HandleWebhook(payload []byte, signature string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, signature, c.webhookSecret)
}
