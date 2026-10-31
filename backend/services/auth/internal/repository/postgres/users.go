package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/domain"
)

const userColumns = `
	u.id::text,
	u.name,
	u.email,
	u.password_hash,
	u.role::text,
	u.status::text,
	COALESCE(u.company_name, ''),
	COALESCE(u.phone, ''),
	COALESCE(u.avatar_url, ''),
	u.email_verified_at,
	u.last_login_at,
	u.created_at,
	u.updated_at`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanUser(row rowScanner) (domain.User, error) {
	var user domain.User
	var emailVerifiedAt, lastLoginAt sql.NullTime
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&user.CompanyName,
		&user.Phone,
		&user.AvatarURL,
		&emailVerifiedAt,
		&lastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, err
	}
	user.EmailVerifiedAt = timePtr(emailVerifiedAt)
	user.LastLoginAt = timePtr(lastLoginAt)
	return user, nil
}

func (s *Store) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, role, status, company_name)
		VALUES ($1, lower($2), $3, $4::user_role, $5::user_status, NULLIF($6, ''))
		RETURNING id::text, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query,
		user.Name,
		strings.TrimSpace(user.Email),
		user.PasswordHash,
		string(user.Role),
		string(user.Status),
		user.CompanyName,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (s *Store) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := "SELECT " + userColumns + " FROM users u WHERE lower(u.email) = lower($1)"
	user, err := scanUser(s.db.QueryRowContext(ctx, query, strings.TrimSpace(email)))
	if err != nil {
		return nil, mapNotFound(err)
	}
	return &user, nil
}

func (s *Store) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := "SELECT " + userColumns + " FROM users u WHERE u.id = $1::uuid"
	user, err := scanUser(s.db.QueryRowContext(ctx, query, id))
	if err != nil {
		return nil, mapNotFound(err)
	}
	return &user, nil
}

func (s *Store) TouchLogin(ctx context.Context, id string, at time.Time) error {
	if _, err := s.db.ExecContext(ctx,
		"UPDATE users SET last_login_at = $2, updated_at = now() WHERE id = $1::uuid",
		id, at); err != nil {
		return fmt.Errorf("update last login: %w", err)
	}
	return nil
}

func (s *Store) UpdatePassword(ctx context.Context, id string, hash string) error {
	result, err := s.db.ExecContext(ctx,
		"UPDATE users SET password_hash = $2, updated_at = now() WHERE id = $1::uuid",
		id, hash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) CreatePasswordResetToken(ctx context.Context, token *domain.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
		VALUES ($1::uuid, $2, $3)
		RETURNING id::text`
	err := s.db.QueryRowContext(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt).Scan(&token.ID)
	if err != nil {
		return fmt.Errorf("create password reset token: %w", err)
	}
	return nil
}
