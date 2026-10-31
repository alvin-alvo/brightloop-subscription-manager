package domain

import (
	"errors"
	"fmt"
)

// Sentinel errors returned by repositories and services. The API layer maps
// these to HTTP status codes.
var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrValidation   = errors.New("validation failed")
	ErrInternal     = errors.New("internal error")
)

// Validationf wraps ErrValidation with a human readable message.
func Validationf(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrValidation, fmt.Sprintf(format, args...))
}

// Conflictf wraps ErrConflict with a human readable message.
func Conflictf(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrConflict, fmt.Sprintf(format, args...))
}

// NotFoundf wraps ErrNotFound with a human readable message.
func NotFoundf(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrNotFound, fmt.Sprintf(format, args...))
}

// Unauthorized wraps ErrUnauthorized with a message-safe reason.
func Unauthorized(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrUnauthorized, fmt.Sprintf(format, args...))
}

// Forbiddenf wraps ErrForbidden with a human readable message.
func Forbiddenf(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrForbidden, fmt.Sprintf(format, args...))
}
