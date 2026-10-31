package dto

import (
	"math"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

// -----------------------------------------------------------------------------
// Request payloads
// -----------------------------------------------------------------------------

// CreateRequest is the POST /api/v1/subscriptions body.
type CreateRequest struct {
	CustomerID         string   `json:"customerId"`
	PlanID             string   `json:"planId"`
	Quantity           int      `json:"quantity"`
	UnitPrice          *float64 `json:"unitPrice"`
	UnitPriceCents     *int64   `json:"unitPriceCents"`
	Discount           *float64 `json:"discount"`
	DiscountCents      int64    `json:"discountCents"`
	DiscountCode       string   `json:"discountCode"`
	Currency           string   `json:"currency"`
	CurrentPeriodStart string   `json:"currentPeriodStart"`
	CurrentPeriodEnd   string   `json:"currentPeriodEnd"`
	Status             string   `json:"status"`
	TrialDays          *int     `json:"trialDays"`
}

// ToInput validates and converts the request to a service command.
func (r CreateRequest) ToInput() (domain.CreateSubscriptionInput, error) {
	input := domain.CreateSubscriptionInput{
		CustomerID:    strings.TrimSpace(r.CustomerID),
		PlanID:        strings.TrimSpace(r.PlanID),
		Quantity:      r.Quantity,
		DiscountCents: r.DiscountCents,
		DiscountCode:  strings.TrimSpace(r.DiscountCode),
		Currency:      strings.TrimSpace(r.Currency),
		TrialDays:     r.TrialDays,
	}
	if r.UnitPriceCents != nil {
		input.UnitPriceCents = r.UnitPriceCents
	} else if r.UnitPrice != nil {
		cents := int64(math.Round(*r.UnitPrice * 100))
		input.UnitPriceCents = &cents
	}
	if r.DiscountCents == 0 && r.Discount != nil {
		input.DiscountCents = int64(math.Round(*r.Discount * 100))
	}
	if r.Status != "" {
		input.Status = domain.SubscriptionStatus(strings.ToLower(strings.TrimSpace(r.Status)))
	}
	if r.CurrentPeriodStart != "" {
		start, err := ParseDate(r.CurrentPeriodStart)
		if err != nil {
			return input, err
		}
		input.CurrentPeriodStart = &start
	}
	if r.CurrentPeriodEnd != "" {
		end, err := ParseDate(r.CurrentPeriodEnd)
		if err != nil {
			return input, err
		}
		input.CurrentPeriodEnd = &end
	}
	return input, nil
}

// UpdateRequest is the PUT /api/v1/subscriptions/:id body.
type UpdateRequest struct {
	PlanID             *string  `json:"planId"`
	Quantity           *int     `json:"quantity"`
	UnitPrice          *float64 `json:"unitPrice"`
	UnitPriceCents     *int64   `json:"unitPriceCents"`
	Discount           *float64 `json:"discount"`
	DiscountCents      *int64   `json:"discountCents"`
	Currency           *string  `json:"currency"`
	CurrentPeriodStart *string  `json:"currentPeriodStart"`
	CurrentPeriodEnd   *string  `json:"currentPeriodEnd"`
	Status             *string  `json:"status"`
	CancelAtPeriodEnd  *bool    `json:"cancelAtPeriodEnd"`
}

// ToInput validates and converts the request to a service command.
func (r UpdateRequest) ToInput() (domain.UpdateSubscriptionInput, error) {
	input := domain.UpdateSubscriptionInput{
		PlanID:            r.PlanID,
		Quantity:          r.Quantity,
		Currency:          r.Currency,
		CancelAtPeriodEnd: r.CancelAtPeriodEnd,
		DiscountCents:     r.DiscountCents,
	}
	if r.UnitPriceCents != nil {
		input.UnitPriceCents = r.UnitPriceCents
	} else if r.UnitPrice != nil {
		cents := int64(math.Round(*r.UnitPrice * 100))
		input.UnitPriceCents = &cents
	}
	if r.DiscountCents == nil && r.Discount != nil {
		cents := int64(math.Round(*r.Discount * 100))
		input.DiscountCents = &cents
	}
	if r.Status != nil {
		status := domain.SubscriptionStatus(strings.ToLower(strings.TrimSpace(*r.Status)))
		input.Status = &status
	}
	if r.CurrentPeriodStart != nil {
		start, err := ParseDate(*r.CurrentPeriodStart)
		if err != nil {
			return input, err
		}
		input.CurrentPeriodStart = &start
	}
	if r.CurrentPeriodEnd != nil {
		end, err := ParseDate(*r.CurrentPeriodEnd)
		if err != nil {
			return input, err
		}
		input.CurrentPeriodEnd = &end
	}
	return input, nil
}

// BillingRunRequest is the POST /api/v1/subscriptions/billing/run body.
type BillingRunRequest struct {
	AsOf   string `json:"asOf"`
	DryRun bool   `json:"dryRun"`
}

// ParseDate accepts ISO dates and RFC3339 timestamps.
func ParseDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, domain.Validationf("date value is required")
	}
	layouts := []string{
		"2006-01-02",
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, domain.Validationf("invalid date %q; expected YYYY-MM-DD", value)
}

