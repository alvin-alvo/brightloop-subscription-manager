package repository

import (
	"context"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/domain"
)

// NotificationRepository persists in-app notifications and templates.
type NotificationRepository interface {
	ListByUser(ctx context.Context, userID string, limit int, unreadOnly bool) ([]domain.Notification, error)
	UnreadCount(ctx context.Context, userID string) (int, error)
	Create(ctx context.Context, notification *domain.Notification) error
	MarkRead(ctx context.Context, id, userID string) error
	MarkAllRead(ctx context.Context, userID string) error
	ListTemplates(ctx context.Context, limit int) ([]domain.NotificationTemplate, error)
	CreateTemplate(ctx context.Context, template *domain.NotificationTemplate) error
}
