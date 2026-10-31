package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/domain"
)

func (s *Store) CreateSession(ctx context.Context, session *domain.Session) error {
	var ip any
	if session.IPAddress != "" {
		ip = session.IPAddress
	}
	query := `
		INSERT INTO auth_sessions (user_id, refresh_token_hash, user_agent, ip_address, expires_at)
		VALUES ($1::uuid, $2, NULLIF($3, ''), $4, $5)
		RETURNING id::text`
	err := s.db.QueryRowContext(ctx, query,
		session.UserID,
		session.RefreshTokenHash,
		session.UserAgent,
		ip,
		session.ExpiresAt,
	).Scan(&session.ID)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (s *Store) GetSessionByID(ctx context.Context, id string) (*domain.Session, error) {
	query := `
		SELECT id::text, user_id::text, refresh_token_hash,
		       COALESCE(user_agent, ''), COALESCE(host(ip_address), ''),
		       expires_at, revoked_at
		FROM auth_sessions WHERE id = $1::uuid`
	var (
		session   domain.Session
		revokedAt sql.NullTime
	)
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshTokenHash,
		&session.UserAgent,
		&session.IPAddress,
		&session.ExpiresAt,
		&revokedAt,
	)
	if err != nil {
		return nil, mapNotFound(err)
	}
	session.RevokedAt = timePtr(revokedAt)
	return &session, nil
}

func (s *Store) RevokeSession(ctx context.Context, id string, at time.Time) error {
	if _, err := s.db.ExecContext(ctx,
		"UPDATE auth_sessions SET revoked_at = $2 WHERE id = $1::uuid AND revoked_at IS NULL",
		id, at); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	return nil
}
