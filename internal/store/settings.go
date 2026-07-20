package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
)

const apiKeySettingKey = "api_key"

// GenerateAPIKey returns a random 64-character hex string suitable for use as an API key.
func GenerateAPIKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GetOrCreateAPIKey returns the current API key, generating and persisting one on first call.
func (s *Store) GetOrCreateAPIKey(ctx context.Context) (string, error) {
	var value string
	err := s.db.QueryRowContext(ctx, `SELECT value FROM app_settings WHERE key = ?`, apiKeySettingKey).Scan(&value)
	if err == nil {
		return value, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	key, err := GenerateAPIKey()
	if err != nil {
		return "", err
	}
	if _, err := s.db.ExecContext(ctx, `
		INSERT INTO app_settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO NOTHING
	`, apiKeySettingKey, key); err != nil {
		return "", err
	}

	// Someone else may have won a concurrent insert; re-read to get the authoritative value.
	if err := s.db.QueryRowContext(ctx, `SELECT value FROM app_settings WHERE key = ?`, apiKeySettingKey).Scan(&value); err != nil {
		return "", err
	}
	return value, nil
}

// RegenerateAPIKey generates a new API key, persists it, and returns it.
func (s *Store) RegenerateAPIKey(ctx context.Context) (string, error) {
	key, err := GenerateAPIKey()
	if err != nil {
		return "", err
	}
	if _, err := s.db.ExecContext(ctx, `
		INSERT INTO app_settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value
	`, apiKeySettingKey, key); err != nil {
		return "", err
	}
	return key, nil
}
