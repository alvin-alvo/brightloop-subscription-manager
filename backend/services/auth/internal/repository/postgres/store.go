package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/auth/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

// Store implements every auth repository interface on top of the shared
// PostgreSQL database.
type Store struct {
	db *sql.DB
}

// New returns a Store backed by the given database handle.
func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func mapNotFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	return false
}

func timePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}

func stringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	s := value.String
	return &s
}
