package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alvin-alvo/brightloop-subscription-manager/backend/services/user/internal/domain"
)

// UpdatePreferences replaces a user's stored JSON preferences.
func (s *Store) UpdatePreferences(ctx context.Context, id string, prefs map[string]any) error {
	if prefs == nil {
		prefs = map[string]any{}
	}
	encoded, err := json.Marshal(prefs)
	if err != nil {
		return fmt.Errorf("encode preferences: %w", err)
	}
	result, err := s.db.ExecContext(ctx,
		"UPDATE users SET preferences = $2::jsonb, updated_at = now() WHERE id = $1::uuid",
		id, string(encoded))
	if err != nil {
		return fmt.Errorf("update preferences: %w", err)
	}
	if affected, err := result.RowsAffected(); err == nil && affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
