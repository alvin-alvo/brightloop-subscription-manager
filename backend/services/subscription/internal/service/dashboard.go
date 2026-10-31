package service

import (
	"context"
	"math"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/repository"
)

// DashboardService implements analytics and reporting read models.
type DashboardService struct {
	analytics repository.AnalyticsRepository
	subs      repository.SubscriptionRepository
	billing   repository.BillingRepository
	now       func() time.Time
}

// NewDashboardService wires the dashboard business rules.
func NewDashboardService(
	analytics repository.AnalyticsRepository,
	subs repository.SubscriptionRepository,
	billing repository.BillingRepository,
) *DashboardService {
	return &DashboardService{analytics: analytics, subs: subs, billing: billing, now: time.Now}
}

// Analytics returns the computed metrics used by the finance dashboard.
func (s *DashboardService) Analytics(ctx context.Context, actor domain.Actor) (domain.ComputedAnalytics, error) {
	if !actor.IsStaff() {
		return domain.ComputedAnalytics{}, domain.Forbiddenf("only staff can read analytics")
	}
	raw, err := s.analytics.Analytics(ctx)
	if err != nil {
		return domain.ComputedAnalytics{}, err
	}
	return computeAnalytics(raw), nil
}

// Dashboard returns every dataset the admin finance dashboard renders.
func (s *DashboardService) Dashboard(ctx context.Context, actor domain.Actor) (*domain.Dashboard, error) {
	computed, err := s.Analytics(ctx, actor)
	if err != nil {
		return nil, err
	}
	transactions, _, err := s.billing.ListTransactions(ctx, domain.TransactionFilter{Page: 1, PageSize: 5})
	if err != nil {
		return nil, err
	}
	renewals, err := s.subs.ListUpcomingRenewals(ctx, utcDate(s.now()), 5)
	if err != nil {
		return nil, err
	}
	revenue, err := s.analytics.MonthlyRevenue(ctx, 6)
	if err != nil {
		return nil, err
	}

	upcoming := make([]domain.UpcomingRenewal, 0, len(renewals))
	for _, sub := range renewals {
		upcoming = append(upcoming, domain.UpcomingRenewal{
			SubscriptionID: sub.ID,
			CustomerID:     sub.CustomerID,
			CustomerName:   sub.CustomerName,
			PlanName:       planName(sub),
			Date:           sub.CurrentPeriodEnd,
			AmountCents:    sub.TotalCents(),
			Currency:       sub.Currency,
		})
	}

	return &domain.Dashboard{
		GeneratedAt:        s.now(),
		Analytics:          computed,
		RecentTransactions: transactions,
		RevenueTrend:       fillMonthly(revenue, 6, s.now()),
		StatusBreakdown:    statusCounts(computed.Analytics),
		UpcomingRenewals:   upcoming,
	}, nil
}

// Report returns the revenue report payload.
func (s *DashboardService) Report(ctx context.Context, actor domain.Actor) (*domain.AnalyticsReport, error) {
	computed, err := s.Analytics(ctx, actor)
	if err != nil {
		return nil, err
	}
	revenue, err := s.analytics.MonthlyRevenue(ctx, 12)
	if err != nil {
		return nil, err
	}
	plans, err := s.analytics.PlanBreakdown(ctx)
	if err != nil {
		return nil, err
	}
	customers, err := s.analytics.TopCustomers(ctx, 5)
	if err != nil {
		return nil, err
	}
	aging, err := s.analytics.InvoiceAging(ctx)
	if err != nil {
		return nil, err
	}

	monthly := fillMonthly(revenue, 12, s.now())
	report := &domain.AnalyticsReport{
		GeneratedAt:     s.now(),
		Analytics:       computed,
		MonthlyRevenue:  monthly,
		PlanBreakdown:   plans,
		TopCustomers:    customers,
		InvoiceAging:    fillAging(aging),
		StatusBreakdown: statusCounts(computed.Analytics),
	}
	if len(monthly) > 0 {
		report.From, _ = time.Parse("2006-01", monthly[0].Month)
		last, _ := time.Parse("2006-01", monthly[len(monthly)-1].Month)
		report.To = addMonthsClamped(last, 1).AddDate(0, 0, -1)
	}
	return report, nil
}

// Transactions returns a paginated settlement stream.
func (s *DashboardService) Transactions(ctx context.Context, filter domain.TransactionFilter, actor domain.Actor) (domain.TransactionPage, error) {
	if !actor.IsStaff() {
		return domain.TransactionPage{}, domain.Forbiddenf("only staff can read transactions")
	}
	filter = filter.Normalize()
	items, total, err := s.billing.ListTransactions(ctx, filter)
	if err != nil {
		return domain.TransactionPage{}, err
	}
	return domain.TransactionPage{
		Items:      items,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: pageCount(total, filter.PageSize),
	}, nil
}

func fillMonthly(points []domain.MonthlyRevenue, months int, now time.Time) []domain.MonthlyRevenue {
	if months < 1 {
		return points
	}
	byMonth := make(map[string]domain.MonthlyRevenue, len(points))
	for _, point := range points {
		byMonth[point.Month] = point
	}

	current := utcDate(now)
	start := time.Date(current.Year(), current.Month()-time.Month(months-1), 1, 0, 0, 0, 0, time.UTC)
	result := make([]domain.MonthlyRevenue, 0, months)
	for i := 0; i < months; i++ {
		month := start.AddDate(0, i, 0)
		key := month.Format("2006-01")
		point, exists := byMonth[key]
		if !exists {
			point = domain.MonthlyRevenue{Month: key}
		}
		point.Label = month.Format("Jan 2006")
		result = append(result, point)
	}
	return result
}

func fillAging(buckets []domain.AgingBucket) []domain.AgingBucket {
	byKey := make(map[string]domain.AgingBucket, len(buckets))
	for _, bucket := range buckets {
		byKey[bucket.Bucket] = bucket
	}
	order := []struct {
		key   string
		label string
	}{
		{"current", "Current"},
		{"1-30", "1–30 days"},
		{"31-60", "31–60 days"},
		{"61-90", "61–90 days"},
		{"90+", "90+ days"},
	}
	result := make([]domain.AgingBucket, 0, len(order))
	for _, item := range order {
		bucket, exists := byKey[item.key]
		if !exists {
			bucket = domain.AgingBucket{Bucket: item.key}
		}
		bucket.Label = item.label
		result = append(result, bucket)
	}
	return result
}

func planName(sub domain.Subscription) string {
	if sub.Plan != nil {
		return sub.Plan.Name
	}
	return ""
}

func percentValue(value float64) float64 {
	return math.Round(value*10) / 10
}
