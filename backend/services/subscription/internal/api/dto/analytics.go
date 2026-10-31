package dto

import (
	"fmt"
	"math"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

// RawAnalytics exposes numeric metric values alongside the formatted cards.
type RawAnalytics struct {
	MRR                       float64        `json:"mrr"`
	MRRCents                  int64          `json:"mrrCents"`
	PreviousMRR               float64        `json:"previousMrr"`
	PreviousMRRCents          int64          `json:"previousMrrCents"`
	ARR                       float64        `json:"arr"`
	ARRCents                  int64          `json:"arrCents"`
	MRRChangePercent          float64        `json:"mrrChangePercent"`
	ActiveSubscriptions       int            `json:"activeSubscriptions"`
	TrialingSubscriptions     int            `json:"trialingSubscriptions"`
	PastDueSubscriptions      int            `json:"pastDueSubscriptions"`
	PausedSubscriptions       int            `json:"pausedSubscriptions"`
	UnpaidSubscriptions       int            `json:"unpaidSubscriptions"`
	CanceledSubscriptions     int            `json:"canceledSubscriptions"`
	TotalSubscriptions        int            `json:"totalSubscriptions"`
	NewSubscriptionsThisMonth int            `json:"newSubscriptionsThisMonth"`
	CanceledThisMonth         int            `json:"canceledThisMonth"`
	NetNewSubscriptions       int            `json:"netNewSubscriptions"`
	ActiveChangePercent       float64        `json:"activeSubscriptionsChangePercent"`
	NetRetentionRate          float64        `json:"netRetentionRate"`
	ChurnRate                 float64        `json:"churnRate"`
	ARPU                      float64        `json:"arpu"`
	ARPUCents                 int64          `json:"arpuCents"`
	TotalCustomers            int            `json:"totalCustomers"`
	ActiveCustomers           int            `json:"activeCustomers"`
	NewCustomersThisMonth     int            `json:"newCustomersThisMonth"`
	CustomersChangePercent    float64        `json:"customersChangePercent"`
	RevenueThisMonth          float64        `json:"revenueThisMonth"`
	RevenueThisMonthCents     int64          `json:"revenueThisMonthCents"`
	RevenueLastMonth          float64        `json:"revenueLastMonth"`
	RevenueLastMonthCents     int64          `json:"revenueLastMonthCents"`
	OutstandingAmount         float64        `json:"outstandingAmount"`
	OutstandingCents          int64          `json:"outstandingCents"`
	OverdueInvoices           int            `json:"overdueInvoices"`
	StatusCounts              map[string]int `json:"statusCounts"`
}

func round1(value float64) float64 {
	return math.Round(value*10) / 10
}

func rawAnalytics(c domain.ComputedAnalytics) RawAnalytics {
	counts := map[string]int{
		"active":   c.ActiveSubscriptions,
		"trialing": c.TrialingSubscriptions,
		"past_due": c.PastDueSubscriptions,
		"paused":   c.PausedSubscriptions,
		"unpaid":   c.UnpaidSubscriptions,
		"canceled": c.CanceledSubscriptions,
	}
	return RawAnalytics{
		MRR:                       Major(c.MRRCents),
		MRRCents:                  c.MRRCents,
		PreviousMRR:               Major(c.PreviousMRRCents),
		PreviousMRRCents:          c.PreviousMRRCents,
		ARR:                       Major(c.ARR),
		ARRCents:                  c.ARR,
		MRRChangePercent:          round1(c.MRRChangePercent),
		ActiveSubscriptions:       c.ActiveSubscriptions,
		TrialingSubscriptions:     c.TrialingSubscriptions,
		PastDueSubscriptions:      c.PastDueSubscriptions,
		PausedSubscriptions:       c.PausedSubscriptions,
		UnpaidSubscriptions:       c.UnpaidSubscriptions,
		CanceledSubscriptions:     c.CanceledSubscriptions,
		TotalSubscriptions:        c.TotalSubscriptions,
		NewSubscriptionsThisMonth: c.NewSubscriptions,
		CanceledThisMonth:         c.CanceledThisMonth,
		NetNewSubscriptions:       c.NetNewSubscriptions,
		ActiveChangePercent:       round1(c.ActiveChangePercent),
		NetRetentionRate:          round1(c.NetRetentionRate),
		ChurnRate:                 round1(c.ChurnRate),
		ARPU:                      Major(c.ARPUCents),
		ARPUCents:                 c.ARPUCents,
		TotalCustomers:            c.TotalCustomers,
		ActiveCustomers:           c.ActiveCustomers,
		NewCustomersThisMonth:     c.NewCustomers,
		CustomersChangePercent:    round1(c.CustomersChangePercent),
		RevenueThisMonth:          Major(c.RevenueThisMonthCents),
		RevenueThisMonthCents:     c.RevenueThisMonthCents,
		RevenueLastMonth:          Major(c.RevenueLastMonthCents),
		RevenueLastMonthCents:     c.RevenueLastMonthCents,
		OutstandingAmount:         Major(c.OutstandingCents),
		OutstandingCents:          c.OutstandingCents,
		OverdueInvoices:           c.OverdueInvoices,
		StatusCounts:              counts,
	}
}

// MetricCard matches the admin dashboard metric card shape.
type MetricCard struct {
	Key      string `json:"key"`
	Title    string `json:"title"`
	Value    string `json:"value"`
	Change   string `json:"change"`
	Trend    string `json:"trend"`
	Subtitle string `json:"subtitle"`
	Icon     string `json:"icon"`
}

func trend(value float64) string {
	if value < 0 {
		return "down"
	}
	return "up"
}

func metricCards(c domain.ComputedAnalytics, currency string) []MetricCard {
	mrrSubtitle := "No prior period baseline"
	if c.PreviousMRRCents > 0 {
		mrrSubtitle = fmt.Sprintf("vs. %s last month", FormatMoney(c.PreviousMRRCents, currency))
	}
	activeSubtitle := fmt.Sprintf("%d net new this month", c.NetNewSubscriptions)
	customerSubtitle := fmt.Sprintf("%d new this month", c.NewCustomers)

	return []MetricCard{
		{
			Key:      "mrr",
			Title:    "Monthly Recurring Revenue",
			Value:    FormatMoney(c.MRRCents, currency),
			Change:   FormatChange(c.MRRChangePercent),
			Trend:    trend(c.MRRChangePercent),
			Subtitle: mrrSubtitle,
			Icon:     "DollarSign",
		},
		{
			Key:      "active_subscriptions",
			Title:    "Active Subscriptions",
			Value:    fmt.Sprintf("%d", c.ActiveSubscriptions),
			Change:   FormatChange(c.ActiveChangePercent),
			Trend:    trend(c.ActiveChangePercent),
			Subtitle: activeSubtitle,
			Icon:     "CreditCard",
		},
		{
			Key:      "net_retention_rate",
			Title:    "Net Retention Rate",
			Value:    FormatPercent(c.NetRetentionRate),
			Change:   FormatChange(c.NetRetentionRate - 100),
			Trend:    trend(c.NetRetentionRate - 100),
			Subtitle: "Target: >105%",
			Icon:     "TrendingUp",
		},
		{
			Key:      "total_customers",
			Title:    "Total Customers",
			Value:    fmt.Sprintf("%d", c.TotalCustomers),
			Change:   FormatChange(c.CustomersChangePercent),
			Trend:    trend(c.CustomersChangePercent),
			Subtitle: customerSubtitle,
			Icon:     "Users",
		},
	}
}

// AnalyticsResponse is returned by /analytics and /metrics.
type AnalyticsResponse struct {
	GeneratedAt string       `json:"generatedAt"`
	Currency    string       `json:"currency"`
	Metrics     []MetricCard `json:"metrics"`
	Raw         RawAnalytics `json:"raw"`
}

// Analytics converts computed analytics.
func Analytics(c domain.ComputedAnalytics, currency string, now time.Time) AnalyticsResponse {
	return AnalyticsResponse{
		GeneratedAt: Timestamp(now),
		Currency:    currency,
		Metrics:     metricCards(c, currency),
		Raw:         rawAnalytics(c),
	}
}

// RevenuePointResponse is one point in a revenue trend chart.
type RevenuePointResponse struct {
	Month        string  `json:"month"`
	Label        string  `json:"label"`
	Revenue      float64 `json:"revenue"`
	RevenueCents int64   `json:"revenueCents"`
	PaymentCount int     `json:"paymentCount"`
}

// StatusBreakdownResponse is one slice of the status donut.
type StatusBreakdownResponse struct {
	Status  string  `json:"status"`
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

// UpcomingRenewalResponse is one row in the renewals list.
type UpcomingRenewalResponse struct {
	SubscriptionID string  `json:"subscriptionId"`
	CustomerID     string  `json:"customerId"`
	CustomerName   string  `json:"customerName"`
	PlanName       string  `json:"planName"`
	Date           string  `json:"date"`
	DaysUntil      int     `json:"daysUntil"`
	Amount         float64 `json:"amount"`
	AmountCents    int64   `json:"amountCents"`
	Currency       string  `json:"currency"`
}

// DashboardResponse is everything /admin/dashboard needs in one payload.
type DashboardResponse struct {
	GeneratedAt        string                    `json:"generatedAt"`
	Currency           string                    `json:"currency"`
	Metrics            []MetricCard              `json:"metrics"`
	Raw                RawAnalytics              `json:"raw"`
	RecentTransactions []TransactionResponse     `json:"recentTransactions"`
	RevenueTrend       []RevenuePointResponse    `json:"revenueTrend"`
	StatusBreakdown    []StatusBreakdownResponse `json:"statusBreakdown"`
	UpcomingRenewals   []UpcomingRenewalResponse `json:"upcomingRenewals"`
}

// Dashboard converts a domain dashboard.
func Dashboard(d domain.Dashboard, currency string) DashboardResponse {
	transactions := make([]TransactionResponse, 0, len(d.RecentTransactions))
	for _, tx := range d.RecentTransactions {
		transactions = append(transactions, Transaction(tx, d.GeneratedAt))
	}
	revenue := make([]RevenuePointResponse, 0, len(d.RevenueTrend))
	for _, point := range d.RevenueTrend {
		revenue = append(revenue, RevenuePointResponse{
			Month:        point.Month,
			Label:        point.Label,
			Revenue:      Major(point.RevenueCents),
			RevenueCents: point.RevenueCents,
			PaymentCount: point.PaymentCount,
		})
	}
	statuses := make([]StatusBreakdownResponse, 0, len(d.StatusBreakdown))
	for _, item := range d.StatusBreakdown {
		statuses = append(statuses, StatusBreakdownResponse{
			Status:  item.Status,
			Count:   item.Count,
			Percent: round1(item.Percent),
		})
	}
	today := time.Date(d.GeneratedAt.Year(), d.GeneratedAt.Month(), d.GeneratedAt.Day(), 0, 0, 0, 0, time.UTC)
	renewals := make([]UpcomingRenewalResponse, 0, len(d.UpcomingRenewals))
	for _, item := range d.UpcomingRenewals {
		days := int(item.Date.Sub(today).Hours() / 24)
		renewals = append(renewals, UpcomingRenewalResponse{
			SubscriptionID: item.SubscriptionID,
			CustomerID:     item.CustomerID,
			CustomerName:   item.CustomerName,
			PlanName:       item.PlanName,
			Date:           Date(item.Date),
			DaysUntil:      days,
			Amount:         Major(item.AmountCents),
			AmountCents:    item.AmountCents,
			Currency:       item.Currency,
		})
	}
	return DashboardResponse{
		GeneratedAt:        Timestamp(d.GeneratedAt),
		Currency:           currency,
		Metrics:            metricCards(d.Analytics, currency),
		Raw:                rawAnalytics(d.Analytics),
		RecentTransactions: transactions,
		RevenueTrend:       revenue,
		StatusBreakdown:    statuses,
		UpcomingRenewals:   renewals,
	}
}

// ReportPeriod is the reporting window.
type ReportPeriod struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// PlanBreakdownResponse is one plan in the revenue report.
type PlanBreakdownResponse struct {
	PlanID              string  `json:"planId"`
	PlanName            string  `json:"planName"`
	ProductName         string  `json:"productName,omitempty"`
	BillingInterval     string  `json:"billingInterval"`
	ActiveSubscriptions int     `json:"activeSubscriptions"`
	MRR                 float64 `json:"mrr"`
	MRRCents            int64   `json:"mrrCents"`
	Revenue             float64 `json:"revenue"`
	RevenueCents        int64   `json:"revenueCents"`
}

// TopCustomerResponse is one customer in the revenue report.
type TopCustomerResponse struct {
	CustomerID          string  `json:"customerId"`
	CustomerName        string  `json:"customerName"`
	CustomerEmail       string  `json:"customerEmail"`
	ActiveSubscriptions int     `json:"activeSubscriptions"`
	MRR                 float64 `json:"mrr"`
	MRRCents            int64   `json:"mrrCents"`
	Revenue             float64 `json:"revenue"`
	RevenueCents        int64   `json:"revenueCents"`
}

// AgingResponse is one accounts-receivable bucket.
type AgingResponse struct {
	Bucket      string  `json:"bucket"`
	Label       string  `json:"label"`
	Count       int     `json:"count"`
	Amount      float64 `json:"amount"`
	AmountCents int64   `json:"amountCents"`
}

// ReportResponse is returned by /reports.
type ReportResponse struct {
	GeneratedAt     string                    `json:"generatedAt"`
	Currency        string                    `json:"currency"`
	Period          ReportPeriod              `json:"period"`
	Metrics         []MetricCard              `json:"metrics"`
	Raw             RawAnalytics              `json:"raw"`
	MonthlyRevenue  []RevenuePointResponse    `json:"monthlyRevenue"`
	PlanBreakdown   []PlanBreakdownResponse   `json:"planBreakdown"`
	TopCustomers    []TopCustomerResponse     `json:"topCustomers"`
	InvoiceAging    []AgingResponse           `json:"invoiceAging"`
	StatusBreakdown []StatusBreakdownResponse `json:"statusBreakdown"`
}

// Report converts a domain report.
func Report(report domain.AnalyticsReport, currency string) ReportResponse {
	revenue := make([]RevenuePointResponse, 0, len(report.MonthlyRevenue))
	for _, point := range report.MonthlyRevenue {
		revenue = append(revenue, RevenuePointResponse{
			Month:        point.Month,
			Label:        point.Label,
			Revenue:      Major(point.RevenueCents),
			RevenueCents: point.RevenueCents,
			PaymentCount: point.PaymentCount,
		})
	}
	plans := make([]PlanBreakdownResponse, 0, len(report.PlanBreakdown))
	for _, plan := range report.PlanBreakdown {
		plans = append(plans, PlanBreakdownResponse{
			PlanID:              plan.PlanID,
			PlanName:            plan.PlanName,
			ProductName:         plan.ProductName,
			BillingInterval:     string(plan.BillingInterval),
			ActiveSubscriptions: plan.ActiveSubscriptions,
			MRR:                 Major(plan.MRRCents),
			MRRCents:            plan.MRRCents,
			Revenue:             Major(plan.RevenueCents),
			RevenueCents:        plan.RevenueCents,
		})
	}
	customers := make([]TopCustomerResponse, 0, len(report.TopCustomers))
	for _, customer := range report.TopCustomers {
		customers = append(customers, TopCustomerResponse{
			CustomerID:          customer.CustomerID,
			CustomerName:        customer.CustomerName,
			CustomerEmail:       customer.CustomerEmail,
			ActiveSubscriptions: customer.ActiveSubscriptions,
			MRR:                 Major(customer.MRRCents),
			MRRCents:            customer.MRRCents,
			Revenue:             Major(customer.RevenueCents),
			RevenueCents:        customer.RevenueCents,
		})
	}
	aging := make([]AgingResponse, 0, len(report.InvoiceAging))
	for _, bucket := range report.InvoiceAging {
		aging = append(aging, AgingResponse{
			Bucket:      bucket.Bucket,
			Label:       bucket.Label,
			Count:       bucket.Count,
			Amount:      Major(bucket.AmountCents),
			AmountCents: bucket.AmountCents,
		})
	}
	statuses := make([]StatusBreakdownResponse, 0, len(report.StatusBreakdown))
	for _, item := range report.StatusBreakdown {
		statuses = append(statuses, StatusBreakdownResponse{
			Status:  item.Status,
			Count:   item.Count,
			Percent: round1(item.Percent),
		})
	}
	period := ReportPeriod{}
	if !report.From.IsZero() {
		period.From = Date(report.From)
	}
	if !report.To.IsZero() {
		period.To = Date(report.To)
	}
	return ReportResponse{
		GeneratedAt:     Timestamp(report.GeneratedAt),
		Currency:        currency,
		Period:          period,
		Metrics:         metricCards(report.Analytics, currency),
		Raw:             rawAnalytics(report.Analytics),
		MonthlyRevenue:  revenue,
		PlanBreakdown:   plans,
		TopCustomers:    customers,
		InvoiceAging:    aging,
		StatusBreakdown: statuses,
	}
}
