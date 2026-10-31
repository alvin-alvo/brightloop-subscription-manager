package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

// Renew atomically advances a subscription period and records its invoice.
func (s *Store) Renew(ctx context.Context, subscription *domain.Subscription, invoice *domain.Invoice, items []domain.InvoiceLineItem) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin renewal: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		UPDATE subscriptions SET
			status = 'active',
			current_period_start = $1,
			current_period_end = $2,
			cancel_at_period_end = FALSE,
			trial_ends_at = NULL,
			canceled_at = NULL,
			ended_at = NULL,
			updated_at = now()
		WHERE id = $3::uuid`,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.ID,
	)
	if err != nil {
		return fmt.Errorf("advance subscription period: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("advance subscription rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrNotFound
	}

	var notes any
	if invoice.Notes != nil {
		notes = *invoice.Notes
	}
	var subscriptionID any
	if invoice.SubscriptionID != nil {
		subscriptionID = *invoice.SubscriptionID
	}
	err = tx.QueryRowContext(ctx, `
		INSERT INTO invoices (
			invoice_number, customer_id, subscription_id, status, currency,
			subtotal_cents, discount_cents, tax_cents, amount_due_cents,
			amount_paid_cents, issue_date, due_date, notes
		) VALUES (
			$1, $2::uuid, $3::uuid, $4::invoice_status, $5,
			$6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id::text, created_at, updated_at`,
		invoice.InvoiceNumber,
		invoice.CustomerID,
		subscriptionID,
		string(invoice.Status),
		invoice.Currency,
		invoice.SubtotalCents,
		invoice.DiscountCents,
		invoice.TaxCents,
		invoice.AmountDueCents,
		invoice.AmountPaidCents,
		invoice.IssueDate,
		invoice.DueDate,
		notes,
	).Scan(&invoice.ID, &invoice.CreatedAt, &invoice.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert renewal invoice: %w", err)
	}

	for i := range items {
		items[i].InvoiceID = invoice.ID
		var planID any
		if items[i].PlanID != nil {
			planID = *items[i].PlanID
		}
		if err := tx.QueryRowContext(ctx, `
			INSERT INTO invoice_line_items (
				invoice_id, plan_id, description, quantity, unit_price_cents, amount_cents
			) VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6)
			RETURNING id::text, created_at`,
			invoice.ID,
			planID,
			items[i].Description,
			items[i].Quantity,
			items[i].UnitPriceCents,
			items[i].AmountCents,
		).Scan(&items[i].ID, &items[i].CreatedAt); err != nil {
			return fmt.Errorf("insert invoice line item: %w", err)
		}
	}
	invoice.Items = items

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit renewal: %w", err)
	}
	return nil
}

func (s *Store) ListInvoicesBySubscription(ctx context.Context, subscriptionID string) ([]domain.Invoice, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id::text, invoice_number, customer_id::text, subscription_id::text,
		       status::text, trim(currency), subtotal_cents, discount_cents, tax_cents,
		       amount_due_cents, amount_paid_cents, issue_date, due_date, paid_at,
		       notes, created_at, updated_at
		FROM invoices
		WHERE subscription_id = $1::uuid
		ORDER BY issue_date ASC, created_at ASC`, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("list invoices: %w", err)
	}
	defer rows.Close()

	invoices := make([]domain.Invoice, 0)
	for rows.Next() {
		var (
			invoice        domain.Invoice
			subscriptionID sql.NullString
			paidAt         sql.NullTime
			notes          sql.NullString
		)
		if err := rows.Scan(
			&invoice.ID,
			&invoice.InvoiceNumber,
			&invoice.CustomerID,
			&subscriptionID,
			&invoice.Status,
			&invoice.Currency,
			&invoice.SubtotalCents,
			&invoice.DiscountCents,
			&invoice.TaxCents,
			&invoice.AmountDueCents,
			&invoice.AmountPaidCents,
			&invoice.IssueDate,
			&invoice.DueDate,
			&paidAt,
			&notes,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan invoice: %w", err)
		}
		invoice.SubscriptionID = stringPtr(subscriptionID)
		invoice.PaidAt = timePtr(paidAt)
		invoice.Notes = stringPtr(notes)
		invoice.Currency = strings.TrimSpace(invoice.Currency)
		invoices = append(invoices, invoice)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate invoices: %w", err)
	}
	return invoices, nil
}

