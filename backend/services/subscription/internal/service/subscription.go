package service

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/repository"
)

// SubscriptionService implements the subscription lifecycle business rules.
type SubscriptionService struct {
	subs    repository.SubscriptionRepository
	refs    repository.ReferenceRepository
	billing repository.BillingRepository
	now     func() time.Time
}

// NewSubscriptionService wires the subscription business rules.
func NewSubscriptionService(
	subs repository.SubscriptionRepository,
	refs repository.ReferenceRepository,
	billing repository.BillingRepository,
) *SubscriptionService {
	return &SubscriptionService{subs: subs, refs: refs, billing: billing, now: time.Now}
}

// List returns subscriptions visible to the actor.
func (s *SubscriptionService) List(ctx context.Context, filter domain.ListFilter, actor domain.Actor) (domain.SubscriptionPage, error) {
	if filter.Status != "" && !filter.Status.Valid() {
		return domain.SubscriptionPage{}, domain.Validationf("invalid status %q", filter.Status)
	}
	if !actor.IsStaff() {
		customerID, err := s.customerIDForActor(ctx, actor)
		if err != nil {
			return domain.SubscriptionPage{}, err
		}
		filter.CustomerID = customerID
	}
	filter = filter.Normalize()

	items, total, err := s.subs.List(ctx, filter)
	if err != nil {
		return domain.SubscriptionPage{}, err
	}
	summary, err := s.subs.Summary(ctx, filter)
	if err != nil {
		return domain.SubscriptionPage{}, err
	}
	return domain.SubscriptionPage{
		Items:      items,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: pageCount(total, filter.PageSize),
		Summary:    summary,
	}, nil
}

// Get returns one subscription if the actor may see it.
func (s *SubscriptionService) Get(ctx context.Context, id string, actor domain.Actor) (*domain.Subscription, error) {
	if strings.TrimSpace(id) == "" {
		return nil, domain.Validationf("subscription id is required")
	}
	sub, err := s.subs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !actor.IsStaff() {
		customerID, err := s.customerIDForActor(ctx, actor)
		if err != nil {
			return nil, err
		}
		if sub.CustomerID != customerID {
			return nil, domain.Forbiddenf("subscription does not belong to the current customer")
		}
	}
	return sub, nil
}

