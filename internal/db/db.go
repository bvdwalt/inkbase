package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// legacyDBName is the default database filename used before the project was
// renamed from Inkbase to Palimpsest.
const legacyDBName = "inkbase.db"

// MigrateLegacyPath renames a database file left over from the pre-rename
// default path (inkbase.db) to path, so upgrading in place doesn't strand
// existing data under the old name. It's a no-op if path already exists or
// no legacy file is found alongside it.
func MigrateLegacyPath(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	legacyPath := filepath.Join(filepath.Dir(path), legacyDBName)
	if _, err := os.Stat(legacyPath); err != nil {
		return nil
	}

	for _, suffix := range []string{"", "-shm", "-wal"} {
		src := legacyPath + suffix
		if _, err := os.Stat(src); err != nil {
			continue
		}
		dst := path + suffix
		if err := os.Rename(src, dst); err != nil {
			return fmt.Errorf("rename legacy db file %s: %w", src, err)
		}
	}

	slog.Info("migrated legacy database file", "from", legacyPath, "to", path)
	return nil
}

func Connect(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	if _, err := db.Exec(`PRAGMA journal_mode = WAL;`); err != nil {
		return nil, fmt.Errorf("enable WAL: %w", err)
	}

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}
