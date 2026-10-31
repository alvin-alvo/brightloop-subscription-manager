package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/dto"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/middleware"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
	"github.com/gin-gonic/gin"
)

// SubscriptionService is the business-rule surface used by subscription
// controllers.
type SubscriptionService interface {
	List(ctx context.Context, filter domain.ListFilter, actor domain.Actor) (domain.SubscriptionPage, error)
	Get(ctx context.Context, id string, actor domain.Actor) (*domain.Subscription, error)
	Create(ctx context.Context, input domain.CreateSubscriptionInput, actor domain.Actor) (*domain.Subscription, error)
	Update(ctx context.Context, id string, input domain.UpdateSubscriptionInput, actor domain.Actor) (*domain.Subscription, error)
	Cancel(ctx context.Context, id string, immediate bool, actor domain.Actor) (*domain.Subscription, error)
	Pause(ctx context.Context, id string, actor domain.Actor) (*domain.Subscription, error)
	Resume(ctx context.Context, id string, actor domain.Actor) (*domain.Subscription, error)
	BillingSchedule(ctx context.Context, id string, actor domain.Actor) (*domain.BillingSchedule, error)
	Plans(ctx context.Context, activeOnly bool, actor domain.Actor) ([]domain.Plan, error)
	Customers(ctx context.Context, search string, actor domain.Actor) ([]domain.Customer, error)
}

// DashboardService is the analytics read-model surface.
type DashboardService interface {
	Analytics(ctx context.Context, actor domain.Actor) (domain.ComputedAnalytics, error)
	Dashboard(ctx context.Context, actor domain.Actor) (*domain.Dashboard, error)
	Report(ctx context.Context, actor domain.Actor) (*domain.AnalyticsReport, error)
	Transactions(ctx context.Context, filter domain.TransactionFilter, actor domain.Actor) (domain.TransactionPage, error)
}

// BillingService is the recurring billing surface.
type BillingService interface {
	Renew(ctx context.Context, id string, actor domain.Actor) (*domain.Subscription, *domain.Invoice, error)
	Overview(ctx context.Context, actor domain.Actor) (*domain.BillingOverview, error)
	Batches(ctx context.Context, limit int, actor domain.Actor) ([]domain.BillingBatch, error)
	Run(ctx context.Context, asOf time.Time, dryRun bool, actor domain.Actor) (*domain.BillingRunResult, error)
}

// Handlers bundles the HTTP controllers.
type Handlers struct {
	Subscriptions SubscriptionService
	DashboardSvc  DashboardService
	Billing       BillingService
	Currency      string
	Now           func() time.Time
}

// New creates the controller bundle.
func New(subscriptions SubscriptionService, dashboard DashboardService, billing BillingService, currency string) *Handlers {
	if currency == "" {
		currency = "USD"
	}
	return &Handlers{
		Subscriptions: subscriptions,
		DashboardSvc:  dashboard,
		Billing:       billing,
		Currency:      currency,
		Now:           time.Now,
	}
}

// Health is the service health probe.
func (h *Handlers) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "subscription"})
}

// ListSubscriptions handles GET /subscriptions.
func (h *Handlers) ListSubscriptions(c *gin.Context) {
	filter, err := parseListFilter(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	page, err := h.Subscriptions.List(c.Request.Context(), filter, middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.SubscriptionPage(page))
}

// CreateSubscription handles POST /subscriptions.
func (h *Handlers) CreateSubscription(c *gin.Context) {
	var request dto.CreateRequest
	if err := response.DecodeJSON(c, &request); err != nil {
		response.Error(c, err)
		return
	}
	input, err := request.ToInput()
	if err != nil {
		response.Error(c, err)
		return
	}
	subscription, err := h.Subscriptions.Create(c.Request.Context(), input, middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.Subscription(*subscription))
}

// GetSubscription handles GET /subscriptions/:id.
func (h *Handlers) GetSubscription(c *gin.Context) {
	subscription, err := h.Subscriptions.Get(c.Request.Context(), c.Param("id"), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Subscription(*subscription))
}

// UpdateSubscription handles PUT /subscriptions/:id.
func (h *Handlers) UpdateSubscription(c *gin.Context) {
	var request dto.UpdateRequest
	if err := response.DecodeJSON(c, &request); err != nil {
		response.Error(c, err)
		return
	}
	input, err := request.ToInput()
	if err != nil {
		response.Error(c, err)
		return
	}
	subscription, err := h.Subscriptions.Update(c.Request.Context(), c.Param("id"), input, middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Subscription(*subscription))
}

// CancelSubscription handles DELETE /subscriptions/:id and POST /:id/cancel.
func (h *Handlers) CancelSubscription(c *gin.Context) {
	immediate := false
	if raw := c.Query("immediate"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			response.Error(c, domain.Validationf("immediate must be true or false"))
			return
		}
		immediate = parsed
	}
	subscription, err := h.Subscriptions.Cancel(c.Request.Context(), c.Param("id"), immediate, middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Subscription(*subscription))
}

// PauseSubscription handles POST /subscriptions/:id/pause.
func (h *Handlers) PauseSubscription(c *gin.Context) {
	subscription, err := h.Subscriptions.Pause(c.Request.Context(), c.Param("id"), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Subscription(*subscription))
}

// ResumeSubscription handles POST /subscriptions/:id/resume.
func (h *Handlers) ResumeSubscription(c *gin.Context) {
	subscription, err := h.Subscriptions.Resume(c.Request.Context(), c.Param("id"), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Subscription(*subscription))
}

// BillingSchedule handles GET /subscriptions/:id/billing.
func (h *Handlers) BillingSchedule(c *gin.Context) {
	schedule, err := h.Subscriptions.BillingSchedule(c.Request.Context(), c.Param("id"), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.BillingSchedule(*schedule))
}

