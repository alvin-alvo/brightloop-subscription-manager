// Package middleware provides CORS handling and bearer-token identity
// resolution.
package middleware

import (
	"net/http"
	"strings"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/auth"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/config"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/domain"
	"github.com/gin-gonic/gin"
)

type principalKey struct{}

// Auth resolves the caller from the bearer token. When no token is present and
// auth is not required, the development default user id is used so the seeded
// demo notifications are visible.
func Auth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		principal := domain.Principal{UserID: cfg.DefaultUserID}
		token := BearerToken(c)
		if token != "" {
			parsed, err := auth.ParsePrincipal(token, cfg.JWTSecret)
			if err == nil {
				principal = parsed
			} else if cfg.AuthRequired {
				response.Error(c, domain.ErrUnauthorized)
				return
			}
		} else if cfg.AuthRequired {
			response.Error(c, domain.ErrUnauthorized)
			return
		}
		c.Set("principal", principal)
		c.Next()
	}
}

// Actor returns the caller resolved by Auth.
func Actor(c *gin.Context) domain.Principal {
	value, ok := c.Get("principal")
	if !ok {
		return domain.Principal{}
	}
	principal, ok := value.(domain.Principal)
	if !ok {
		return domain.Principal{}
	}
	return principal
}

// CORS allows the Vite frontend to call the service cross-origin.
func CORS(origins []string) gin.HandlerFunc {
	allowAll := false
	allowed := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		if origin == "*" {
			allowAll = true
		}
		allowed[origin] = struct{}{}
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if allowAll {
				c.Header("Access-Control-Allow-Origin", "*")
			} else if _, ok := allowed[origin]; ok {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Vary", "Origin")
			}
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Header("Access-Control-Max-Age", "86400")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// BearerToken extracts a token from the Authorization header.
func BearerToken(c *gin.Context) string {
	header := strings.TrimSpace(c.GetHeader("Authorization"))
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
