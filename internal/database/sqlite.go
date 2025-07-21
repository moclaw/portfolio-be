package database

import (
	"portfolio-be/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

func InitSQLite(databaseURL string) (*gorm.DB, error) {
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

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Content{},
		&models.Upload{},
		&models.Resource{},
		&models.Experience{},
		&models.Service{},
		&models.Technology{},
		&models.Project{},
		&models.Testimonial{},
		&models.Contact{},
	)
}
