package dto

import (
	"math"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

func toCents(value float64) int64 {
	return int64(math.Round(value * 100))
}

// MessageResponse is a generic success payload.
type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// OKMessage builds a success message response.
func OKMessage(message string) MessageResponse {
	return MessageResponse{Success: true, Message: message}
}

// ProductResponse is an admin catalogue product.
type ProductResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	SKU         string   `json:"sku,omitempty"`
	Active      bool     `json:"active"`
	Features    []string `json:"features"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

// Product converts a domain product.
func Product(p domain.Product) ProductResponse {
	features := p.Features
	if features == nil {
		features = []string{}
	}
	return ProductResponse{
		ID: p.ID, Name: p.Name, Description: p.Description, SKU: p.SKU,
		Active: p.Active, Features: features,
		CreatedAt: Timestamp(p.CreatedAt), UpdatedAt: Timestamp(p.UpdatedAt),
	}
}

// ProductRequest is a product create/update payload.
type ProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	SKU         string   `json:"sku"`
	Active      *bool    `json:"active"`
	Features    []string `json:"features"`
}

func (r ProductRequest) ToCreate() domain.CreateProductInput {
	return domain.CreateProductInput{Name: r.Name, Description: r.Description, SKU: r.SKU, Features: r.Features}
}

func (r ProductRequest) ToUpdate() domain.UpdateProductInput {
	return domain.UpdateProductInput{Name: r.Name, Description: r.Description, SKU: r.SKU, Active: r.Active, Features: r.Features}
}

// DiscountResponse is an admin discount.
type DiscountResponse struct {
	ID                 string  `json:"id"`
	Code               string  `json:"code"`
	Name               string  `json:"name"`
	Type               string  `json:"type"`
	Value              float64 `json:"value"`
	PercentBasisPoints int     `json:"percentBasisPoints"`
	AmountCents        int64   `json:"amountCents"`
	Duration           string  `json:"duration"`
	DurationMonths     int     `json:"durationMonths"`
	UsageLimit         *int    `json:"usageLimit"`
	Redemptions        int     `json:"redemptions"`
	StartsAt           *string `json:"startsAt,omitempty"`
	EndsAt             *string `json:"endsAt,omitempty"`
	Active             bool    `json:"active"`
}

// Discount converts a domain discount.
func Discount(d domain.Discount) DiscountResponse {
	value := Major(d.AmountCents)
	if d.Kind == domain.DiscountPercentage {
		value = float64(d.PercentBasisPoints) / 100
	}
	resp := DiscountResponse{
		ID: d.ID, Code: d.Code, Name: d.Code, Type: string(d.Kind), Value: value,
		PercentBasisPoints: d.PercentBasisPoints, AmountCents: d.AmountCents,
		Duration: string(d.Duration), DurationMonths: d.DurationMonths,
		UsageLimit: d.MaxRedemptions, Redemptions: d.TimesRedeemed,
		StartsAt: TimestampPtr(d.StartsAt), EndsAt: TimestampPtr(d.EndsAt), Active: d.Active,
	}
	return resp
}

// DiscountRequest is a discount create payload.
type DiscountRequest struct {
	Code           string  `json:"code"`
	Description    string  `json:"description"`
	Type           string  `json:"type"`
	Value          float64 `json:"value"`
	Duration       string  `json:"duration"`
	DurationMonths int     `json:"durationMonths"`
	UsageLimit     *int    `json:"usageLimit"`
	StartsAt       *string `json:"startsAt"`
	EndsAt         *string `json:"endsAt"`
	Active         bool    `json:"active"`
}

func (r DiscountRequest) ToDomain() domain.Discount {
	d := domain.Discount{
		Code: r.Code, Description: r.Description,
		Kind: domain.DiscountKind(r.Type), Duration: domain.DiscountDuration(r.Duration),
		DurationMonths: r.DurationMonths, MaxRedemptions: r.UsageLimit, Active: r.Active,
	}
	if r.Type == string(domain.DiscountPercentage) {
		d.PercentBasisPoints = int(toCents(r.Value))
	} else {
		d.AmountCents = toCents(r.Value)
	}
	d.StartsAt = parseTimePtr(r.StartsAt)
	d.EndsAt = parseTimePtr(r.EndsAt)
	return d
}

// OverdueInvoiceResponse is an AR row.
type OverdueInvoiceResponse struct {
	ID            string  `json:"id"`
	InvoiceNumber string  `json:"invoiceNumber"`
	CustomerID    string  `json:"customerId"`
	CustomerName  string  `json:"customerName"`
	CustomerEmail string  `json:"customerEmail"`
	PlanName      string  `json:"planName"`
	DueDate       string  `json:"dueDate"`
	DaysOverdue   int     `json:"daysOverdue"`
	Amount        float64 `json:"amount"`
	AmountCents   int64   `json:"amountCents"`
	Status        string  `json:"status"`
	RemindersSent int     `json:"remindersSent"`
}

// OverdueInvoice converts a domain overdue invoice.
func OverdueInvoice(inv domain.OverdueInvoice) OverdueInvoiceResponse {
	return OverdueInvoiceResponse{
		ID: inv.ID, InvoiceNumber: inv.InvoiceNumber, CustomerID: inv.CustomerID,
		CustomerName: inv.CustomerName, CustomerEmail: inv.CustomerEmail, PlanName: inv.PlanName,
		DueDate: Date(inv.DueDate), DaysOverdue: inv.DaysOverdue,
		Amount: Major(inv.AmountCents), AmountCents: inv.AmountCents,
		Status: inv.Status, RemindersSent: inv.RemindersSent,
	}
}

// CustomerInvoiceResponse is a self-service invoice.
type CustomerInvoiceResponse struct {
	ID             string  `json:"id"`
	InvoiceNumber  string  `json:"invoiceNumber"`
	CustomerID     string  `json:"customerId"`
	SubscriptionID string  `json:"subscriptionId,omitempty"`
	Status         string  `json:"status"`
	Date           string  `json:"date"`
	DueDate        string  `json:"dueDate"`
	Subtotal       float64 `json:"subtotalAmount"`
	Tax            float64 `json:"taxAmount"`
	Discount       float64 `json:"discountAmount"`
	Total          float64 `json:"totalAmount"`
	AmountDue      float64 `json:"amountDue"`
	AmountPaid     float64 `json:"amountPaid"`
	Currency       string  `json:"currency"`
	PaymentMethod  string  `json:"paymentMethod,omitempty"`
	PaidAt         *string `json:"paidAt,omitempty"`
}

// CustomerInvoice converts a domain customer invoice.
func CustomerInvoice(inv domain.CustomerInvoice) CustomerInvoiceResponse {
	return CustomerInvoiceResponse{
		ID: inv.ID, InvoiceNumber: inv.InvoiceNumber, CustomerID: inv.CustomerID,
		SubscriptionID: inv.SubscriptionID, Status: inv.Status,
		Date: Date(inv.Date), DueDate: Date(inv.DueDate),
		Subtotal: Major(inv.SubtotalCents), Tax: Major(inv.TaxCents), Discount: Major(inv.DiscountCents),
		Total: Major(inv.TotalCents), AmountDue: Major(inv.AmountDueCents), AmountPaid: Major(inv.AmountPaidCents),
		Currency: inv.Currency, PaymentMethod: inv.PaymentMethod, PaidAt: TimestampPtr(inv.PaidAt),
	}
}

// NextPaymentResponse is the next renewal for the overview.
type NextPaymentResponse struct {
	Amount             float64 `json:"amount"`
	Date               string  `json:"date"`
	SubscriptionNumber string  `json:"subscriptionNumber"`
	PlanName           string  `json:"planName"`
}

// CustomerOverviewResponse is the customer portal dashboard summary.
type CustomerOverviewResponse struct {
	ActiveSubscriptionsCount int                  `json:"activeSubscriptionsCount"`
	TotalMonthlyValue        float64              `json:"totalMonthlyValue"`
	TotalOutstandingAmount   float64              `json:"totalOutstandingAmount"`
	OutstandingInvoicesCount int                  `json:"outstandingInvoicesCount"`
	NextPaymentDue           *NextPaymentResponse `json:"nextPaymentDue"`
}

// CustomerOverview converts a domain overview.
func CustomerOverview(o domain.CustomerOverview) CustomerOverviewResponse {
	resp := CustomerOverviewResponse{
		ActiveSubscriptionsCount: o.ActiveSubscriptionsCount,
		TotalMonthlyValue:        Major(o.TotalMonthlyValueCents),
		TotalOutstandingAmount:   Major(o.TotalOutstandingCents),
		OutstandingInvoicesCount: o.OutstandingInvoicesCount,
	}
	if o.NextPaymentDue != nil {
		resp.NextPaymentDue = &NextPaymentResponse{
			Amount: Major(o.NextPaymentDue.AmountCents), Date: Date(o.NextPaymentDue.Date),
			SubscriptionNumber: o.NextPaymentDue.SubscriptionNumber, PlanName: o.NextPaymentDue.PlanName,
		}
	}
	return resp
}

func parseTimePtr(s *string) *time.Time {
	if s == nil {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil
	}
	return &t
}
