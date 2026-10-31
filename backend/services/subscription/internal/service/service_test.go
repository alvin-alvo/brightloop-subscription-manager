package service

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

var testNow = time.Date(2026, 10, 31, 12, 0, 0, 0, time.UTC)

func testFixture() (*fakeStore, *SubscriptionService, *BillingService) {
	store := newFakeStore()
	store.customer = domain.Customer{
		ID:       "cus_1",
		Name:     "Acme Enterprises",
		Email:    "billing@acme.test",
		Currency: "USD",
		Status:   "active",
	}
	store.plan = domain.Plan{
		ID:              "plan_1",
		ProductID:       "prod_1",
		Name:            "Growth",
		BillingInterval: domain.IntervalMonthly,
		PriceCents:      10000,
		Currency:        "USD",
		Active:          true,
	}
	subscriptions := NewSubscriptionService(store, store, store)
	subscriptions.now = func() time.Time { return testNow }
	billing := NewBillingService(store, store)
	billing.now = func() time.Time { return testNow }
	return store, subscriptions, billing
}

func staff() domain.Actor {
	return domain.Actor{UserID: "usr_staff", Role: "staff"}
}

func TestCreateSubscriptionDefaults(t *testing.T) {
	store, subscriptions, _ := testFixture()

	sub, err := subscriptions.Create(context.Background(), domain.CreateSubscriptionInput{
		CustomerID: "cus_1",
		PlanID:     "plan_1",
		Quantity:   2,
	}, staff())
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if sub.Status != domain.StatusActive {
		t.Fatalf("status = %s, want active", sub.Status)
	}
	if sub.UnitPriceCents != 10000 {
		t.Fatalf("unit price = %d, want plan price 10000", sub.UnitPriceCents)
	}
	if sub.TotalCents() != 20000 {
		t.Fatalf("total = %d, want 20000", sub.TotalCents())
	}
	if got, want := sub.CurrentPeriodEnd.Format("2006-01-02"), "2026-11-30"; got != want {
		t.Fatalf("period end = %s, want %s", got, want)
	}
	if store.appliedDiscountID != nil {
		t.Fatalf("unexpected discount id")
	}
}

func TestCreateSubscriptionAppliesDiscountCode(t *testing.T) {
	store, subscriptions, _ := testFixture()
	store.discount = &domain.Discount{
		ID:                 "disc_1",
		Code:               "SAVE20",
		Kind:               domain.DiscountPercentage,
		PercentBasisPoints: 2000,
		Active:             true,
	}

	sub, err := subscriptions.Create(context.Background(), domain.CreateSubscriptionInput{
		CustomerID:   "cus_1",
		PlanID:       "plan_1",
		Quantity:     2,
		DiscountCode: "SAVE20",
	}, staff())
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if sub.DiscountCents != 4000 {
		t.Fatalf("discount = %d, want 4000", sub.DiscountCents)
	}
	if store.appliedDiscountID == nil || *store.appliedDiscountID != "disc_1" {
		t.Fatalf("discount was not linked to the subscription")
	}
}

func TestCreateRejectsDiscountAboveSubtotal(t *testing.T) {
	_, subscriptions, _ := testFixture()
	_, err := subscriptions.Create(context.Background(), domain.CreateSubscriptionInput{
		CustomerID:    "cus_1",
		PlanID:        "plan_1",
		Quantity:      1,
		DiscountCents: 50000,
	}, staff())
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestCancelAtPeriodEndAndImmediately(t *testing.T) {
	store, subscriptions, _ := testFixture()
	sub := domain.Subscription{
		ID:                 "sub_1",
		CustomerID:         "cus_1",
		PlanID:             "plan_1",
		Plan:               &store.plan,
		Status:             domain.StatusActive,
		Quantity:           1,
		UnitPriceCents:     10000,
		Currency:           "USD",
		CurrentPeriodStart: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC),
		CurrentPeriodEnd:   time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC),
	}
	store.subscriptions[sub.ID] = clone(sub)

	scheduled, err := subscriptions.Cancel(context.Background(), sub.ID, false, staff())
	if err != nil {
		t.Fatalf("Cancel returned error: %v", err)
	}
	if !scheduled.CancelAtPeriodEnd || scheduled.Status != domain.StatusActive {
		t.Fatalf("scheduled cancel = cancelAtPeriodEnd %v status %s", scheduled.CancelAtPeriodEnd, scheduled.Status)
	}
	if scheduled.CanceledAt == nil {
		t.Fatal("expected canceledAt to be recorded")
	}

	immediate, err := subscriptions.Cancel(context.Background(), sub.ID, true, staff())
	if err != nil {
		t.Fatalf("immediate Cancel returned error: %v", err)
	}
	if immediate.Status != domain.StatusCanceled || immediate.EndedAt == nil {
		t.Fatalf("immediate cancel = status %s endedAt %v", immediate.Status, immediate.EndedAt)
	}
}

