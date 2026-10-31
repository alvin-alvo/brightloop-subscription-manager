package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/domain"
)

// fakeStore implements the user repository for unit tests.
type fakeStore struct {
	users   map[string]*domain.User // by id
	byEmail map[string]*domain.User
}

func newFakeStore() *fakeStore {
	return &fakeStore{users: map[string]*domain.User{}, byEmail: map[string]*domain.User{}}
}

var idCounter int

func nextID(prefix string) string {
	idCounter++
	return fmt.Sprintf("%s_%03d", prefix, idCounter)
}

func (f *fakeStore) seed(u domain.User) {
	u.ID = nextID("usr")
	f.users[u.ID] = &u
	f.byEmail[strings.ToLower(u.Email)] = &u
}

func (f *fakeStore) Clone(u domain.User) *domain.User {
	c := u
	if u.LastLoginAt != nil {
		v := *u.LastLoginAt
		c.LastLoginAt = &v
	}
	return &c
}

func (f *fakeStore) GetByID(_ context.Context, id string) (*domain.User, error) {
	u, ok := f.users[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return f.Clone(*u), nil
}

func (f *fakeStore) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	u, ok := f.byEmail[strings.ToLower(email)]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return f.Clone(*u), nil
}

func (f *fakeStore) ListStaff(_ context.Context, _ string, _ int) ([]domain.User, error) {
	var items []domain.User
	for _, u := range f.users {
		if u.Role.StaffRole() {
			items = append(items, *f.Clone(*u))
		}
	}
	return items, nil
}

func (f *fakeStore) Create(_ context.Context, u *domain.User) error {
	if _, ok := f.byEmail[strings.ToLower(u.Email)]; ok {
		return domain.ErrConflict
	}
	c := *u
	c.ID = nextID("usr")
	c.CreatedAt = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	c.UpdatedAt = c.CreatedAt
	f.users[c.ID] = &c
	f.byEmail[strings.ToLower(c.Email)] = &c
	*u = c
	return nil
}

func (f *fakeStore) UpdateProfile(_ context.Context, u *domain.User) (*domain.User, error) {
	return f.update(u, false)
}

func (f *fakeStore) UpdateAdmin(_ context.Context, u *domain.User) (*domain.User, error) {
	return f.update(u, true)
}

func (f *fakeStore) update(u *domain.User, _ bool) (*domain.User, error) {
	cur, ok := f.users[u.ID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	*cur = *u
	return f.Clone(*cur), nil
}

func (f *fakeStore) SetStatus(_ context.Context, id string, status domain.UserStatus, _ bool) error {
	u, ok := f.users[id]
	if !ok {
		return domain.ErrNotFound
	}
	u.Status = status
	return nil
}

func (f *fakeStore) UpdatePreferences(_ context.Context, id string, prefs map[string]any) error {
	u, ok := f.users[id]
	if !ok {
		return domain.ErrNotFound
	}
	u.Preferences = prefs
	return nil
}

var _ = errors.Is
