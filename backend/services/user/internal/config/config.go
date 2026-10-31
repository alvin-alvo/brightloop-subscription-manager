package config

import (
	"os"
	"strconv"
	"strings"
)

// Config is the user service runtime configuration.
type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	AuthRequired   bool
	DefaultRole    string
	DefaultUserID  string
	AllowedOrigins []string
}

// Load reads configuration from environment variables with development
// friendly defaults.
func Load() Config {
	return Config{
		Port:           env("PORT", "8082"),
		DatabaseURL:    env("DATABASE_URL", "postgres://brightloop:brightloop@localhost:5432/brightloop?sslmode=disable"),
		JWTSecret:      env("JWT_SECRET", ""),
		AuthRequired:   envBool("AUTH_REQUIRED", false),
		DefaultRole:    strings.ToLower(env("AUTH_DEFAULT_ROLE", "admin")),
		DefaultUserID:  env("AUTH_DEFAULT_USER_ID", ""),
		AllowedOrigins: splitCSV(env("CORS_ALLOWED_ORIGINS", "*")),
	}
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
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
