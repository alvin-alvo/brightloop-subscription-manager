package middleware

import (
	"net/http"
	"strings"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/config"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type principalKey struct{}

// Auth accepts a bearer token. When JWT_SECRET is configured, tokens are
// verified; otherwise claims are read unverified so the mock frontend token
// continues to work in development. Set AUTH_REQUIRED=true to reject missing or
// invalid tokens.
func Auth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		principal := domain.Actor{Role: cfg.DefaultRole}
		if principal.Role == "" {
			principal.Role = "admin"
		}

		token := bearerToken(c.GetHeader("Authorization"))
		if token == "" {
			if cfg.AuthRequired {
				response.Error(c, domain.ErrUnauthorized)
				return
			}
			c.Set("principal", principal)
			c.Next()
			return
		}

		if role, ok := mockRole(token); ok {
			principal.Role = role
			c.Set("principal", principal)
			c.Next()
			return
		}

		claims := jwt.MapClaims{}
		var tokenObj *jwt.Token
		var err error
		if cfg.JWTSecret != "" {
			tokenObj, err = jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
				return []byte(cfg.JWTSecret), nil
			}, jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}))
		} else {
			tokenObj, _, err = jwt.NewParser().ParseUnverified(token, claims)
		}
		if err != nil || tokenObj == nil || !tokenObj.Valid && cfg.JWTSecret != "" {
			if cfg.AuthRequired {
				response.Error(c, domain.ErrUnauthorized)
				return
			}
			c.Set("principal", principal)
			c.Next()
			return
		}

		principal.UserID = firstClaim(claims, "sub", "user_id", "uid", "id")
		principal.CustomerID = firstClaim(claims, "customer_id", "customerId")
		if role := strings.ToLower(firstClaim(claims, "role")); role != "" {
			principal.Role = role
		}
		if principal.Role != "admin" && principal.Role != "staff" && principal.Role != "customer" {
			principal.Role = cfg.DefaultRole
		}
		c.Set("principal", principal)
		c.Next()
	}
}

// RequireStaff blocks customers from staff-only endpoints.
func RequireStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		principal := Actor(c)
		if !principal.IsStaff() {
			response.Error(c, domain.ErrForbidden)
			return
		}
		c.Next()
	}
}

// Actor returns the authenticated caller from the request context.
func Actor(c *gin.Context) domain.Actor {
	value, exists := c.Get("principal")
	if !exists {
		return domain.Actor{Role: "admin"}
	}
	principal, ok := value.(domain.Actor)
	if !ok {
		return domain.Actor{Role: "admin"}
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

func bearerToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func mockRole(token string) (string, bool) {
	if !strings.HasPrefix(token, "mock_token_") {
		return "", false
	}
	role := strings.TrimPrefix(token, "mock_token_")
	if role == "admin" || role == "staff" || role == "customer" {
		return role, true
	}
	return "", false
}

func firstClaim(claims jwt.MapClaims, keys ...string) string {
	for _, key := range keys {
		value, ok := claims[key]
		if !ok || value == nil {
			continue
		}
		if text, ok := value.(string); ok && strings.TrimSpace(text) != "" {
			return strings.TrimSpace(text)
		}
	}
	return ""
}
