package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/subscription/internal/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Store implements every repository interface on top of the shared
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

func likePattern(search string) string {
	search = strings.TrimSpace(strings.ToLower(search))
	search = strings.ReplaceAll(search, "\\", "\\\\")
	search = strings.ReplaceAll(search, "%", "\\%")
	search = strings.ReplaceAll(search, "_", "\\_")
	return "%" + search + "%"
}

func nullString(value string) sql.NullString {
	return sql.NullString{String: value, Valid: value != ""}
}

func nullTime(value *time.Time) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *value, Valid: true}
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

func placeholders(start, count int) string {
	parts := make([]string, 0, count)
	for i := 0; i < count; i++ {
		parts = append(parts, fmt.Sprintf("$%d", start+i))
	}
	return strings.Join(parts, ", ")
}
