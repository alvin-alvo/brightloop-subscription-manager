// Package service implements the notification business rules.
package service

import (
	"context"
	"errors"
	"strings"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/domain"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/repository"
)

// NotificationService combines the notification repository.
type NotificationService struct {
	notifications repository.NotificationRepository
}

// NewNotificationService wires the notification business rules.
func NewNotificationService(notifications repository.NotificationRepository) *NotificationService {
	return &NotificationService{notifications: notifications}
}

// List returns a user's notifications and unread count.
func (s *NotificationService) List(ctx context.Context, actor domain.Principal, limit int, unreadOnly bool) ([]domain.Notification, int, error) {
	if actor.UserID == "" {
		return nil, 0, domain.ErrUnauthorized
	}
	items, err := s.notifications.ListByUser(ctx, actor.UserID, limit, unreadOnly)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.notifications.UnreadCount(ctx, actor.UserID)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}

// UnreadCount returns the unread total for a user.
func (s *NotificationService) UnreadCount(ctx context.Context, actor domain.Principal) (int, error) {
	if actor.UserID == "" {
		return 0, domain.ErrUnauthorized
	}
	return s.notifications.UnreadCount(ctx, actor.UserID)
}

// MarkRead marks one notification read.
func (s *NotificationService) MarkRead(ctx context.Context, actor domain.Principal, id string) error {
	if actor.UserID == "" {
		return domain.ErrUnauthorized
	}
	return s.notifications.MarkRead(ctx, id, actor.UserID)
}

// MarkAllRead marks every notification read.
func (s *NotificationService) MarkAllRead(ctx context.Context, actor domain.Principal) error {
	if actor.UserID == "" {
		return domain.ErrUnauthorized
	}
	return s.notifications.MarkAllRead(ctx, actor.UserID)
}

// Send creates an in-app notification for a user. Used by the in-app, email and
// push endpoints; delivery integrations can be layered on later.
func (s *NotificationService) Send(ctx context.Context, input domain.CreateNotificationInput) (*domain.Notification, error) {
	input.Sanitize()
	if err := validateNotification(input); err != nil {
		return nil, err
	}
	n := &domain.Notification{
		UserID:  input.UserID,
		Channel: input.Channel,
		Title:   input.Title,
		Body:    input.Body,
	}
	if err := s.notifications.Create(ctx, n); err != nil {
		return nil, err
	}
	return n, nil
}

// CreateTemplate stores a reusable template.
func (s *NotificationService) CreateTemplate(ctx context.Context, input domain.CreateTemplateInput) (*domain.NotificationTemplate, error) {
	input.Sanitize()
	if input.Code == "" {
		return nil, domain.Validationf("code is required")
	}
	if input.Channel == "" {
		input.Channel = domain.ChannelInApp
	}
	if !input.Channel.Valid() {
		return nil, domain.Validationf("channel must be one of: email, push, in_app")
	}
	if input.Body == "" {
		return nil, domain.Validationf("body is required")
	}
	t := &domain.NotificationTemplate{
		Code:    input.Code,
		Channel: input.Channel,
		Subject: input.Subject,
		Body:    input.Body,
	}
	if err := s.notifications.CreateTemplate(ctx, t); err != nil {
		if isConflict(err) {
			return nil, domain.Conflictf("a template with code %q already exists", input.Code)
		}
		return nil, err
	}
	return t, nil
}

// ListTemplates returns stored templates.
func (s *NotificationService) ListTemplates(ctx context.Context, limit int) ([]domain.NotificationTemplate, error) {
	return s.notifications.ListTemplates(ctx, limit)
}

func validateNotification(input domain.CreateNotificationInput) error {
	if input.UserID == "" {
		return domain.Validationf("userId is required")
	}
	if input.Channel == "" {
		return domain.Validationf("channel is required")
	}
	if !input.Channel.Valid() {
		return domain.Validationf("channel must be one of: email, push, in_app")
	}
	if strings.TrimSpace(input.Title) == "" {
		return domain.Validationf("title is required")
	}
	if strings.TrimSpace(input.Body) == "" {
		return domain.Validationf("body is required")
	}
	return nil
}

func isConflict(err error) bool {
	return errors.Is(err, domain.ErrConflict)
}
