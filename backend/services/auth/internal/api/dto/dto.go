// Package dto converts domain values to and from API request/response shapes.
package dto

import (
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/domain"
)

// UserResponse mirrors the frontend User interface (camelCase).
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

// AuthSessionResponse is the payload returned by register and login.
type AuthSessionResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// AuthSession builds a session response.
func AuthSession(u domain.User, token string) AuthSessionResponse {
	return AuthSessionResponse{User: User(u), Token: token}
}

// MessageResponse is a generic success payload (logout, password reset).
type MessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// OKMessage builds a success message response.
func OKMessage(message string) MessageResponse {
	return MessageResponse{Success: true, Message: message}
}

// RegisterRequest is the register request body.
type RegisterRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	CompanyName string `json:"companyName"`
	Role        string `json:"role"`
}

// ToInput converts to the domain input.
func (r RegisterRequest) ToInput() domain.RegisterInput {
	return domain.RegisterInput{
		Name:        r.Name,
		Email:       r.Email,
		Password:    r.Password,
		CompanyName: r.CompanyName,
		Role:        domain.UserRole(r.Role),
	}
}

// LoginRequest is the login request body.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ToInput converts to the domain input.
func (r LoginRequest) ToInput() domain.LoginInput {
	return domain.LoginInput{Email: r.Email, Password: r.Password}
}

// ResetRequest is the password-reset request body.
type ResetRequest struct {
	Email string `json:"email"`
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
