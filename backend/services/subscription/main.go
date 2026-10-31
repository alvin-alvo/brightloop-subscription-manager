package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Subscription service draft.
//
// Responsibilities (to implement):
//   - Subscription CRUD
//   - Billing cycles
//   - Renewals and reminders
//   - Analytics and reports
func main() {
	port := env("PORT", "8083")

	router := gin.Default()
	router.GET("/health", health("subscription"))

	api := router.Group("/api/v1/subscriptions")
	{
		api.GET("", stub("list subscriptions"))
		api.POST("", stub("create subscription"))
		api.GET("/analytics", stub("analytics"))
		api.GET("/reports", stub("reports"))
		api.GET("/:id", stub("get subscription"))
		api.PUT("/:id", stub("update subscription"))
		api.DELETE("/:id", stub("cancel subscription"))
		api.GET("/:id/billing", stub("billing cycles"))
		api.POST("/:id/renew", stub("renew subscription"))
	}

	log.Printf("subscription service listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("subscription service: %v", err)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func health(service string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": service})
	}
}

func stub(feature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"error":   "not implemented",
			"feature": feature,
		})
	}
}
