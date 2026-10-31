package domain

import (
	"strings"
	"time"
)

// UserRole mirrors the user_role enum in PostgreSQL.
type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleStaff    UserRole = "staff"
	RoleCustomer UserRole = "customer"
)

// Valid reports whether r is a value accepted by the database enum.
func (r UserRole) Valid() bool {
	switch r {
	case RoleAdmin, RoleStaff, RoleCustomer:
		return true
	default:
		return false
	}
}

// UserStatus mirrors the user_status enum in PostgreSQL.
type UserStatus string

const (
	UserActive    UserStatus = "active"
	UserInactive  UserStatus = "inactive"
	UserSuspended UserStatus = "suspended"
)

// Valid reports whether s is a value accepted by the database enum.
func (s UserStatus) Valid() bool {
	switch s {
	case UserActive, UserInactive, UserSuspended:
		return true
	default:
		return false
	}
}

// User is an account that can sign in. PasswordHash is never serialized.
type User struct {
	ID              string
	Name            string
	Email           string
	PasswordHash    string
	Role            UserRole
	Status          UserStatus
	CompanyName     string
	Phone           string
	AvatarURL       string
	EmailVerifiedAt *time.Time
	LastLoginAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Principal carries the identity resolved from an access token.
type Principal struct {
	UserID string
	Name   string
	Email  string
	Role   UserRole
}

// RegisterInput is the validated register request.
type RegisterInput struct {
	Name        string
	Email       string
	Password    string
	CompanyName string
	Role        UserRole
}

// LoginInput is the validated login request.
type LoginInput struct {
	Email    string
	Password string
}

// ResetInput is the validated password-reset request.
type ResetInput struct {
	Email string
}

// Session is one issued authentication session; it can be revoked on logout.
type Session struct {
	ID               string
	UserID           string
	RefreshTokenHash string
	UserAgent        string
	IPAddress        string
	ExpiresAt        time.Time
	RevokedAt        *time.Time
}

// PasswordResetToken is a hashed password-reset code issued for a user.
type PasswordResetToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
}

// Sanitize trims whitespace and normalizes case-sensitive fields.
func (r *RegisterInput) Sanitize() {
	r.Name = strings.TrimSpace(r.Name)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.CompanyName = strings.TrimSpace(r.CompanyName)
	r.Role = UserRole(strings.ToLower(strings.TrimSpace(string(r.Role))))
}