// Create validates and persists a new subscription.
func (s *SubscriptionService) Create(ctx context.Context, input domain.CreateSubscriptionInput, actor domain.Actor) (*domain.Subscription, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can create subscriptions")
	}
	if strings.TrimSpace(input.CustomerID) == "" {
		return nil, domain.Validationf("customerId is required")
	}
	if strings.TrimSpace(input.PlanID) == "" {
		return nil, domain.Validationf("planId is required")
	}

	customer, err := s.refs.GetCustomer(ctx, input.CustomerID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.Validationf("customer %s does not exist", input.CustomerID)
		}
		return nil, err
	}
	if customer.Status == "suspended" {
		return nil, domain.Validationf("customer is suspended")
	}
	plan, err := s.refs.GetPlan(ctx, input.PlanID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.Validationf("plan %s does not exist", input.PlanID)
		}
		return nil, err
	}
	if !plan.Active {
		return nil, domain.Validationf("plan %s is not active", plan.ID)
	}

	now := s.now()
	quantity := input.Quantity
	if quantity == 0 {
		quantity = 1
	}
	if quantity < 1 {
		return nil, domain.Validationf("quantity must be greater than zero")
	}

	unitPrice := plan.PriceCents
	if input.UnitPriceCents != nil {
		unitPrice = *input.UnitPriceCents
	}
	if unitPrice < 0 {
		return nil, domain.Validationf("unitPrice must not be negative")
	}

	currency := upperCurrency(input.Currency)
	if currency == "" {
		currency = upperCurrency(plan.Currency)
	}
	if currency == "" {
		currency = upperCurrency(customer.Currency)
	}
	if currency == "" {
		currency = "USD"
	}

	start := utcDate(now)
	if input.CurrentPeriodStart != nil {
		start = utcDate(*input.CurrentPeriodStart)
	}
	end := addMonthsClamped(start, plan.BillingInterval.Months())
	if input.CurrentPeriodEnd != nil {
		end = utcDate(*input.CurrentPeriodEnd)
	}
	if err := validatePeriod(start, end); err != nil {
		return nil, err
	}

	status := input.Status
	if status != "" && !status.Valid() {
		return nil, domain.Validationf("invalid status %q", status)
	}
	trialDays := 0
	if input.TrialDays != nil {
		trialDays = *input.TrialDays
	} else if status == "" {
		trialDays = plan.TrialPeriodDays
	}
	if trialDays < 0 {
		return nil, domain.Validationf("trialDays must not be negative")
	}
	if status == "" {
		if trialDays > 0 {
			status = domain.StatusTrialing
		} else {
			status = domain.StatusActive
		}
	}
	var trialEndsAt *time.Time
	if status == domain.StatusTrialing && trialDays > 0 {
		ends := now.AddDate(0, 0, trialDays)
		trialEndsAt = &ends
	}

	subtotal := unitPrice * int64(quantity)
	discountCents := input.DiscountCents
	var discountID *string
	if strings.TrimSpace(input.DiscountCode) != "" {
		if discountCents > 0 {
			return nil, domain.Validationf("provide either discountCode or discountCents, not both")
		}
		discount, err := s.refs.GetDiscountByCode(ctx, input.DiscountCode)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return nil, domain.Validationf("discount %q does not exist", input.DiscountCode)
			}
			return nil, err
		}
		if err := discountUsable(*discount, currency, now); err != nil {
			return nil, err
		}
		discountCents = discountAmountCents(*discount, subtotal)
		discountID = &discount.ID
	}
	if discountCents < 0 {
		return nil, domain.Validationf("discountCents must not be negative")
	}
	if discountCents > subtotal {
		return nil, domain.Validationf("discountCents cannot exceed the subtotal")
	}

	createdBy := input.CreatedBy
	if createdBy == nil && actor.UserID != "" {
		userID := actor.UserID
		createdBy = &userID
	}

	subscription := domain.Subscription{
		CustomerID:         customer.ID,
		PlanID:             plan.ID,
		Plan:               plan,
		Status:             status,
		Quantity:           quantity,
		UnitPriceCents:     unitPrice,
		DiscountCents:      discountCents,
		Currency:           currency,
		CurrentPeriodStart: start,
		CurrentPeriodEnd:   end,
		CancelAtPeriodEnd:  false,
		TrialEndsAt:        trialEndsAt,
		StartedAt:          now,
		CreatedBy:          createdBy,
	}
	return s.subs.Create(ctx, &subscription, discountID)
}

