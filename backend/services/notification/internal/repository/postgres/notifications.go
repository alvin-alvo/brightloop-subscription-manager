package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

const notificationColumns = `
	n.id::text,
	n.user_id::text,
	n.channel::text,
	n.title,
	n.body,
	n.read_at,
	n.created_at`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanNotification(row rowScanner) (domain.Notification, error) {
	var n domain.Notification
	var readAt sql.NullTime
	err := row.Scan(
		&n.ID,
		&n.UserID,
		&n.Channel,
		&n.Title,
		&n.Body,
		&readAt,
		&n.CreatedAt,
	)
	if err != nil {
		return domain.Notification{}, err
	}
	n.ReadAt = timePtr(readAt)
	return n, nil
}

func (s *Store) ListByUser(ctx context.Context, userID string, limit int, unreadOnly bool) ([]domain.Notification, error) {
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	query := "SELECT " + notificationColumns + `
		FROM notifications n
		WHERE n.user_id = $1::uuid`
	args := []any{userID}
	if unreadOnly {
		query += " AND n.read_at IS NULL"
	}
	query += fmt.Sprintf(" ORDER BY n.created_at DESC LIMIT $%d", len(args)+1)
	args = append(args, limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()

	items := make([]domain.Notification, 0, limit)
	for rows.Next() {
		n, err := scanNotification(rows)
		if err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		items = append(items, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate notifications: %w", err)
	}
	return items, nil
}

func (s *Store) UnreadCount(ctx context.Context, userID string) (int, error) {
	var count int
	if err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*)::int FROM notifications WHERE user_id = $1::uuid AND read_at IS NULL",
		userID).Scan(&count); err != nil {
		return 0, fmt.Errorf("count unread notifications: %w", err)
	}
	return count, nil
}

func (s *Store) Create(ctx context.Context, n *domain.Notification) error {
	query := `
		INSERT INTO notifications (user_id, channel, title, body)
		VALUES ($1::uuid, $2::notification_channel, $3, $4)
		RETURNING id::text, created_at`
	err := s.db.QueryRowContext(ctx, query,
		n.UserID,
		string(n.Channel),
		n.Title,
		n.Body,
	).Scan(&n.ID, &n.CreatedAt)
	if err != nil {
		return fmt.Errorf("create notification: %w", err)
	}
	return nil
}

func (s *Store) MarkRead(ctx context.Context, id, userID string) error {
	result, err := s.db.ExecContext(ctx,
		"UPDATE notifications SET read_at = now() WHERE id = $1::uuid AND user_id = $2::uuid AND read_at IS NULL",
		id, userID)
	if err != nil {
		return fmt.Errorf("mark notification read: %w", err)
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) MarkAllRead(ctx context.Context, userID string) error {
	if _, err := s.db.ExecContext(ctx,
		"UPDATE notifications SET read_at = now() WHERE user_id = $1::uuid AND read_at IS NULL",
		userID); err != nil {
		return fmt.Errorf("mark all notifications read: %w", err)
	}
	return nil
}

func (s *Store) ListTemplates(ctx context.Context, limit int) ([]domain.NotificationTemplate, error) {
	if limit < 1 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id::text, code, channel::text, COALESCE(subject, ''), body, created_at, updated_at
		FROM notification_templates
		ORDER BY code ASC
		LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("list templates: %w", err)
	}
	defer rows.Close()

	items := make([]domain.NotificationTemplate, 0)
	for rows.Next() {
		var t domain.NotificationTemplate
		if err := rows.Scan(&t.ID, &t.Code, &t.Channel, &t.Subject, &t.Body, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan template: %w", err)
		}
		items = append(items, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate templates: %w", err)
	}
	return items, nil
}

func (s *Store) CreateTemplate(ctx context.Context, t *domain.NotificationTemplate) error {
	query := `
		INSERT INTO notification_templates (code, channel, subject, body)
		VALUES ($1, $2::notification_channel, NULLIF($3, ''), $4)
		RETURNING id::text, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query, t.Code, string(t.Channel), t.Subject, t.Body).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("create template: %w", err)
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	return false
}
