package dto

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

// -----------------------------------------------------------------------------
// Money and formatting helpers
// -----------------------------------------------------------------------------

// Major converts minor units to the decimal number expected by the frontend.
func Major(cents int64) float64 {
	return float64(cents) / 100
}

// FormatMoney renders minor units as a human readable currency string.
func FormatMoney(cents int64, currency string) string {
	sign := ""
	if cents < 0 {
		sign = "-"
		cents = -cents
	}
	whole := cents / 100
	fraction := cents % 100
	grouped := groupDigits(whole)
	return fmt.Sprintf("%s%s%s.%02d", sign, currencySymbol(currency), grouped, fraction)
}

func currencySymbol(currency string) string {
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "", "USD":
		return "$"
	case "EUR":
		return "€"
	case "GBP":
		return "£"
	case "JPY":
		return "¥"
	case "INR":
		return "₹"
	case "CAD":
		return "CA$"
	case "AUD":
		return "A$"
	default:
		return strings.ToUpper(currency) + " "
	}
}

func groupDigits(value int64) string {
	raw := strconv.FormatInt(value, 10)
	if len(raw) <= 3 {
		return raw
	}
	var builder strings.Builder
	first := len(raw) % 3
	if first == 0 {
		first = 3
	}
	builder.WriteString(raw[:first])
	for i := first; i < len(raw); i += 3 {
		builder.WriteByte(',')
		builder.WriteString(raw[i : i+3])
	}
	return builder.String()
}

// FormatChange renders a signed percentage for dashboard metric cards.
func FormatChange(percent float64) string {
	rounded := math.Round(percent*10) / 10
	if rounded > 0 {
		return fmt.Sprintf("+%.1f%%", rounded)
	}
	return fmt.Sprintf("%.1f%%", rounded)
}

// FormatPercent renders an unsigned percentage with one decimal.
func FormatPercent(percent float64) string {
	return fmt.Sprintf("%.1f%%", math.Round(percent*10)/10)
}

// FormatRate renders a percentage used by billing runs.
func FormatRate(percent float64) string {
	return fmt.Sprintf("%.1f%%", math.Round(percent*10)/10)
}

// Date formats billing dates as YYYY-MM-DD.
func Date(value time.Time) string {
	return value.UTC().Format("2006-01-02")
}

// Timestamp formats audit timestamps as RFC3339 UTC.
func Timestamp(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}

// TimestampPtr formats an optional timestamp.
func TimestampPtr(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := Timestamp(*value)
	return &formatted
}

// TransactionDate renders the relative date shown in the settlement stream.
func TransactionDate(value, now time.Time) string {
	value = value.UTC()
	now = now.UTC()
	valueDate := time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	days := int(nowDate.Sub(valueDate).Hours() / 24)
	switch days {
	case 0:
		return "Today, " + value.Format("3:04 PM")
	case 1:
		return "Yesterday"
	default:
		return value.Format("Jan 2, 2006")
	}
}

func parseInt(value string) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0
	}
	return parsed
}

// -----------------------------------------------------------------------------
// Resource responses
// -----------------------------------------------------------------------------

// PlanResponse mirrors the frontend Plan interface, with cents fields added for
// lossless consumers.
type PlanResponse struct {
	ID              string   `json:"id"`
	ProductID       string   `json:"productId"`
	ProductName     string   `json:"productName,omitempty"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	BillingInterval string   `json:"billingInterval"`
	Price           float64  `json:"price"`
	PriceCents      int64    `json:"priceCents"`
	Currency        string   `json:"currency"`
	Features        []string `json:"features"`
	TrialPeriodDays int      `json:"trialPeriodDays"`
	Active          bool     `json:"active"`
	CreatedAt       string   `json:"createdAt"`
	UpdatedAt       string   `json:"updatedAt"`
}

// Plan converts a domain plan.
func Plan(plan domain.Plan) PlanResponse {
	features := plan.Features
	if features == nil {
		features = []string{}
	}
	return PlanResponse{
		ID:              plan.ID,
		ProductID:       plan.ProductID,
		ProductName:     plan.ProductName,
		Name:            plan.Name,
		Description:     plan.Description,
		BillingInterval: string(plan.BillingInterval),
		Price:           Major(plan.PriceCents),
		PriceCents:      plan.PriceCents,
		Currency:        plan.Currency,
		Features:        features,
		TrialPeriodDays: plan.TrialPeriodDays,
		Active:          plan.Active,
		CreatedAt:       Timestamp(plan.CreatedAt),
		UpdatedAt:       Timestamp(plan.UpdatedAt),
	}
}

// CustomerResponse is the compact customer shape used by assignment forms.
type CustomerResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CompanyName string `json:"companyName,omitempty"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
}

// Customer converts a domain customer.
func Customer(customer domain.Customer) CustomerResponse {
	return CustomerResponse{
		ID:          customer.ID,
		Name:        customer.Name,
		Email:       customer.Email,
		CompanyName: customer.CompanyName,
		Currency:    customer.Currency,
		Status:      customer.Status,
	}
}

