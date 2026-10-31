package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/domain"
)

// UserRepository is the persistence contract for users. Implementations live
// in the repository/postgres package.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	ByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	ByEmail(ctx context.Context, email string) (*domain.User, error)
}
