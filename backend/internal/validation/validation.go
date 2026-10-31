package validation

import (
	"fmt"
	"net/mail"
	"unicode"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/domain"
)

// Email checks that the value is a single, plain email address.
func Email(value string) error {
	address, err := mail.ParseAddress(value)
	if err != nil || address.Address != value {
		return fmt.Errorf("%w: invalid email address", domain.ErrValidation)
	}
	return nil
}

// Password enforces the product rule: more than eight characters, with an
// uppercase letter, a lowercase letter and a special character.
func Password(value string) error {
	if len(value) <= 8 {
		return fmt.Errorf("%w: password must be more than 8 characters", domain.ErrValidation)
	}

	var hasUpper, hasLower, hasSpecial bool
	for _, r := range value {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case !unicode.IsLetter(r) && !unicode.IsDigit(r):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasSpecial {
		return fmt.Errorf(
			"%w: password must contain an uppercase letter, a lowercase letter and a special character",
			domain.ErrValidation,
		)
	}

	return nil
}
