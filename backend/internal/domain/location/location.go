package location

import (
	"context"
	"time"
)

// Location represents a physical business location.
type Location struct {
	ID            uint   `gorm:"primaryKey"`
	AccountID     uint
	Name          string `gorm:"not null"`
	Address       string
	Phone         string
	Website       string
	OpeningHours  string
	Description   string
	GoogleProfile *GoogleProfile `gorm:"foreignKey:LocationID"`
	LastReviewCheck *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// GoogleProfile is a value object linking to a Google Business Profile.
type GoogleProfile struct {
	ID           uint   `gorm:"primaryKey"`
	LocationID   uint   `gorm:"uniqueIndex"`
	GoogleID     string `gorm:"uniqueIndex"` // Google Place ID (e.g., "accounts/123/locations/456")
	RefreshToken string // AES-256-GCM encrypted, Base64 encoded
	Scope        string
	TokenType    string
	Expiry       time.Time
	LastSyncAt   *time.Time
}

// LocationRepository defines the port for location persistence.
type LocationRepository interface {
	FindByID(ctx context.Context, id uint) (*Location, error)
	FindByAccount(ctx context.Context, accountID uint) ([]Location, error)
	CountByAccount(ctx context.Context, accountID uint) (int, error)
	Create(ctx context.Context, location *Location) error
	Update(ctx context.Context, location *Location) error
	Delete(ctx context.Context, id uint) error
	UpdateGoogleProfile(ctx context.Context, gp *GoogleProfile) error
	FindGoogleProfileByLocation(ctx context.Context, locationID uint) (*GoogleProfile, error)
}
