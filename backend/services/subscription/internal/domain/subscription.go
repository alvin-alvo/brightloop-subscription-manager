package domain

import (
	"strings"
	"time"
)

// SubscriptionStatus mirrors the subscription_status enum in PostgreSQL.
type SubscriptionStatus string

const (
	StatusActive   SubscriptionStatus = "active"
	StatusTrialing SubscriptionStatus = "trialing"
	StatusPastDue  SubscriptionStatus = "past_due"
	StatusCanceled SubscriptionStatus = "canceled"
	StatusPaused   SubscriptionStatus = "paused"
	StatusUnpaid   SubscriptionStatus = "unpaid"
)

// Valid reports whether s is a value accepted by the database enum.
func (s SubscriptionStatus) Valid() bool {
	switch s {
	case StatusActive, StatusTrialing, StatusPastDue, StatusCanceled, StatusPaused, StatusUnpaid:
		return true
	default:
		return false
	}
}

// Billable reports whether a subscription in this state should generate
// recurring revenue.
func (s SubscriptionStatus) Billable() bool {
	switch s {
	case StatusActive, StatusTrialing, StatusPastDue:
		return true
	default:
		return false
	}
}

// BillingInterval mirrors the billing_interval enum in PostgreSQL.
type BillingInterval string

const (
	IntervalMonthly   BillingInterval = "monthly"
	IntervalQuarterly BillingInterval = "quarterly"
	IntervalYearly    BillingInterval = "yearly"
)

// Valid reports whether b is a value accepted by the database enum.
func (b BillingInterval) Valid() bool {
	switch b {
	case IntervalMonthly, IntervalQuarterly, IntervalYearly:
		return true
	default:
		return false
	}
}

// Months returns how many months a billing interval spans.
func (b BillingInterval) Months() int {
	switch b {
	case IntervalQuarterly:
		return 3
	case IntervalYearly:
		return 12
	default:
		return 1
	}
}

// Title returns the interval formatted for display ("Monthly", "Annual"...).
func (b BillingInterval) Title() string {
	switch b {
	case IntervalQuarterly:
		return "Quarterly"
	case IntervalYearly:
		return "Annual"
	default:
		return "Monthly"
	}
}

// DiscountKind mirrors discount_kind.
type DiscountKind string

const (
	DiscountPercentage  DiscountKind = "percentage"
	DiscountFixedAmount DiscountKind = "fixed_amount"
)

// DiscountDuration mirrors discount_duration.
type DiscountDuration string

const (
	DiscountOnce      DiscountDuration = "once"
	DiscountRepeating DiscountDuration = "repeating"
	DiscountForever   DiscountDuration = "forever"
)

// InvoiceStatus mirrors invoice_status.
type InvoiceStatus string

const (
	InvoiceDraft         InvoiceStatus = "draft"
	InvoiceOpen          InvoiceStatus = "open"
	InvoicePaid          InvoiceStatus = "paid"
	InvoiceUncollectible InvoiceStatus = "uncollectible"
	InvoiceVoid          InvoiceStatus = "void"
)

// PaymentStatus mirrors payment_status.
type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentSucceeded PaymentStatus = "succeeded"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
)

// Customer is the subset of the customer record required by subscriptions.
type Customer struct {
	ID          string
	UserID      *string
	Name        string
	Email       string
	CompanyName string
	Currency    string
	Status      string
}

