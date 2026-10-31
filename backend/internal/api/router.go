package api

import (
	"net/http"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/api/handler"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/api/middleware"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/domain"
)

// Dependencies are the handlers and middleware needed to build the router.
type Dependencies struct {
	Health *handler.HealthHandler
	Auth   *handler.AuthHandler
	AuthMW *middleware.Authenticator
}

// NewRouter registers every HTTP route. Add new modules here as their
// handlers are implemented.
func NewRouter(deps Dependencies) http.Handler {
	mux := http.NewServeMux()

	// Public routes.
	mux.HandleFunc("GET /health/live", deps.Health.Live)
	mux.HandleFunc("GET /health/ready", deps.Health.Ready)
	mux.HandleFunc("POST /api/v1/auth/register", deps.Auth.Register)
	mux.HandleFunc("POST /api/v1/auth/login", deps.Auth.Login)

	// Authenticated routes.
	mux.Handle("GET /api/v1/auth/me", middleware.Chain(
		http.HandlerFunc(deps.Auth.Me),
		deps.AuthMW.Require,
	))

	// Example admin-only route. Replace with real admin handlers as modules
	// are added.
	mux.Handle("GET /api/v1/admin/ping", middleware.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response.JSON(w, http.StatusOK, map[string]string{"scope": "admin"})
		}),
		deps.AuthMW.Require,
		deps.AuthMW.RequireRole(domain.RoleAdmin),
	))

	return mux
}
