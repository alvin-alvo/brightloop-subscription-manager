package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/domain"
)

const maxBodyBytes = 1 << 20 // 1 MiB

// ErrorBody is the standard JSON error payload.
type ErrorBody struct {
	Error string `json:"error"`
}

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if payload == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("encode response", "error", err)
	}
}

// Error maps domain errors to HTTP status codes. Unexpected errors are logged
// and returned as a generic 500 so internals are not leaked to clients.
func Error(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrValidation):
		JSON(w, http.StatusUnprocessableEntity, ErrorBody{Error: err.Error()})
	case errors.Is(err, domain.ErrUnauthorized):
		JSON(w, http.StatusUnauthorized, ErrorBody{Error: "unauthorized"})
	case errors.Is(err, domain.ErrForbidden):
		JSON(w, http.StatusForbidden, ErrorBody{Error: "forbidden"})
	case errors.Is(err, domain.ErrNotFound):
		JSON(w, http.StatusNotFound, ErrorBody{Error: "not found"})
	case errors.Is(err, domain.ErrConflict):
		JSON(w, http.StatusConflict, ErrorBody{Error: "resource already exists"})
	default:
		slog.Error("unhandled error", "error", err)
		JSON(w, http.StatusInternalServerError, ErrorBody{Error: "internal server error"})
	}
}

// DecodeJSON decodes a single JSON object from the request body.
func DecodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, maxBodyBytes))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("%w: invalid JSON body", domain.ErrValidation)
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w: request body must contain a single JSON object", domain.ErrValidation)
	}

	return nil
}