// Plans handles GET /subscriptions/plans.
func (h *Handlers) Plans(c *gin.Context) {
	activeOnly := true
	if raw := c.Query("activeOnly"); raw != "" {
		parsed, err := strconv.ParseBool(raw)
		if err != nil {
			response.Error(c, domain.Validationf("activeOnly must be true or false"))
			return
		}
		activeOnly = parsed
	}
	plans, err := h.Subscriptions.Plans(c.Request.Context(), activeOnly, middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	items := make([]dto.PlanResponse, 0, len(plans))
	for _, plan := range plans {
		items = append(items, dto.Plan(plan))
	}
	response.JSON(c, http.StatusOK, items)
}

// Customers handles GET /subscriptions/customers.
func (h *Handlers) Customers(c *gin.Context) {
	customers, err := h.Subscriptions.Customers(c.Request.Context(), c.Query("search"), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	items := make([]dto.CustomerResponse, 0, len(customers))
	for _, customer := range customers {
		items = append(items, dto.Customer(customer))
	}
	response.JSON(c, http.StatusOK, items)
}

// Analytics handles GET /subscriptions/analytics and /metrics.
func (h *Handlers) Analytics(c *gin.Context) {
	computed, err := h.DashboardSvc.Analytics(c.Request.Context(), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Analytics(computed, h.Currency, h.Now()))
}

// Dashboard handles GET /subscriptions/dashboard.
func (h *Handlers) Dashboard(c *gin.Context) {
	dashboard, err := h.DashboardSvc.Dashboard(c.Request.Context(), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Dashboard(*dashboard, h.Currency))
}

// Reports handles GET /subscriptions/reports.
func (h *Handlers) Reports(c *gin.Context) {
	report, err := h.DashboardSvc.Report(c.Request.Context(), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Report(*report, h.Currency))
}

// Transactions handles GET /subscriptions/transactions.
func (h *Handlers) Transactions(c *gin.Context) {
	page, err := queryInt(c, "page", 1, 0)
	if err != nil {
		response.Error(c, err)
		return
	}
	pageSize, err := queryInt(c, "pageSize", 20, 100)
	if err != nil {
		response.Error(c, err)
		return
	}
	filter := domain.TransactionFilter{
		CustomerID: strings.TrimSpace(c.Query("customerId")),
		Page:       page,
		PageSize:   pageSize,
	}
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		filter.Status = domain.PaymentStatus(strings.ToLower(raw))
		switch filter.Status {
		case domain.PaymentPending, domain.PaymentSucceeded, domain.PaymentFailed, domain.PaymentRefunded:
		default:
			response.Error(c, domain.Validationf("invalid payment status %q", raw))
			return
		}
	}
	result, err := h.DashboardSvc.Transactions(c.Request.Context(), filter, middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.TransactionPage(result, h.Now()))
}

// BillingOverview handles GET /subscriptions/billing/overview.
func (h *Handlers) BillingOverview(c *gin.Context) {
	overview, err := h.Billing.Overview(c.Request.Context(), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.BillingOverview(*overview, h.Now()))
}

// BillingRuns handles GET /subscriptions/billing/runs.
func (h *Handlers) BillingRuns(c *gin.Context) {
	limit, err := queryInt(c, "limit", 20, 100)
	if err != nil {
		response.Error(c, err)
		return
	}
	batches, err := h.Billing.Batches(c.Request.Context(), limit, middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.BillingBatches(batches, h.Now()))
}

// RunBilling handles POST /subscriptions/billing/run.
func (h *Handlers) RunBilling(c *gin.Context) {
	request := dto.BillingRunRequest{}
	if c.Request.ContentLength != 0 {
		if err := response.DecodeJSON(c, &request); err != nil {
			response.Error(c, err)
			return
		}
	}
	var asOf time.Time
	if request.AsOf != "" {
		parsed, err := dto.ParseDate(request.AsOf)
		if err != nil {
			response.Error(c, err)
			return
		}
		asOf = parsed
	}
	result, err := h.Billing.Run(c.Request.Context(), asOf, request.DryRun, middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.BillingRun(*result))
}

// RenewSubscription handles POST /subscriptions/:id/renew.
func (h *Handlers) RenewSubscription(c *gin.Context) {
	subscription, invoice, err := h.Billing.Renew(c.Request.Context(), c.Param("id"), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.RenewResponse{
		Subscription: dto.Subscription(*subscription),
		Invoice:      dto.Invoice(*invoice),
	})
}

func parseListFilter(c *gin.Context) (domain.ListFilter, error) {
	page, err := queryInt(c, "page", 1, 0)
	if err != nil {
		return domain.ListFilter{}, err
	}
	pageSize, err := queryInt(c, "pageSize", 20, 100)
	if err != nil {
		return domain.ListFilter{}, err
	}
	filter := domain.ListFilter{
		CustomerID: strings.TrimSpace(c.Query("customerId")),
		PlanID:     strings.TrimSpace(c.Query("planId")),
		Search:     strings.TrimSpace(c.Query("search")),
		SortBy:     strings.TrimSpace(c.Query("sortBy")),
		SortOrder:  strings.TrimSpace(c.Query("sortOrder")),
		Page:       page,
		PageSize:   pageSize,
	}
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		filter.Status = domain.SubscriptionStatus(strings.ToLower(raw))
	}
	return filter, nil
}

func queryInt(c *gin.Context, key string, fallback, max int) (int, error) {
	raw := strings.TrimSpace(c.Query(key))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return 0, domain.Validationf("%s must be a positive integer", key)
	}
	if max > 0 && value > max {
		value = max
	}
	return value, nil
}
