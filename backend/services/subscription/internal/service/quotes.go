package service

import (
	"context"
	"errors"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/repository"
)

// QuoteService implements quotations, quote templates and the staff CRM
// workspace.
type QuoteService struct {
	quotes  repository.QuoteRepository
	catalog repository.CatalogRepository
	subs    repository.SubscriptionRepository
}

// NewQuoteService wires the quote business rules.
func NewQuoteService(quotes repository.QuoteRepository, catalog repository.CatalogRepository, subs repository.SubscriptionRepository) *QuoteService {
	return &QuoteService{quotes: quotes, catalog: catalog, subs: subs}
}

// --- Quotes ----------------------------------------------------------------

func (s *QuoteService) ListQuotes(ctx context.Context, actor domain.Actor) ([]domain.Quote, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can view quotes")
	}
	return s.quotes.ListQuotes(ctx)
}

func (s *QuoteService) GetQuote(ctx context.Context, actor domain.Actor, id string) (*domain.Quote, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can view quotes")
	}
	return s.quotes.GetQuote(ctx, id)
}

func (s *QuoteService) CreateQuote(ctx context.Context, actor domain.Actor, quote *domain.Quote) (*domain.Quote, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can create quotes")
	}
	if quote.CustomerID == "" {
		return nil, domain.Validationf("customerId is required")
	}
	if quote.Currency == "" {
		quote.Currency = "USD"
	}
	if err := s.quotes.CreateQuote(ctx, quote); err != nil {
		return nil, err
	}
	return s.quotes.GetQuote(ctx, quote.ID)
}

func (s *QuoteService) UpdateQuoteStatus(ctx context.Context, actor domain.Actor, id, status string) (*domain.Quote, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can update quotes")
	}
	if err := validQuoteStatus(status); err != nil {
		return nil, err
	}
	if err := s.quotes.UpdateQuoteStatus(ctx, id, status); err != nil {
		return nil, err
	}
	return s.quotes.GetQuote(ctx, id)
}

func (s *QuoteService) DeleteQuote(ctx context.Context, actor domain.Actor, id string) error {
	if !actor.IsStaff() {
		return domain.Forbiddenf("only staff can update quotes")
	}
	return s.quotes.DeleteQuote(ctx, id)
}

// --- Templates -------------------------------------------------------------

func (s *QuoteService) ListTemplates(ctx context.Context, actor domain.Actor) ([]domain.QuoteTemplate, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can view templates")
	}
	return s.quotes.ListTemplates(ctx)
}

func (s *QuoteService) GetTemplate(ctx context.Context, actor domain.Actor, id string) (*domain.QuoteTemplate, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can view templates")
	}
	return s.quotes.GetTemplate(ctx, id)
}

func (s *QuoteService) CreateTemplate(ctx context.Context, actor domain.Actor, template *domain.QuoteTemplate) (*domain.QuoteTemplate, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can manage templates")
	}
	if template.Name == "" {
		return nil, domain.Validationf("template name is required")
	}
	if err := s.quotes.CreateTemplate(ctx, template); err != nil {
		if errors.Is(err, domain.ErrConflict) {
			return nil, domain.Conflictf("a template with this configuration already exists")
		}
		return nil, err
	}
	return s.quotes.GetTemplate(ctx, template.ID)
}

// Workspace returns the aggregated staff CRM payload.
func (s *QuoteService) Workspace(ctx context.Context, actor domain.Actor) (domain.Workspace, error) {
	if !actor.IsStaff() {
		return domain.Workspace{}, domain.Forbiddenf("only staff can use the workspace")
	}
	customers, err := s.quotes.ListCustomerWorkspace(ctx)
	if err != nil {
		return domain.Workspace{}, err
	}
	templates, err := s.quotes.ListTemplates(ctx)
	if err != nil {
		return domain.Workspace{}, err
	}
	quotes, err := s.quotes.ListQuotes(ctx)
	if err != nil {
		return domain.Workspace{}, err
	}
	subs, _, err := s.subs.List(ctx, domain.ListFilter{Page: 1, PageSize: 200}.Normalize())
	if err != nil {
		return domain.Workspace{}, err
	}
	invoices, err := s.catalog.ListOverdueInvoices(ctx, false, 100)
	if err != nil {
		return domain.Workspace{}, err
	}
	return domain.Workspace{
		Customers:     customers,
		Templates:     templates,
		Quotes:        quotes,
		Subscriptions: subs,
		Invoices:      invoices,
	}, nil
}

func validQuoteStatus(status string) error {
	switch status {
	case "draft", "sent", "accepted", "declined", "expired":
		return nil
	default:
		return domain.Validationf("quote status must be draft, sent, accepted, declined or expired")
	}
}
