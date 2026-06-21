package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"lokalsicht/internal/domain/account"
	"lokalsicht/internal/domain/insight"
	"lokalsicht/internal/domain/location"
	"lokalsicht/internal/domain/notification"
	"lokalsicht/internal/domain/review"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5174"
	}

	// Database
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{
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

	// Router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{os.Getenv("FRONTEND_URL")},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}))

	// Health check
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Internal cron endpoints
	r.Group(func(r chi.Router) {
		r.Post("/api/internal/cron/check-reviews", func(w http.ResponseWriter, r *http.Request) {
			// TODO: Implement notification service
			w.WriteHeader(http.StatusOK)
		})
		r.Post("/api/internal/cron/sync-insights", func(w http.ResponseWriter, r *http.Request) {
			// TODO: Implement analytics service
			w.WriteHeader(http.StatusOK)
		})
	})

	slog.Info("server starting", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