func buildTransactionWhere(filter domain.TransactionFilter) (string, []any) {
	var conditions []string
	var args []any
	if filter.CustomerID != "" {
		args = append(args, filter.CustomerID)
		conditions = append(conditions, fmt.Sprintf("p.customer_id = $%d::uuid", len(args)))
	}
	if filter.Status != "" {
		args = append(args, string(filter.Status))
		conditions = append(conditions, fmt.Sprintf("p.status = $%d::payment_status", len(args)))
	}
	if len(conditions) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

const transactionFrom = `
	FROM payments p
	JOIN invoices i ON i.id = p.invoice_id
	JOIN customers c ON c.id = p.customer_id
	LEFT JOIN subscriptions s ON s.id = i.subscription_id
	LEFT JOIN plans pl ON pl.id = s.plan_id`

func (s *Store) ListTransactions(ctx context.Context, filter domain.TransactionFilter) ([]domain.Transaction, int, error) {
	filter = filter.Normalize()
	where, args := buildTransactionWhere(filter)

	var total int
	if err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) "+transactionFrom+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count transactions: %w", err)
	}

	offset := (filter.Page - 1) * filter.PageSize
	query := `
		SELECT p.id::text, p.invoice_id::text, i.invoice_number,
		       p.customer_id::text, c.name, c.email,
		       COALESCE(i.subscription_id::text, ''),
		       COALESCE(pl.id::text, ''), COALESCE(pl.name, ''),
		       COALESCE(pl.billing_interval::text, 'monthly'),
		       p.amount_cents, trim(p.currency), p.method::text, p.status::text,
		       COALESCE(p.reference, ''), p.paid_at, p.created_at
		` + transactionFrom + where + fmt.Sprintf(`
		ORDER BY COALESCE(p.paid_at, p.created_at) DESC
		LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	queryArgs := append(append([]any{}, args...), filter.PageSize, offset)

	rows, err := s.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list transactions: %w", err)
	}
	defer rows.Close()

	items := make([]domain.Transaction, 0, filter.PageSize)
	for rows.Next() {
		var (
			tx     domain.Transaction
			paidAt sql.NullTime
		)
		if err := rows.Scan(
			&tx.ID,
			&tx.InvoiceID,
			&tx.InvoiceNumber,
			&tx.CustomerID,
			&tx.CustomerName,
			&tx.CustomerEmail,
			&tx.SubscriptionID,
			&tx.PlanID,
			&tx.PlanName,
			&tx.BillingInterval,
			&tx.AmountCents,
			&tx.Currency,
			&tx.Method,
			&tx.Status,
			&tx.Reference,
			&paidAt,
			&tx.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan transaction: %w", err)
		}
		tx.PaidAt = timePtr(paidAt)
		tx.Currency = strings.TrimSpace(tx.Currency)
		items = append(items, tx)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate transactions: %w", err)
	}
	return items, total, nil
}

func (s *Store) BillingBatches(ctx context.Context, limit int) ([]domain.BillingBatch, error) {
	if limit < 1 {
		limit = 20
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT issue_date, COUNT(*)::int, COALESCE(SUM(amount_due_cents), 0)::bigint,
		       COUNT(*) FILTER (WHERE status = 'paid')::int
		FROM invoices
		GROUP BY issue_date
		ORDER BY issue_date DESC
		LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("list billing batches: %w", err)
	}
	defer rows.Close()

	batches := make([]domain.BillingBatch, 0)
	for rows.Next() {
		var (
			batch     domain.BillingBatch
			issueDate time.Time
		)
		if err := rows.Scan(&issueDate, &batch.ItemsCount, &batch.TotalBilledCents, &batch.PaidCount); err != nil {
			return nil, fmt.Errorf("scan billing batch: %w", err)
		}
		batch.Date = issueDate
		batch.ID = "BATCH-" + issueDate.Format("2006-01-02")
		if batch.ItemsCount > 0 {
			batch.SuccessRatePercent = float64(batch.PaidCount) / float64(batch.ItemsCount) * 100
		}
		batches = append(batches, batch)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate billing batches: %w", err)
	}
	return batches, nil
}

func (s *Store) PendingInvoices(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*)::int FROM invoices WHERE status IN ('open', 'draft')`,
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("count pending invoices: %w", err)
	}
	return count, nil
}

func (s *Store) OverdueInvoices(ctx context.Context) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx,
		`SELECT COUNT(*)::int FROM invoices WHERE status = 'open' AND due_date < CURRENT_DATE`,
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("count overdue invoices: %w", err)
	}
	return count, nil
}
