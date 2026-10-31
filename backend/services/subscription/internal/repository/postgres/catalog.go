package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

// ---------------------------------------------------------------------------
// Products
// ---------------------------------------------------------------------------

func (s *Store) ListProducts(ctx context.Context) ([]domain.Product, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id::text, name, description, COALESCE(sku, ''), active,
		       COALESCE(array_to_json(features)::text, '[]'), created_at, updated_at
		FROM products ORDER BY name ASC`)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()
	items := make([]domain.Product, 0)
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		items = append(items, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}
	return items, nil
}

func (s *Store) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	p, err := scanProduct(s.db.QueryRowContext(ctx, `
		SELECT id::text, name, description, COALESCE(sku, ''), active,
		       COALESCE(array_to_json(features)::text, '[]'), created_at, updated_at
		FROM products WHERE id = $1::uuid`, id))
	if err != nil {
		return nil, mapNotFound(err)
	}
	return &p, nil
}

func scanProduct(row rowScanner) (domain.Product, error) {
	var (
		product      domain.Product
		featuresJSON string
	)
	if err := row.Scan(
		&product.ID, &product.Name, &product.Description, &product.SKU,
		&product.Active, &featuresJSON, &product.CreatedAt, &product.UpdatedAt,
	); err != nil {
		return domain.Product{}, err
	}
	if featuresJSON != "" && featuresJSON != "[]" {
		_ = json.Unmarshal([]byte(featuresJSON), &product.Features)
	}
	if product.Features == nil {
		product.Features = []string{}
	}
	return product, nil
}

func (s *Store) CreateProduct(ctx context.Context, p *domain.Product) error {
	features, _ := json.Marshal(p.Features)
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO products (name, description, sku, active, features)
		VALUES ($1, $2, NULLIF($3, ''), $4, $5::jsonb)
		RETURNING id::text, created_at, updated_at`,
		p.Name, p.Description, p.SKU, p.Active, string(features),
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("create product: %w", err)
	}
	return nil
}

