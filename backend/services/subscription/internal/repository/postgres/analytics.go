package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
)

func (s *Store) Analytics(ctx context.Context) (domain.Analytics, error) {
	query := `
		WITH month AS (
			SELECT date_trunc('month', now()) AS start_day
		),
		subs AS (
			SELECT s.*,
			       GREATEST(s.unit_price_cents * s.quantity - s.discount_cents, 0)::numeric AS gross_cents,
			       p.billing_interval
			FROM subscriptions s
			JOIN plans p ON p.id = s.plan_id
		),
		norm AS (
			SELECT *,
			       ROUND(gross_cents / CASE billing_interval
			           WHEN 'quarterly' THEN 3
			           WHEN 'yearly' THEN 12
			           ELSE 1
			       END)::bigint AS mrr_cents
			FROM subs
		)
		SELECT
			(SELECT COUNT(*) FROM subscriptions WHERE status = 'active')::int,
			(SELECT COUNT(*) FROM subscriptions WHERE status = 'trialing')::int,
			(SELECT COUNT(*) FROM subscriptions WHERE status = 'past_due')::int,
			(SELECT COUNT(*) FROM subscriptions WHERE status = 'canceled')::int,
			(SELECT COUNT(*) FROM subscriptions WHERE status = 'paused')::int,
			(SELECT COUNT(*) FROM subscriptions WHERE status = 'unpaid')::int,
			(SELECT COUNT(*) FROM subscriptions)::int,
			(SELECT COUNT(*) FROM subscriptions WHERE started_at >= (SELECT start_day FROM month))::int,
			(SELECT COUNT(*) FROM subscriptions WHERE canceled_at IS NOT NULL AND canceled_at >= (SELECT start_day FROM month))::int,
			(SELECT COALESCE(SUM(mrr_cents), 0) FROM norm WHERE status IN ('active', 'trialing', 'past_due'))::bigint,
			(SELECT COALESCE(SUM(mrr_cents), 0) FROM norm
			  WHERE started_at < (SELECT start_day FROM month)
			    AND (canceled_at IS NULL OR canceled_at >= (SELECT start_day FROM month)))::bigint,
			(SELECT COALESCE(SUM(mrr_cents), 0) FROM norm
			  WHERE status IN ('active', 'trialing', 'past_due')
			    AND started_at < (SELECT start_day FROM month)
			    AND (canceled_at IS NULL OR canceled_at >= (SELECT start_day FROM month)))::bigint,
			(SELECT COUNT(*) FROM customers)::int,
			(SELECT COUNT(*) FROM customers WHERE status = 'active')::int,
			(SELECT COUNT(*) FROM customers WHERE created_at >= (SELECT start_day FROM month))::int,
			(SELECT COALESCE(SUM(amount_cents), 0) FROM payments
			  WHERE status = 'succeeded'
			    AND COALESCE(paid_at, created_at) >= (SELECT start_day FROM month))::bigint,
			(SELECT COALESCE(SUM(amount_cents), 0) FROM payments
			  WHERE status = 'succeeded'
			    AND COALESCE(paid_at, created_at) >= (SELECT start_day FROM month) - interval '1 month'
			    AND COALESCE(paid_at, created_at) < (SELECT start_day FROM month))::bigint,
			(SELECT COALESCE(SUM(amount_due_cents - amount_paid_cents), 0) FROM invoices WHERE status = 'open')::bigint,
			(SELECT COUNT(*) FROM invoices WHERE status = 'open' AND due_date < CURRENT_DATE)::int
	`

	var a domain.Analytics
	err := s.db.QueryRowContext(ctx, query).Scan(
		&a.ActiveSubscriptions,
		&a.TrialingSubscriptions,
		&a.PastDueSubscriptions,
		&a.CanceledSubscriptions,
		&a.PausedSubscriptions,
		&a.UnpaidSubscriptions,
		&a.TotalSubscriptions,
		&a.NewSubscriptions,
		&a.CanceledThisMonth,
		&a.MRRCents,
		&a.PreviousMRRCents,
		&a.NetRetainedMRRCents,
		&a.TotalCustomers,
		&a.ActiveCustomers,
		&a.NewCustomers,
		&a.RevenueThisMonthCents,
		&a.RevenueLastMonthCents,
		&a.OutstandingCents,
		&a.OverdueInvoices,
	)
	if err != nil {
		return domain.Analytics{}, fmt.Errorf("read analytics: %w", err)
	}
	return a, nil
}

