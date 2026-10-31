// Package handler implements the user HTTP controllers.
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/api/dto"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/api/middleware"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/domain"
	"github.com/gin-gonic/gin"
)

// UserService is the business-rule surface used by the controllers.
type UserService interface {
	Me(ctx context.Context, actor domain.Principal) (*domain.User, error)
	UpdateProfile(ctx context.Context, actor domain.Principal, patch domain.UpdateProfilePatch) (*domain.User, error)
	Preferences(ctx context.Context, actor domain.Principal) (map[string]any, error)
	UpdatePreferences(ctx context.Context, actor domain.Principal, updates map[string]any) (map[string]any, error)
	DeleteAccount(ctx context.Context, actor domain.Principal) error
	ListStaff(ctx context.Context, actor domain.Principal, search string) ([]domain.User, error)
	InviteStaff(ctx context.Context, actor domain.Principal, input domain.InviteInput) (*domain.User, error)
	GetUser(ctx context.Context, actor domain.Principal, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, actor domain.Principal, id string, patch domain.UpdateUserPatch) (*domain.User, error)
}

// Handlers bundles the user HTTP controllers.
type Handlers struct {
	Users UserService
	Now   func() time.Time
}

// New creates the controller bundle.
func New(users UserService) *Handlers {
	return &Handlers{Users: users, Now: time.Now}
}

// Health is the service health probe.
func (h *Handlers) Health(c *gin.Context) {
	response.JSON(c, http.StatusOK, gin.H{"status": "ok", "service": "user"})
}

// Me handles GET /api/v1/users/me.
func (h *Handlers) Me(c *gin.Context) {
	user, err := h.Users.Me(c.Request.Context(), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.User(*user))
}

// UpdateMe handles PUT /api/v1/users/me.
func (h *Handlers) UpdateMe(c *gin.Context) {
	var request dto.UpdateProfileRequest
	if err := response.DecodeJSON(c, &request); err != nil {
		response.Error(c, err)
		return
	}
	user, err := h.Users.UpdateProfile(c.Request.Context(), middleware.Actor(c), request.ToPatch())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.User(*user))
}

// DeleteMe handles DELETE /api/v1/users/me.
func (h *Handlers) DeleteMe(c *gin.Context) {
	if err := h.Users.DeleteAccount(c.Request.Context(), middleware.Actor(c)); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.OKMessage("Account deactivated"))
}

// Preferences handles GET /api/v1/users/me/preferences.
func (h *Handlers) Preferences(c *gin.Context) {
	prefs, err := h.Users.Preferences(c.Request.Context(), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, prefs)
}

// UpdatePreferences handles PUT /api/v1/users/me/preferences.
func (h *Handlers) UpdatePreferences(c *gin.Context) {
	var updates map[string]any
	if err := response.DecodeJSON(c, &updates); err != nil {
		response.Error(c, err)
		return
	}
	prefs, err := h.Users.UpdatePreferences(c.Request.Context(), middleware.Actor(c), updates)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, prefs)
}

// ListUsers handles GET /api/v1/users (admin/staff list).
func (h *Handlers) ListUsers(c *gin.Context) {
	users, err := h.Users.ListStaff(c.Request.Context(), middleware.Actor(c), c.Query("search"))
	if err != nil {
		response.Error(c, err)
		return
	}
	items := make([]dto.StaffResponse, 0, len(users))
	for _, user := range users {
		items = append(items, dto.Staff(user, h.Now()))
	}
	response.JSON(c, http.StatusOK, items)
}

// InviteUser handles POST /api/v1/users (admin invite).
func (h *Handlers) InviteUser(c *gin.Context) {
	var request dto.InviteRequest
	if err := response.DecodeJSON(c, &request); err != nil {
		response.Error(c, err)
		return
	}
	user, err := h.Users.InviteStaff(c.Request.Context(), middleware.Actor(c), request.ToInput())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.Staff(*user, h.Now()))
}

// GetUser handles GET /api/v1/users/:id (admin).
func (h *Handlers) GetUser(c *gin.Context) {
	user, err := h.Users.GetUser(c.Request.Context(), middleware.Actor(c), c.Param("id"))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Staff(*user, h.Now()))
}

// UpdateUser handles PUT /api/v1/users/:id (admin).
func (h *Handlers) UpdateUser(c *gin.Context) {
	var request dto.UpdateUserRequest
	if err := response.DecodeJSON(c, &request); err != nil {
		response.Error(c, err)
		return
	}
	user, err := h.Users.UpdateUser(c.Request.Context(), middleware.Actor(c), c.Param("id"), request.ToPatch())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.Staff(*user, h.Now()))
}
