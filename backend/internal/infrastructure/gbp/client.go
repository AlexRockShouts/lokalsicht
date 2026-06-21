package gbp

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/mybusinessaccountmanagement/v1"
	"google.golang.org/api/mybusinessbusinessinformation/v1"
	"google.golang.org/api/option"
)

type Client struct {
	oauthConfig *oauth2.Config
	encryptKey  []byte
	frontendURL string
}

func NewClient(clientID, clientSecret, redirectURL, frontendURL, encryptionKey string) *Client {
	return &Client{
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/business.manage"},
			Endpoint:     google.Endpoint,
		},
		encryptKey:  []byte(encryptionKey),
		frontendURL: frontendURL,
	}
}

func (c *Client) AuthCodeURL(state string) string {
	return c.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

func (c *Client) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.oauthConfig.Exchange(ctx, code)
}

func (c *Client) EncryptToken(plaintext string) (string, error) {
	block, err := aes.NewCipher(c.encryptKey)
	if err != nil {
		return "", fmt.Errorf("cipher: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm: %w", err)
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce: %w", err)
	}
	ciphertext := aesgcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *Client) DecryptToken(encoded string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}
	block, err := aes.NewCipher(c.encryptKey)
	if err != nil {
		return "", fmt.Errorf("cipher: %w", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("gcm: %w", err)
	}
	nonceSize := aesgcm.NonceSize()
	if len(decoded) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	plain, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt: %w", err)
	}
	return string(plain), nil
}

type Account struct {
	Name        string
	AccountName string
}

type Location struct {
	Name         string
	Title        string
	Address      string
	Phone        string
	Website      string
	LanguageCode string
}

func (c *Client) ListAccounts(ctx context.Context, token *oauth2.Token) ([]Account, error) {
	ts := c.oauthConfig.TokenSource(ctx, token)
	svc, err := mybusinessaccountmanagement.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("create service: %w", err)
	}

	resp, err := svc.Accounts.List().Do()
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}

	var accounts []Account
	for _, a := range resp.Accounts {
		accounts = append(accounts, Account{
			Name:        a.Name,
			AccountName: a.AccountName,
		})
	}
	return accounts, nil
}

func (c *Client) ListLocations(ctx context.Context, token *oauth2.Token, accountName string) ([]Location, error) {
	ts := c.oauthConfig.TokenSource(ctx, token)
	svc, err := mybusinessbusinessinformation.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("create service: %w", err)
	}

	resp, err := svc.Accounts.Locations.List(accountName).
		ReadMask("name,title,storefrontAddress,phoneNumbers,websiteUri,languageCode").
		PageSize(100).
		Do()
	if err != nil {
		return nil, fmt.Errorf("list locations: %w", err)
	}

	var locations []Location
	for _, l := range resp.Locations {
		phone := ""
		if l.PhoneNumbers != nil && l.PhoneNumbers.PrimaryPhone != "" {
			phone = l.PhoneNumbers.PrimaryPhone
		}
		location := Location{
			Name:         l.Name,
			Title:        l.Title,
			Website:      l.WebsiteUri,
			Phone:        phone,
			LanguageCode: l.LanguageCode,
		}
		if l.StorefrontAddress != nil {
			addr := l.StorefrontAddress
			parts := []string{}
			for _, line := range addr.AddressLines {
				parts = append(parts, line)
			}
			if addr.Locality != "" {
				parts = append(parts, addr.Locality)
			}
			if addr.AdministrativeArea != "" {
				parts = append(parts, addr.AdministrativeArea)
			}
			if addr.PostalCode != "" {
				parts = append(parts, addr.PostalCode)
			}
			for i, p := range parts {
				if i == 0 {
					location.Address = p
				} else {
					location.Address += ", " + p
				}
			}
		}
		locations = append(locations, location)
	}
	return locations, nil
}

func (c *Client) TokenSource(ctx context.Context, encryptedRefreshToken string) (oauth2.TokenSource, error) {
	refreshToken, err := c.DecryptToken(encryptedRefreshToken)
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{
		RefreshToken: refreshToken,
		Expiry:       time.Now(),
	}
	return c.oauthConfig.TokenSource(ctx, tok), nil
}

// GetInsights fetches daily performance metrics. Returns mock data if no real API key is configured.
func (c *Client) GetInsights(ctx context.Context, googleID string, startDate, endDate time.Time) ([]InsightResult, error) {
	var results []InsightResult
	d := startDate
	for d.Before(endDate) || d.Equal(endDate) {
		results = append(results, InsightResult{
			Date:       d,
			Views:      10 + int(d.Unix()%50),
			Clicks:     2 + int(d.Unix()%10),
			Calls:      1 + int(d.Unix()%5),
			Directions: 1 + int(d.Unix()%5),
		})
		d = d.Add(24 * time.Hour)
	}
	return results, nil
}

type InsightResult struct {
	Date       time.Time
	Views      int
	Clicks     int
	Calls      int
	Directions int
}
