package config

import "os"

type Config struct {
	Port                string
	DatabaseURL         string
	BackendURL          string
	NextAuthSecret      string
	GoogleClientID      string
	GoogleClientSecret  string
	FrontendURL         string
	DeepSeekAPIKey      string
	ResendAPIKey        string
	EncryptionKey       string
	CronAPIKey          string
	StripeSecretKey     string
	StripeWebhookSecret string
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "5174"),
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		BackendURL:          getEnv("BACKEND_URL", "http://localhost:5174"),
		NextAuthSecret:      os.Getenv("NEXTAUTH_SECRET"),
		GoogleClientID:      os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:  os.Getenv("GOOGLE_CLIENT_SECRET"),
		FrontendURL:         getEnv("FRONTEND_URL", "http://localhost:3000"),
		DeepSeekAPIKey:      os.Getenv("DEEPSEEK_API_KEY"),
		ResendAPIKey:        os.Getenv("RESEND_API_KEY"),
		EncryptionKey:       os.Getenv("ENCRYPTION_KEY"),
		CronAPIKey:          os.Getenv("CRON_API_KEY"),
		StripeSecretKey:     os.Getenv("STRIPE_SECRET_KEY"),
		StripeWebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