// Update patches mutable subscription fields.
func (s *SubscriptionService) Update(ctx context.Context, id string, input domain.UpdateSubscriptionInput, actor domain.Actor) (*domain.Subscription, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can update subscriptions")
	}
	if strings.TrimSpace(id) == "" {
		return nil, domain.Validationf("subscription id is required")
	}
	sub, err := s.subs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.PlanID != nil && *input.PlanID != "" && *input.PlanID != sub.PlanID {
		plan, err := s.refs.GetPlan(ctx, *input.PlanID)
		if err != nil {
			if errors.Is(err, domain.ErrNotFound) {
				return nil, domain.Validationf("plan %s does not exist", *input.PlanID)
			}
			return nil, err
		}
		if !plan.Active {
			return nil, domain.Validationf("plan %s is not active", plan.ID)
		}
		sub.PlanID = plan.ID
		sub.Plan = plan
		if input.UnitPriceCents == nil {
			sub.UnitPriceCents = plan.PriceCents
		}
	}
	if input.Quantity != nil {
		if *input.Quantity < 1 {
			return nil, domain.Validationf("quantity must be greater than zero")
		}
		sub.Quantity = *input.Quantity
	}
	if input.UnitPriceCents != nil {
		if *input.UnitPriceCents < 0 {
			return nil, domain.Validationf("unitPrice must not be negative")
		}
		sub.UnitPriceCents = *input.UnitPriceCents
	}
	if input.Currency != nil {
		currency := upperCurrency(*input.Currency)
		if currency == "" {
			return nil, domain.Validationf("currency must be a 3-letter ISO code")
		}
		sub.Currency = currency
	}
	if input.CurrentPeriodStart != nil {
		sub.CurrentPeriodStart = utcDate(*input.CurrentPeriodStart)
		if input.CurrentPeriodEnd == nil {
			sub.CurrentPeriodEnd = addMonthsClamped(sub.CurrentPeriodStart, sub.BillingInterval().Months())
		}
	}
	if input.CurrentPeriodEnd != nil {
		sub.CurrentPeriodEnd = utcDate(*input.CurrentPeriodEnd)
	}
	if err := validatePeriod(sub.CurrentPeriodStart, sub.CurrentPeriodEnd); err != nil {
		return nil, err
	}
	if input.DiscountCents != nil {
		if *input.DiscountCents < 0 {
			return nil, domain.Validationf("discountCents must not be negative")
		}
		sub.DiscountCents = *input.DiscountCents
	}
	if sub.DiscountCents > sub.SubtotalCents() {
		return nil, domain.Validationf("discountCents cannot exceed the subtotal")
	}
	if input.CancelAtPeriodEnd != nil {
		sub.CancelAtPeriodEnd = *input.CancelAtPeriodEnd
	}
	if input.Status != nil {
		if !input.Status.Valid() {
			return nil, domain.Validationf("invalid status %q", *input.Status)
		}
		now := s.now()
		sub.Status = *input.Status
		switch sub.Status {
		case domain.StatusCanceled:
			if sub.CanceledAt == nil {
				sub.CanceledAt = &now
			}
			if sub.EndedAt == nil {
				sub.EndedAt = &now
			}
			sub.CancelAtPeriodEnd = false
		case domain.StatusActive:
			sub.CanceledAt = nil
			sub.EndedAt = nil
		}
	}
	return s.subs.Update(ctx, sub)
}

// Cancel schedules or immediately applies a cancellation.
func (s *SubscriptionService) Cancel(ctx context.Context, id string, immediate bool, actor domain.Actor) (*domain.Subscription, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can cancel subscriptions")
	}
	sub, err := s.subs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sub.Status == domain.StatusCanceled {
		return sub, nil
	}
	now := s.now()
	if immediate {
		sub.Status = domain.StatusCanceled
		sub.CancelAtPeriodEnd = false
		sub.CanceledAt = &now
		sub.EndedAt = &now
	} else {
		sub.CancelAtPeriodEnd = true
		sub.CanceledAt = &now
	}
	return s.subs.Update(ctx, sub)
}

// Pause suspends billing without terminating the subscription.
func (s *SubscriptionService) Pause(ctx context.Context, id string, actor domain.Actor) (*domain.Subscription, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can pause subscriptions")
	}
	sub, err := s.subs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sub.Status == domain.StatusPaused {
		return sub, nil
	}
	if sub.Status == domain.StatusCanceled {
		return nil, domain.Conflictf("a canceled subscription cannot be paused")
	}
	sub.Status = domain.StatusPaused
	return s.subs.Update(ctx, sub)
}

// Resume reactivates a paused subscription.
func (s *SubscriptionService) Resume(ctx context.Context, id string, actor domain.Actor) (*domain.Subscription, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can resume subscriptions")
	}
	sub, err := s.subs.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sub.Status == domain.StatusActive {
		return sub, nil
	}
	if sub.Status != domain.StatusPaused {
		return nil, domain.Conflictf("only paused subscriptions can be resumed")
	}
	now := s.now()
	today := utcDate(now)
	if sub.CurrentPeriodEnd.Before(today) || sub.CurrentPeriodEnd.Equal(today) {
		sub.CurrentPeriodStart = today
		sub.CurrentPeriodEnd = addMonthsClamped(today, sub.BillingInterval().Months())
	}
	sub.Status = domain.StatusActive
	sub.CanceledAt = nil
	sub.EndedAt = nil
	return s.subs.Update(ctx, sub)
}