// -----------------------------------------------------------------------------
// Invoice and billing responses
// -----------------------------------------------------------------------------

// InvoiceLineItemResponse mirrors the frontend InvoiceLineItem interface.
type InvoiceLineItemResponse struct {
	ID          string  `json:"id"`
	PlanID      string  `json:"planId,omitempty"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	Amount      float64 `json:"amount"`
}

// InvoiceResponse mirrors the frontend Invoice interface.
type InvoiceResponse struct {
	ID              string                    `json:"id"`
	InvoiceNumber   string                    `json:"invoiceNumber"`
	CustomerID      string                    `json:"customerId"`
	SubscriptionID  string                    `json:"subscriptionId,omitempty"`
	Status          string                    `json:"status"`
	Currency        string                    `json:"currency"`
	Subtotal        float64                   `json:"subtotal"`
	SubtotalCents   int64                     `json:"subtotalCents"`
	Discount        float64                   `json:"discount"`
	DiscountCents   int64                     `json:"discountCents"`
	Tax             float64                   `json:"tax"`
	TaxCents        int64                     `json:"taxCents"`
	AmountDue       float64                   `json:"amountDue"`
	AmountDueCents  int64                     `json:"amountDueCents"`
	AmountPaid      float64                   `json:"amountPaid"`
	AmountPaidCents int64                     `json:"amountPaidCents"`
	IssueDate       string                    `json:"issueDate"`
	DueDate         string                    `json:"dueDate"`
	PaidAt          *string                   `json:"paidAt,omitempty"`
	CreatedAt       string                    `json:"createdAt"`
	UpdatedAt       string                    `json:"updatedAt"`
	LineItems       []InvoiceLineItemResponse `json:"lineItems"`
}

// Invoice converts a domain invoice.
func Invoice(invoice domain.Invoice) InvoiceResponse {
	items := make([]InvoiceLineItemResponse, 0, len(invoice.Items))
	for _, item := range invoice.Items {
		planID := ""
		if item.PlanID != nil {
			planID = *item.PlanID
		}
		items = append(items, InvoiceLineItemResponse{
			ID:          item.ID,
			PlanID:      planID,
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   Major(item.UnitPriceCents),
			Amount:      Major(item.AmountCents),
		})
	}
	var subscriptionID string
	if invoice.SubscriptionID != nil {
		subscriptionID = *invoice.SubscriptionID
	}
	return InvoiceResponse{
		ID:              invoice.ID,
		InvoiceNumber:   invoice.InvoiceNumber,
		CustomerID:      invoice.CustomerID,
		SubscriptionID:  subscriptionID,
		Status:          string(invoice.Status),
		Currency:        invoice.Currency,
		Subtotal:        Major(invoice.SubtotalCents),
		SubtotalCents:   invoice.SubtotalCents,
		Discount:        Major(invoice.DiscountCents),
		DiscountCents:   invoice.DiscountCents,
		Tax:             Major(invoice.TaxCents),
		TaxCents:        invoice.TaxCents,
		AmountDue:       Major(invoice.AmountDueCents),
		AmountDueCents:  invoice.AmountDueCents,
		AmountPaid:      Major(invoice.AmountPaidCents),
		AmountPaidCents: invoice.AmountPaidCents,
		IssueDate:       Date(invoice.IssueDate),
		DueDate:         Date(invoice.DueDate),
		PaidAt:          TimestampPtr(invoice.PaidAt),
		CreatedAt:       Timestamp(invoice.CreatedAt),
		UpdatedAt:       Timestamp(invoice.UpdatedAt),
		LineItems:       items,
	}
}

// RenewResponse returns the advanced subscription and its new invoice.
type RenewResponse struct {
	Subscription SubscriptionResponse `json:"subscription"`
	Invoice      InvoiceResponse      `json:"invoice"`
}

// BillingCycleResponse is one cycle in a subscription's schedule.
type BillingCycleResponse struct {
	PeriodStart   string  `json:"periodStart"`
	PeriodEnd     string  `json:"periodEnd"`
	Amount        float64 `json:"amount"`
	AmountCents   int64   `json:"amountCents"`
	Status        string  `json:"status"`
	InvoiceID     string  `json:"invoiceId,omitempty"`
	InvoiceNumber string  `json:"invoiceNumber,omitempty"`
	IssuedAt      *string `json:"issuedAt,omitempty"`
	PaidAt        *string `json:"paidAt,omitempty"`
}

// BillingScheduleResponse is returned by /:id/billing.
type BillingScheduleResponse struct {
	SubscriptionID     string                 `json:"subscriptionId"`
	PlanID             string                 `json:"planId"`
	BillingInterval    string                 `json:"billingInterval"`
	BillingCycle       string                 `json:"billingCycle"`
	Currency           string                 `json:"currency"`
	CurrentPeriodStart string                 `json:"currentPeriodStart"`
	CurrentPeriodEnd   string                 `json:"currentPeriodEnd"`
	NextBillingDate    string                 `json:"nextBillingDate"`
	Cycles             []BillingCycleResponse `json:"cycles"`
	TotalPaid          float64                `json:"totalPaid"`
	TotalPaidCents     int64                  `json:"totalPaidCents"`
}

// BillingSchedule converts a domain schedule.
func BillingSchedule(schedule domain.BillingSchedule) BillingScheduleResponse {
	cycles := make([]BillingCycleResponse, 0, len(schedule.Cycles))
	for _, cycle := range schedule.Cycles {
		cycles = append(cycles, BillingCycleResponse{
			PeriodStart:   Date(cycle.PeriodStart),
			PeriodEnd:     Date(cycle.PeriodEnd),
			Amount:        Major(cycle.AmountCents),
			AmountCents:   cycle.AmountCents,
			Status:        cycle.Status,
			InvoiceID:     cycle.InvoiceID,
			InvoiceNumber: cycle.InvoiceNumber,
			IssuedAt:      TimestampPtr(cycle.IssuedAt),
			PaidAt:        TimestampPtr(cycle.PaidAt),
		})
	}
	return BillingScheduleResponse{
		SubscriptionID:     schedule.SubscriptionID,
		PlanID:             schedule.PlanID,
		BillingInterval:    string(schedule.BillingInterval),
		BillingCycle:       schedule.BillingInterval.Title(),
		Currency:           schedule.Currency,
		CurrentPeriodStart: Date(schedule.CurrentPeriodStart),
		CurrentPeriodEnd:   Date(schedule.CurrentPeriodEnd),
		NextBillingDate:    Date(schedule.NextBillingDate),
		Cycles:             cycles,
		TotalPaid:          Major(schedule.TotalPaidCents),
		TotalPaidCents:     schedule.TotalPaidCents,
	}
}

// BillingBatchResponse is one recurring billing run row.
type BillingBatchResponse struct {
	ID                 string  `json:"id"`
	Date               string  `json:"date"`
	Label              string  `json:"label"`
	ItemsCount         int     `json:"itemsCount"`
	TotalBilled        string  `json:"totalBilled"`
	TotalBilledCents   int64   `json:"totalBilledCents"`
	PaidCount          int     `json:"paidCount"`
	SuccessRate        string  `json:"successRate"`
	SuccessRatePercent float64 `json:"successRatePercent"`
	Status             string  `json:"status"`
}

// SchedulerResponse describes the recurring billing scheduler.
type SchedulerResponse struct {
	Enabled   bool   `json:"enabled"`
	Schedule  string `json:"schedule"`
	NextRunAt string `json:"nextRunAt"`
}

// DunningResponse describes retry policy.
type DunningResponse struct {
	Attempts     int     `json:"attempts"`
	RecoveryRate float64 `json:"recoveryRate"`
	Description  string  `json:"description"`
}

// BillingOverviewResponse is returned by /billing/overview.
type BillingOverviewResponse struct {
	GeneratedAt     string                 `json:"generatedAt"`
	PendingInvoices int                    `json:"pendingInvoices"`
	OverdueInvoices int                    `json:"overdueInvoices"`
	Scheduler       SchedulerResponse      `json:"scheduler"`
	Dunning         DunningResponse        `json:"dunning"`
	Batches         []BillingBatchResponse `json:"batches"`
}

// BillingBatches converts recurring billing runs.
func BillingBatches(batches []domain.BillingBatch, now time.Time) []BillingBatchResponse {
	items := make([]BillingBatchResponse, 0, len(batches))
	for _, batch := range batches {
		items = append(items, BillingBatchResponse{
			ID:                 batch.ID,
			Date:               Date(batch.Date),
			Label:              TransactionDate(batch.Date, now),
			ItemsCount:         batch.ItemsCount,
			TotalBilled:        FormatMoney(batch.TotalBilledCents, "USD"),
			TotalBilledCents:   batch.TotalBilledCents,
			PaidCount:          batch.PaidCount,
			SuccessRate:        FormatRate(batch.SuccessRatePercent),
			SuccessRatePercent: round1(batch.SuccessRatePercent),
			Status:             "completed",
		})
	}
	return items
}

// BillingOverview converts a domain overview.
func BillingOverview(overview domain.BillingOverview, now time.Time) BillingOverviewResponse {
	batches := BillingBatches(overview.Batches, now)
	nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
	return BillingOverviewResponse{
		GeneratedAt:     Timestamp(overview.GeneratedAt),
		PendingInvoices: overview.PendingInvoices,
		OverdueInvoices: overview.OverdueInvoices,
		Scheduler: SchedulerResponse{
			Enabled:   true,
			Schedule:  "daily at 00:00 UTC",
			NextRunAt: Timestamp(nextDay),
		},
		Dunning: DunningResponse{
			Attempts:     3,
			RecoveryRate: 0,
			Description:  "Smart retries on failed invoices",
		},
		Batches: batches,
	}
}

// BillingRunItemResponse is one subscription outcome.
type BillingRunItemResponse struct {
	SubscriptionID string  `json:"subscriptionId"`
	CustomerName   string  `json:"customerName"`
	Action         string  `json:"action"`
	InvoiceID      string  `json:"invoiceId,omitempty"`
	InvoiceNumber  string  `json:"invoiceNumber,omitempty"`
	Amount         float64 `json:"amount"`
	AmountCents    int64   `json:"amountCents"`
	Error          string  `json:"error,omitempty"`
}

// BillingRunResponse is returned by POST /billing/run.
type BillingRunResponse struct {
	AsOf             string                   `json:"asOf"`
	DryRun           bool                     `json:"dryRun"`
	Processed        int                      `json:"processed"`
	Renewed          int                      `json:"renewed"`
	Canceled         int                      `json:"canceled"`
	InvoicesCreated  int                      `json:"invoicesCreated"`
	TotalBilled      float64                  `json:"totalBilled"`
	TotalBilledCents int64                    `json:"totalBilledCents"`
	Failed           int                      `json:"failed"`
	Items            []BillingRunItemResponse `json:"items"`
}

// BillingRun converts a domain billing run.
func BillingRun(result domain.BillingRunResult) BillingRunResponse {
	items := make([]BillingRunItemResponse, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, BillingRunItemResponse{
			SubscriptionID: item.SubscriptionID,
			CustomerName:   item.CustomerName,
			Action:         item.Action,
			InvoiceID:      item.InvoiceID,
			InvoiceNumber:  item.InvoiceNumber,
			Amount:         Major(item.AmountCents),
			AmountCents:    item.AmountCents,
			Error:          item.Error,
		})
	}
	return BillingRunResponse{
		AsOf:             Date(result.AsOf),
		DryRun:           result.DryRun,
		Processed:        result.Processed,
		Renewed:          result.Renewed,
		Canceled:         result.Canceled,
		InvoicesCreated:  result.InvoicesCreated,
		TotalBilled:      Major(result.TotalBilledCents),
		TotalBilledCents: result.TotalBilledCents,
		Failed:           result.Failed,
		Items:            items,
	}
}
