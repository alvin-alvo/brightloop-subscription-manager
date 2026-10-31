package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

const subscriptionColumns = `
	s.id::text,
	s.customer_id::text,
	c.name,
	c.email,
	s.plan_id::text,
	p.name,
	COALESCE(p.description, ''),
	p.product_id::text,
	COALESCE(pr.name, ''),
	p.billing_interval::text,
	p.price_cents,
	trim(p.currency),
	COALESCE(array_to_json(p.features)::text, '[]'),
	p.trial_period_days,
	p.active,
	p.created_at,
	p.updated_at,
	s.status::text,
	s.quantity,
	s.unit_price_cents,
	s.discount_cents,
	trim(s.currency),
	s.current_period_start,
	s.current_period_end,
	s.cancel_at_period_end,
	s.trial_ends_at,
	s.started_at,
	s.canceled_at,
	s.ended_at,
	s.created_by::text,
	s.created_at,
	s.updated_at`

const subscriptionFrom = `
	FROM subscriptions s
	JOIN customers c ON c.id = s.customer_id
	JOIN plans p ON p.id = s.plan_id
	JOIN products pr ON pr.id = p.product_id`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanSubscription(row rowScanner) (domain.Subscription, error) {
	var (
		sub          domain.Subscription
		plan         domain.Plan
		featuresJSON string
		trialEndsAt  sql.NullTime
		canceledAt   sql.NullTime
		endedAt      sql.NullTime
		createdBy    sql.NullString
		currency     string
		planCurrency string
	)

	err := row.Scan(
		&sub.ID,
		&sub.CustomerID,
		&sub.CustomerName,
		&sub.CustomerEmail,
		&sub.PlanID,
		&plan.Name,
		&plan.Description,
		&plan.ProductID,
		&plan.ProductName,
		&plan.BillingInterval,
		&plan.PriceCents,
		&planCurrency,
		&featuresJSON,
		&plan.TrialPeriodDays,
		&plan.Active,
		&plan.CreatedAt,
		&plan.UpdatedAt,
		&sub.Status,
		&sub.Quantity,
		&sub.UnitPriceCents,
		&sub.DiscountCents,
		&currency,
		&sub.CurrentPeriodStart,
		&sub.CurrentPeriodEnd,
		&sub.CancelAtPeriodEnd,
		&trialEndsAt,
		&sub.StartedAt,
		&canceledAt,
		&endedAt,
		&createdBy,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		return domain.Subscription{}, err
	}

	plan.ID = sub.PlanID
	plan.Currency = strings.TrimSpace(planCurrency)
	if err := json.Unmarshal([]byte(featuresJSON), &plan.Features); err != nil || plan.Features == nil {
		plan.Features = []string{}
	}
	sub.Plan = &plan
	sub.Currency = strings.TrimSpace(currency)
	sub.TrialEndsAt = timePtr(trialEndsAt)
	sub.CanceledAt = timePtr(canceledAt)
	sub.EndedAt = timePtr(endedAt)
	sub.CreatedBy = stringPtr(createdBy)
	return sub, nil
}

func (s *Store) List(ctx context.Context, filter domain.ListFilter) ([]domain.Subscription, int, error) {
	filter = filter.Normalize()
	where, args := buildSubscriptionWhere(filter)

	var total int
	countQuery := "SELECT COUNT(*) " + subscriptionFrom + where
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count subscriptions: %w", err)
	}

	sortColumn := "s.created_at"
	switch filter.SortBy {
	case "updated_at":
		sortColumn = "s.updated_at"
	case "current_period_end":
		sortColumn = "s.current_period_end"
	case "unit_price_cents":
		sortColumn = "s.unit_price_cents"
	case "customer_name":
		sortColumn = "c.name"
	case "status":
		sortColumn = "s.status"
	}
	direction := "DESC"
	if filter.SortOrder == "asc" {
		direction = "ASC"
	}

	offset := (filter.Page - 1) * filter.PageSize
	query := "SELECT " + subscriptionColumns + subscriptionFrom + where +
		fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", sortColumn, direction, len(args)+1, len(args)+2)
	queryArgs := append(append([]any{}, args...), filter.PageSize, offset)

	rows, err := s.db.QueryContext(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list subscriptions: %w", err)
	}
	defer rows.Close()

	items := make([]domain.Subscription, 0, filter.PageSize)
	for rows.Next() {
		sub, err := scanSubscription(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan subscription: %w", err)
		}
		items = append(items, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate subscriptions: %w", err)
	}
	return items, total, nil
}

