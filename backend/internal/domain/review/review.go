package review

import (
	"context"
	"time"
)

// Review represents a customer review.
type Review struct {
	ID         uint      `gorm:"primaryKey"`
	LocationID uint
	GoogleID   string    `gorm:"uniqueIndex"`
	Author     string
	Rating     int       // 1–5 stars
	Text       string
	Language   string    // detected language (de, fr, it, en)
	ReplyText  *string   // the business's published reply (nil if unanswered)
	RepliedAt  *time.Time
	CreatedAt  time.Time
}

// Reply is a value object representing a business reply to a review.
type Reply struct {
	Text        string    // the final reply text sent to Google
	GeneratedBy string    // "ai" or "manual"
	Variant     int       // which AI variant was chosen (0 if manual)
	CreatedAt   time.Time
}

// ReviewReceived is a domain event emitted when a new review is discovered.
type ReviewReceived struct {
	Review    Review
	LocationID uint
	AccountID  uint
	Timestamp time.Time
}

// ReviewRepository defines the port for review persistence.
type ReviewRepository interface {
	FindByID(ctx context.Context, id uint) (*Review, error)
	FindByLocation(ctx context.Context, locationID uint, since *time.Time) ([]Review, error)
	FindByGoogleID(ctx context.Context, googleID string) (*Review, error)
	Save(ctx context.Context, review *Review) error
	SaveReply(ctx context.Context, reviewID uint, reply *Reply) error
	CountByLocation(ctx context.Context, locationID uint) (int64, error)
}

// AIReplyGenerator defines the port for AI-generated reply suggestions.
type AIReplyGenerator interface {
	GenerateReply(ctx context.Context, reviewText string, language string, businessContext string) ([]string, error)
}

// GBPClient defines the port for Google Business Profile API operations.
type GBPClient interface {
	ListReviews(ctx context.Context, googleID string, since *time.Time) ([]Review, error)
	ReplyToReview(ctx context.Context, googleReviewID string, replyText string) error
}
