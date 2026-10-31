package repository

import (
	"context"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/domain"
)

// UserRepository persists accounts and password-reset tokens.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	TouchLogin(ctx context.Context, id string, at time.Time) error
	UpdatePassword(ctx context.Context, id string, hash string) error
	CreatePasswordResetToken(ctx context.Context, token *domain.PasswordResetToken) error
}

// SessionRepository persists issued authentication sessions so they can be
// revoked on logout and checked during validation.
type SessionRepository interface {
	CreateSession(ctx context.Context, session *domain.Session) error
	GetSessionByID(ctx context.Context, id string) (*domain.Session, error)
	RevokeSession(ctx context.Context, id string, at time.Time) error
}
