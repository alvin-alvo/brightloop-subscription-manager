package domain

import (
	"time"

	"github.com/google/uuid"
)

// Role is the access level of a user.
type Role string

const (
	RoleAdmin    Role = "admin"
	RoleStaff    Role = "staff"
	RoleCustomer Role = "customer"
)

// Valid reports whether the role is one of the supported values.
func (r Role) Valid() bool {
	switch r {
	case RoleAdmin, RoleStaff, RoleCustomer:
		return true
	default:
		return false
	}
}

// User is an account that can sign in. Customers are linked to subscriptions
// and invoices; staff and admins operate the system.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FullName     string
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
