package postgres

import (
	"database/sql"
	"errors"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/notification/internal/domain"
)

// Store implements the notification repository on top of PostgreSQL.
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

func timePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}