// BillingSchedule returns past/current/upcoming billing cycles.
func (s *SubscriptionService) BillingSchedule(ctx context.Context, id string, actor domain.Actor) (*domain.BillingSchedule, error) {
	sub, err := s.Get(ctx, id, actor)
	if err != nil {
		return nil, err
	}
	invoices, err := s.billing.ListInvoicesBySubscription(ctx, sub.ID)
	if err != nil {
		return nil, err
	}

	months := sub.BillingInterval().Months()
	cycles := make([]domain.BillingCycle, 0, len(invoices)+2)
	totalPaid := int64(0)
	currentCovered := false
	for _, invoice := range invoices {
		start := utcDate(invoice.IssueDate)
		end := addMonthsClamped(start, months)
		cycle := domain.BillingCycle{
			PeriodStart:   start,
			PeriodEnd:     end,
			AmountCents:   invoice.AmountDueCents,
			Status:        invoiceCycleStatus(invoice.Status),
			InvoiceID:     invoice.ID,
			InvoiceNumber: invoice.InvoiceNumber,
			IssuedAt:      &invoice.IssueDate,
			PaidAt:        invoice.PaidAt,
		}
		cycles = append(cycles, cycle)
		totalPaid += invoice.AmountPaidCents
		if !start.Before(sub.CurrentPeriodStart) && start.Before(sub.CurrentPeriodEnd) {
			currentCovered = true
		}
	}

	if !currentCovered && sub.Status != domain.StatusCanceled {
		status := "current"
		if sub.CurrentPeriodStart.After(utcDate(s.now())) {
			status = "upcoming"
		}
		cycles = append(cycles, domain.BillingCycle{
			PeriodStart: sub.CurrentPeriodStart,
			PeriodEnd:   sub.CurrentPeriodEnd,
			AmountCents: sub.TotalCents(),
			Status:      status,
		})
	}

	if sub.Status == domain.StatusActive || sub.Status == domain.StatusTrialing || sub.Status == domain.StatusPastDue {
		cycles = append(cycles, domain.BillingCycle{
			PeriodStart: sub.CurrentPeriodEnd,
			PeriodEnd:   addMonthsClamped(sub.CurrentPeriodEnd, months),
			AmountCents: sub.TotalCents(),
			Status:      "upcoming",
		})
	}

	sort.SliceStable(cycles, func(i, j int) bool {
		return cycles[i].PeriodStart.Before(cycles[j].PeriodStart)
	})

	return &domain.BillingSchedule{
		SubscriptionID:     sub.ID,
		PlanID:             sub.PlanID,
		BillingInterval:    sub.BillingInterval(),
		Currency:           sub.Currency,
		CurrentPeriodStart: sub.CurrentPeriodStart,
		CurrentPeriodEnd:   sub.CurrentPeriodEnd,
		NextBillingDate:    sub.CurrentPeriodEnd,
		Cycles:             cycles,
		TotalPaidCents:     totalPaid,
	}, nil
}

// Plans returns active plans for the assignment form.
func (s *SubscriptionService) Plans(ctx context.Context, activeOnly bool, _ domain.Actor) ([]domain.Plan, error) {
	return s.refs.ListPlans(ctx, activeOnly)
}

// Customers returns customers eligible for a subscription (staff only).
func (s *SubscriptionService) Customers(ctx context.Context, search string, actor domain.Actor) ([]domain.Customer, error) {
	if !actor.IsStaff() {
		return nil, domain.Forbiddenf("only staff can list customers")
	}
	return s.refs.ListCustomers(ctx, search, 200)
}

func (s *SubscriptionService) customerIDForActor(ctx context.Context, actor domain.Actor) (string, error) {
	if actor.CustomerID != "" {
		return actor.CustomerID, nil
	}
	if actor.UserID == "" {
		return "", domain.Forbiddenf("customer identity is not available")
	}
	customer, err := s.refs.GetCustomerByUserID(ctx, actor.UserID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.Forbiddenf("no customer profile is linked to this account")
		}
		return "", err
	}
	return customer.ID, nil
}

func invoiceCycleStatus(status domain.InvoiceStatus) string {
	switch status {
	case domain.InvoicePaid:
		return "paid"
	case domain.InvoiceOpen:
		return "open"
	case domain.InvoiceVoid:
		return "void"
	case domain.InvoiceUncollectible:
		return "uncollectible"
	default:
		return "draft"
	}
}
