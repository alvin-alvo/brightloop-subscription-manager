// Package api builds the notification service router.
package api

import (
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/api/handler"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/api/middleware"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/config"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/domain"
	"github.com/gin-gonic/gin"
)

// NewRouter builds the notification HTTP router.
func NewRouter(cfg config.Config, handlers *handler.Handlers) *gin.Engine {
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.AllowedOrigins))

	router.GET("/health", handlers.Health)

	api := router.Group("/api/v1/notifications")
	api.Use(middleware.Auth(cfg))
	{
		api.GET("", handlers.List)
		api.POST("", handlers.SendInApp)
		api.GET("/unread-count", handlers.UnreadCount)
		api.POST("/read-all", handlers.MarkAllRead)
		api.GET("/templates", handlers.ListTemplates)
		api.POST("/templates", handlers.CreateTemplate)
		api.POST("/email", handlers.SendEmail)
		api.POST("/push", handlers.SendPush)
		api.POST("/:id/read", handlers.MarkRead)
	}

	router.NoRoute(func(c *gin.Context) {
		response.Error(c, domain.NotFoundf("route %s %s not found", c.Request.Method, c.Request.URL.Path))
	})
	return router
}
