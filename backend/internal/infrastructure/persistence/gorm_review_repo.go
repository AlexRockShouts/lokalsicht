package persistence

import (
	"context"
	"time"

	"gorm.io/gorm"

	"lokalsicht/internal/domain/review"
)

type GormReviewRepo struct {
	db *gorm.DB
}

func NewGormReviewRepo(db *gorm.DB) *GormReviewRepo {
	return &GormReviewRepo{db: db}
}

func (r *GormReviewRepo) FindByID(ctx context.Context, id uint) (*review.Review, error) {
	var rev review.Review
	if err := r.db.WithContext(ctx).First(&rev, id).Error; err != nil {
		return nil, err
	}
	return &rev, nil
}

func (r *GormReviewRepo) FindByLocation(ctx context.Context, locationID uint, since *time.Time) ([]review.Review, error) {
	var reviews []review.Review
	query := r.db.WithContext(ctx).Where("location_id = ?", locationID).Order("created_at DESC")
	if since != nil {
		query = query.Where("created_at > ?", since)
	}
	if err := query.Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *GormReviewRepo) FindByGoogleID(ctx context.Context, googleID string) (*review.Review, error) {
	var rev review.Review
	err := r.db.WithContext(ctx).Where("google_id = ?", googleID).First(&rev).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rev, nil
}

func (r *GormReviewRepo) Save(ctx context.Context, rev *review.Review) error {
	if rev.ID != 0 {
		return r.db.WithContext(ctx).Save(rev).Error
	}
	// Upsert by GoogleID
	var existing review.Review
	err := r.db.WithContext(ctx).Where("google_id = ?", rev.GoogleID).First(&existing).Error
	if err == nil {
		rev.ID = existing.ID
		return r.db.WithContext(ctx).Save(rev).Error
	}
	return r.db.WithContext(ctx).Create(rev).Error
}

func (r *GormReviewRepo) SaveReply(ctx context.Context, reviewID uint, reply *review.Reply) error {
	return r.db.WithContext(ctx).Model(&review.Review{}).Where("id = ?", reviewID).Updates(map[string]interface{}{
		"reply_text": reply.Text,
		"replied_at": reply.CreatedAt,
	}).Error
}

func (r *GormReviewRepo) CountByLocation(ctx context.Context, locationID uint) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&review.Review{}).Where("location_id = ?", locationID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
