package domain

import "time"

// Analytics contains the raw aggregates read from PostgreSQL. Money is stored
// in minor units; the API layer converts it for the frontend.
type Analytics struct {
	ActiveSubscriptions   int
	TrialingSubscriptions int
	PastDueSubscriptions  int
	CanceledSubscriptions int
	PausedSubscriptions   int
	UnpaidSubscriptions   int
	TotalSubscriptions    int
	NewSubscriptions      int
	CanceledThisMonth     int

	MRRCents            int64
	PreviousMRRCents    int64
	NetRetainedMRRCents int64

	TotalCustomers  int
	ActiveCustomers int
	NewCustomers    int

	RevenueThisMonthCents int64
	RevenueLastMonthCents int64
	OutstandingCents      int64
	OverdueInvoices       int
}

// ComputedAnalytics decorates raw analytics with derived business metrics.
type ComputedAnalytics struct {
	Analytics

	ARR                    int64
	ARPUCents              int64
	MRRChangePercent       float64
	ActiveChangePercent    float64
	CustomersChangePercent float64
	NetNewSubscriptions    int
	NetRetentionRate       float64
	ChurnRate              float64
}

// MonthlyRevenue is one point in the revenue trend.
type MonthlyRevenue struct {
	Month        string
	Label        string
	RevenueCents int64
	PaymentCount int
}

// PlanBreakdown is revenue and subscription distribution for a plan.
type PlanBreakdown struct {
	PlanID              string
	PlanName            string
	ProductName         string
	BillingInterval     BillingInterval
	ActiveSubscriptions int
	MRRCents            int64
	RevenueCents        int64
}

// StatusBreakdown is one subscription status count.
type StatusBreakdown struct {
	Status  string
	Count   int
	Percent float64
}

// TopCustomer ranks a customer by recurring and collected revenue.
type TopCustomer struct {
	CustomerID          string
	CustomerName        string
	CustomerEmail       string
	ActiveSubscriptions int
	MRRCents            int64
	RevenueCents        int64
}

// AgingBucket is an accounts-receivable aging bucket.
type AgingBucket struct {
	Bucket      string
	Label       string
	Count       int
	AmountCents int64
}

// UpcomingRenewal is a subscription approaching its next renewal.
type UpcomingRenewal struct {
	SubscriptionID string
	CustomerID     string
	CustomerName   string
	PlanName       string
	Date           time.Time
	AmountCents    int64
	Currency       string
}

// Dashboard is the payload consumed by /admin/dashboard.
type Dashboard struct {
	GeneratedAt        time.Time
	Analytics          ComputedAnalytics
	RecentTransactions []Transaction
	RevenueTrend       []MonthlyRevenue
	StatusBreakdown    []StatusBreakdown
	UpcomingRenewals   []UpcomingRenewal
}

// AnalyticsReport is the payload consumed by the reports view.
type AnalyticsReport struct {
	GeneratedAt     time.Time
	From            time.Time
	To              time.Time
	Analytics       ComputedAnalytics
	MonthlyRevenue  []MonthlyRevenue
	PlanBreakdown   []PlanBreakdown
	TopCustomers    []TopCustomer
	InvoiceAging    []AgingBucket
	StatusBreakdown []StatusBreakdown
}

// BillingBatch is one day of generated invoices, shown as a billing run.
type BillingBatch struct {
	ID                 string
	Date               time.Time
	ItemsCount         int
	TotalBilledCents   int64
	PaidCount          int
	SuccessRatePercent float64
}

// BillingOverview is the billing operations dashboard payload.
type BillingOverview struct {
	PendingInvoices int
	OverdueInvoices int
	Batches         []BillingBatch
	GeneratedAt     time.Time
}

// BillingRunResult reports what a recurring billing pass changed.
type BillingRunResult struct {
	AsOf             time.Time
	DryRun           bool
	Processed        int
	Renewed          int
	Canceled         int
	InvoicesCreated  int
	TotalBilledCents int64
	Failed           int
	Items            []BillingRunItem
}

// BillingRunItem is the per-subscription outcome of a billing pass.
type BillingRunItem struct {
	SubscriptionID string
	CustomerName   string
	Action         string // renewed, canceled, skipped, failed
	InvoiceID      string
	InvoiceNumber  string
	AmountCents    int64
	Error          string
}
