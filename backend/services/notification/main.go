package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Notification service draft.
//
// Responsibilities (to implement):
//   - Email notifications
//   - In-app notifications
//   - Push notifications
//   - Notification templates
func main() {
	port := env("PORT", "8084")

	router := gin.Default()
	router.GET("/health", health("notification"))

	api := router.Group("/api/v1/notifications")
	{
		api.GET("", stub("list in-app notifications"))
		api.POST("/email", stub("send email notification"))
		api.POST("/push", stub("send push notification"))
		api.GET("/templates", stub("list notification templates"))
		api.POST("/templates", stub("create notification template"))
	}

	log.Printf("notification service listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("notification service: %v", err)
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
