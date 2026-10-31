package service

import (
	"context"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

// fakeStore implements every repository interface for unit tests.
type fakeStore struct {
	subscriptions     map[string]*domain.Subscription
	customer          domain.Customer
	plan              domain.Plan
	discount          *domain.Discount
	appliedDiscountID *string
	invoice           *domain.Invoice
	invoiceItems      []domain.InvoiceLineItem
	invoicesBySub     []domain.Invoice
	analytics         domain.Analytics
	monthly           []domain.MonthlyRevenue
	plans             []domain.PlanBreakdown
	customers         []domain.TopCustomer
	aging             []domain.AgingBucket
	transactions      []domain.Transaction
	batches           []domain.BillingBatch
	pending           int
	overdue           int
}

func newFakeStore() *fakeStore {
	return &fakeStore{subscriptions: map[string]*domain.Subscription{}}
}

func clone(sub domain.Subscription) *domain.Subscription {
	copy := sub
	if sub.Plan != nil {
		plan := *sub.Plan
		copy.Plan = &plan
	}
	return &copy
}

func (f *fakeStore) List(_ context.Context, filter domain.ListFilter) ([]domain.Subscription, int, error) {
	items := make([]domain.Subscription, 0, len(f.subscriptions))
	for _, sub := range f.subscriptions {
		items = append(items, *clone(*sub))
	}
	return items, len(items), nil
}

func (f *fakeStore) Summary(_ context.Context, _ domain.ListFilter) (domain.SubscriptionSummary, error) {
	summary := domain.SubscriptionSummary{Currency: "USD"}
	for _, sub := range f.subscriptions {
		summary.Total++
		switch sub.Status {
		case domain.StatusActive:
			summary.Active++
		case domain.StatusTrialing:
			summary.Trialing++
		case domain.StatusPastDue:
			summary.PastDue++
		case domain.StatusCanceled:
			summary.Canceled++
		case domain.StatusPaused:
			summary.Paused++
		case domain.StatusUnpaid:
			summary.Unpaid++
		}
		if sub.Status.Billable() {
			summary.MRRCents += sub.MonthlyRecurringCents()
		}
	}
	return summary, nil
}

func (f *fakeStore) GetByID(_ context.Context, id string) (*domain.Subscription, error) {
	sub, ok := f.subscriptions[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return clone(*sub), nil
}

func (f *fakeStore) Create(_ context.Context, sub *domain.Subscription, discountID *string) (*domain.Subscription, error) {
	sub.ID = "sub_test"
	sub.CreatedAt = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	sub.UpdatedAt = sub.CreatedAt
	f.appliedDiscountID = discountID
	f.subscriptions[sub.ID] = clone(*sub)
	return clone(*sub), nil
}

func (f *fakeStore) Update(_ context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	if _, ok := f.subscriptions[sub.ID]; !ok {
		return nil, domain.ErrNotFound
	}
	f.subscriptions[sub.ID] = clone(*sub)
	return clone(*sub), nil
}

func (f *fakeStore) ListDueForBilling(_ context.Context, asOf time.Time) ([]domain.Subscription, error) {
	items := make([]domain.Subscription, 0)
	for _, sub := range f.subscriptions {
		if sub.Status.Billable() && !sub.CurrentPeriodEnd.After(asOf) {
			items = append(items, *clone(*sub))
		}
	}
	return items, nil
}

func (f *fakeStore) ListUpcomingRenewals(_ context.Context, _ time.Time, _ int) ([]domain.Subscription, error) {
	items := make([]domain.Subscription, 0, len(f.subscriptions))
	for _, sub := range f.subscriptions {
		items = append(items, *clone(*sub))
	}
	return items, nil
}

func (f *fakeStore) GetCustomer(_ context.Context, _ string) (*domain.Customer, error) {
	customer := f.customer
	return &customer, nil
}

func (f *fakeStore) GetCustomerByUserID(_ context.Context, _ string) (*domain.Customer, error) {
	customer := f.customer
	return &customer, nil
}

func (f *fakeStore) ListCustomers(_ context.Context, _ string, _ int) ([]domain.Customer, error) {
	return []domain.Customer{f.customer}, nil
}

func (f *fakeStore) GetPlan(_ context.Context, _ string) (*domain.Plan, error) {
	plan := f.plan
	return &plan, nil
}

func (f *fakeStore) ListPlans(_ context.Context, _ bool) ([]domain.Plan, error) {
	return []domain.Plan{f.plan}, nil
}

func (f *fakeStore) GetDiscountByCode(_ context.Context, _ string) (*domain.Discount, error) {
	if f.discount == nil {
		return nil, domain.ErrNotFound
	}
	discount := *f.discount
	return &discount, nil
}

func (f *fakeStore) Renew(_ context.Context, sub *domain.Subscription, invoice *domain.Invoice, items []domain.InvoiceLineItem) error {
	invoice.ID = "inv_test"
	f.invoice = invoice
	f.invoiceItems = items
	f.subscriptions[sub.ID] = clone(*sub)
	return nil
}

func (f *fakeStore) ListInvoicesBySubscription(_ context.Context, _ string) ([]domain.Invoice, error) {
	return f.invoicesBySub, nil
}

func (f *fakeStore) ListTransactions(_ context.Context, filter domain.TransactionFilter) ([]domain.Transaction, int, error) {
	items := f.transactions
	if len(items) > filter.PageSize {
		items = items[:filter.PageSize]
	}
	return items, len(f.transactions), nil
}

func (f *fakeStore) BillingBatches(_ context.Context, _ int) ([]domain.BillingBatch, error) {
	return f.batches, nil
}

func (f *fakeStore) PendingInvoices(_ context.Context) (int, error) { return f.pending, nil }
func (f *fakeStore) OverdueInvoices(_ context.Context) (int, error) { return f.overdue, nil }

func (f *fakeStore) Analytics(_ context.Context) (domain.Analytics, error) {
	return f.analytics, nil
}

func (f *fakeStore) MonthlyRevenue(_ context.Context, _ int) ([]domain.MonthlyRevenue, error) {
	return f.monthly, nil
}

func (f *fakeStore) PlanBreakdown(_ context.Context) ([]domain.PlanBreakdown, error) {
	return f.plans, nil
}

func (f *fakeStore) TopCustomers(_ context.Context, _ int) ([]domain.TopCustomer, error) {
	return f.customers, nil
}

func (f *fakeStore) InvoiceAging(_ context.Context) ([]domain.AgingBucket, error) {
	return f.aging, nil
}
