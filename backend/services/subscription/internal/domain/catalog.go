package domain

import "time"

// Product is a catalogue item a plan belongs to.
type Product struct {
	ID          string
	Name        string
	Description string
	SKU         string
	Active      bool
	Features    []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CreateProductInput is a validated product.
type CreateProductInput struct {
	Name        string
	Description string
	SKU         string
	Features    []string
}

// UpdateProductInput patches a product.
type UpdateProductInput struct {
	Name        string
	Description string
	SKU         string
	Active      *bool
	Features    []string
}

// OverdueInvoice is one open/overdue invoice with AR-follow-up context.
type OverdueInvoice struct {
	ID             string
	InvoiceNumber  string
	CustomerID     string
	CustomerName   string
	CustomerEmail  string
	SubscriptionID string
	PlanName       string
	DueDate        time.Time
	DaysOverdue    int
	AmountCents    int64
	Status         string
	RemindersSent  int
}

// CustomerInvoice is an invoice scoped to a customer for the self-service portal.
type CustomerInvoice struct {
	ID              string
	InvoiceNumber   string
	CustomerID      string
	SubscriptionID  string
	Status          string
	Date            time.Time
	DueDate         time.Time
	SubtotalCents   int64
	TaxCents        int64
	DiscountCents   int64
	TotalCents      int64
	AmountDueCents  int64
	AmountPaidCents int64
	Currency        string
	PaymentMethod   string
	PaidAt          *time.Time
	Items           []InvoiceLineItem
}

// CustomerOverview is the customer portal dashboard summary.
type CustomerOverview struct {
	ActiveSubscriptionsCount int
	TotalMonthlyValueCents   int64
	TotalOutstandingCents    int64
	OutstandingInvoicesCount int
	NextPaymentDue           *NextPayment
}

// NextPayment is the next upcoming renewal.
type NextPayment struct {
	AmountCents        int64
	Date               time.Time
	SubscriptionNumber string
	PlanName           string
}
