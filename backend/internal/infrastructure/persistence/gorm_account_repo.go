package persistence

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"lokalsicht/internal/domain/account"
)

type GormAccountRepo struct {
	db *gorm.DB
}

func NewGormAccountRepo(db *gorm.DB) *GormAccountRepo {
	return &GormAccountRepo{db: db}
}

func (r *GormAccountRepo) FindByID(ctx context.Context, id uint) (*account.Account, error) {
	var a account.Account
	if err := r.db.WithContext(ctx).Preload("Users").First(&a, id).Error; err != nil {
		return nil, fmt.Errorf("account not found: %w", err)
	}
	return &a, nil
}

func (r *GormAccountRepo) FindByEmail(ctx context.Context, email string) (*account.User, error) {
	var u account.User
	err := r.db.WithContext(ctx).Where("email = ?", email).Preload("Account").First(&u).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("user lookup failed: %w", err)
	}
	return &u, nil
}

func (r *GormAccountRepo) Create(ctx context.Context, a *account.Account) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *GormAccountRepo) CreateUser(ctx context.Context, u *account.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *GormAccountRepo) Update(ctx context.Context, a *account.Account) error {
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *GormAccountRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&account.Account{}, id).Error
}
