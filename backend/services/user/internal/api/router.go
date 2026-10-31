// Package api builds the user service router.
package api

import (
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/api/handler"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/api/middleware"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/config"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/domain"
	"github.com/gin-gonic/gin"
)

// NewRouter builds the user HTTP router.
func NewRouter(cfg config.Config, handlers *handler.Handlers) *gin.Engine {
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.AllowedOrigins))

	router.GET("/health", handlers.Health)

	api := router.Group("/api/v1/users")
	api.Use(middleware.Auth(cfg))
	{
		api.GET("/me", handlers.Me)
		api.PUT("/me", handlers.UpdateMe)
		api.DELETE("/me", handlers.DeleteMe)
		api.GET("/me/preferences", handlers.Preferences)
		api.PUT("/me/preferences", handlers.UpdatePreferences)
		api.GET("", handlers.ListUsers)
		api.POST("", handlers.InviteUser)
		api.GET("/:id", handlers.GetUser)
		api.PUT("/:id", handlers.UpdateUser)
	}

	router.NoRoute(func(c *gin.Context) {
		response.Error(c, domain.NotFoundf("route %s %s not found", c.Request.Method, c.Request.URL.Path))
	})
	return router
}