func (s *Store) UpdateProduct(ctx context.Context, p *domain.Product) error {
	features, _ := json.Marshal(p.Features)
	result, err := s.db.ExecContext(ctx, `
		UPDATE products SET
			name = COALESCE(NULLIF($1, ''), name),
			description = COALESCE(NULLIF($2, ''), description),
			sku = NULLIF($3, ''),
			active = COALESCE($4, active),
			features = COALESCE($5::jsonb, features),
			updated_at = now()
		WHERE id = $6::uuid`, p.Name, p.Description, p.SKU, p.Active, string(features), p.ID)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) DeleteProduct(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1::uuid", id)
	if err != nil {
		if isForeignKeyViolation(err) {
			return domain.Conflictf("product has attached plans and cannot be deleted")
		}
		return fmt.Errorf("delete product: %w", err)
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ---------------------------------------------------------------------------
// Discounts
// ---------------------------------------------------------------------------

const discountColumns = `
	d.id::text, d.code, d.description, d.kind::text,
	COALESCE(d.percent_basis_points, 0), COALESCE(d.amount_cents, 0),
	trim(COALESCE(d.currency, '')), d.duration::text,
	COALESCE(d.duration_months, 0), d.max_redemptions, d.times_redeemed,
	d.starts_at, d.ends_at, d.active`

func scanDiscount(row rowScanner) (domain.Discount, error) {
	var (
		d              domain.Discount
		maxRedemptions sql.NullInt64
		startsAt       sql.NullTime
		endsAt         sql.NullTime
	)
	if err := row.Scan(
		&d.ID, &d.Code, &d.Description, &d.Kind,
		&d.PercentBasisPoints, &d.AmountCents, &d.Currency, &d.Duration,
		&d.DurationMonths, &maxRedemptions, &d.TimesRedeemed,
		&startsAt, &endsAt, &d.Active,
	); err != nil {
		return domain.Discount{}, err
	}
	if maxRedemptions.Valid {
		v := int(maxRedemptions.Int64)
		d.MaxRedemptions = &v
	}
	d.StartsAt = timePtr(startsAt)
	d.EndsAt = timePtr(endsAt)
	return d, nil
}

func (s *Store) ListDiscounts(ctx context.Context, activeOnly bool) ([]domain.Discount, error) {
	query := "SELECT " + discountColumns + " FROM discounts d"
	if activeOnly {
		query += " WHERE d.active = TRUE"
	}
	query += " ORDER BY d.code ASC"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list discounts: %w", err)
	}
	defer rows.Close()
	items := make([]domain.Discount, 0)
	for rows.Next() {
		d, err := scanDiscount(rows)
		if err != nil {
			return nil, fmt.Errorf("scan discount: %w", err)
		}
		items = append(items, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate discounts: %w", err)
	}
	return items, nil
}

func (s *Store) GetDiscount(ctx context.Context, id string) (*domain.Discount, error) {
	d, err := scanDiscount(s.db.QueryRowContext(ctx, "SELECT "+discountColumns+" FROM discounts d WHERE d.id = $1::uuid", id))
	if err != nil {
		return nil, mapNotFound(err)
	}
	return &d, nil
}

func (s *Store) CreateDiscount(ctx context.Context, d *domain.Discount) error {
	var maxRedemptions any
	if d.MaxRedemptions != nil {
		maxRedemptions = *d.MaxRedemptions
	}
	var currency any
	if d.Currency != "" {
		currency = d.Currency
	}
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO discounts (
			code, description, kind, percent_basis_points, amount_cents,
			currency, duration, duration_months, max_redemptions, starts_at, ends_at, active
		) VALUES ($1, $2, $3::discount_kind, $4, $5, $6, $7::discount_duration, $8, $9, $10, $11, $12)
		RETURNING id::text`,
		strings.ToLower(strings.TrimSpace(d.Code)),
		d.Description,
		string(d.Kind),
		nullableInt(d.PercentBasisPoints),
		nullableInt64(d.AmountCents),
		currency,
		string(d.Duration),
		nullableInt(d.DurationMonths),
		maxRedemptions,
		nullTime(d.StartsAt),
		nullTime(d.EndsAt),
		d.Active,
	).Scan(&d.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("create discount: %w", err)
	}
	return nil
}

func (s *Store) UpdateDiscount(ctx context.Context, d *domain.Discount) error {
	result, err := s.db.ExecContext(ctx, `
		UPDATE discounts SET active = $2, updated_at = now() WHERE id = $1::uuid`,
		d.ID, d.Active)
	if err != nil {
		return fmt.Errorf("update discount: %w", err)
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ---------------------------------------------------------------------------
// AR & customer portal
// ---------------------------------------------------------------------------

func (s *Store) ListOverdueInvoices(ctx context.Context, overdueOnly bool, limit int) ([]domain.OverdueInvoice, error) {
	if limit < 1 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	query := `
		SELECT i.id::text, i.invoice_number, i.customer_id::text, c.name, c.email,
		       COALESCE(i.subscription_id::text, ''), COALESCE(pl.name, ''),
		       i.due_date, GREATEST((CURRENT_DATE - i.due_date), 0)::int,
		       COALESCE(i.amount_due_cents - i.amount_paid_cents, 0),
		       i.status::text, 0
		FROM invoices i
		JOIN customers c ON c.id = i.customer_id
		LEFT JOIN subscriptions sub ON sub.id = i.subscription_id
		LEFT JOIN plans pl ON pl.id = sub.plan_id
		WHERE i.status IN ('open', 'draft')`
	if overdueOnly {
		query += " AND i.due_date < CURRENT_DATE"
	}
	query += " ORDER BY i.due_date ASC LIMIT $1"
	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("list overdue invoices: %w", err)
	}
	defer rows.Close()
	items := make([]domain.OverdueInvoice, 0)
	for rows.Next() {
		var inv domain.OverdueInvoice
		if err := rows.Scan(
			&inv.ID, &inv.InvoiceNumber, &inv.CustomerID, &inv.CustomerName, &inv.CustomerEmail,
			&inv.SubscriptionID, &inv.PlanName, &inv.DueDate, &inv.DaysOverdue,
			&inv.AmountCents, &inv.Status, &inv.RemindersSent,
		); err != nil {
			return nil, fmt.Errorf("scan overdue invoice: %w", err)
		}
		items = append(items, inv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate overdue invoices: %w", err)
	}
	return items, nil
}

func (s *Store) CustomerInvoices(ctx context.Context, customerID string) ([]domain.CustomerInvoice, error) {
	rows, err := s.db.QueryContext(ctx, `
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
		WHERE i.customer_id = $1::uuid
		ORDER BY i.issue_date DESC, i.created_at DESC`,
		customerID)
	if err != nil {
		return nil, fmt.Errorf("list customer invoices: %w", err)
	}
	defer rows.Close()
	items := make([]domain.CustomerInvoice, 0)
	for rows.Next() {
		var inv domain.CustomerInvoice
		var paymentMethod sql.NullString
		var paidAt sql.NullTime
		if err := rows.Scan(
			&inv.ID, &inv.InvoiceNumber, &inv.CustomerID, &inv.SubscriptionID, &inv.Status,
			&inv.Date, &inv.DueDate, &inv.SubtotalCents, &inv.TaxCents, &inv.DiscountCents,
			&inv.TotalCents, &inv.AmountDueCents, &inv.AmountPaidCents, &inv.Currency,
			&paymentMethod, &paidAt,
		); err != nil {
			return nil, fmt.Errorf("scan customer invoice: %w", err)
		}
		inv.PaymentMethod = paymentMethod.String
		inv.PaidAt = timePtr(paidAt)
		items = append(items, inv)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate customer invoices: %w", err)
	}
	return items, nil
}

func (s *Store) CustomerOverview(ctx context.Context, customerID string) (domain.CustomerOverview, error) {
	var overview domain.CustomerOverview
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)::int,
		       COALESCE(SUM(
		           GREATEST(s.unit_price_cents * s.quantity - s.discount_cents, 0)::numeric /
		           CASE p.billing_interval WHEN 'quarterly' THEN 3 WHEN 'yearly' THEN 12 ELSE 1 END
		       ), 0)::bigint
		FROM subscriptions s JOIN plans p ON p.id = s.plan_id
		WHERE s.customer_id = $1::uuid AND s.status IN ('active', 'trialing', 'past_due')`,
		customerID).Scan(&overview.ActiveSubscriptionsCount, &overview.TotalMonthlyValueCents); err != nil {
		return domain.CustomerOverview{}, fmt.Errorf("customer overview subscriptions: %w", err)
	}
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*)::int, COALESCE(SUM(amount_due_cents - amount_paid_cents), 0)::bigint
		FROM invoices WHERE customer_id = $1::uuid AND status IN ('open', 'draft')`,
		customerID).Scan(&overview.OutstandingInvoicesCount, &overview.TotalOutstandingCents); err != nil {
		return domain.CustomerOverview{}, fmt.Errorf("customer overview invoices: %w", err)
	}
	var (
		amountCents int64
		dueDate     time.Time
		subID       string
		planName    string
	)
	err := s.db.QueryRowContext(ctx, `
		SELECT s.id::text, p.name, s.current_period_end,
		       GREATEST(s.unit_price_cents * s.quantity - s.discount_cents, 0)
		FROM subscriptions s JOIN plans p ON p.id = s.plan_id
		WHERE s.customer_id = $1::uuid
		  AND s.status IN ('active', 'trialing', 'past_due')
		  AND s.current_period_end >= CURRENT_DATE
		ORDER BY s.current_period_end ASC
		LIMIT 1`, customerID).Scan(&subID, &planName, &dueDate, &amountCents)
	if err == nil {
		overview.NextPaymentDue = &domain.NextPayment{
			AmountCents:        amountCents,
			Date:               dueDate,
			SubscriptionNumber: "SUB-" + subID[:8],
			PlanName:           planName,
		}
	}
	return overview, nil
}

// small helpers -------------------------------------------------------------

func nullableInt(v int) any {
	if v == 0 {
		return nil
	}
	return v
}

func nullableInt64(v int64) any {
	if v == 0 {
		return nil
	}
	return v
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
