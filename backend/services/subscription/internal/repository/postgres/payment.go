package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

// GetInvoice returns one invoice (AR / customer portal row).
func (s *Store) GetInvoice(ctx context.Context, id string) (*domain.CustomerInvoice, error) {
	var inv domain.CustomerInvoice
	var paymentMethod sql.NullString
	var paidAt sql.NullTime
	err := s.db.QueryRowContext(ctx, `
		SELECT i.id::text, i.invoice_number, i.customer_id::text,
		       COALESCE(i.subscription_id::text, ''), i.status::text,
		       i.issue_date, i.due_date,
		       i.subtotal_cents, i.tax_cents, i.discount_cents,
		       i.total_cents, i.amount_due_cents, i.amount_paid_cents,
		       trim(i.currency), COALESCE(p.method::text, ''), i.paid_at
		FROM invoices i
		LEFT JOIN LATERAL (
			SELECT method FROM payments WHERE invoice_id = i.id AND status = 'succeeded' ORDER BY created_at DESC LIMIT 1
		) p ON TRUE
		WHERE i.id = $1::uuid`, id).Scan(
		&inv.ID, &inv.InvoiceNumber, &inv.CustomerID, &inv.SubscriptionID, &inv.Status,
		&inv.Date, &inv.DueDate, &inv.SubtotalCents, &inv.TaxCents, &inv.DiscountCents,
		&inv.TotalCents, &inv.AmountDueCents, &inv.AmountPaidCents, &inv.Currency,
		&paymentMethod, &paidAt,
	)
	if err != nil {
		return nil, mapNotFound(err)
	}
	inv.PaymentMethod = paymentMethod.String
	inv.PaidAt = timePtr(paidAt)
	return &inv, nil
}

// RecordPayment settles an open invoice: inserts a succeeded payment and marks
// the invoice paid. Returns the updated invoice.
func (s *Store) RecordPayment(ctx context.Context, invoiceID, customerID, method string) (*domain.CustomerInvoice, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin payment: %w", err)
	}
	defer tx.Rollback()

	if method == "" {
		method = "other"
	}

	var amountDue, total int64
	var currency string
	err = tx.QueryRowContext(ctx,
		`SELECT amount_due_cents, total_cents, trim(currency) FROM invoices WHERE id = $1::uuid`,
		invoiceID).Scan(&amountDue, &total, &currency)
	if err != nil {
		return nil, mapNotFound(err)
	}
	if amountDue <= 0 {
		return nil, domain.Conflictf("invoice is already settled")
	}

	now := time.Now().UTC()
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO payments (invoice_id, customer_id, amount_cents, currency, method, status, paid_at)
		VALUES ($1::uuid, $2::uuid, $3, $4, $5::payment_method, 'succeeded', $6)`,
		invoiceID, customerID, amountDue, currency, method, now); err != nil {
		return nil, fmt.Errorf("insert payment: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE invoices SET status = 'paid', amount_paid_cents = amount_due_cents, paid_at = $2, updated_at = now()
		WHERE id = $1::uuid`, invoiceID, now); err != nil {
		return nil, fmt.Errorf("mark invoice paid: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit payment: %w", err)
	}
	return s.GetInvoice(ctx, invoiceID)
}
