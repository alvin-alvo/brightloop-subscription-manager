// Package auth provides bearer-token parsing and password hashing used by the
// user service to identify callers and provision invited staff.
package auth

import (
	"errors"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Claims are the typed claims in a token issued by the auth service.
type Claims struct {
	Role      string `json:"role"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	UserID    string `json:"user_id"`
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

// ErrInvalidToken is returned when a token cannot be parsed.
var ErrInvalidToken = errors.New("invalid token")

// ParsePrincipal parses a bearer token and extracts the actor. When secret is
// empty the token is read unverified (development matching the subscription
// service's behaviour); otherwise the signature is checked.
func ParsePrincipal(tokenString, secret string) (domain.Principal, error) {
	claims := &Claims{}
	var err error
	if secret != "" {
		var tokenObj *jwt.Token
		tokenObj, err = jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}))
		if err != nil || tokenObj == nil || !tokenObj.Valid {
			return domain.Principal{}, ErrInvalidToken
		}
	} else {
		var tokenObj *jwt.Token
		tokenObj, _, err = jwt.NewParser().ParseUnverified(tokenString, claims)
		if err != nil || tokenObj == nil {
			return domain.Principal{}, ErrInvalidToken
		}
	}
	principal := domain.Principal{Role: domain.UserRole(claims.Role)}
	if claims.UserID != "" {
		principal.UserID = claims.UserID
	} else if claims.Subject != "" {
		principal.UserID = claims.Subject
	}
	return principal, nil
}

// HashPassword hashes a plaintext password with bcrypt.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