// Plan is a billable plan in the product catalogue.
type Plan struct {
	ID              string
	ProductID       string
	ProductName     string
	Name            string
	Description     string
	BillingInterval BillingInterval
	PriceCents      int64
	Currency        string
	Features        []string
	TrialPeriodDays int
	Active          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Discount is a redeemable promotion.
type Discount struct {
	ID                 string
	Code               string
	Description        string
	Kind               DiscountKind
	PercentBasisPoints int
	AmountCents        int64
	Currency           string
	Duration           DiscountDuration
	DurationMonths     int
	MaxRedemptions     *int
	TimesRedeemed      int
	StartsAt           *time.Time
	EndsAt             *time.Time
	Active             bool
}

// Subscription is the core aggregate of this service.
type Subscription struct {
	ID                 string
	CustomerID         string
	CustomerName       string
	CustomerEmail      string
	PlanID             string
	Plan               *Plan
	Status             SubscriptionStatus
	Quantity           int
	UnitPriceCents     int64
	DiscountCents      int64
	Currency           string
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	CancelAtPeriodEnd  bool
	TrialEndsAt        *time.Time
	StartedAt          time.Time
	CanceledAt         *time.Time
	EndedAt            *time.Time
	CreatedBy          *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// SubtotalCents is the undiscounted amount for one billing period.
func (s Subscription) SubtotalCents() int64 {
	return s.UnitPriceCents * int64(s.Quantity)
}

// TotalCents is the amount charged for one billing period.
func (s Subscription) TotalCents() int64 {
	total := s.SubtotalCents() - s.DiscountCents
	if total < 0 {
		return 0
	}
	return total
}

// MonthlyRecurringCents normalises one billing period to a monthly amount.
func (s Subscription) MonthlyRecurringCents() int64 {
	months := 1
	if s.Plan != nil {
		months = s.Plan.BillingInterval.Months()
	}
	if months <= 0 {
		months = 1
	}
	return (s.TotalCents() + int64(months)/2) / int64(months)
}

// BillingInterval returns the interval of the attached plan.
func (s Subscription) BillingInterval() BillingInterval {
	if s.Plan != nil && s.Plan.BillingInterval.Valid() {
		return s.Plan.BillingInterval
	}
	return IntervalMonthly
}

// SubscriptionSummary aggregates a filtered subscription list.
type SubscriptionSummary struct {
	Total    int   `json:"total"`
	Active   int   `json:"active"`
	Trialing int   `json:"trialing"`
	PastDue  int   `json:"pastDue"`
	Canceled int   `json:"canceled"`
	Paused   int   `json:"paused"`
	Unpaid   int   `json:"unpaid"`
	MRRCents int64 `json:"mrrCents"`
	Currency string
}

// ListFilter controls subscription listing and searching.
type ListFilter struct {
	Status     SubscriptionStatus
	CustomerID string
	PlanID     string
	Search     string
	SortBy     string
	SortOrder  string
	Page       int
	PageSize   int
}

// Normalize applies pagination defaults and bounds.
func (f ListFilter) Normalize() ListFilter {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	if f.PageSize > 100 {
		f.PageSize = 100
	}
	switch strings.ToLower(f.SortBy) {
	case "created_at", "updated_at", "current_period_end", "unit_price_cents", "customer_name", "status":
		f.SortBy = strings.ToLower(f.SortBy)
	default:
		f.SortBy = "created_at"
	}
	if strings.EqualFold(f.SortOrder, "asc") {
		f.SortOrder = "asc"
	} else {
		f.SortOrder = "desc"
	}
	return f
}

// SubscriptionPage is a paginated result with aggregate counts.
type SubscriptionPage struct {
	Items      []Subscription
	Total      int
	Page       int
	PageSize   int
	TotalPages int
	Summary    SubscriptionSummary
}

// CreateSubscriptionInput is the validated command used to create a
// subscription. Money is expressed in minor units.
type CreateSubscriptionInput struct {
	CustomerID         string
	PlanID             string
	Quantity           int
	UnitPriceCents     *int64
	DiscountCents      int64
	DiscountCode       string
	Currency           string
	CurrentPeriodStart *time.Time
	CurrentPeriodEnd   *time.Time
	Status             SubscriptionStatus
	TrialDays          *int
	CreatedBy          *string
}

// UpdateSubscriptionInput is the validated command used to patch a
// subscription. Nil pointers mean "leave unchanged".
type UpdateSubscriptionInput struct {
	PlanID             *string
	Quantity           *int
	UnitPriceCents     *int64
	DiscountCents      *int64
	Currency           *string
	CurrentPeriodStart *time.Time
	CurrentPeriodEnd   *time.Time
	Status             *SubscriptionStatus
	CancelAtPeriodEnd  *bool
}

// Actor is the authenticated caller. Role is one of admin, staff or customer.
type Actor struct {
	UserID     string
	Role       string
	CustomerID string
}

// IsStaff reports whether the actor may mutate subscriptions.
func (a Actor) IsStaff() bool {
	return a.Role == "admin" || a.Role == "staff"
}

// BillingCycle is one period in a subscription's billing schedule.
type BillingCycle struct {
	PeriodStart   time.Time
	PeriodEnd     time.Time
	AmountCents   int64
	Status        string
	InvoiceID     string
	InvoiceNumber string
	IssuedAt      *time.Time
	PaidAt        *time.Time
}

// BillingSchedule is the billing history / upcoming cycle for a subscription.
type BillingSchedule struct {
	SubscriptionID     string
	PlanID             string
	BillingInterval    BillingInterval
	Currency           string
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	NextBillingDate    time.Time
	Cycles             []BillingCycle
	TotalPaidCents     int64
}

// Invoice is a stored invoice.
type Invoice struct {
	ID              string
	InvoiceNumber   string
	CustomerID      string
	SubscriptionID  *string
	Status          InvoiceStatus
	Currency        string
	SubtotalCents   int64
	DiscountCents   int64
	TaxCents        int64
	AmountDueCents  int64
	AmountPaidCents int64
	IssueDate       time.Time
	DueDate         time.Time
	PaidAt          *time.Time
	Notes           *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Items           []InvoiceLineItem
}

// InvoiceLineItem is one line on an invoice.
type InvoiceLineItem struct {
	ID             string
	InvoiceID      string
	PlanID         *string
	Description    string
	Quantity       int
	UnitPriceCents int64
	AmountCents    int64
	CreatedAt      time.Time
}

// Transaction is a payment joined with its invoice, customer and plan. It is
// what the admin dashboard settlement stream renders.
type Transaction struct {
	ID              string
	InvoiceID       string
	InvoiceNumber   string
	CustomerID      string
	CustomerName    string
	CustomerEmail   string
	SubscriptionID  string
	PlanID          string
	PlanName        string
	BillingInterval BillingInterval
	AmountCents     int64
	Currency        string
	Method          string
	Status          PaymentStatus
	Reference       string
	PaidAt          *time.Time
	CreatedAt       time.Time
}

// TransactionFilter filters the settlement stream.
type TransactionFilter struct {
	CustomerID string
	Status     PaymentStatus
	Page       int
	PageSize   int
}

// Normalize applies pagination defaults and bounds.
func (f TransactionFilter) Normalize() TransactionFilter {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	if f.PageSize > 100 {
		f.PageSize = 100
	}
	return f
}

// TransactionPage is a paginated settlement stream.
type TransactionPage struct {
	Items      []Transaction
	Total      int
	Page       int
	PageSize   int
	TotalPages int
}
