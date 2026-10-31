package domain

import (
	"time"

	"github.com/google/uuid"
)

// DiscountKind selects which value field a discount uses.
type DiscountKind string

const (
	DiscountFixed      DiscountKind = "fixed"
	DiscountPercentage DiscountKind = "percentage"
)

// DiscountAppliesTo selects what a discount can be attached to.
type DiscountAppliesTo string

const (
	DiscountAppliesToProduct      DiscountAppliesTo = "product"
	DiscountAppliesToSubscription DiscountAppliesTo = "subscription"
)

// Discount is an admin-created pricing rule. Fixed discounts use
// FixedAmountCents; percentage discounts use PercentageBasisPoints.
type Discount struct {
	ID                    uuid.UUID
	Name                  string
	Kind                  DiscountKind
	FixedAmountCents      int64
	PercentageBasisPoints int
	MinPurchaseCents      int64
	MinQuantity           int
	StartsAt              *time.Time
	EndsAt                *time.Time
	MaxUses               *int
	UsedCount             int
	AppliesTo             DiscountAppliesTo
	ProductID             *uuid.UUID
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
