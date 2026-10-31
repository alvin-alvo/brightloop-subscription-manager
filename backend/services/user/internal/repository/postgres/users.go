package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/domain"
)

const userColumns = `
	u.id::text,
	u.name,
	u.email,
	u.password_hash,
	u.role::text,
	u.status::text,
	COALESCE(u.company_name, ''),
	COALESCE(u.department, ''),
	COALESCE(u.phone, ''),
	COALESCE(u.avatar_url, ''),
	COALESCE(u.preferences::text, '{}'),
	u.email_verified_at,
	u.last_login_at,
	u.created_at,
	u.updated_at`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanUser(row rowScanner) (domain.User, error) {
	var (
		user          domain.User
		prefsJSON     string
		emailVerified sql.NullTime
		lastLogin     sql.NullTime
	)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Status,
		&user.CompanyName,
		&user.Department,
		&user.Phone,
		&user.AvatarURL,
		&prefsJSON,
		&emailVerified,
		&lastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, err
	}
	user.EmailVerifiedAt = timePtr(emailVerified)
	user.LastLoginAt = timePtr(lastLogin)
	user.Preferences = map[string]any{}
	if prefsJSON != "" && prefsJSON != "{}" {
		_ = json.Unmarshal([]byte(prefsJSON), &user.Preferences)
	}
	return user, nil
}

func (s *Store) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := "SELECT " + userColumns + " FROM users u WHERE u.id = $1::uuid"
	user, err := scanUser(s.db.QueryRowContext(ctx, query, id))
	if err != nil {
		return nil, mapNotFound(err)
	}
	return &user, nil
}

func (s *Store) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := "SELECT " + userColumns + " FROM users u WHERE lower(u.email) = lower($1)"
	user, err := scanUser(s.db.QueryRowContext(ctx, query, strings.TrimSpace(email)))
	if err != nil {
		return nil, mapNotFound(err)
	}
	return &user, nil
}

func (s *Store) ListStaff(ctx context.Context, search string, limit int) ([]domain.User, error) {
	if limit < 1 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	query := "SELECT " + userColumns + `
		FROM users u
		WHERE u.role IN ('admin', 'staff')`
	args := []any{}
	if strings.TrimSpace(search) != "" {
		args = append(args, likePattern(search))
		query += fmt.Sprintf(" AND (lower(u.name) LIKE $%d ESCAPE '\\' OR lower(u.email) LIKE $%d ESCAPE '\\' OR lower(COALESCE(u.department, '')) LIKE $%d ESCAPE '\\')", 1, 1, 1)
	}
	query += fmt.Sprintf(" ORDER BY u.name ASC LIMIT $%d", len(args)+1)
	args = append(args, limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list staff: %w", err)
	}
	defer rows.Close()

	items := make([]domain.User, 0)
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("scan staff: %w", err)
		}
		items = append(items, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate staff: %w", err)
	}
	return items, nil
}

func (s *Store) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, role, status, company_name, department)
		VALUES ($1, lower($2), $3, $4::user_role, $5::user_status, NULLIF($6, ''), $7)
		RETURNING id::text, created_at, updated_at`
	err := s.db.QueryRowContext(ctx, query,
		user.Name,
		user.Email,
		user.PasswordHash,
		string(user.Role),
		string(user.Status),
		user.CompanyName,
		user.Department,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (s *Store) UpdateProfile(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		UPDATE users SET
			name = $1,
			company_name = NULLIF($2, ''),
			phone = NULLIF($3, ''),
			avatar_url = NULLIF($4, ''),
			updated_at = now()
		WHERE id = $5::uuid
		RETURNING updated_at`
	err := s.db.QueryRowContext(ctx, query,
		user.Name,
		user.CompanyName,
		user.Phone,
		user.AvatarURL,
		user.ID,
	).Scan(&user.UpdatedAt)
	if err != nil {
		return nil, mapNotFound(err)
	}
	return s.GetByID(ctx, user.ID)
}

func (s *Store) UpdateAdmin(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		UPDATE users SET
			name = COALESCE(NULLIF($1, ''), name),
			role = COALESCE(NULLIF($2, '')::user_role, role),
			status = COALESCE(NULLIF($3, '')::user_status, status),
			department = COALESCE(NULLIF($4, ''), department),
			updated_at = now()
		WHERE id = $5::uuid
		RETURNING updated_at`
	err := s.db.QueryRowContext(ctx, query,
		user.Name,
		user.Role,
		user.Status,
		user.Department,
		user.ID,
	).Scan(&user.UpdatedAt)
	if err != nil {
		return nil, mapNotFound(err)
	}
	return s.GetByID(ctx, user.ID)
}

func (s *Store) SetStatus(ctx context.Context, id string, status domain.UserStatus, _ bool) error {
	result, err := s.db.ExecContext(ctx,
		"UPDATE users SET status = $2::user_status, updated_at = now() WHERE id = $1::uuid",
		id, string(status))
	if err != nil {
		return fmt.Errorf("set user status: %w", err)
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func likePattern(search string) string {
	search = strings.TrimSpace(strings.ToLower(search))
	search = strings.ReplaceAll(search, "\\", "\\\\")
	search = strings.ReplaceAll(search, "%", "\\%")
	search = strings.ReplaceAll(search, "_", "\\_")
	return "%" + search + "%"
}

func timePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}