func TestRenewAdvancesPeriodAndCreatesInvoice(t *testing.T) {
	store, _, billing := testFixture()
	sub := domain.Subscription{
		ID:                 "sub_1",
		CustomerID:         "cus_1",
		CustomerName:       "Acme Enterprises",
		PlanID:             "plan_1",
		Plan:               &store.plan,
		Status:             domain.StatusActive,
		Quantity:           1,
		UnitPriceCents:     10000,
		Currency:           "USD",
		CurrentPeriodStart: time.Date(2026, 10, 30, 0, 0, 0, 0, time.UTC),
		CurrentPeriodEnd:   time.Date(2026, 11, 30, 0, 0, 0, 0, time.UTC),
	}
	store.subscriptions[sub.ID] = clone(sub)

	updated, invoice, err := billing.Renew(context.Background(), sub.ID, staff())
	if err != nil {
		t.Fatalf("Renew returned error: %v", err)
	}
	if got, want := updated.CurrentPeriodStart.Format("2006-01-02"), "2026-11-30"; got != want {
		t.Fatalf("new period start = %s, want %s", got, want)
	}
	if got, want := updated.CurrentPeriodEnd.Format("2006-01-02"), "2026-12-30"; got != want {
		t.Fatalf("new period end = %s, want %s", got, want)
	}
	if invoice.AmountDueCents != 10000 || invoice.Status != domain.InvoiceOpen {
		t.Fatalf("invoice = amount %d status %s", invoice.AmountDueCents, invoice.Status)
	}
	if len(store.invoiceItems) != 1 {
		t.Fatalf("invoice items = %d, want 1", len(store.invoiceItems))
	}
}

func TestBillingRunRenewsAndCancelsDueSubscriptions(t *testing.T) {
	store, subscriptions, billing := testFixture()
	renew := domain.Subscription{
		ID:                 "sub_renew",
		CustomerID:         "cus_1",
		PlanID:             "plan_1",
		Plan:               &store.plan,
		Status:             domain.StatusActive,
		Quantity:           1,
		UnitPriceCents:     10000,
		Currency:           "USD",
		CurrentPeriodStart: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC),
		CurrentPeriodEnd:   time.Date(2026, 10, 30, 0, 0, 0, 0, time.UTC),
	}
	cancel := renew
	cancel.ID = "sub_cancel"
	cancel.CancelAtPeriodEnd = true
	store.subscriptions[renew.ID] = clone(renew)
	store.subscriptions[cancel.ID] = clone(cancel)

	result, err := billing.Run(context.Background(), testNow, false, staff())
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if result.Processed != 2 || result.Renewed != 1 || result.Canceled != 1 || result.InvoicesCreated != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
	if result.TotalBilledCents != 10000 {
		t.Fatalf("total billed = %d, want 10000", result.TotalBilledCents)
	}
	if store.subscriptions["sub_cancel"].Status != domain.StatusCanceled {
		t.Fatalf("scheduled cancel was not finalized")
	}
	if store.subscriptions["sub_renew"].CurrentPeriodEnd.Format("2006-01-02") != "2026-11-30" {
		t.Fatalf("renewed subscription did not advance")
	}
	if _, err := subscriptions.Get(context.Background(), "sub_renew", staff()); err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
}

func TestBillingRunDryRunDoesNotPersist(t *testing.T) {
	store, _, billing := testFixture()
	sub := domain.Subscription{
		ID:                 "sub_1",
		CustomerID:         "cus_1",
		PlanID:             "plan_1",
		Plan:               &store.plan,
		Status:             domain.StatusActive,
		Quantity:           1,
		UnitPriceCents:     10000,
		Currency:           "USD",
		CurrentPeriodStart: time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC),
		CurrentPeriodEnd:   time.Date(2026, 10, 30, 0, 0, 0, 0, time.UTC),
	}
	store.subscriptions[sub.ID] = clone(sub)

	result, err := billing.Run(context.Background(), testNow, true, staff())
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if result.Renewed != 1 || result.InvoicesCreated != 1 {
		t.Fatalf("unexpected dry run result: %+v", result)
	}
	if store.invoice != nil {
		t.Fatal("dry run persisted an invoice")
	}
	if store.subscriptions[sub.ID].CurrentPeriodEnd != sub.CurrentPeriodEnd {
		t.Fatal("dry run advanced the subscription")
	}
}

func TestComputeAnalytics(t *testing.T) {
	computed := computeAnalytics(domain.Analytics{
		ActiveSubscriptions: 100,
		NewSubscriptions:    10,
		CanceledThisMonth:   5,
		MRRCents:            100000,
		PreviousMRRCents:    90000,
		NetRetainedMRRCents: 95000,
		TotalCustomers:      50,
		NewCustomers:        4,
	})
	assertClose(t, computed.MRRChangePercent, 11.111111111111112)
	assertClose(t, computed.NetRetentionRate, 105.55555555555556)
	assertClose(t, computed.ChurnRate, 5.263157894736842)
	if computed.ARR != 1200000 {
		t.Fatalf("arr = %d, want 1200000", computed.ARR)
	}
}

func assertClose(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestAddMonthsClamped(t *testing.T) {
	jan31 := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	got := addMonthsClamped(jan31, 1).Format("2006-01-02")
	if got != "2026-02-28" {
		t.Fatalf("Jan 31 + 1 month = %s, want 2026-02-28", got)
	}
	leap := addMonthsClamped(time.Date(2028, 1, 31, 0, 0, 0, 0, time.UTC), 1).Format("2006-01-02")
	if leap != "2028-02-29" {
		t.Fatalf("leap year = %s, want 2028-02-29", leap)
	}
}

func TestCustomerScoping(t *testing.T) {
	store, subscriptions, _ := testFixture()
	sub := domain.Subscription{
		ID:         "sub_1",
		CustomerID: "cus_1",
		PlanID:     "plan_1",
		Plan:       &store.plan,
		Status:     domain.StatusActive,
	}
	store.subscriptions[sub.ID] = clone(sub)
	actor := domain.Actor{Role: "customer", CustomerID: "cus_1"}
	got, err := subscriptions.Get(context.Background(), sub.ID, actor)
	if err != nil || got.ID != sub.ID {
		t.Fatalf("customer Get = %v, %v", got, err)
	}
	actor.CustomerID = "cus_other"
	if _, err := subscriptions.Get(context.Background(), sub.ID, actor); err == nil {
		t.Fatal("expected forbidden for another customer's subscription")
	}
}
