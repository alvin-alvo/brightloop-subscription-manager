package api

import (
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/handler"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/middleware"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/config"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
	"github.com/gin-gonic/gin"
)

// NewRouter builds the HTTP router. The same routes are served under the
// canonical /api/v1/subscriptions prefix and the shorter /api/subscriptions
// alias used by the current frontend default base URL.
func NewRouter(cfg config.Config, handlers *handler.Handlers) *gin.Engine {
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.AllowedOrigins))

	router.GET("/health", handlers.Health)
	registerRoutes(router, "/api/v1/subscriptions", cfg, handlers)
	registerRoutes(router, "/api/subscriptions", cfg, handlers)

	router.NoRoute(func(c *gin.Context) {
		response.Error(c, domain.NotFoundf("route %s %s not found", c.Request.Method, c.Request.URL.Path))
	})
	return router
}

func registerRoutes(router *gin.Engine, prefix string, cfg config.Config, handlers *handler.Handlers) {
	group := router.Group(prefix)
	group.Use(middleware.Auth(cfg))

	group.GET("/health", handlers.Health)

	// Static routes must be registered before the /:id family.
	group.GET("", handlers.ListSubscriptions)
	group.POST("", middleware.RequireStaff(), handlers.CreateSubscription)
	group.GET("/plans", handlers.Plans)
	group.GET("/customers", middleware.RequireStaff(), handlers.Customers)
	group.GET("/dashboard", handlers.Dashboard)
	group.GET("/analytics", handlers.Analytics)
	group.GET("/metrics", handlers.Analytics)
	group.GET("/reports", handlers.Reports)
	group.GET("/transactions", middleware.RequireStaff(), handlers.Transactions)
	group.GET("/billing/overview", middleware.RequireStaff(), handlers.BillingOverview)
	group.GET("/billing/runs", middleware.RequireStaff(), handlers.BillingRuns)
	group.POST("/billing/run", middleware.RequireStaff(), handlers.RunBilling)

	group.GET("/:id", handlers.GetSubscription)
	group.PUT("/:id", middleware.RequireStaff(), handlers.UpdateSubscription)
	group.DELETE("/:id", middleware.RequireStaff(), handlers.CancelSubscription)
	group.POST("/:id/cancel", middleware.RequireStaff(), handlers.CancelSubscription)
	group.POST("/:id/pause", middleware.RequireStaff(), handlers.PauseSubscription)
	group.POST("/:id/resume", middleware.RequireStaff(), handlers.ResumeSubscription)
	group.POST("/:id/renew", middleware.RequireStaff(), handlers.RenewSubscription)
	group.GET("/:id/billing", handlers.BillingSchedule)
}
