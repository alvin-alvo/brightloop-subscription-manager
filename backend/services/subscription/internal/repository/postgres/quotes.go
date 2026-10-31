package postgres

import (
	"context"
	"crypto/rand"
	"fmt"
	"strconv"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

func randomCompact(n int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err == nil {
		for i := range buf {
			buf[i] = chars[int(buf[i])%len(chars)]
		}
		return string(buf)
	}
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func quoteNumberGen() string { return "QT-" + randomCompact(8) }

func (s *Store) ListQuotes(ctx context.Context) ([]domain.Quote, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT q.id::text, q.quote_number, q.customer_id::text, c.name,
		       q.status::text, trim(q.currency),
		       q.subtotal_cents, q.discount_cents, q.tax_cents, q.total_cents,
		       COALESCE(q.valid_until, CURRENT_DATE), COALESCE(q.notes, ''),
		       q.created_at, q.updated_at
		FROM quotations q JOIN customers c ON c.id = q.customer_id
		ORDER BY q.created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list quotes: %w", err)
	}
	defer rows.Close()
	items := make([]domain.Quote, 0)
	for rows.Next() {
		q, err := scanQuote(rows)
		if err != nil {
			return nil, fmt.Errorf("scan quote: %w", err)
		}
		items = append(items, q)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate quotes: %w", err)
	}
	for i := range items {
		items[i].Items = s.quoteItems(ctx, items[i].ID)
	}
	return items, nil
}

func (s *Store) GetQuote(ctx context.Context, id string) (*domain.Quote, error) {
	q, err := scanQuote(s.db.QueryRowContext(ctx, `
		SELECT q.id::text, q.quote_number, q.customer_id::text, c.name,
		       q.status::text, trim(q.currency),
		       q.subtotal_cents, q.discount_cents, q.tax_cents, q.total_cents,
		       COALESCE(q.valid_until, CURRENT_DATE), COALESCE(q.notes, ''),
		       q.created_at, q.updated_at
		FROM quotations q JOIN customers c ON c.id = q.customer_id
		WHERE q.id = $1::uuid`, id))
	if err != nil {
		return nil, mapNotFound(err)
	}
	q.Items = s.quoteItems(ctx, q.ID)
	return &q, nil
}

func scanQuote(row rowScanner) (domain.Quote, error) {
	var q domain.Quote
	if err := row.Scan(
		&q.ID, &q.QuoteNumber, &q.CustomerID, &q.CustomerName,
		&q.Status, &q.Currency, &q.SubtotalCents, &q.DiscountCents, &q.TaxCents, &q.TotalCents,
		&q.ValidUntil, &q.Notes, &q.CreatedAt, &q.UpdatedAt,
	); err != nil {
		return domain.Quote{}, err
	}
	q.Items = []domain.QuoteItem{}
	return q, nil
}

func (s *Store) quoteItems(ctx context.Context, quoteID string) []domain.QuoteItem {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id::text, COALESCE(plan_id::text, ''), description, quantity,
		       unit_price_cents, tax_basis_points, line_total_cents
		FROM quotation_items WHERE quotation_id = $1::uuid ORDER BY created_at ASC`, quoteID)
	if err != nil {
		return []domain.QuoteItem{}
	}
	defer rows.Close()
	items := make([]domain.QuoteItem, 0)
	for rows.Next() {
		var item domain.QuoteItem
		var planID string
		if err := rows.Scan(&item.ID, &planID, &item.Description, &item.Quantity,
			&item.UnitPriceCents, &item.TaxBasisPoints, &item.LineTotalCents); err != nil {
			continue
		}
		if planID != "" {
			item.PlanID = &planID
		}
		items = append(items, item)
	}
	return items
}

func (s *Store) CreateQuote(ctx context.Context, q *domain.Quote) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin create quote: %w", err)
	}
	defer tx.Rollback()

	if q.QuoteNumber == "" {
		q.QuoteNumber = quoteNumberGen()
	}
	err = tx.QueryRowContext(ctx, `
		INSERT INTO quotations (
			quote_number, customer_id, status, currency,
			subtotal_cents, discount_cents, tax_cents, total_cents, valid_until, notes
		) VALUES ($1, $2::uuid, $3::quotation_status, $4, $5, $6, $7, $8, $9, NULLIF($10, ''))
		RETURNING id::text, created_at, updated_at`,
		q.QuoteNumber, q.CustomerID, defaultQuoteStatus(q.Status), q.Currency,
		q.SubtotalCents, q.DiscountCents, q.TaxCents, q.TotalCents, q.ValidUntil, q.Notes,
	).Scan(&q.ID, &q.CreatedAt, &q.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("insert quote: %w", err)
	}

	for i := range q.Items {
		var planID any
		if q.Items[i].PlanID != nil {
			planID = *q.Items[i].PlanID
		}
		if err := tx.QueryRowContext(ctx, `
			INSERT INTO quotation_items (
				quotation_id, plan_id, description, quantity, unit_price_cents,
				discount_basis_points, tax_basis_points, line_total_cents
			) VALUES ($1::uuid, $2::uuid, $3, $4, $5, 0, $6, $7)
			RETURNING id::text`,
			q.ID, planID, q.Items[i].Description, q.Items[i].Quantity,
			q.Items[i].UnitPriceCents, q.Items[i].TaxBasisPoints, q.Items[i].LineTotalCents,
		).Scan(&q.Items[i].ID); err != nil {
			return fmt.Errorf("insert quote item: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit quote: %w", err)
	}
	return nil
}

func (s *Store) UpdateQuoteStatus(ctx context.Context, id, status string) error {
	result, err := s.db.ExecContext(ctx,
		"UPDATE quotations SET status = $2::quotation_status, updated_at = now() WHERE id = $1::uuid",
		id, status)
	if err != nil {
		return fmt.Errorf("update quote status: %w", err)
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) DeleteQuote(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM quotations WHERE id = $1::uuid", id)
	if err != nil {
		return fmt.Errorf("delete quote: %w", err)
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func defaultQuoteStatus(status string) string {
	if status == "" {
		return "draft"
	}
	return status
}

// --- Templates -------------------------------------------------------------

func (s *Store) ListTemplates(ctx context.Context) ([]domain.QuoteTemplate, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT t.id::text, t.name, COALESCE(t.plan_id::text, ''), COALESCE(p.name, ''),
		       t.validity_days, t.created_at, t.updated_at
		FROM quote_templates t LEFT JOIN plans p ON p.id = t.plan_id
		ORDER BY t.name ASC`)
	if err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}
	defer rows.Close()
	items := make([]domain.QuoteTemplate, 0)
	for rows.Next() {
		t, err := scanTemplate(rows)
		if err != nil {
			return nil, fmt.Errorf("scan template: %w", err)
		}
		t.Items = s.templateItems(ctx, t.ID)
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate templates: %w", err)
	}
	return items, nil
}

