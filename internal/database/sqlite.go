package database

import (
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

func InitSQLite(databaseURL string) (*gorm.DB, error) {
	// Handle special SQLite database URLs
	if databaseURL != ":memory:" && databaseURL != "" {
		// Ensure the directory exists for the database file
		dir := filepath.Dir(databaseURL)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, err
			}
		}
	}

	// Use the pure Go SQLite driver by specifying the driver name
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        databaseURL,
	}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
