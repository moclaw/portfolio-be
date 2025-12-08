package database

import (
	"fmt"
	"log"
	"portfolio-be/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitPostgres initializes PostgreSQL database connection
func InitPostgres(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("✓ Successfully connected to PostgreSQL database")
	return db, nil
}

// InitDatabase initializes PostgreSQL database connection
func InitDatabase(dbType, databaseURL string) (*gorm.DB, error) {
	return InitPostgres(databaseURL)
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
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
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	log.Println("✓ Database migrations completed successfully")
	return nil
}

// IsEmpty checks if the database has been seeded with initial data
func IsEmpty(db *gorm.DB) bool {
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	return userCount == 0
}
