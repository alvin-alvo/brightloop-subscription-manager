package middleware

import (
	"log/slog"
	"net/http"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/api/response"
)

// Recover turns a panic into a 500 response instead of crashing the server.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				slog.Error("panic recovered",
					"panic", recovered,
					"method", r.Method,
					"path", r.URL.Path,
				)
				response.JSON(w, http.StatusInternalServerError, response.ErrorBody{Error: "internal server error"})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
