package handler

import (
	"net/http"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/api/middleware"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/domain"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/service"
)

// AuthHandler exposes registration, login and profile endpoints.
type AuthHandler struct {
	auth *service.AuthService
}

// NewAuthHandler builds auth handlers around the auth service.
func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type authResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

// Register creates a customer account and returns an access token.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	user, err := h.auth.Register(r.Context(), service.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
	})
	if err != nil {
		response.Error(w, err)
		return
	}

	token, err := h.auth.IssueToken(user)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, authResponse{Token: token, User: toUserResponse(user)})
}

// Login verifies credentials and returns an access token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.Error(w, err)
		return
	}

	token, user, err := h.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, authResponse{Token: token, User: toUserResponse(user)})
}

// Me returns the authenticated user's profile.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, domain.ErrUnauthorized)
		return
	}

	user, err := h.auth.ByID(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.JSON(w, http.StatusOK, toUserResponse(user))
}

func toUserResponse(user *domain.User) userResponse {
	return userResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}
}