func (s *Store) MonthlyRevenue(ctx context.Context, months int) ([]domain.MonthlyRevenue, error) {
	if months < 1 {
		months = 12
	}
	if months > 36 {
		months = 36
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT to_char(date_trunc('month', COALESCE(paid_at, created_at)), 'YYYY-MM') AS month,
		       COALESCE(SUM(amount_cents), 0)::bigint,
		       COUNT(*)::int
		FROM payments
		WHERE status = 'succeeded'
		  AND COALESCE(paid_at, created_at) >= date_trunc('month', now()) - ($1::int - 1) * interval '1 month'
		GROUP BY 1
		ORDER BY 1`, months)
	if err != nil {
		return nil, fmt.Errorf("read monthly revenue: %w", err)
	}
	defer rows.Close()

	points := make([]domain.MonthlyRevenue, 0)
	for rows.Next() {
		var point domain.MonthlyRevenue
		if err := rows.Scan(&point.Month, &point.RevenueCents, &point.PaymentCount); err != nil {
			return nil, fmt.Errorf("scan monthly revenue: %w", err)
		}
		points = append(points, point)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate monthly revenue: %w", err)
	}
	return points, nil
}

func (s *Store) PlanBreakdown(ctx context.Context) ([]domain.PlanBreakdown, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT p.id::text,
		       p.name,
		       COALESCE(pr.name, ''),
		       p.billing_interval::text,
		       COUNT(s.id) FILTER (WHERE s.status IN ('active', 'trialing', 'past_due'))::int,
		       COALESCE(SUM(
		           ROUND(
		               GREATEST(s.unit_price_cents * s.quantity - s.discount_cents, 0)::numeric /
		               CASE p.billing_interval
		                   WHEN 'quarterly' THEN 3
		                   WHEN 'yearly' THEN 12
		                   ELSE 1
		               END
		           )
		       ) FILTER (WHERE s.status IN ('active', 'trialing', 'past_due')), 0)::bigint,
		       COALESCE(rev.revenue_cents, 0)::bigint
		FROM plans p
		JOIN products pr ON pr.id = p.product_id
		LEFT JOIN subscriptions s ON s.plan_id = p.id
		LEFT JOIN (
			SELECT ili.plan_id, SUM(ili.amount_cents) AS revenue_cents
			FROM invoice_line_items ili
			JOIN invoices i ON i.id = ili.invoice_id
			WHERE i.status = 'paid'
			GROUP BY ili.plan_id
		) rev ON rev.plan_id = p.id
		GROUP BY p.id, p.name, pr.name, p.billing_interval, rev.revenue_cents
		ORDER BY 6 DESC, p.name ASC`)
	if err != nil {
		return nil, fmt.Errorf("read plan breakdown: %w", err)
	}
	defer rows.Close()

	plans := make([]domain.PlanBreakdown, 0)
	for rows.Next() {
		var plan domain.PlanBreakdown
		if err := rows.Scan(
			&plan.PlanID,
			&plan.PlanName,
			&plan.ProductName,
			&plan.BillingInterval,
			&plan.ActiveSubscriptions,
			&plan.MRRCents,
			&plan.RevenueCents,
		); err != nil {
			return nil, fmt.Errorf("scan plan breakdown: %w", err)
		}
		plans = append(plans, plan)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate plan breakdown: %w", err)
	}
	return plans, nil
}

func (s *Store) TopCustomers(ctx context.Context, limit int) ([]domain.TopCustomer, error) {
	if limit < 1 {
		limit = 5
	}
	if limit > 50 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id::text,
		       c.name,
		       c.email,
		       COUNT(s.id) FILTER (WHERE s.status IN ('active', 'trialing', 'past_due'))::int,
		       COALESCE(SUM(
		           ROUND(
		               GREATEST(s.unit_price_cents * s.quantity - s.discount_cents, 0)::numeric /
		               CASE COALESCE(p.billing_interval, 'monthly')
		                   WHEN 'quarterly' THEN 3
		                   WHEN 'yearly' THEN 12
		                   ELSE 1
		               END
		           )
		       ) FILTER (WHERE s.status IN ('active', 'trialing', 'past_due')), 0)::bigint,
		       COALESCE((
		           SELECT SUM(pay.amount_cents) FROM payments pay
		           WHERE pay.customer_id = c.id AND pay.status = 'succeeded'
		       ), 0)::bigint
		FROM customers c
		LEFT JOIN subscriptions s ON s.customer_id = c.id
		LEFT JOIN plans p ON p.id = s.plan_id
		GROUP BY c.id, c.name, c.email
		ORDER BY 5 DESC, 6 DESC
		LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("read top customers: %w", err)
	}
	defer rows.Close()

	customers := make([]domain.TopCustomer, 0)
	for rows.Next() {
		var customer domain.TopCustomer
		if err := rows.Scan(
			&customer.CustomerID,
			&customer.CustomerName,
			&customer.CustomerEmail,
			&customer.ActiveSubscriptions,
			&customer.MRRCents,
			&customer.RevenueCents,
		); err != nil {
			return nil, fmt.Errorf("scan top customer: %w", err)
		}
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate top customers: %w", err)
	}
	return customers, nil
}

func (s *Store) InvoiceAging(ctx context.Context) ([]domain.AgingBucket, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT CASE
		           WHEN due_date >= CURRENT_DATE THEN 'current'
		           WHEN CURRENT_DATE - due_date <= 30 THEN '1-30'
		           WHEN CURRENT_DATE - due_date <= 60 THEN '31-60'
		           WHEN CURRENT_DATE - due_date <= 90 THEN '61-90'
		           ELSE '90+'
		       END AS bucket,
		       COUNT(*)::int,
		       COALESCE(SUM(amount_due_cents - amount_paid_cents), 0)::bigint
		FROM invoices
		WHERE status = 'open'
		GROUP BY 1`)
	if err != nil {
		return nil, fmt.Errorf("read invoice aging: %w", err)
	}
	defer rows.Close()

	buckets := make([]domain.AgingBucket, 0)
	for rows.Next() {
		var bucket domain.AgingBucket
		if err := rows.Scan(&bucket.Bucket, &bucket.Count, &bucket.AmountCents); err != nil {
			return nil, fmt.Errorf("scan invoice aging: %w", err)
		}
		bucket.Label = agingLabel(bucket.Bucket)
		buckets = append(buckets, bucket)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate invoice aging: %w", err)
	}
	return buckets, nil
}

func agingLabel(bucket string) string {
	switch strings.ToLower(bucket) {
	case "current":
		return "Current"
	case "1-30":
		return "1–30 days"
	case "31-60":
		return "31–60 days"
	case "61-90":
		return "61–90 days"
	case "90+":
		return "90+ days"
	default:
		return bucket
	}
}
