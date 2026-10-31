package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Auth service draft.
//
// Responsibilities (to implement):
//   - Register / login / logout
//   - JWT generation and validation
//   - Password reset
//   - Role-based access control
func main() {
	port := env("PORT", "8081")

	router := gin.Default()
	router.GET("/health", health("auth"))

	api := router.Group("/api/v1/auth")
	{
		api.POST("/register", stub("register"))
		api.POST("/login", stub("login"))
		api.POST("/logout", stub("logout"))
		api.POST("/password-reset", stub("password reset"))
		api.GET("/validate", stub("JWT validation"))
	}

	log.Printf("auth service listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("auth service: %v", err)
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