func buildSubscriptionWhere(filter domain.ListFilter) (string, []any) {
	var conditions []string
	var args []any

	if filter.Status != "" {
		args = append(args, string(filter.Status))
		conditions = append(conditions, fmt.Sprintf("s.status = $%d::subscription_status", len(args)))
	}
	if filter.CustomerID != "" {
		args = append(args, filter.CustomerID)
		conditions = append(conditions, fmt.Sprintf("s.customer_id = $%d::uuid", len(args)))
	}
	if filter.PlanID != "" {
		args = append(args, filter.PlanID)
		conditions = append(conditions, fmt.Sprintf("s.plan_id = $%d::uuid", len(args)))
	}
	if filter.Search != "" {
		args = append(args, likePattern(filter.Search))
		idx := len(args)
		conditions = append(conditions, fmt.Sprintf(
			"(lower(c.name) LIKE $%d ESCAPE '\\' OR lower(c.email) LIKE $%d ESCAPE '\\' OR lower(p.name) LIKE $%d ESCAPE '\\' OR s.id::text LIKE $%d)",
			idx, idx, idx, idx,
		))
	}

	if len(conditions) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func (s *Store) Summary(ctx context.Context, filter domain.ListFilter) (domain.SubscriptionSummary, error) {
	filter = filter.Normalize()
	where, args := buildSubscriptionWhere(filter)

	query := `
		SELECT
			COUNT(*)::int,
			COUNT(*) FILTER (WHERE s.status = 'active')::int,
			COUNT(*) FILTER (WHERE s.status = 'trialing')::int,
			COUNT(*) FILTER (WHERE s.status = 'past_due')::int,
			COUNT(*) FILTER (WHERE s.status = 'canceled')::int,
			COUNT(*) FILTER (WHERE s.status = 'paused')::int,
			COUNT(*) FILTER (WHERE s.status = 'unpaid')::int,
			COALESCE(SUM(
				CASE WHEN s.status IN ('active', 'trialing', 'past_due') THEN
					ROUND(
						GREATEST(s.unit_price_cents * s.quantity - s.discount_cents, 0)::numeric /
						CASE p.billing_interval
							WHEN 'quarterly' THEN 3
							WHEN 'yearly' THEN 12
							ELSE 1
						END
					)
				ELSE 0 END
			), 0)::bigint,
			COALESCE(MAX(trim(s.currency)), 'USD')
		` + subscriptionFrom + where

	var summary domain.SubscriptionSummary
	err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&summary.Total,
		&summary.Active,
		&summary.Trialing,
		&summary.PastDue,
		&summary.Canceled,
		&summary.Paused,
		&summary.Unpaid,
		&summary.MRRCents,
		&summary.Currency,
	)
	if err != nil {
		return domain.SubscriptionSummary{}, fmt.Errorf("summarise subscriptions: %w", err)
	}
	return summary, nil
}

func (s *Store) GetByID(ctx context.Context, id string) (*domain.Subscription, error) {
	query := "SELECT " + subscriptionColumns + subscriptionFrom + " WHERE s.id = $1::uuid"
	sub, err := scanSubscription(s.db.QueryRowContext(ctx, query, id))
	if err != nil {
		return nil, mapNotFound(err)
	}
	return &sub, nil
}

