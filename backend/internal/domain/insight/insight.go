package insight

import (
	"context"
	"time"
)

type InsightSnapshot struct {
	ID         uint      `gorm:"primaryKey"`
	LocationID uint      `gorm:"uniqueIndex:idx_loc_date"`
	Date       time.Time `gorm:"uniqueIndex:idx_loc_date;type:date"`
	Views      int
	Clicks     int
	Calls      int
	Directions int
	CreatedAt  time.Time
}

type SyncResult struct {
	LocationsChecked int
	Synced           int
	Errors           int
	CheckedAt        time.Time
}

type DailyPoint struct {
	Date       string `json:"date"`
	Views      int    `json:"views"`
	Clicks     int    `json:"clicks"`
	Calls      int    `json:"calls"`
	Directions int    `json:"directions"`
}

type Analytics struct {
	Days            int          `json:"days"`
	TotalViews      int          `json:"totalViews"`
	TotalClicks     int          `json:"totalClicks"`
	TotalCalls      int          `json:"totalCalls"`
	TotalDirections int          `json:"totalDirections"`
	Snapshots       []DailyPoint `json:"snapshots"`
}

type InsightRepository interface {
	Save(ctx context.Context, snapshot *InsightSnapshot) error
	FindByLocation(ctx context.Context, locationID uint, from, to time.Time) ([]InsightSnapshot, error)
	LastSyncDate(ctx context.Context, locationID uint) (*time.Time, error)
}

type InsightsClient interface {
	GetInsights(ctx context.Context, googleID string, startDate, endDate time.Time) ([]InsightSnapshot, error)
}
