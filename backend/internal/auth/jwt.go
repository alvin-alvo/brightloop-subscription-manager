package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/domain"
)

// Claims is the payload carried by an access token.
type Claims struct {
	Role domain.Role `json:"role"`
	jwt.RegisteredClaims
}

// TokenManager issues and verifies HMAC-signed access tokens.
type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

// NewTokenManager builds a token manager from a shared secret and token TTL.
func NewTokenManager(secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{secret: []byte(secret), ttl: ttl}
}

// Generate creates a signed access token for a user.
func (m *TokenManager) Generate(userID uuid.UUID, role domain.Role) (string, error) {
	now := time.Now().UTC()

	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// Parse verifies a token signature and expiry, returning its claims.
func (m *TokenManager) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(_ *jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil || !token.Valid {
		return nil, domain.ErrUnauthorized
	}

	return claims, nil
}
