package response

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
	"github.com/gin-gonic/gin"
)

const maxBodyBytes = 1 << 20

type errorBody struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// JSON writes a successful response.
func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

// DecodeJSON decodes and validates a JSON request body.
func DecodeJSON(c *gin.Context, target any) error {
	if c.Request.Body == nil {
		return domain.Validationf("request body is required")
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodyBytes)
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(target); err != nil {
		if errors.Is(err, io.EOF) {
			return domain.Validationf("request body is required")
		}
		return domain.Validationf("invalid JSON body: %v", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != nil && !errors.Is(err, io.EOF) {
		return domain.Validationf("invalid JSON body: %v", err)
	}
	return nil
}

// Error maps service/repository errors to HTTP status codes and a consistent
// JSON body. The frontend's api client reads the `message` field.
func Error(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	code := "internal_error"
	message := "internal server error"

	switch {
	case errors.Is(err, domain.ErrValidation):
		status, code, message = http.StatusBadRequest, "validation_error", cleanMessage(err)
	case errors.Is(err, domain.ErrNotFound):
		status, code, message = http.StatusNotFound, "not_found", cleanMessage(err)
	case errors.Is(err, domain.ErrConflict):
		status, code, message = http.StatusConflict, "conflict", cleanMessage(err)
	case errors.Is(err, domain.ErrUnauthorized):
		status, code, message = http.StatusUnauthorized, "unauthorized", cleanMessage(err)
	case errors.Is(err, domain.ErrForbidden):
		status, code, message = http.StatusForbidden, "forbidden", cleanMessage(err)
	default:
		log.Printf("unhandled error: %v", err)
	}

	c.AbortWithStatusJSON(status, errorBody{
		Success: false,
		Error:   code,
		Message: message,
	})
}

func cleanMessage(err error) string {
	message := err.Error()
	for _, prefix := range []string{
		"validation failed: ",
		"not found: ",
		"conflict: ",
		"unauthorized: ",
		"forbidden: ",
		"internal error: ",
	} {
		message = strings.TrimPrefix(message, prefix)
	}
	return message
}
