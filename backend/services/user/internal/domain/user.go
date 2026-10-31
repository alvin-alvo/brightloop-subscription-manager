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

// StaffRole reports whether r is an internal (non-customer) role.
func (r UserRole) StaffRole() bool {
	return r == RoleAdmin || r == RoleStaff
}

// UserStatus mirrors the user_status enum in PostgreSQL.
type UserStatus string

const (
	UserActive    UserStatus = "active"
	UserInvited   UserStatus = "invited"
	UserInactive  UserStatus = "inactive"
	UserSuspended UserStatus = "suspended"
)

// Valid reports whether s is a value accepted by the database enum.
func (s UserStatus) Valid() bool {
	switch s {
	case UserActive, UserInvited, UserInactive, UserSuspended:
		return true
	default:
		return false
	}
}

// User is an account managed by the user service. PasswordHash is never
// serialized by the API layer.
type User struct {
	ID              string
	Name            string
	Email           string
	PasswordHash    string
	Role            UserRole
	Status          UserStatus
	CompanyName     string
	Department      string
	Phone           string
	AvatarURL       string
	Preferences     map[string]any
	EmailVerifiedAt *time.Time
	LastLoginAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Principal carries the identity resolved from the bearer token.
type Principal struct {
	UserID string
	Role   UserRole
}

// IsStaff reports whether the principal holds an internal role.
func (p Principal) IsStaff() bool {
	return p.Role == RoleAdmin || p.Role == RoleStaff
}

// UpdateProfilePatch is an allowed set of self-service profile fields.
type UpdateProfilePatch struct {
	Name        string
	CompanyName string
	Phone       string
	AvatarURL   string
}

// InviteInput is the validated admin invite request.
type InviteInput struct {
	Name       string
	Email      string
	Role       UserRole
	Department string
}

// UpdateUserPatch is the admin edit request.
type UpdateUserPatch struct {
	Name       string
	Role       UserRole
	Status     UserStatus
	Department string
}

func cleanEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// Sanitize normalizes an invite.
func (i *InviteInput) Sanitize() {
	i.Name = strings.TrimSpace(i.Name)
	i.Email = cleanEmail(i.Email)
	i.Department = strings.TrimSpace(i.Department)
	i.Role = UserRole(strings.ToLower(strings.TrimSpace(string(i.Role))))
}
