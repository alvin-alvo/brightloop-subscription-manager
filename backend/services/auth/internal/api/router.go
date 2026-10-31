// Package api builds the auth service router.
package api

import (
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/api/handler"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/api/middleware"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/config"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/domain"
	"github.com/gin-gonic/gin"
)

// NewRouter builds the auth HTTP router. Routes mirror the scaffold contract
// expected by the gateway at /api/v1/auth.
func NewRouter(cfg config.Config, handlers *handler.Handlers) *gin.Engine {
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.AllowedOrigins))

	router.GET("/health", handlers.Health)

	api := router.Group("/api/v1/auth")
	{
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)
		api.POST("/logout", handlers.Logout)
		api.GET("/validate", handlers.Validate)
		api.POST("/password-reset", handlers.ResetPassword)
	}

	router.NoRoute(func(c *gin.Context) {
		response.Error(c, domain.NotFoundf("route %s %s not found", c.Request.Method, c.Request.URL.Path))
	})
	return router
}
