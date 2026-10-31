// Package handler implements the notification HTTP controllers.
package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/api/dto"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/api/middleware"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/api/response"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/domain"
	"github.com/gin-gonic/gin"
)

// NotificationService is the business-rule surface used by the controllers.
type NotificationService interface {
	List(ctx context.Context, actor domain.Principal, limit int, unreadOnly bool) ([]domain.Notification, int, error)
	UnreadCount(ctx context.Context, actor domain.Principal) (int, error)
	MarkRead(ctx context.Context, actor domain.Principal, id string) error
	MarkAllRead(ctx context.Context, actor domain.Principal) error
	Send(ctx context.Context, input domain.CreateNotificationInput) (*domain.Notification, error)
	CreateTemplate(ctx context.Context, input domain.CreateTemplateInput) (*domain.NotificationTemplate, error)
	ListTemplates(ctx context.Context, limit int) ([]domain.NotificationTemplate, error)
}

// Handlers bundles the notification HTTP controllers.
type Handlers struct {
	Notifications NotificationService
	Now           func() time.Time
}

// New creates the controller bundle.
func New(notifications NotificationService) *Handlers {
	return &Handlers{Notifications: notifications, Now: time.Now}
}

// Health is the service health probe.
func (h *Handlers) Health(c *gin.Context) {
	response.JSON(c, http.StatusOK, gin.H{"status": "ok", "service": "notification"})
}

// List handles GET /api/v1/notifications.
func (h *Handlers) List(c *gin.Context) {
	limit, err := queryInt(c, "limit", 20, 100)
	if err != nil {
		response.Error(c, err)
		return
	}
	unreadOnly := c.Query("unreadOnly") == "true"
	items, count, err := h.Notifications.List(c.Request.Context(), middleware.Actor(c), limit, unreadOnly)
	if err != nil {
		response.Error(c, err)
		return
	}
	list := createList(items, count, h.Now())
	response.JSON(c, http.StatusOK, list)
}

// UnreadCount handles GET /api/v1/notifications/unread-count.
func (h *Handlers) UnreadCount(c *gin.Context) {
	count, err := h.Notifications.UnreadCount(c.Request.Context(), middleware.Actor(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.UnreadCountResponse{Count: count})
}

// MarkRead handles POST /api/v1/notifications/:id/read.
func (h *Handlers) MarkRead(c *gin.Context) {
	if err := h.Notifications.MarkRead(c.Request.Context(), middleware.Actor(c), c.Param("id")); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.OKMessage("Marked as read"))
}

// MarkAllRead handles POST /api/v1/notifications/read-all.
func (h *Handlers) MarkAllRead(c *gin.Context) {
	if err := h.Notifications.MarkAllRead(c.Request.Context(), middleware.Actor(c)); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.OKMessage("All notifications marked as read"))
}

// SendInApp handles POST /api/v1/notifications.
func (h *Handlers) SendInApp(c *gin.Context) {
	var request dto.SendRequest
	if err := response.DecodeJSON(c, &request); err != nil {
		response.Error(c, err)
		return
	}
	if request.Channel == "" {
		request.Channel = string(domain.ChannelInApp)
	}
	n, err := h.Notifications.Send(c.Request.Context(), request.ToInput())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.Notification(*n, h.Now()))
}

// SendEmail handles POST /api/v1/notifications/email.
func (h *Handlers) SendEmail(c *gin.Context) {
	h.sendDelivered(c, domain.ChannelEmail)
}

// SendPush handles POST /api/v1/notifications/push.
func (h *Handlers) SendPush(c *gin.Context) {
	h.sendDelivered(c, domain.ChannelPush)
}

func (h *Handlers) sendDelivered(c *gin.Context, channel domain.Channel) {
	var request dto.SendRequest
	if err := response.DecodeJSON(c, &request); err != nil {
		response.Error(c, err)
		return
	}
	request.Channel = string(channel)
	n, err := h.Notifications.Send(c.Request.Context(), request.ToInput())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.Notification(*n, h.Now()))
}

// ListTemplates handles GET /api/v1/notifications/templates.
func (h *Handlers) ListTemplates(c *gin.Context) {
	limit, err := queryInt(c, "limit", 50, 200)
	if err != nil {
		response.Error(c, err)
		return
	}
	templates, err := h.Notifications.ListTemplates(c.Request.Context(), limit)
	if err != nil {
		response.Error(c, err)
		return
	}
	items := make([]dto.TemplateResponse, 0, len(templates))
	for _, t := range templates {
		items = append(items, dto.Template(t))
	}
	response.JSON(c, http.StatusOK, items)
}

// CreateTemplate handles POST /api/v1/notifications/templates.
func (h *Handlers) CreateTemplate(c *gin.Context) {
	var request dto.CreateTemplateRequest
	if err := response.DecodeJSON(c, &request); err != nil {
		response.Error(c, err)
		return
	}
	t, err := h.Notifications.CreateTemplate(c.Request.Context(), request.ToInput())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.Template(*t))
}

func createList(items []domain.Notification, count int, now time.Time) dto.NotificationListResponse {
	list := make([]dto.NotificationResponse, 0, len(items))
	for _, n := range items {
		list = append(list, dto.Notification(n, now))
	}
	return dto.NotificationListResponse{Items: list, UnreadCount: count}
}

func queryInt(c *gin.Context, key string, fallback, max int) (int, error) {
	raw := c.Query(key)
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return 0, domain.Validationf("%s must be a positive integer", key)
	}
	if max > 0 && value > max {
		value = max
	}
	return value, nil
}
