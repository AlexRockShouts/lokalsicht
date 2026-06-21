package application

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"lokalsicht/internal/domain/insight"
	"lokalsicht/internal/domain/location"
	"lokalsicht/internal/infrastructure/gbp"
)

type InsightService struct {
	db        *gorm.DB
	gbpClient *gbp.Client
}

func NewInsightService(db *gorm.DB, gbpClient *gbp.Client) *InsightService {
	return &InsightService{db: db, gbpClient: gbpClient}
}

func (s *InsightService) SyncAll(ctx context.Context) (*insight.SyncResult, error) {
	var locations []location.Location
	s.db.Preload("GoogleProfile").Where("google_profile_id IS NOT NULL").Find(&locations)

	result := &insight.SyncResult{CheckedAt: time.Now()}

	for _, loc := range locations {
		if loc.GoogleProfile == nil || loc.GoogleProfile.GoogleID == "" {
			continue
		}
		result.LocationsChecked++

		// Get last sync date
		var lastDate *time.Time
		var lastSnap insight.InsightSnapshot
		err := s.db.Where("location_id = ?", loc.ID).Order("date DESC").First(&lastSnap).Error
		if err == nil {
			lastDate = &lastSnap.Date
		}
		if lastDate == nil {
			t := time.Now().AddDate(0, 0, -30)
			lastDate = &t
		}

		// Fetch insights from Google
		startDate := lastDate.Add(24 * time.Hour)
		endDate := time.Now()

		dailyInsights, err := s.gbpClient.GetInsights(ctx, loc.GoogleProfile.GoogleID, startDate, endDate)
		if err != nil {
			slog.Warn("insight sync failed for location", "location", loc.ID, "error", err)
			result.Errors++
			continue
		}

		// Save snapshots
		for _, ins := range dailyInsights {
			snap := insight.InsightSnapshot{
				LocationID: loc.ID,
				Date:       ins.Date,
				Views:      ins.Views,
				Clicks:     ins.Clicks,
				Calls:      ins.Calls,
				Directions: ins.Directions,
			}
			var existing insight.InsightSnapshot
			err := s.db.Where("location_id = ? AND date = ?", loc.ID, ins.Date).First(&existing).Error
			if err == nil {
				continue
			}
			s.db.Create(&snap)
			result.Synced++
		}
	}

	return result, nil
}

func (s *InsightService) GetAnalytics(ctx context.Context, locationID uint, days int) (*insight.Analytics, error) {
	from := time.Now().AddDate(0, 0, -days)
	var snapshots []insight.InsightSnapshot
	s.db.Where("location_id = ? AND date >= ?", locationID, from).Order("date ASC").Find(&snapshots)

	analytics := &insight.Analytics{
		Days:      days,
		Snapshots: make([]insight.DailyPoint, 0, len(snapshots)),
	}

	for _, snap := range snapshots {
		analytics.TotalViews += snap.Views
		analytics.TotalClicks += snap.Clicks
		analytics.TotalCalls += snap.Calls
		analytics.TotalDirections += snap.Directions
		analytics.Snapshots = append(analytics.Snapshots, insight.DailyPoint{
			Date:       snap.Date.Format("2006-01-02"),
			Views:      snap.Views,
			Clicks:     snap.Clicks,
			Calls:      snap.Calls,
			Directions: snap.Directions,
		})
	}

	return analytics, nil
}
