package domain

import (
	"time"

	"github.com/google/uuid"
)

// TaxRate is an admin-managed tax rule. The rate is expressed in basis points
// (2000 = 20.00%) to keep calculations exact.
type TaxRate struct {
	ID              uuid.UUID
	Name            string
	RateBasisPoints int
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
