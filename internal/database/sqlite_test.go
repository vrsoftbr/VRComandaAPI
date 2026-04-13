package database

import (
	"os"
	"path/filepath"
	"testing"

	"vrcomandaapi/internal/config"
)

func TestEnsureSQLiteDirNoopForCurrentDir(t *testing.T) {
	if err := ensureSQLiteDir("db.sqlite"); err != nil {
		t.Fatalf("ensureSQLiteDir returned error: %v", err)
	}
}

func TestEnsureSQLiteDirCreatesNestedPath(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "a", "b", "db.sqlite")

	if err := ensureSQLiteDir(path); err != nil {
		t.Fatalf("ensureSQLiteDir returned error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(root, "a", "b")); err != nil {
		t.Fatalf("expected nested directory to exist: %v", err)
	}
}

func TestEnsureSQLiteDirFailsWhenParentIsFile(t *testing.T) {
	root := t.TempDir()
	filePath := filepath.Join(root, "file-parent")
	if err := os.WriteFile(filePath, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed creating fixture file: %v", err)
	}

	err := ensureSQLiteDir(filepath.Join(filePath, "db.sqlite"))
	if err == nil {
		t.Fatal("expected error when parent path is file")
	}
}

func TestConnectSQLiteEnsureDirError(t *testing.T) {
	root := t.TempDir()
	blocker := filepath.Join(root, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed creating fixture file: %v", err)
	}

	cfg := config.Config{SQLitePath: filepath.Join(blocker, "db.sqlite")}
	db, err := ConnectSQLite(cfg)
	if err == nil {
		t.Fatal("expected ConnectSQLite to return error when dir cannot be created")
	}
	if db != nil {
		t.Fatal("expected nil db on error")
	}
}

func TestConnectSQLiteSuccess(t *testing.T) {
	cfg := config.Config{SQLitePath: ":memory:"}

	db, err := ConnectSQLite(cfg)
	if err != nil {
		t.Fatalf("ConnectSQLite returned error: %v", err)
	}

	if db == nil {
		t.Fatal("ConnectSQLite returned nil db")
	}
}
