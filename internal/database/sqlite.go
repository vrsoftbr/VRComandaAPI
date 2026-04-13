package database

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"vrcomandaapi/internal/config"
)

// ConnectSQLite opens the local SQLite database used for operational writes.
// The directory is created beforehand to avoid runtime failures on first boot.
func ConnectSQLite(cfg config.Config) (*gorm.DB, error) {
	if err := ensureSQLiteDir(cfg.SQLitePath); err != nil {
		return nil, err
	}

	return gorm.Open(sqlite.Open(cfg.SQLitePath), &gorm.Config{})
}

// ensureSQLiteDir guarantees the database parent directory exists.
func ensureSQLiteDir(sqlitePath string) error {
	dir := filepath.Dir(sqlitePath)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}
