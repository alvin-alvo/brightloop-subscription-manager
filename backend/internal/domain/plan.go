package domain

import (
	"time"

	"github.com/google/uuid"
)

// BillingInterval is how often a plan bills.
type BillingInterval string

const (
	BillingDaily   BillingInterval = "daily"
	BillingWeekly  BillingInterval = "weekly"
	BillingMonthly BillingInterval = "monthly"
	BillingYearly  BillingInterval = "yearly"
)

// Valid reports whether the interval is supported.
func (b BillingInterval) Valid() bool {
	switch b {
	case BillingDaily, BillingWeekly, BillingMonthly, BillingYearly:
		return true
	default:
		return false
	}
}

// Plan defines the price and billing behaviour of a product.
type Plan struct {
	ID              uuid.UUID
	ProductID       uuid.UUID
	Name            string
	PriceCents      int64
	BillingInterval BillingInterval
	MinQuantity     int
	AvailableFrom   *time.Time
	AvailableTo     *time.Time
	AutoClose       bool
	AllowEarlyClose bool
	AllowPause      bool
	AllowRenewal    bool
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
