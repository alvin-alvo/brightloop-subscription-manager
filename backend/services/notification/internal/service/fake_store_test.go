package service

import (
	"context"
	"fmt"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/domain"
)

// fakeStore implements the notification repository for unit tests.
type fakeStore struct {
	items     []*domain.Notification
	templates []*domain.NotificationTemplate
	seq       int
}

func newFakeStore() *fakeStore {
	return &fakeStore{}
}

func itoaID(n int) string {
	return fmt.Sprintf("n_%03d", n)
}

func (f *fakeStore) ListByUser(_ context.Context, userID string, limit int, unreadOnly bool) ([]domain.Notification, error) {
	var out []domain.Notification
	count := 0
	for _, n := range f.items {
		if n.UserID != userID {
			continue
		}
		if unreadOnly && n.Read() {
			continue
		}
		if count >= limit {
			break
		}
		count++
		out = append(out, *n)
	}
	return out, nil
}

func (f *fakeStore) UnreadCount(_ context.Context, userID string) (int, error) {
	count := 0
	for _, n := range f.items {
		if n.UserID == userID && !n.Read() {
			count++
		}
	}
	return count, nil
}

func (f *fakeStore) Create(_ context.Context, n *domain.Notification) error {
	f.seq++
	now := time.Now()
	copy := *n
	copy.ID = itoaID(f.seq)
	copy.CreatedAt = now
	f.items = append(f.items, &copy)
	*n = copy
	return nil
}

func (f *fakeStore) MarkRead(_ context.Context, id, userID string) error {
	for _, n := range f.items {
		if n.ID == id && n.UserID == userID {
			if n.Read() {
				return domain.ErrNotFound
			}
			at := time.Now()
			n.ReadAt = &at
			return nil
		}
	}
	return domain.ErrNotFound
}

func (f *fakeStore) MarkAllRead(_ context.Context, userID string) error {
	for _, n := range f.items {
		if n.UserID == userID && !n.Read() {
			at := time.Now()
			n.ReadAt = &at
		}
	}
	return nil
}

func (f *fakeStore) ListTemplates(_ context.Context, _ int) ([]domain.NotificationTemplate, error) {
	out := make([]domain.NotificationTemplate, 0, len(f.templates))
	for _, t := range f.templates {
		out = append(out, *t)
	}
	return out, nil
}

func (f *fakeStore) CreateTemplate(_ context.Context, t *domain.NotificationTemplate) error {
	for _, existing := range f.templates {
		if existing.Code == t.Code {
			return domain.ErrConflict
		}
	}
	f.seq++
	copy := *t
	copy.ID = itoaID(f.seq)
	f.templates = append(f.templates, &copy)
	*t = copy
	return nil
}
