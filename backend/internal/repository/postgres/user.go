package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/domain"
	"github.com/alvin-alvo/brightloop-subscription-manager/backend/internal/repository"
)

// UserRepository stores users in PostgreSQL.
type UserRepository struct {
	db *sql.DB
}

var _ repository.UserRepository = (*UserRepository)(nil)

// NewUserRepository builds a PostgreSQL-backed user repository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user. A duplicate email is reported as
// domain.ErrConflict.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	const query = `
		INSERT INTO users (id, email, password_hash, full_name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		user.ID.String(),
		user.Email,
		user.PasswordHash,
		user.FullName,
		string(user.Role),
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrConflict
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// ByID loads a user by primary key.
func (r *UserRepository) ByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const query = `
		SELECT id, email, password_hash, full_name, role, created_at, updated_at
		FROM users
		WHERE id = $1`

	return r.findOne(ctx, query, id.String())
}

// ByEmail loads a user by their unique email address.
func (r *UserRepository) ByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = `
		SELECT id, email, password_hash, full_name, role, created_at, updated_at
		FROM users
		WHERE email = $1`

	return r.findOne(ctx, query, email)
}

func (r *UserRepository) findOne(ctx context.Context, query string, arg string) (*domain.User, error) {
	var (
		user domain.User
		id   string
		role string
	)

	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&id,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	user.ID, err = uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}
	user.Role = domain.Role(role)

	return &user, nil
}
