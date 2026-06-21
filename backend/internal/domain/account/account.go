package account

import (
	"context"
	"time"
)

// Plan represents a subscription tier.
type Plan string

const (
	PlanBasic    Plan = "basic"
	PlanStandard Plan = "standard"
	PlanPro      Plan = "pro"
)

// Account is the aggregate root for account management.
// Locations are referenced by foreign key (Location.AccountID) — queried via LocationRepository.
type Account struct {
	ID                   uint   `gorm:"primaryKey"`
	Name                 string `gorm:"not null"`
	Plan                 Plan   `gorm:"default:basic"`
	ResellerID           *uint
	StripeCustomerID     *string
	StripeSubscriptionID *string
	TrialEndsAt          *time.Time
	Users                []User
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// User belongs to an Account.
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	Name      string
	AccountID uint
	Account   Account `gorm:"constraint:OnDelete:CASCADE"`
	Role      string  `gorm:"default:owner"` // owner | member
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AccountRepository defines the port for account persistence.
type AccountRepository interface {
	FindByID(ctx context.Context, id uint) (*Account, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, account *Account) error
	CreateUser(ctx context.Context, user *User) error
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, id uint) error
}
