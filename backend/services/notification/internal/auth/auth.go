// Package auth provides bearer-token parsing used by the notification service
// to identify the target user.
package auth

import (
	"errors"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken is returned when a token cannot be parsed.
var ErrInvalidToken = errors.New("invalid token")

// ParsePrincipal extracts the user id (and role) from a bearer token. When
// secret is empty the token is read unverified (development).
func ParsePrincipal(tokenString, secret string) (domain.Principal, error) {
	claims := &jwt.MapClaims{}
	if secret != "" {
		token, terr := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}))
		if terr != nil || token == nil || !token.Valid {
			return domain.Principal{}, ErrInvalidToken
		}
	} else {
		token, _, terr := jwt.NewParser().ParseUnverified(tokenString, claims)
		if terr != nil || token == nil {
			return domain.Principal{}, ErrInvalidToken
		}
	}
	principal := domain.Principal{UserID: firstClaim(*claims, "user_id", "sub", "uid", "id")}
	principal.Role = firstClaim(*claims, "role")
	return principal, nil
}

func firstClaim(claims jwt.MapClaims, keys ...string) string {
	for _, key := range keys {
		value, ok := claims[key]
		if !ok || value == nil {
			continue
		}
		if text, ok := value.(string); ok && text != "" {
			return text
		}
	}
	return ""
}
