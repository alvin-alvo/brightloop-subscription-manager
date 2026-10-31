package service

import (
	"context"
	"fmt"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/repository"
)

// BillingService implements renewals, recurring billing runs and billing
// operations read models.
type BillingService struct {
	subs    repository.SubscriptionRepository
	billing repository.BillingRepository
	now     func() time.Time
}

// NewBillingService wires the recurring billing engine.
func NewBillingService(subs repository.SubscriptionRepository, billing repository.BillingRepository) *BillingService {
	return &BillingService{subs: subs, billing: billing, now: time.Now}
}

// Renew advances one subscription by a full billing interval and issues an
// invoice for the new period.
func (s *BillingService) Renew(ctx context.Context, id string, actor domain.Actor) (*domain.Subscription, *domain.Invoice, error) {
	if !actor.IsStaff() {
		return nil, nil, domain.Forbiddenf("only staff can renew subscriptions")
	}
	sub, err := s.subs.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if sub.Status == domain.StatusCanceled {
		return nil, nil, domain.Conflictf("a canceled subscription cannot be renewed")
	}
	now := s.now()
	invoice, items, err := s.buildRenewal(sub, sub.CurrentPeriodEnd, now)
	if err != nil {
		return nil, nil, err
	}
	if err := s.billing.Renew(ctx, sub, invoice, items); err != nil {
		return nil, nil, err
	}
	updated, err := s.subs.GetByID(ctx, sub.ID)
	if err != nil {
		return nil, nil, err
	}
	return updated, invoice, nil
}

// Overview returns the billing operations dashboard payload.
func (s *BillingService) Overview(ctx context.Context, actor domain.Actor) (*domain.BillingOverview, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can read billing operations")
	}
	batches, err := s.billing.BillingBatches(ctx, 20)
	if err != nil {
		return nil, err
	}
	pending, err := s.billing.PendingInvoices(ctx)
	if err != nil {
		return nil, err
	}
	overdue, err := s.billing.OverdueInvoices(ctx)
	if err != nil {
		return nil, err
	}
	return &domain.BillingOverview{
		PendingInvoices: pending,
		OverdueInvoices: overdue,
		Batches:         batches,
		GeneratedAt:     s.now(),
	}, nil
}

// Batches returns recurring billing runs derived from invoice issue dates.
func (s *BillingService) Batches(ctx context.Context, limit int, actor domain.Actor) ([]domain.BillingBatch, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can read billing runs")
	}
	return s.billing.BillingBatches(ctx, limit)
}

// Run executes nightly-style recurring billing: due subscriptions are renewed
// and invoiced, or finalized when they were scheduled to cancel.
func (s *BillingService) Run(ctx context.Context, asOf time.Time, dryRun bool, actor domain.Actor) (*domain.BillingRunResult, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can run recurring billing")
	}
	if asOf.IsZero() {
		asOf = s.now()
	}
	asOf = utcDate(asOf)

	due, err := s.subs.ListDueForBilling(ctx, asOf)
	if err != nil {
		return nil, err
	}

	result := &domain.BillingRunResult{
		AsOf:   asOf,
		DryRun: dryRun,
		Items:  make([]domain.BillingRunItem, 0, len(due)),
	}
	result.Processed = len(due)

	for _, sub := range due {
		item := domain.BillingRunItem{
			SubscriptionID: sub.ID,
			CustomerName:   sub.CustomerName,
		}

		if sub.CancelAtPeriodEnd {
			if !dryRun {
				canceledAt := asOf
				sub.Status = domain.StatusCanceled
				sub.CancelAtPeriodEnd = false
				sub.CanceledAt = &canceledAt
				sub.EndedAt = &canceledAt
				if _, err := s.subs.Update(ctx, &sub); err != nil {
					item.Action = "failed"
					item.Error = err.Error()
					result.Failed++
					result.Items = append(result.Items, item)
					continue
				}
			}
			item.Action = "canceled"
			result.Canceled++
			result.Items = append(result.Items, item)
			continue
		}

		invoice, items, err := s.buildRenewal(&sub, sub.CurrentPeriodEnd, asOf)
		if err != nil {
			item.Action = "failed"
			item.Error = err.Error()
			result.Failed++
			result.Items = append(result.Items, item)
			continue
		}
		item.Action = "renewed"
		item.AmountCents = invoice.AmountDueCents
		item.InvoiceNumber = invoice.InvoiceNumber
		result.Renewed++
		result.InvoicesCreated++
		result.TotalBilledCents += invoice.AmountDueCents

		if !dryRun {
			if err := s.billing.Renew(ctx, &sub, invoice, items); err != nil {
				item.Action = "failed"
				item.Error = err.Error()
				result.Renewed--
				result.InvoicesCreated--
				result.TotalBilledCents -= invoice.AmountDueCents
				result.Failed++
			} else {
				item.InvoiceID = invoice.ID
			}
		}
		result.Items = append(result.Items, item)
	}
	return result, nil
}

func (s *BillingService) buildRenewal(sub *domain.Subscription, periodStart, now time.Time) (*domain.Invoice, []domain.InvoiceLineItem, error) {
	start := utcDate(periodStart)
	months := sub.BillingInterval().Months()
	end := addMonthsClamped(start, months)
	if err := validatePeriod(start, end); err != nil {
		return nil, nil, err
	}

	subtotal := sub.SubtotalCents()
	discount := sub.DiscountCents
	if discount > subtotal {
		discount = subtotal
	}
	amountDue := subtotal - discount
	issueDate := utcDate(now)
	dueDate := start
	if dueDate.Before(issueDate) {
		dueDate = issueDate
	}
	currency := sub.Currency
	if currency == "" && sub.Plan != nil {
		currency = sub.Plan.Currency
	}
	if currency == "" {
		currency = "USD"
	}
	name := planName(*sub)
	description := fmt.Sprintf("%s — %s to %s", name, start.Format("Jan 2, 2006"), end.AddDate(0, 0, -1).Format("Jan 2, 2006"))
	if discount > 0 {
		description += " (discount applied)"
	}
	subscriptionID := sub.ID
	invoice := &domain.Invoice{
		InvoiceNumber:  invoiceNumber(now),
		CustomerID:     sub.CustomerID,
		SubscriptionID: &subscriptionID,
		Status:         domain.InvoiceOpen,
		Currency:       currency,
		SubtotalCents:  subtotal,
		DiscountCents:  discount,
		TaxCents:       0,
		AmountDueCents: amountDue,
		IssueDate:      issueDate,
		DueDate:        dueDate,
	}
	planID := sub.PlanID
	items := []domain.InvoiceLineItem{{
		PlanID:         &planID,
		Description:    description,
		Quantity:       sub.Quantity,
		UnitPriceCents: sub.UnitPriceCents,
		AmountCents:    amountDue,
	}}

	// Advance the in-memory aggregate so the repository update is consistent.
	sub.Status = domain.StatusActive
	sub.CurrentPeriodStart = start
	sub.CurrentPeriodEnd = end
	sub.CancelAtPeriodEnd = false
	sub.TrialEndsAt = nil
	sub.CanceledAt = nil
	sub.EndedAt = nil
	return invoice, items, nil
}
