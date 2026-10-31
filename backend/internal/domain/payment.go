package domain

import (
	"time"

	"github.com/google/uuid"
)

// Payment records money received against an invoice.
type Payment struct {
	ID          uuid.UUID
	InvoiceID   uuid.UUID
	Method      string
	AmountCents int64
	PaidAt      time.Time
	Note        string
	CreatedAt   time.Time
}
