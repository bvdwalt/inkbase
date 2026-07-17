package db_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bvdwalt/palimpsest/internal/db"
)

func TestConnectPingFailureOnDirectoryPath(t *testing.T) {
	_, err := db.Connect(t.TempDir())
	if err == nil {
		t.Fatal("Connect with a directory path: err = nil, want an error")
	}
}

func TestConnectAppliesMigrationsOnce(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	first, err := db.Connect(path)
	if err != nil {
		t.Fatalf("first Connect: %v", err)
	}
	first.Close()

	// Reconnecting to the same file should skip already-applied migrations without error.
	second, err := db.Connect(path)
	if err != nil {
		t.Fatalf("second Connect: %v", err)
	}
	defer second.Close()

	var count int
	if err := second.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&count); err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if count == 0 {
		t.Error("schema_migrations is empty after two Connect calls, want recorded migrations")
	}
}

func TestMigrateLegacyPathRenamesLegacyFileAndSidecars(t *testing.T) {
	dir := t.TempDir()
	legacyPath := filepath.Join(dir, "inkbase.db")
	newPath := filepath.Join(dir, "palimpsest.db")

	for _, suffix := range []string{"", "-shm", "-wal"} {
		if err := os.WriteFile(legacyPath+suffix, []byte("data"), 0o644); err != nil {
			t.Fatalf("write legacy file %s: %v", suffix, err)
		}
	}

	if err := db.MigrateLegacyPath(newPath); err != nil {
		t.Fatalf("MigrateLegacyPath: %v", err)
	}

	for _, suffix := range []string{"", "-shm", "-wal"} {
		if _, err := os.Stat(legacyPath + suffix); !os.IsNotExist(err) {
			t.Errorf("legacy file %s still exists after migration", legacyPath+suffix)
		}
		if _, err := os.Stat(newPath + suffix); err != nil {
			t.Errorf("expected %s to exist after migration: %v", newPath+suffix, err)
		}
	}
}

func TestMigrateLegacyPathNoOpWhenNewPathExists(t *testing.T) {
	dir := t.TempDir()
	legacyPath := filepath.Join(dir, "inkbase.db")
	newPath := filepath.Join(dir, "palimpsest.db")

	if err := os.WriteFile(legacyPath, []byte("legacy"), 0o644); err != nil {
		t.Fatalf("write legacy file: %v", err)
	}
	if err := os.WriteFile(newPath, []byte("current"), 0o644); err != nil {
		t.Fatalf("write new file: %v", err)
	}

	if err := db.MigrateLegacyPath(newPath); err != nil {
		t.Fatalf("MigrateLegacyPath: %v", err)
	}

	got, err := os.ReadFile(newPath)
	if err != nil {
		t.Fatalf("read new path: %v", err)
	}
	if string(got) != "current" {
		t.Errorf("newPath content = %q, want unchanged %q", got, "current")
	}
	if _, err := os.Stat(legacyPath); err != nil {
		t.Errorf("legacy file should be left untouched: %v", err)
	}
}

func TestMigrateLegacyPathNoOpWhenNoLegacyFile(t *testing.T) {
	dir := t.TempDir()
	newPath := filepath.Join(dir, "palimpsest.db")

	if err := db.MigrateLegacyPath(newPath); err != nil {
		t.Fatalf("MigrateLegacyPath: %v", err)
	}

	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Errorf("newPath should not have been created")
	}
}