func (s *Store) GetTemplate(ctx context.Context, id string) (*domain.QuoteTemplate, error) {
	t, err := scanTemplate(s.db.QueryRowContext(ctx, `
		SELECT t.id::text, t.name, COALESCE(t.plan_id::text, ''), COALESCE(p.name, ''),
		       t.validity_days, t.created_at, t.updated_at
		FROM quote_templates t LEFT JOIN plans p ON p.id = t.plan_id
		WHERE t.id = $1::uuid`, id))
	if err != nil {
		return nil, mapNotFound(err)
	}
	t.Items = s.templateItems(ctx, t.ID)
	return &t, nil
}

func scanTemplate(row rowScanner) (domain.QuoteTemplate, error) {
	var (
		t      domain.QuoteTemplate
		planID string
	)
	if err := row.Scan(&t.ID, &t.Name, &planID, &t.PlanName, &t.ValidityDays, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return domain.QuoteTemplate{}, err
	}
	if planID != "" {
		t.PlanID = &planID
	}
	t.Items = []domain.QuoteTemplateItem{}
	return t, nil
}

func (s *Store) templateItems(ctx context.Context, templateID string) []domain.QuoteTemplateItem {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id::text, COALESCE(product_id::text, ''), product_name, quantity,
		       unit_price_cents, tax_rate_basis_points
		FROM quote_template_items WHERE template_id = $1::uuid ORDER BY created_at ASC`, templateID)
	if err != nil {
		return []domain.QuoteTemplateItem{}
	}
	defer rows.Close()
	items := make([]domain.QuoteTemplateItem, 0)
	for rows.Next() {
		var item domain.QuoteTemplateItem
		var productID string
		if err := rows.Scan(&item.ID, &productID, &item.ProductName, &item.Quantity,
			&item.UnitPriceCents, &item.TaxBasisPoints); err != nil {
			continue
		}
		if productID != "" {
			item.ProductID = &productID
		}
		items = append(items, item)
	}
	return items
}

func (s *Store) CreateTemplate(ctx context.Context, t *domain.QuoteTemplate) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin create template: %w", err)
	}
	defer tx.Rollback()

	if t.ValidityDays == 0 {
		t.ValidityDays = 30
	}
	var planID any
	if t.PlanID != nil {
		planID = *t.PlanID
	}
	err = tx.QueryRowContext(ctx, `
		INSERT INTO quote_templates (name, plan_id, validity_days)
		VALUES ($1, $2::uuid, $3) RETURNING id::text, created_at, updated_at`,
		t.Name, planID, t.ValidityDays).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("insert template: %w", err)
	}
	for i := range t.Items {
		var productID any
		if t.Items[i].ProductID != nil {
			productID = *t.Items[i].ProductID
		}
		if err := tx.QueryRowContext(ctx, `
			INSERT INTO quote_template_items (
				template_id, product_id, product_name, quantity, unit_price_cents, tax_rate_basis_points
			) VALUES ($1::uuid, $2::uuid, $3, $4, $5, $6) RETURNING id::text`,
			t.ID, productID, t.Items[i].ProductName, t.Items[i].Quantity,
			t.Items[i].UnitPriceCents, t.Items[i].TaxBasisPoints,
		).Scan(&t.Items[i].ID); err != nil {
			return fmt.Errorf("insert template item: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit template: %w", err)
	}
	return nil
}

// --- CRM workspace customer view -------------------------------------------

func (s *Store) ListCustomerWorkspace(ctx context.Context) ([]domain.CustomerWorkspace, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id::text, c.name, COALESCE(c.email, '') AS contact, c.email,
		       COUNT(s.id) FILTER (WHERE s.status IN ('active', 'trialing', 'past_due'))::int,
		       COALESCE(SUM(
		           GREATEST(s.unit_price_cents * s.quantity - s.discount_cents, 0)::numeric /
		           CASE p.billing_interval WHEN 'quarterly' THEN 3 WHEN 'yearly' THEN 12 ELSE 1 END
		       ) FILTER (WHERE s.status IN ('active', 'trialing', 'past_due')), 0)::bigint,
		       c.status::text
		FROM customers c
		LEFT JOIN subscriptions s ON s.customer_id = c.id
		LEFT JOIN plans p ON p.id = s.plan_id
		GROUP BY c.id, c.name, c.email, c.status
		ORDER BY c.name ASC`)
	if err != nil {
		return nil, fmt.Errorf("list customer workspace: %w", err)
	}
	defer rows.Close()
	items := make([]domain.CustomerWorkspace, 0)
	for rows.Next() {
		var cw domain.CustomerWorkspace
		if err := rows.Scan(&cw.ID, &cw.Name, &cw.Contact, &cw.Email,
			&cw.ActiveSubscriptions, &cw.MonthlyValueCents, &cw.Status); err != nil {
			return nil, fmt.Errorf("scan customer workspace: %w", err)
		}
		cw.Tier = customerTier(cw.MonthlyValueCents)
		items = append(items, cw)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate customer workspace: %w", err)
	}
	return items, nil
}

func customerTier(mrr int64) string {
	switch {
	case mrr >= 500000:
		return "Enterprise"
	case mrr >= 50000:
		return "Growth"
	case mrr > 0:
		return "Pro"
	default:
		return "Trial"
	}
}