// SubscriptionResponse mirrors the frontend Subscription interface.
type SubscriptionResponse struct {
	ID                 string        `json:"id"`
	CustomerID         string        `json:"customerId"`
	CustomerName       string        `json:"customerName,omitempty"`
	CustomerEmail      string        `json:"customerEmail,omitempty"`
	PlanID             string        `json:"planId"`
	Plan               *PlanResponse `json:"plan,omitempty"`
	Status             string        `json:"status"`
	CurrentPeriodStart string        `json:"currentPeriodStart"`
	CurrentPeriodEnd   string        `json:"currentPeriodEnd"`
	CancelAtPeriodEnd  bool          `json:"cancelAtPeriodEnd"`
	Quantity           int           `json:"quantity"`
	UnitPrice          float64       `json:"unitPrice"`
	UnitPriceCents     int64         `json:"unitPriceCents"`
	Discount           float64       `json:"discount"`
	DiscountCents      int64         `json:"discountCents"`
	TotalAmount        float64       `json:"totalAmount"`
	TotalAmountCents   int64         `json:"totalAmountCents"`
	MRR                float64       `json:"mrr"`
	MRRCents           int64         `json:"mrrCents"`
	Currency           string        `json:"currency"`
	BillingInterval    string        `json:"billingInterval"`
	BillingCycle       string        `json:"billingCycle"`
	NextBillingDate    string        `json:"nextBillingDate"`
	TrialEndsAt        *string       `json:"trialEndsAt,omitempty"`
	StartedAt          string        `json:"startedAt"`
	CanceledAt         *string       `json:"canceledAt,omitempty"`
	EndedAt            *string       `json:"endedAt,omitempty"`
	CreatedAt          string        `json:"createdAt"`
	UpdatedAt          string        `json:"updatedAt"`
}

// Subscription converts a domain subscription.
func Subscription(sub domain.Subscription) SubscriptionResponse {
	response := SubscriptionResponse{
		ID:                 sub.ID,
		CustomerID:         sub.CustomerID,
		CustomerName:       sub.CustomerName,
		CustomerEmail:      sub.CustomerEmail,
		PlanID:             sub.PlanID,
		Status:             string(sub.Status),
		CurrentPeriodStart: Date(sub.CurrentPeriodStart),
		CurrentPeriodEnd:   Date(sub.CurrentPeriodEnd),
		CancelAtPeriodEnd:  sub.CancelAtPeriodEnd,
		Quantity:           sub.Quantity,
		UnitPrice:          Major(sub.UnitPriceCents),
		UnitPriceCents:     sub.UnitPriceCents,
		Discount:           Major(sub.DiscountCents),
		DiscountCents:      sub.DiscountCents,
		TotalAmount:        Major(sub.TotalCents()),
		TotalAmountCents:   sub.TotalCents(),
		MRR:                Major(sub.MonthlyRecurringCents()),
		MRRCents:           sub.MonthlyRecurringCents(),
		Currency:           sub.Currency,
		BillingInterval:    string(sub.BillingInterval()),
		BillingCycle:       sub.BillingInterval().Title(),
		NextBillingDate:    Date(sub.CurrentPeriodEnd),
		TrialEndsAt:        TimestampPtr(sub.TrialEndsAt),
		StartedAt:          Timestamp(sub.StartedAt),
		CanceledAt:         TimestampPtr(sub.CanceledAt),
		EndedAt:            TimestampPtr(sub.EndedAt),
		CreatedAt:          Timestamp(sub.CreatedAt),
		UpdatedAt:          Timestamp(sub.UpdatedAt),
	}
	if sub.Plan != nil {
		plan := Plan(*sub.Plan)
		response.Plan = &plan
	}
	return response
}

// SubscriptionSummaryResponse aggregates a subscription list.
type SubscriptionSummaryResponse struct {
	Total    int     `json:"total"`
	Active   int     `json:"active"`
	Trialing int     `json:"trialing"`
	PastDue  int     `json:"pastDue"`
	Canceled int     `json:"canceled"`
	Paused   int     `json:"paused"`
	Unpaid   int     `json:"unpaid"`
	MRR      float64 `json:"mrr"`
	MRRCents int64   `json:"mrrCents"`
	Currency string  `json:"currency"`
}

// SubscriptionSummary converts a domain summary.
func SubscriptionSummary(summary domain.SubscriptionSummary) SubscriptionSummaryResponse {
	currency := summary.Currency
	if currency == "" {
		currency = "USD"
	}
	return SubscriptionSummaryResponse{
		Total:    summary.Total,
		Active:   summary.Active,
		Trialing: summary.Trialing,
		PastDue:  summary.PastDue,
		Canceled: summary.Canceled,
		Paused:   summary.Paused,
		Unpaid:   summary.Unpaid,
		MRR:      Major(summary.MRRCents),
		MRRCents: summary.MRRCents,
		Currency: currency,
	}
}

