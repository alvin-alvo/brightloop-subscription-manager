package repository

import (
	"context"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

// QuoteRepository persists quotations, quote templates and the CRM customer view.
type QuoteRepository interface {
	ListQuotes(ctx context.Context) ([]domain.Quote, error)
	GetQuote(ctx context.Context, id string) (*domain.Quote, error)
	CreateQuote(ctx context.Context, quote *domain.Quote) error
	UpdateQuoteStatus(ctx context.Context, id, status string) error
	DeleteQuote(ctx context.Context, id string) error

	ListTemplates(ctx context.Context) ([]domain.QuoteTemplate, error)
	GetTemplate(ctx context.Context, id string) (*domain.QuoteTemplate, error)
	CreateTemplate(ctx context.Context, template *domain.QuoteTemplate) error

	ListCustomerWorkspace(ctx context.Context) ([]domain.CustomerWorkspace, error)
}
