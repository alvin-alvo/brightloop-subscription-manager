package domain

import (
	"time"

	"github.com/google/uuid"
)

// InvoiceStatus is the lifecycle stage of an invoice.
type InvoiceStatus string

const (
	InvoiceDraft     InvoiceStatus = "draft"
	InvoiceConfirmed InvoiceStatus = "confirmed"
	InvoicePaid      InvoiceStatus = "paid"
)

// Valid reports whether the status is supported.
func (s InvoiceStatus) Valid() bool {
	switch s {
	case InvoiceDraft, InvoiceConfirmed, InvoicePaid:
		return true
	default:
		return false
	}
}

// Invoice is a bill generated from a subscription.
type Invoice struct {
	ID                 uuid.UUID
	Number             string
	SubscriptionID     *uuid.UUID
	CustomerID         uuid.UUID
	Status             InvoiceStatus
	IssueDate          time.Time
	DueDate            time.Time
	SubtotalCents      int64
	TaxTotalCents      int64
	DiscountTotalCents int64
	TotalCents         int64
	Lines              []InvoiceLine
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// InvoiceLine is a priced line on an invoice, with a snapshot of the tax rate
// so historical invoices remain correct if tax rules change later.
type InvoiceLine struct {
	ID                 uuid.UUID
	InvoiceID          uuid.UUID
	ProductID          *uuid.UUID
	Description        string
	Quantity           int
	UnitPriceCents     int64
	TaxRateBasisPoints int
	AmountCents        int64
}
