package database

import (
	"portfolio-be/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitSQLite(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{
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
		&models.Experience{},
		&models.Service{},
		&models.Technology{},
		&models.Project{},
		&models.Testimonial{},
		&models.Contact{},
	)
}
