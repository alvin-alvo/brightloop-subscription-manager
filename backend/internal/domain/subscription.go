package domain

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionStatus is the lifecycle stage of a subscription.
type SubscriptionStatus string

const (
	SubscriptionDraft     SubscriptionStatus = "draft"
	SubscriptionQuotation SubscriptionStatus = "quotation"
	SubscriptionConfirmed SubscriptionStatus = "confirmed"
	SubscriptionActive    SubscriptionStatus = "active"
	SubscriptionClosed    SubscriptionStatus = "closed"
)

// Valid reports whether the status is supported.
func (s SubscriptionStatus) Valid() bool {
	switch s {
	case SubscriptionDraft, SubscriptionQuotation, SubscriptionConfirmed, SubscriptionActive, SubscriptionClosed:
		return true
	default:
		return false
	}
}

// Subscription is a customer's agreement to pay for a plan over time.
type Subscription struct {
	ID           uuid.UUID
	Number       string
	CustomerID   uuid.UUID
	PlanID       uuid.UUID
	Status       SubscriptionStatus
	StartDate    time.Time
	ExpiryDate   *time.Time
	PaymentTerms string
	Notes        string
	Lines        []SubscriptionLine
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// SubscriptionLine is one priced product line on a subscription.
type SubscriptionLine struct {
	ID             uuid.UUID
	SubscriptionID uuid.UUID
	ProductID      uuid.UUID
	Quantity       int
	UnitPriceCents int64
	TaxRateID      *uuid.UUID
	AmountCents    int64
}
