package config

import (
	"os"
	"strings"
	"time"
)

// Config is the auth service runtime configuration.
type Config struct {
	Port             string
	DatabaseURL      string
	JWTSecret        string
	JWTIssuer        string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	AllowedOrigins   []string
	PasswordResetTTL time.Duration
}

// Load reads configuration from environment variables with development
// friendly defaults.
func Load() Config {
	return Config{
		Port:             env("PORT", "8081"),
		DatabaseURL:      env("DATABASE_URL", "postgres://brightloop:brightloop@localhost:5432/brightloop?sslmode=disable"),
		JWTSecret:        env("JWT_SECRET", ""),
		JWTIssuer:        env("JWT_ISSUER", "brightloop-auth"),
		AccessTokenTTL:   envDuration("ACCESS_TOKEN_TTL", 7*24*time.Hour),
		RefreshTokenTTL:  envDuration("REFRESH_TOKEN_TTL", 30*24*time.Hour),
		PasswordResetTTL: envDuration("PASSWORD_RESET_TTL", 30*time.Minute),
		AllowedOrigins:   splitCSV(env("CORS_ALLOWED_ORIGINS", "*")),
	}
}

// EffectiveSecret returns the signing key, falling back to a stable development
// secret when none is configured so tokens remain verifiable locally.
func (c Config) EffectiveSecret() string {
	if strings.TrimSpace(c.JWTSecret) != "" {
		return strings.TrimSpace(c.JWTSecret)
	}
	return "brightloop-dev-secret-key-change-me-1234567890"
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	if parsed <= 0 {
		return fallback
	}
	return parsed
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	if len(result) == 0 {
		return []string{"*"}
	}
	return result
}
