package repository

import (
	"context"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

// CatalogRepository provides admin catalogue management (products, discounts)
// and the AR detail / customer self-service views.
type CatalogRepository interface {
	ListProducts(ctx context.Context) ([]domain.Product, error)
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
	CreateProduct(ctx context.Context, product *domain.Product) error
	UpdateProduct(ctx context.Context, product *domain.Product) error
	DeleteProduct(ctx context.Context, id string) error

	ListDiscounts(ctx context.Context, activeOnly bool) ([]domain.Discount, error)
	GetDiscount(ctx context.Context, id string) (*domain.Discount, error)
	CreateDiscount(ctx context.Context, discount *domain.Discount) error
	UpdateDiscount(ctx context.Context, discount *domain.Discount) error

	ListOverdueInvoices(ctx context.Context, overdueOnly bool, limit int) ([]domain.OverdueInvoice, error)
	CustomerInvoices(ctx context.Context, customerID string) ([]domain.CustomerInvoice, error)
	CustomerOverview(ctx context.Context, customerID string) (domain.CustomerOverview, error)
	GetInvoice(ctx context.Context, id string) (*domain.CustomerInvoice, error)
	RecordPayment(ctx context.Context, invoiceID, customerID, method string) (*domain.CustomerInvoice, error)
}
