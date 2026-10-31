package domain

import (
	"time"

	"github.com/google/uuid"
)

// Product is a sellable item. Money is stored in minor units (cents) to avoid
// floating-point rounding errors.
type Product struct {
	ID              uuid.UUID
	Name            string
	Type            string
	SalesPriceCents int64
	CostPriceCents  int64
	Notes           string
	IsActive        bool
	Variants        []ProductVariant
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// ProductVariant is an option of a product (for example a brand) with an
// optional price adjustment.
type ProductVariant struct {
	ID              uuid.UUID
	ProductID       uuid.UUID
	Name            string
	PriceDeltaCents int64
}
