package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"lokalsicht/internal/application"
	"lokalsicht/internal/config"
	"lokalsicht/internal/domain/account"
	"lokalsicht/internal/domain/insight"
	"lokalsicht/internal/domain/location"
	"lokalsicht/internal/domain/notification"
	"lokalsicht/internal/domain/review"
	"lokalsicht/internal/infrastructure/ai"
	"lokalsicht/internal/infrastructure/gbp"
	"lokalsicht/internal/infrastructure/persistence"
	"lokalsicht/internal/infrastructure/stripe"
	httphandler "lokalsicht/internal/interfaces/http"
	"lokalsicht/internal/interfaces/middleware"
)

func main() {
	// Load .env from project root (or current dir for production)
	godotenv.Load("../.env")
	godotenv.Load(".env")

	cfg := config.Load()

	// Database
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	db.AutoMigrate(
		&account.Account{},
		&account.User{},
		&location.Location{},
		&location.GoogleProfile{},
		&review.Review{},
		&insight.InsightSnapshot{},
		&notification.Preference{},
	)

	// Persistence
	accountRepo := persistence.NewGormAccountRepo(db)
	reviewRepo := persistence.NewGormReviewRepo(db)

	// DeepSeek
	aiClient := ai.NewDeepSeekClient(cfg.DeepSeekAPIKey)

	// Handlers
	reviewHandler := httphandler.NewReviewHandler(reviewRepo, aiClient)
	stripeClient := stripe.NewClient(cfg.StripeSecretKey, cfg.StripeWebhookSecret, "", "", cfg.FrontendURL)
	billingHandler := httphandler.NewBillingHandler(db, stripeClient)

	// GBP Client
	redirectURI := cfg.BackendURL + "/api/gbp/callback"
	gbpClient := gbp.NewClient(cfg.GoogleClientID, cfg.GoogleClientSecret, redirectURI, cfg.FrontendURL, cfg.EncryptionKey)
	gbpHandler := httphandler.NewGBPHandler(gbpClient, db, cfg.FrontendURL)

	// Insight service + cron
	insightSvc := application.NewInsightService(db, gbpClient)
	analyticsHandler := httphandler.NewAnalyticsHandler(insightSvc)
	cronHandler := httphandler.NewCronHandler(insightSvc)

	// Router
	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{cfg.FrontendURL},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}))

	// Health check
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// GBP OAuth callback (no auth required — Google redirects here)
	r.Get("/api/gbp/callback", gbpHandler.Callback)

	// Public
	r.Post("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"login endpoint ready"}`))
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(cfg.NextAuthSecret, accountRepo))

		r.Get("/api/me", httphandler.MeHandler)

		// GBP connect
		r.Post("/api/gbp/connect", gbpHandler.Connect)

		// Locations
		locHandler := httphandler.NewLocationHandler(db)
		r.Get("/api/locations", locHandler.List)
		r.Put("/api/locations/{id}", locHandler.Update)

		// Reviews
		r.Get("/api/locations/{id}/reviews", reviewHandler.List)
		r.Post("/api/locations/{id}/reviews/{rid}/generate", reviewHandler.GenerateReply)
		r.Post("/api/locations/{id}/reviews/{rid}/reply", reviewHandler.Reply)

		// Billing
		r.Post("/api/billing/create-checkout", billingHandler.CreateCheckout)

		// Analytics
		r.Get("/api/locations/{id}/analytics", analyticsHandler.Analytics)
	})

	// Stripe webhook (no auth — Stripe signs it)
	r.Post("/api/billing/webhook", billingHandler.Webhook)

	// Internal cron endpoints
	r.Group(func(r chi.Router) {
		r.Post("/api/internal/cron/check-reviews", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		r.Post("/api/internal/cron/sync-insights", cronHandler.SyncInsights)
	})

	slog.Info("server starting", "port", cfg.Port, "redirect_uri", redirectURI)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
