package notification

import (
	"context"
	"time"
)

// Preference controls notification delivery per user.
type Preference struct {
	ID                     uint   `gorm:"primaryKey"`
	UserID                 uint   `gorm:"uniqueIndex"`
	NewReviewEmail         bool   `gorm:"default:true"`
	NewReviewDigest        bool   `gorm:"default:false"`
	AIResponseReadyEmail   bool   `gorm:"default:true"`
	MonthlyReportEmail     bool   `gorm:"default:true"`
}

// EmailClient defines the port for sending transactional emails.
type EmailClient interface {
	SendNewReviewEmail(ctx context.Context, to string, reviewText string, locationName string) error
	SendDigestEmail(ctx context.Context, to string, newReviews int) error
	SendAIResponseReady(ctx context.Context, to string, locationName string) error
}

// NotificationService orchestrates notification logic.
type NotificationService struct {
	EmailClient EmailClient
	Preferences func(ctx context.Context, userID uint) (*Preference, error)
}

// CheckNewReviewsResult is returned after a review check cycle.
type CheckNewReviewsResult struct {
	LocationsChecked int
	NewReviews       int
	EmailsSent       int
	Errors           int
	CheckedAt        time.Time
}
