// Package dto converts domain values to and from API request/response shapes.
package dto

import (
	"strconv"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/domain"
)

// UserResponse is the session-user shape (matches the frontend User type).
type UserResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Role            string  `json:"role"`
	Status          string  `json:"status"`
	CompanyName     string  `json:"companyName,omitempty"`
	Phone           string  `json:"phone,omitempty"`
	AvatarURL       string  `json:"avatarUrl,omitempty"`
	EmailVerifiedAt *string `json:"emailVerifiedAt,omitempty"`
	LastLoginAt     *string `json:"lastLoginAt,omitempty"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// User converts a domain user, never exposing the password hash.
func User(u domain.User) UserResponse {
	return UserResponse{
		ID:              u.ID,
		Name:            u.Name,
		Email:           u.Email,
		Role:            string(u.Role),
		Status:          string(u.Status),
		CompanyName:     u.CompanyName,
		Phone:           u.Phone,
		AvatarURL:       u.AvatarURL,
		EmailVerifiedAt: timestampPtr(u.EmailVerifiedAt),
		LastLoginAt:     timestampPtr(u.LastLoginAt),
		CreatedAt:       timestamp(u.CreatedAt),
		UpdatedAt:       timestamp(u.UpdatedAt),
	}
}

// StaffResponse is the staff-management shape (matches AdminStaff).
type StaffResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
	Status     string `json:"status"`
	LastActive string `json:"lastActive"`
	CreatedAt  string `json:"createdAt"`
}

// Staff converts a domain user for the staff-management page.
func Staff(u domain.User, now time.Time) StaffResponse {
	return StaffResponse{
		ID:         u.ID,
		Name:       u.Name,
		Email:      u.Email,
		Role:       string(u.Role),
		Department: u.Department,
		Status:     string(u.Status),
		LastActive: relativeLastActive(u, now),
		CreatedAt:  timestamp(u.CreatedAt),
	}
}

// MessageResponse is a generic success payload.
type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// OKMessage builds a success message response.
func OKMessage(message string) MessageResponse {
	return MessageResponse{Success: true, Message: message}
}

// UpdateProfileRequest is the /me update payload.
type UpdateProfileRequest struct {
	Name        string `json:"name"`
	CompanyName string `json:"companyName"`
	Phone       string `json:"phone"`
	AvatarURL   string `json:"avatarUrl"`
}

// ToPatch converts to the domain patch.
func (r UpdateProfileRequest) ToPatch() domain.UpdateProfilePatch {
	return domain.UpdateProfilePatch{
		Name:        r.Name,
		CompanyName: r.CompanyName,
		Phone:       r.Phone,
		AvatarURL:   r.AvatarURL,
	}
}

// InviteRequest is the admin invite payload.
type InviteRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Department string `json:"department"`
}

// ToInput converts to the domain input.
func (r InviteRequest) ToInput() domain.InviteInput {
	return domain.InviteInput{
		Name:       r.Name,
		Email:      r.Email,
		Role:       domain.UserRole(r.Role),
		Department: r.Department,
	}
}

// UpdateUserRequest is the admin edit payload.
type UpdateUserRequest struct {
	Name       string `json:"name"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	Department string `json:"department"`
}

// ToPatch converts to the domain patch.
func (r UpdateUserRequest) ToPatch() domain.UpdateUserPatch {
	return domain.UpdateUserPatch{
		Name:       r.Name,
		Role:       domain.UserRole(r.Role),
		Status:     domain.UserStatus(r.Status),
		Department: r.Department,
	}
}

func timestamp(value time.Time) string {
	return value.UTC().Format(time.RFC3339)
}

func timestampPtr(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := timestamp(*value)
	return &formatted
}

// relativeLastActive renders a display-ready relative time like the frontend
// mock ("Just now", "2h ago", "Yesterday").
func relativeLastActive(u domain.User, now time.Time) string {
	if u.Status == domain.UserInvited {
		return "Just invited"
	}
	base := u.CreatedAt
	if u.LastLoginAt != nil {
		base = *u.LastLoginAt
	}
	base = base.UTC()
	now = now.UTC()
	diff := now.Sub(base)
	switch {
	case diff < time.Minute:
		return "Just now"
	case diff < time.Hour:
		return strconv.Itoa(int(diff/time.Minute)) + "m ago"
	case diff < 24*time.Hour:
		return strconv.Itoa(int(diff/time.Hour)) + "h ago"
	case diff < 48*time.Hour:
		return "Yesterday"
	default:
		return base.Format("Jan 2, 2006")
	}
}
