package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// User service draft.
//
// Responsibilities (to implement):
//   - Profile management
//   - Preferences
//   - Account settings
func main() {
	port := env("PORT", "8082")

	router := gin.Default()
	router.GET("/health", health("user"))

	api := router.Group("/api/v1/users")
	{
		api.GET("/me", stub("get profile"))
		api.PUT("/me", stub("update profile"))
		api.GET("/me/preferences", stub("get preferences"))
		api.PUT("/me/preferences", stub("update preferences"))
		api.DELETE("/me", stub("delete account"))
		api.GET("", stub("list users"))
		api.GET("/:id", stub("get user"))
	}

	log.Printf("user service listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("user service: %v", err)
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
