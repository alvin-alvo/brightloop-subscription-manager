// Package auth provides JWT issuing/verification, password hashing and token
// helpers used by the auth service. It does not depend on the API layer.
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// TokenManager signs and verifies access tokens.
type TokenManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
	now    func() time.Time
}

// NewTokenManager builds a TokenManager with the given signing secret.
func NewTokenManager(secret, issuer string, ttl time.Duration) *TokenManager {
	now := time.Now
	return &TokenManager{secret: []byte(secret), issuer: issuer, ttl: ttl, now: now}
}

// Claims are the parsed, typed contents of an access token.
type Claims struct {
	Role      string `json:"role"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	UserID    string `json:"user_id"`
	SessionID string `json:"sid"`
	jwt.RegisteredClaims
}

// Issue signs an access token for the user, embedding the session id so logout
// can revoke the session. It returns the signed token string.
func (m *TokenManager) Issue(user domain.User, sessionID string) (string, error) {
	now := m.now().UTC()
	claims := Claims{
		Role:      string(user.Role),
		Name:      user.Name,
		Email:     user.Email,
		UserID:    user.ID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ErrInvalidToken is returned when a token cannot be verified or is malformed.
var ErrInvalidToken = errors.New("invalid token")

// Parse verifies the token signature and expiry against the manager settings
// and returns the typed claims.
func (m *TokenManager) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	}, jwt.WithValidMethods([]string{"HS256", "HS384", "HS512"}), jwt.WithIssuer(m.issuer))
	if err != nil || token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || claims.UserID == "" {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// ToPrincipal maps claims to a domain principal.
func (c *Claims) ToPrincipal() domain.Principal {
	return domain.Principal{
		UserID: c.UserID,
		Name:   c.Name,
		Email:  c.Email,
		Role:   domain.UserRole(c.Role),
	}
}

// HashPassword hashes a plaintext password with bcrypt.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword reports whether the plaintext matches the stored bcrypt hash.
func VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// RandomToken returns a URL-safe random token suitable for refresh tokens or
// password-reset codes.
func RandomToken(bytes int) (string, error) {
	buffer := make([]byte, bytes)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

// HashToken hashes a token value before storage (SHA-256 hex digest).
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