func (s *Store) Create(ctx context.Context, subscription *domain.Subscription, discountID *string) (*domain.Subscription, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin create subscription: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO subscriptions (
			customer_id, plan_id, status, quantity, unit_price_cents,
			discount_cents, currency, current_period_start, current_period_end,
			cancel_at_period_end, trial_ends_at, started_at, created_by
		) VALUES (
			$1::uuid, $2::uuid, $3::subscription_status, $4, $5,
			$6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id::text`

	var createdBy any
	if subscription.CreatedBy != nil && *subscription.CreatedBy != "" {
		createdBy = *subscription.CreatedBy
	}
	var trialEndsAt any
	if subscription.TrialEndsAt != nil {
		trialEndsAt = *subscription.TrialEndsAt
	}

	var id string
	err = tx.QueryRowContext(ctx, query,
		subscription.CustomerID,
		subscription.PlanID,
		string(subscription.Status),
		subscription.Quantity,
		subscription.UnitPriceCents,
		subscription.DiscountCents,
		subscription.Currency,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.CancelAtPeriodEnd,
		trialEndsAt,
		subscription.StartedAt,
		createdBy,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("insert subscription: %w", err)
	}

	if discountID != nil && *discountID != "" {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO subscription_discounts (subscription_id, discount_id) VALUES ($1::uuid, $2::uuid)
			 ON CONFLICT (subscription_id, discount_id) DO NOTHING`,
			id, *discountID,
		); err != nil {
			return nil, fmt.Errorf("link subscription discount: %w", err)
		}
		if _, err := tx.ExecContext(ctx,
			`UPDATE discounts SET times_redeemed = times_redeemed + 1, updated_at = now() WHERE id = $1::uuid`,
			*discountID,
		); err != nil {
			return nil, fmt.Errorf("redeem discount: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit create subscription: %w", err)
	}
	return s.GetByID(ctx, id)
}

func (s *Store) Update(ctx context.Context, subscription *domain.Subscription) (*domain.Subscription, error) {
	query := `
		UPDATE subscriptions SET
			plan_id = $1::uuid,
			status = $2::subscription_status,
			quantity = $3,
			unit_price_cents = $4,
			discount_cents = $5,
			currency = $6,
			current_period_start = $7,
			current_period_end = $8,
			cancel_at_period_end = $9,
			trial_ends_at = $10,
			canceled_at = $11,
			ended_at = $12,
			updated_at = now()
		WHERE id = $13::uuid`

	result, err := s.db.ExecContext(ctx, query,
		subscription.PlanID,
		string(subscription.Status),
		subscription.Quantity,
		subscription.UnitPriceCents,
		subscription.DiscountCents,
		subscription.Currency,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.CancelAtPeriodEnd,
		nullTime(subscription.TrialEndsAt),
		nullTime(subscription.CanceledAt),
		nullTime(subscription.EndedAt),
		subscription.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("update subscription: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("update subscription rows affected: %w", err)
	}
	if affected == 0 {
		return nil, domain.ErrNotFound
	}
	return s.GetByID(ctx, subscription.ID)
}

func (s *Store) ListDueForBilling(ctx context.Context, asOf time.Time) ([]domain.Subscription, error) {
	query := "SELECT " + subscriptionColumns + subscriptionFrom + `
		WHERE s.current_period_end <= $1::date
		  AND s.status IN ('active', 'trialing', 'past_due')
		ORDER BY s.current_period_end ASC, s.created_at ASC`
	rows, err := s.db.QueryContext(ctx, query, asOf)
	if err != nil {
		return nil, fmt.Errorf("list due subscriptions: %w", err)
	}
	defer rows.Close()
	return s.collectSubscriptions(rows)
}

func (s *Store) ListUpcomingRenewals(ctx context.Context, asOf time.Time, limit int) ([]domain.Subscription, error) {
	if limit < 1 {
		limit = 5
	}
	query := "SELECT " + subscriptionColumns + subscriptionFrom + `
		WHERE s.status IN ('active', 'trialing', 'past_due')
		  AND s.current_period_end >= $1::date
		ORDER BY s.current_period_end ASC
		LIMIT $2`
	rows, err := s.db.QueryContext(ctx, query, asOf, limit)
	if err != nil {
		return nil, fmt.Errorf("list upcoming renewals: %w", err)
	}
	defer rows.Close()
	return s.collectSubscriptions(rows)
}

func (s *Store) collectSubscriptions(rows *sql.Rows) ([]domain.Subscription, error) {
	items := make([]domain.Subscription, 0)
	for rows.Next() {
		sub, err := scanSubscription(rows)
		if err != nil {
			return nil, fmt.Errorf("scan subscription: %w", err)
		}
		items = append(items, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate subscriptions: %w", err)
	}
	return items, nil
}