// SubscriptionPageResponse is the paginated list contract.
type SubscriptionPageResponse struct {
	Items      []SubscriptionResponse      `json:"items"`
	Total      int                         `json:"total"`
	Page       int                         `json:"page"`
	PageSize   int                         `json:"pageSize"`
	TotalPages int                         `json:"totalPages"`
	Summary    SubscriptionSummaryResponse `json:"summary"`
}

// SubscriptionPage converts a domain page.
func SubscriptionPage(page domain.SubscriptionPage) SubscriptionPageResponse {
	items := make([]SubscriptionResponse, 0, len(page.Items))
	for _, sub := range page.Items {
		items = append(items, Subscription(sub))
	}
	return SubscriptionPageResponse{
		Items:      items,
		Total:      page.Total,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: page.TotalPages,
		Summary:    SubscriptionSummary(page.Summary),
	}
}

// TransactionResponse matches the dashboard settlement stream and adds the
// structured fields required by other consumers.
type TransactionResponse struct {
	ID             string  `json:"id"`
	InvoiceID      string  `json:"invoiceId"`
	InvoiceNumber  string  `json:"invoiceNumber"`
	CustomerID     string  `json:"customerId"`
	CustomerName   string  `json:"customerName"`
	CustomerEmail  string  `json:"customerEmail,omitempty"`
	SubscriptionID string  `json:"subscriptionId,omitempty"`
	PlanID         string  `json:"planId,omitempty"`
	PlanName       string  `json:"planName,omitempty"`
	Amount         float64 `json:"amountValue"`
	AmountCents    int64   `json:"amountCents"`
	AmountLabel    string  `json:"amount"`
	Currency       string  `json:"currency"`
	Method         string  `json:"method"`
	PaymentStatus  string  `json:"paymentStatus"`
	Status         string  `json:"status"`
	Reference      string  `json:"reference,omitempty"`
	PaidAt         *string `json:"paidAt,omitempty"`
	CreatedAt      string  `json:"createdAt"`
	Date           string  `json:"date"`
	Customer       string  `json:"customer"`
	Plan           string  `json:"plan"`
}

// Transaction converts a domain transaction.
func Transaction(tx domain.Transaction, now time.Time) TransactionResponse {
	paidOrCreated := tx.CreatedAt
	if tx.PaidAt != nil {
		paidOrCreated = *tx.PaidAt
	}
	plan := tx.PlanName
	if tx.BillingInterval.Valid() {
		plan = fmt.Sprintf("%s (%s)", tx.PlanName, tx.BillingInterval.Title())
	}
	shortID := tx.ID
	if len(shortID) > 8 {
		shortID = shortID[:8]
	}
	return TransactionResponse{
		ID:             "tx_" + strings.ReplaceAll(shortID, "-", ""),
		InvoiceID:      tx.InvoiceID,
		InvoiceNumber:  tx.InvoiceNumber,
		CustomerID:     tx.CustomerID,
		CustomerName:   tx.CustomerName,
		CustomerEmail:  tx.CustomerEmail,
		SubscriptionID: tx.SubscriptionID,
		PlanID:         tx.PlanID,
		PlanName:       tx.PlanName,
		Amount:         Major(tx.AmountCents),
		AmountCents:    tx.AmountCents,
		AmountLabel:    FormatMoney(tx.AmountCents, tx.Currency),
		Currency:       tx.Currency,
		Method:         tx.Method,
		PaymentStatus:  string(tx.Status),
		Status:         paymentDisplayStatus(tx.Status),
		Reference:      tx.Reference,
		PaidAt:         TimestampPtr(tx.PaidAt),
		CreatedAt:      Timestamp(tx.CreatedAt),
		Date:           TransactionDate(paidOrCreated, now),
		Customer:       tx.CustomerName,
		Plan:           plan,
	}
}

func paymentDisplayStatus(status domain.PaymentStatus) string {
	switch status {
	case domain.PaymentSucceeded:
		return "success"
	case domain.PaymentPending:
		return "warning"
	case domain.PaymentFailed:
		return "danger"
	case domain.PaymentRefunded:
		return "neutral"
	default:
		return "neutral"
	}
}

// TransactionPageResponse is the paginated settlement stream.
type TransactionPageResponse struct {
	Items      []TransactionResponse `json:"items"`
	Total      int                   `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"pageSize"`
	TotalPages int                   `json:"totalPages"`
}

// TransactionPage converts a domain page.
func TransactionPage(page domain.TransactionPage, now time.Time) TransactionPageResponse {
	items := make([]TransactionResponse, 0, len(page.Items))
	for _, item := range page.Items {
		items = append(items, Transaction(item, now))
	}
	return TransactionPageResponse{
		Items:      items,
		Total:      page.Total,
		Page:       page.Page,
		PageSize:   page.PageSize,
		TotalPages: page.TotalPages,
	}
}
