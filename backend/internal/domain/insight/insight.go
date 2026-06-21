package insight

import (
	"context"
	"time"
)

// InsightSnapshot is a daily snapshot of Google Business Profile performance metrics.
type InsightSnapshot struct {
	ID         uint      `gorm:"primaryKey"`
	LocationID uint      `gorm:"uniqueIndex:idx_loc_date"`
	Date       time.Time `gorm:"uniqueIndex:idx_loc_date"` // Date of snapshot (midnight UTC)
	Views      int       // Search views
	Clicks     int       // Website clicks
	Calls      int       // Phone call clicks
	Directions int       // Direction requests
	CreatedAt  time.Time
}

// InsightRepository defines the port for insight persistence.
type InsightRepository interface {
	Save(ctx context.Context, snapshot *InsightSnapshot) error
	FindByLocation(ctx context.Context, locationID uint, from, to time.Time) ([]InsightSnapshot, error)
	LastSyncDate(ctx context.Context, locationID uint) (*time.Time, error)
}

// InsightsClient defines the port for fetching performance data from Google.
type InsightsClient interface {
	GetInsights(ctx context.Context, googleID string, startDate, endDate time.Time) (*InsightSnapshot, error)
}
