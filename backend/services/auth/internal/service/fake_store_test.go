package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/domain"
)

// fakeStore implements every repository interface for unit tests and records
// calls so tests can assert side effects.
type fakeStore struct {
	users    map[string]*domain.User // keyed by lower(email)
	byID     map[string]*domain.User
	sessions map[string]*domain.Session
	resets   []*domain.PasswordResetToken
	created  []*domain.User
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		users:    map[string]*domain.User{},
		byID:     map[string]*domain.User{},
		sessions: map[string]*domain.Session{},
		resets:   []*domain.PasswordResetToken{},
		created:  []*domain.User{},
	}
}

var fakeIDCounter int

func nextID() string {
	fakeIDCounter++
	return fmt.Sprintf("usr_%03d", fakeIDCounter)
}

func cloneUser(u domain.User) *domain.User {
	c := u
	if u.EmailVerifiedAt != nil {
		v := *u.EmailVerifiedAt
		c.EmailVerifiedAt = &v
	}
	if u.LastLoginAt != nil {
		v := *u.LastLoginAt
		c.LastLoginAt = &v
	}
	return &c
}

func (f *fakeStore) seed(user domain.User) {
	user.ID = nextID()
	f.users[strings.ToLower(user.Email)] = &user
	f.byID[user.ID] = &user
}

func (f *fakeStore) Create(_ context.Context, user *domain.User) error {
	email := strings.ToLower(user.Email)
	if _, ok := f.users[email]; ok {
		return domain.ErrConflict
	}
	copy := *user
	copy.ID = nextID()
	copy.CreatedAt = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	copy.UpdatedAt = copy.CreatedAt
	f.users[email] = &copy
	f.byID[copy.ID] = &copy
	f.created = append(f.created, cloneUser(copy))
	*user = copy
	return nil
}

func (f *fakeStore) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	user, ok := f.users[strings.ToLower(email)]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return cloneUser(*user), nil
}

func (f *fakeStore) GetByID(_ context.Context, id string) (*domain.User, error) {
	user, ok := f.byID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return cloneUser(*user), nil
}

func (f *fakeStore) TouchLogin(_ context.Context, id string, at time.Time) error {
	user, ok := f.byID[id]
	if !ok {
		return domain.ErrNotFound
	}
	value := at
	user.LastLoginAt = &value
	return nil
}

func (f *fakeStore) UpdatePassword(_ context.Context, id string, hash string) error {
	user, ok := f.byID[id]
	if !ok {
		return domain.ErrNotFound
	}
	user.PasswordHash = hash
	return nil
}

func (f *fakeStore) CreatePasswordResetToken(_ context.Context, token *domain.PasswordResetToken) error {
	token.ID = nextID()
	f.resets = append(f.resets, token)
	return nil
}

func (f *fakeStore) CreateSession(_ context.Context, session *domain.Session) error {
	copy := *session
	copy.ID = nextID()
	f.sessions[copy.ID] = &copy
	*session = copy
	return nil
}

func (f *fakeStore) GetSessionByID(_ context.Context, id string) (*domain.Session, error) {
	session, ok := f.sessions[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	copy := *session
	return &copy, nil
}

func (f *fakeStore) RevokeSession(_ context.Context, id string, at time.Time) error {
	session, ok := f.sessions[id]
	if !ok {
		return domain.ErrNotFound
	}
	value := at
	session.RevokedAt = &value
	return nil
}
