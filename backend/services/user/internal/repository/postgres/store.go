package postgres

import (
	"database/sql"
	"errors"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

// Store implements the user repository on top of the shared PostgreSQL database.
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
