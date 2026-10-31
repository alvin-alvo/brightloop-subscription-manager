package repository

import (
	"context"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/domain"
)

// UserRepository persists accounts.
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	ListStaff(ctx context.Context, search string, limit int) ([]domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	UpdateProfile(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateAdmin(ctx context.Context, user *domain.User) (*domain.User, error)
	SetStatus(ctx context.Context, id string, status domain.UserStatus, deactivate bool) error
	UpdatePreferences(ctx context.Context, id string, prefs map[string]any) error
}
