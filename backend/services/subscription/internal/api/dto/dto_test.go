package dto

import (
	"testing"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

func TestFormatMoney(t *testing.T) {
	cases := []struct {
		cents    int64
		currency string
		want     string
	}{
		{4825000, "USD", "$48,250.00"},
		{120000, "USD", "$1,200.00"},
		{99, "USD", "$0.99"},
		{-250, "USD", "-$2.50"},
		{100000, "EUR", "€1,000.00"},
	}
	for _, test := range cases {
		if got := FormatMoney(test.cents, test.currency); got != test.want {
			t.Fatalf("FormatMoney(%d, %s) = %q, want %q", test.cents, test.currency, got, test.want)
		}
	}
}

func TestFormatChange(t *testing.T) {
	if got := FormatChange(14.23); got != "+14.2%" {
		t.Fatalf("FormatChange = %q", got)
	}
	if got := FormatChange(-3.15); got != "-3.2%" {
		t.Fatalf("FormatChange negative = %q", got)
	}
	if got := FormatChange(0); got != "0.0%" {
		t.Fatalf("FormatChange zero = %q", got)
	}
}

func TestDashboardMetricsMatchFrontendCards(t *testing.T) {
	now := time.Date(2026, 10, 31, 10, 0, 0, 0, time.UTC)
	computed := domain.ComputedAnalytics{Analytics: domain.Analytics{
		ActiveSubscriptions: 1284,
		NewSubscriptions:    94,
		MRRCents:            4825000,
		PreviousMRRCents:    4224000,
		NetRetainedMRRCents: 4500000,
		TotalCustomers:      842,
		NewCustomers:        42,
	}}
	computed.NetRetentionRate = 108.4
	response := Dashboard(domain.Dashboard{
		GeneratedAt: now,
		Analytics:   computed,
	}, "USD")

	if len(response.Metrics) != 4 {
		t.Fatalf("metrics = %d, want 4", len(response.Metrics))
	}
	titles := []string{
		"Monthly Recurring Revenue",
		"Active Subscriptions",
		"Net Retention Rate",
		"Total Customers",
	}
	for index, title := range titles {
		if response.Metrics[index].Title != title {
			t.Fatalf("metric[%d] title = %q, want %q", index, response.Metrics[index].Title, title)
		}
	}
	if response.Metrics[0].Value != "$48,250.00" {
		t.Fatalf("mrr value = %q", response.Metrics[0].Value)
	}
	if response.Metrics[2].Value != "108.4%" {
		t.Fatalf("nrr value = %q", response.Metrics[2].Value)
	}
	if response.Raw.TotalCustomers != 842 {
		t.Fatalf("raw total customers = %d", response.Raw.TotalCustomers)
	}
}

func TestTransactionShape(t *testing.T) {
	now := time.Date(2026, 10, 31, 10, 0, 0, 0, time.UTC)
	paid := now
	tx := Transaction(domain.Transaction{
		ID:              "12345678-aaaa-bbbb-cccc-ddddeeeeffff",
		InvoiceID:       "inv_1",
		InvoiceNumber:   "INV-2026-1",
		CustomerID:      "cus_1",
		CustomerName:    "Acme",
		PlanName:        "Growth",
		BillingInterval: domain.IntervalMonthly,
		AmountCents:     49900,
		Currency:        "USD",
		Status:          domain.PaymentSucceeded,
		PaidAt:          &paid,
		CreatedAt:       paid,
	}, now)
	if tx.Status != "success" {
		t.Fatalf("status = %q", tx.Status)
	}
	if tx.AmountLabel != "$499.00" {
		t.Fatalf("amount = %q", tx.AmountLabel)
	}
	if tx.Plan != "Growth (Monthly)" {
		t.Fatalf("plan = %q", tx.Plan)
	}
	if tx.Customer != "Acme" {
		t.Fatalf("customer = %q", tx.Customer)
	}
	if tx.Date != "Today, 10:00 AM" {
		t.Fatalf("date = %q", tx.Date)
	}
}
