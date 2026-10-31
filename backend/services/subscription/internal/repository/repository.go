package repository

import (
	"context"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

// SubscriptionRepository persists subscriptions.
type SubscriptionRepository interface {
	List(ctx context.Context, filter domain.ListFilter) ([]domain.Subscription, int, error)
	Summary(ctx context.Context, filter domain.ListFilter) (domain.SubscriptionSummary, error)
	GetByID(ctx context.Context, id string) (*domain.Subscription, error)
	Create(ctx context.Context, subscription *domain.Subscription, discountID *string) (*domain.Subscription, error)
	Update(ctx context.Context, subscription *domain.Subscription) (*domain.Subscription, error)
	ListDueForBilling(ctx context.Context, asOf time.Time) ([]domain.Subscription, error)
	ListUpcomingRenewals(ctx context.Context, asOf time.Time, limit int) ([]domain.Subscription, error)
}

// ReferenceRepository reads related aggregates owned by other services
// (customers, plans and discounts) from the shared database.
type ReferenceRepository interface {
	GetCustomer(ctx context.Context, id string) (*domain.Customer, error)
	GetCustomerByUserID(ctx context.Context, userID string) (*domain.Customer, error)
	ListCustomers(ctx context.Context, search string, limit int) ([]domain.Customer, error)
	GetPlan(ctx context.Context, id string) (*domain.Plan, error)
	ListPlans(ctx context.Context, activeOnly bool) ([]domain.Plan, error)
	GetDiscountByCode(ctx context.Context, code string) (*domain.Discount, error)
}

// BillingRepository persists invoices and reads the settlement stream.
type BillingRepository interface {
	Renew(ctx context.Context, subscription *domain.Subscription, invoice *domain.Invoice, items []domain.InvoiceLineItem) error
	ListInvoicesBySubscription(ctx context.Context, subscriptionID string) ([]domain.Invoice, error)
	ListTransactions(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, int, error)
	BillingBatches(ctx context.Context, limit int) ([]domain.BillingBatch, error)
	PendingInvoices(ctx context.Context) (int, error)
	OverdueInvoices(ctx context.Context) (int, error)
}

// AnalyticsRepository reads revenue and subscription aggregates.
type AnalyticsRepository interface {
	Analytics(ctx context.Context) (domain.Analytics, error)
	MonthlyRevenue(ctx context.Context, months int) ([]domain.MonthlyRevenue, error)
	PlanBreakdown(ctx context.Context) ([]domain.PlanBreakdown, error)
	TopCustomers(ctx context.Context, limit int) ([]domain.TopCustomer, error)
	InvoiceAging(ctx context.Context) ([]domain.AgingBucket, error)
}
