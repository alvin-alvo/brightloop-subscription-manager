package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/api/response"
)

// Pinger is the subset of the database handle needed for readiness checks.
type Pinger interface {
	PingContext(ctx context.Context) error
}

// HealthHandler serves liveness and readiness probes.
type HealthHandler struct {
	db Pinger
}

// NewHealthHandler builds health handlers around a database connection.
func NewHealthHandler(db Pinger) *HealthHandler {
	return &HealthHandler{db: db}
}

// Live reports that the process is running.
func (h *HealthHandler) Live(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready reports that the process can reach its dependencies.
func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		response.JSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
