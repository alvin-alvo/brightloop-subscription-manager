package dto

import (
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

// QuoteResponse matches the frontend StaffQuote shape plus structured fields.
type QuoteResponse struct {
	ID            string              `json:"id"`
	QuoteNumber   string              `json:"quoteNumber"`
	CustomerID    string              `json:"customerId"`
	Customer      string              `json:"customer"`
	Plan          string              `json:"plan"`
	AmountCents   int64               `json:"amountCents"`
	ValidUntil    string              `json:"validUntil"`
	Status        string              `json:"status"`
	Currency      string              `json:"currency"`
	SubtotalCents int64               `json:"subtotalCents"`
	DiscountCents int64               `json:"discountCents"`
	TaxCents      int64               `json:"taxCents"`
	Items         []QuoteItemResponse `json:"items"`
	CreatedAt     string              `json:"createdAt"`
	UpdatedAt     string              `json:"updatedAt"`
}

// QuoteItemResponse is one quote line.
type QuoteItemResponse struct {
	ID             string  `json:"id"`
	PlanID         string  `json:"planId,omitempty"`
	Description    string  `json:"description"`
	Quantity       int     `json:"quantity"`
	UnitPrice      float64 `json:"unitPrice"`
	UnitPriceCents int64   `json:"unitPriceCents"`
	TaxBasisPoints int64   `json:"taxRateBasisPoints"`
	LineTotalCents int64   `json:"lineTotalCents"`
}

// Quote converts a domain quote.
func Quote(q domain.Quote) QuoteResponse {
	items := make([]QuoteItemResponse, 0, len(q.Items))
	for _, item := range q.Items {
		var planID string
		if item.PlanID != nil {
			planID = *item.PlanID
		}
		items = append(items, QuoteItemResponse{
			ID: item.ID, PlanID: planID, Description: item.Description, Quantity: item.Quantity,
			UnitPrice: Major(item.UnitPriceCents), UnitPriceCents: item.UnitPriceCents,
			TaxBasisPoints: item.TaxBasisPoints, LineTotalCents: item.LineTotalCents,
		})
	}
	return QuoteResponse{
		ID: q.ID, QuoteNumber: q.QuoteNumber, CustomerID: q.CustomerID, Customer: q.CustomerName,
		Plan: q.Plan(), AmountCents: q.TotalCents, ValidUntil: Date(q.ValidUntil),
		Status: quoteStatus(q.Status), Currency: q.Currency,
		SubtotalCents: q.SubtotalCents, DiscountCents: q.DiscountCents, TaxCents: q.TaxCents,
		Items: items, CreatedAt: Timestamp(q.CreatedAt), UpdatedAt: Timestamp(q.UpdatedAt),
	}
}

func quoteStatus(status string) string {
	switch status {
	case "accepted":
		return "Approved"
	case "expired":
		return "Expired"
	default:
		return "Pending"
	}
}

// QuoteRequest is a quote create payload.
type QuoteRequest struct {
	CustomerID string           `json:"customerId"`
	Notes      string           `json:"notes"`
	ValidUntil string           `json:"validUntil"`
	Items      []QuoteItemInput `json:"items"`
}

// QuoteItemInput is one line on create.
type QuoteItemInput struct {
	PlanID         string `json:"planId"`
	Description    string `json:"description"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents int64  `json:"unitPriceCents"`
	TaxBasisPoints int64  `json:"taxRateBasisPoints"`
}

// ToDomain converts a quote request.
func (r QuoteRequest) ToDomain() domain.Quote {
	q := domain.Quote{CustomerID: r.CustomerID, Status: "draft", Currency: "USD", Notes: r.Notes}
	if r.ValidUntil != "" {
		if t, err := time.Parse("2006-01-02", r.ValidUntil); err == nil {
			q.ValidUntil = t
		}
	}
	subtotal := int64(0)
	for _, in := range r.Items {
		item := domain.QuoteItem{
			Description: in.Description, Quantity: in.Quantity,
			UnitPriceCents: in.UnitPriceCents, TaxBasisPoints: in.TaxBasisPoints,
			LineTotalCents: in.UnitPriceCents * int64(in.Quantity),
		}
		if in.PlanID != "" {
			item.PlanID = &in.PlanID
		}
		if item.Description == "" {
			item.Description = "Subscription"
		}
		if item.Quantity < 1 {
			item.Quantity = 1
		}
		subtotal += item.LineTotalCents
		q.Items = append(q.Items, item)
	}
	q.SubtotalCents = subtotal
	q.TotalCents = subtotal
	return q
}

// TemplateResponse matches the frontend QuoteTemplate shape.
type TemplateResponse struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Plan         string                 `json:"plan"`
	ValidityDays int                    `json:"validityDays"`
	Lines        []TemplateItemResponse `json:"lines"`
	CreatedAt    string                 `json:"createdAt"`
	UpdatedAt    string                 `json:"updatedAt"`
}

// TemplateItemResponse is one template line.
type TemplateItemResponse struct {
	ID                 string `json:"id"`
	ProductID          string `json:"productId,omitempty"`
	Product            string `json:"product"`
	Quantity           int    `json:"quantity"`
	UnitPriceCents     int64  `json:"unitPriceCents"`
	TaxRateBasisPoints int64  `json:"taxRateBasisPoints"`
}

// Template converts a domain template.
func Template(t domain.QuoteTemplate) TemplateResponse {
	lines := make([]TemplateItemResponse, 0, len(t.Items))
	for _, item := range t.Items {
		var productID string
		if item.ProductID != nil {
			productID = *item.ProductID
		}
		lines = append(lines, TemplateItemResponse{
			ID: item.ID, ProductID: productID, Product: item.ProductName,
			Quantity: item.Quantity, UnitPriceCents: item.UnitPriceCents,
			TaxRateBasisPoints: item.TaxBasisPoints,
		})
	}
	return TemplateResponse{
		ID: t.ID, Name: t.Name, Plan: t.PlanName, ValidityDays: t.ValidityDays,
		Lines: lines, CreatedAt: Timestamp(t.CreatedAt), UpdatedAt: Timestamp(t.UpdatedAt),
	}
}

// TemplateRequest is a template create payload.
type TemplateRequest struct {
	Name         string              `json:"name"`
	Plan         string              `json:"plan"`
	ValidityDays int                 `json:"validityDays"`
	Lines        []TemplateItemInput `json:"lines"`
}

// TemplateItemInput is one template line on create.
type TemplateItemInput struct {
	ProductID          string `json:"productId"`
	Product            string `json:"product"`
	Quantity           int    `json:"quantity"`
	UnitPriceCents     int64  `json:"unitPriceCents"`
	TaxRateBasisPoints int64  `json:"taxRateBasisPoints"`
}

// ToDomain converts a template request.
func (r TemplateRequest) ToDomain() domain.QuoteTemplate {
	t := domain.QuoteTemplate{Name: r.Name, ValidityDays: r.ValidityDays}
	if r.Plan != "" {
		t.PlanName = r.Plan
	}
	for _, in := range r.Lines {
		t.Items = append(t.Items, domain.QuoteTemplateItem{
			ProductName: in.Product, Quantity: in.Quantity,
			UnitPriceCents: in.UnitPriceCents, TaxBasisPoints: in.TaxRateBasisPoints,
		})
	}
	return t
}

// CustomerWorkspaceResponse is a CRM row.
type CustomerWorkspaceResponse struct {
	ID                  string  `json:"id"`
	Name                string  `json:"name"`
	Contact             string  `json:"contact"`
	Email               string  `json:"email"`
	ActiveSubscriptions int     `json:"activeSubscriptions"`
	MonthlyValue        float64 `json:"monthlyValue"`
	MonthlyValueCents   int64   `json:"monthlyValueCents"`
	Tier                string  `json:"tier"`
	Status              string  `json:"status"`
}

func customerWorkspace(cw domain.CustomerWorkspace) CustomerWorkspaceResponse {
	return CustomerWorkspaceResponse{
		ID: cw.ID, Name: cw.Name, Contact: cw.Contact, Email: cw.Email,
		ActiveSubscriptions: cw.ActiveSubscriptions,
		MonthlyValue:        Major(cw.MonthlyValueCents), MonthlyValueCents: cw.MonthlyValueCents,
		Tier: cw.Tier, Status: cw.Status,
	}
}

// WorkspaceResponse is the aggregated staff CRM payload.
type WorkspaceResponse struct {
	Customers     []CustomerWorkspaceResponse `json:"customers"`
	Templates     []TemplateResponse          `json:"templates"`
	Quotes        []QuoteResponse             `json:"quotes"`
	Subscriptions []SubscriptionResponse      `json:"subscriptions"`
	Invoices      []OverdueInvoiceResponse    `json:"invoices"`
}

// Workspace converts a domain workspace.
func Workspace(w domain.Workspace) WorkspaceResponse {
	resp := WorkspaceResponse{
		Customers:     make([]CustomerWorkspaceResponse, 0, len(w.Customers)),
		Templates:     make([]TemplateResponse, 0, len(w.Templates)),
		Quotes:        make([]QuoteResponse, 0, len(w.Quotes)),
		Subscriptions: make([]SubscriptionResponse, 0, len(w.Subscriptions)),
		Invoices:      make([]OverdueInvoiceResponse, 0, len(w.Invoices)),
	}
	for _, cw := range w.Customers {
		resp.Customers = append(resp.Customers, customerWorkspace(cw))
	}
	for _, t := range w.Templates {
		resp.Templates = append(resp.Templates, Template(t))
	}
	for _, q := range w.Quotes {
		resp.Quotes = append(resp.Quotes, Quote(q))
	}
	for _, s := range w.Subscriptions {
		resp.Subscriptions = append(resp.Subscriptions, Subscription(s))
	}
	for _, inv := range w.Invoices {
		resp.Invoices = append(resp.Invoices, OverdueInvoice(inv))
	}
	return resp
}
