package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/auth"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/domain"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

// Middleware is the standard net/http middleware shape.
type Middleware func(http.Handler) http.Handler

// Chain applies middleware from left to right, so the first entry is the
// outermost wrapper.
func Chain(next http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}
	return next
}

// Authenticator verifies bearer tokens and stores the caller in the request
// context.
type Authenticator struct {
	tokens *auth.TokenManager
}

// NewAuthenticator builds bearer-token middleware.
func NewAuthenticator(tokens *auth.TokenManager) *Authenticator {
	return &Authenticator{tokens: tokens}
}

// Require rejects requests without a valid bearer token.
func (a *Authenticator) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "Bearer "

		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, prefix) {
			response.Error(w, domain.ErrUnauthorized)
			return
		}

		claims, err := a.tokens.Parse(strings.TrimSpace(strings.TrimPrefix(header, prefix)))
		if err != nil {
			response.Error(w, err)
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			response.Error(w, domain.ErrUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, roleKey, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole rejects authenticated callers whose role is not allowed. It must
// run after Require.
func (a *Authenticator) RequireRole(roles ...domain.Role) Middleware {
	allowed := make(map[domain.Role]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := RoleFromContext(r.Context())
			if !ok {
				response.Error(w, domain.ErrUnauthorized)
				return
			}
			if _, ok := allowed[role]; !ok {
				response.Error(w, domain.ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// UserIDFromContext returns the authenticated user's ID.
func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

// RoleFromContext returns the authenticated user's role.
func RoleFromContext(ctx context.Context) (domain.Role, bool) {
	role, ok := ctx.Value(roleKey).(domain.Role)
	return role, ok
}
