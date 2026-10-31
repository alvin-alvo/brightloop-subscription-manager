package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

func (s *Store) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return s.getCustomer(ctx, `WHERE c.id = $1::uuid`, id)
}

func (s *Store) GetCustomerByUserID(ctx context.Context, userID string) (*domain.Customer, error) {
	return s.getCustomer(ctx, `WHERE c.user_id = $1::uuid`, userID)
}

func (s *Store) getCustomer(ctx context.Context, where string, arg any) (*domain.Customer, error) {
	query := `
		SELECT c.id::text, c.user_id::text, c.name, c.email,
		       COALESCE(c.company_name, ''), trim(c.currency), c.status::text
		FROM customers c ` + where
	var customer domain.Customer
	var userID sql.NullString
	err := s.db.QueryRowContext(ctx, query, arg).Scan(
		&customer.ID,
		&userID,
		&customer.Name,
		&customer.Email,
		&customer.CompanyName,
		&customer.Currency,
		&customer.Status,
	)
	if err != nil {
		return nil, mapNotFound(err)
	}
	customer.UserID = stringPtr(userID)
	return &customer, nil
}

func (s *Store) ListCustomers(ctx context.Context, search string, limit int) ([]domain.Customer, error) {
	if limit < 1 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	query := `
		SELECT c.id::text, c.user_id::text, c.name, c.email,
		       COALESCE(c.company_name, ''), trim(c.currency), c.status::text
		FROM customers c`
	args := []any{}
	if strings.TrimSpace(search) != "" {
		args = append(args, likePattern(search))
		query += " WHERE lower(c.name) LIKE $1 ESCAPE '\\' OR lower(c.email) LIKE $1 ESCAPE '\\' OR lower(COALESCE(c.company_name, '')) LIKE $1 ESCAPE '\\'"
	}
	query += fmt.Sprintf(" ORDER BY c.name ASC LIMIT $%d", len(args)+1)
	args = append(args, limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}
	defer rows.Close()

	customers := make([]domain.Customer, 0)
	for rows.Next() {
		var customer domain.Customer
		var userID sql.NullString
		if err := rows.Scan(
			&customer.ID,
			&userID,
			&customer.Name,
			&customer.Email,
			&customer.CompanyName,
			&customer.Currency,
			&customer.Status,
		); err != nil {
			return nil, fmt.Errorf("scan customer: %w", err)
		}
		customer.UserID = stringPtr(userID)
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate customers: %w", err)
	}
	return customers, nil
}

const planColumns = `
	p.id::text,
	p.product_id::text,
	COALESCE(pr.name, ''),
	p.name,
	p.description,
	p.billing_interval::text,
	p.price_cents,
	trim(p.currency),
	COALESCE(array_to_json(p.features)::text, '[]'),
	p.trial_period_days,
	p.active,
	p.created_at,
	p.updated_at`

func scanPlan(row rowScanner) (domain.Plan, error) {
	var (
		plan         domain.Plan
		featuresJSON string
	)
	err := row.Scan(
		&plan.ID,
		&plan.ProductID,
		&plan.ProductName,
		&plan.Name,
		&plan.Description,
		&plan.BillingInterval,
		&plan.PriceCents,
		&plan.Currency,
		&featuresJSON,
		&plan.TrialPeriodDays,
		&plan.Active,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err != nil {
		return domain.Plan{}, err
	}
	plan.Currency = strings.TrimSpace(plan.Currency)
	if err := json.Unmarshal([]byte(featuresJSON), &plan.Features); err != nil || plan.Features == nil {
		plan.Features = []string{}
	}
	return plan, nil
}

func (s *Store) GetPlan(ctx context.Context, id string) (*domain.Plan, error) {
	query := "SELECT " + planColumns + `
		FROM plans p JOIN products pr ON pr.id = p.product_id
		WHERE p.id = $1::uuid`
	plan, err := scanPlan(s.db.QueryRowContext(ctx, query, id))
	if err != nil {
		return nil, mapNotFound(err)
	}
	return &plan, nil
}

func (s *Store) ListPlans(ctx context.Context, activeOnly bool) ([]domain.Plan, error) {
	query := "SELECT " + planColumns + `
		FROM plans p JOIN products pr ON pr.id = p.product_id`
	if activeOnly {
		query += " WHERE p.active = TRUE AND pr.active = TRUE"
	}
	query += " ORDER BY pr.name ASC, p.price_cents ASC"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list plans: %w", err)
	}
	defer rows.Close()

	plans := make([]domain.Plan, 0)
	for rows.Next() {
		plan, err := scanPlan(rows)
		if err != nil {
			return nil, fmt.Errorf("scan plan: %w", err)
		}
		plans = append(plans, plan)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate plans: %w", err)
	}
	return plans, nil
}

func (s *Store) GetDiscountByCode(ctx context.Context, code string) (*domain.Discount, error) {
	query := `
		SELECT id::text, code, description, kind::text,
		       COALESCE(percent_basis_points, 0), COALESCE(amount_cents, 0),
		       trim(COALESCE(currency, '')), duration::text,
		       COALESCE(duration_months, 0), max_redemptions, times_redeemed,
		       starts_at, ends_at, active
		FROM discounts
		WHERE lower(code) = lower($1)`
	var (
		discount       domain.Discount
		maxRedemptions sql.NullInt64
		startsAt       sql.NullTime
		endsAt         sql.NullTime
	)
	err := s.db.QueryRowContext(ctx, query, strings.TrimSpace(code)).Scan(
		&discount.ID,
		&discount.Code,
		&discount.Description,
		&discount.Kind,
		&discount.PercentBasisPoints,
		&discount.AmountCents,
		&discount.Currency,
		&discount.Duration,
		&discount.DurationMonths,
		&maxRedemptions,
		&discount.TimesRedeemed,
		&startsAt,
		&endsAt,
		&discount.Active,
	)
	if err != nil {
		return nil, mapNotFound(err)
	}
	if maxRedemptions.Valid {
		value := int(maxRedemptions.Int64)
		discount.MaxRedemptions = &value
	}
	discount.StartsAt = timePtr(startsAt)
	discount.EndsAt = timePtr(endsAt)
	return &discount, nil
}
