package domain

import "time"

// QuoteItem is one line in a quotation.
type QuoteItem struct {
	ID             string
	PlanID         *string
	Description    string
	Quantity       int
	UnitPriceCents int64
	TaxBasisPoints int64
	LineTotalCents int64
}

// Quote is a sales quotation.
type Quote struct {
	ID            string
	QuoteNumber   string
	CustomerID    string
	CustomerName  string
	Status        string // draft/sent/accepted/declined/expired
	Currency      string
	SubtotalCents int64
	DiscountCents int64
	TaxCents      int64
	TotalCents    int64
	ValidUntil    time.Time
	Notes         string
	Items         []QuoteItem
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Plan returns the display plan name (first non-empty line description).
func (q Quote) Plan() string {
	for _, item := range q.Items {
		if item.Description != "" {
			return item.Description
		}
	}
	return q.Currency
}

// QuoteTemplate is a reusable quote template.
type QuoteTemplate struct {
	ID           string
	Name         string
	PlanID       *string
	PlanName     string
	ValidityDays int
	Items        []QuoteTemplateItem
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// QuoteTemplateItem is one line in a quote template.
type QuoteTemplateItem struct {
	ID             string
	ProductID      *string
	ProductName    string
	Quantity       int
	UnitPriceCents int64
	TaxBasisPoints int64
}

// CustomerWorkspace is a CRM row with subscription/value insight.
type CustomerWorkspace struct {
	ID                  string
	Name                string
	Contact             string
	Email               string
	ActiveSubscriptions int
	MonthlyValueCents   int64
	Tier                string
	Status              string
}

// Workspace is the aggregated staff CRM payload.
type Workspace struct {
	Customers     []CustomerWorkspace `json:"customers"`
	Templates     []QuoteTemplate     `json:"templates"`
	Quotes        []Quote             `json:"quotes"`
	Subscriptions []Subscription      `json:"subscriptions"`
	Invoices      []OverdueInvoice    `json:"invoices"`
}
