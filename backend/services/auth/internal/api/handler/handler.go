// Package handler implements the auth HTTP controllers.
package handler

import (
	"context"
	"net/http"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/api/dto"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/api/middleware"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/domain"
	"github.com/gin-gonic/gin"
)

// AuthService is the business-rule surface used by the controllers.
type AuthService interface {
	Register(ctx context.Context, input domain.RegisterInput) (domain.User, string, error)
	Login(ctx context.Context, input domain.LoginInput) (domain.User, string, error)
	Validate(ctx context.Context, token string) (domain.User, error)
	Logout(ctx context.Context, token string) error
	ResetPassword(ctx context.Context, email string) error
}

// Handlers bundles the auth HTTP controllers.
type Handlers struct {
	Auth AuthService
}

// New creates the controller bundle.
func New(auth AuthService) *Handlers {
	return &Handlers{Auth: auth}
}

// Health is the service health probe.
func (h *Handlers) Health(c *gin.Context) {
	response.JSON(c, http.StatusOK, gin.H{"status": "ok", "service": "auth"})
}

// Register handles POST /api/v1/auth/register.
func (h *Handlers) Register(c *gin.Context) {
	var request dto.RegisterRequest
	if err := response.DecodeJSON(c, &request); err != nil {
		response.Error(c, err)
		return
	}
	user, token, err := h.Auth.Register(c.Request.Context(), request.ToInput())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.AuthSession(user, token))
}

// Login handles POST /api/v1/auth/login.
func (h *Handlers) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := response.DecodeJSON(c, &request); err != nil {
		response.Error(c, err)
		return
	}
	user, token, err := h.Auth.Login(c.Request.Context(), request.ToInput())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.AuthSession(user, token))
}

// Logout handles POST /api/v1/auth/logout.
func (h *Handlers) Logout(c *gin.Context) {
	if err := h.Auth.Logout(c.Request.Context(), middleware.BearerToken(c)); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.OKMessage("Signed out"))
}

// Validate handles GET /api/v1/auth/validate.
func (h *Handlers) Validate(c *gin.Context) {
	token := middleware.BearerToken(c)
	if token == "" {
		response.Error(c, domain.ErrUnauthorized)
		return
	}
	user, err := h.Auth.Validate(c.Request.Context(), token)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.User(user))
}

// ResetPassword handles POST /api/v1/auth/password-reset.
func (h *Handlers) ResetPassword(c *gin.Context) {
	var request dto.ResetRequest
	if err := response.DecodeJSON(c, &request); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.Auth.ResetPassword(c.Request.Context(), request.Email); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.OKMessage("If an account exists, a reset link has been sent to your email."))
}
