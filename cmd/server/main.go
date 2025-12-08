// Package main Portfolio Backend API
// @title Portfolio Backend API
// @version 1.0
// @description This is a portfolio backend server with content management and file upload capabilities.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3001
// @BasePath /
// @schemes http https
package main

import (
	"log"
	"net/http"

	"portfolio-be/internal/api"
	"portfolio-be/internal/config"
	"portfolio-be/internal/database"
	"portfolio-be/internal/services"

	"portfolio-be/docs" // This is required for swagger
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Update swagger info with dynamic host
	docs.SwaggerInfo.Host = cfg.Host + ":" + cfg.Port

	// Initialize database (supports both PostgreSQL and SQLite)
	db, err := database.InitDatabase(cfg.DatabaseType, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	log.Printf("Database initialized with type: %s", cfg.DatabaseType)

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Check if database needs seeding and seed if empty
	if database.IsEmpty(db) {
		// Seed database with initial data
		if err := database.Seed(db); err != nil {
			log.Fatal("Failed to seed database:", err)
		}
		log.Println("Database seeded successfully!")
	}

	// Initialize Redis if enabled
	if cfg.RedisConfig.Enabled {
		redisClient, err := database.NewRedisClient(database.RedisConfig{
			Host:     cfg.RedisConfig.Host,
			Port:     cfg.RedisConfig.Port,
			Password: cfg.RedisConfig.Password,
			DB:       cfg.RedisConfig.DB,
		})
		if err != nil {
			log.Printf("Warning: Failed to initialize Redis: %v", err)
		} else {
			defer redisClient.Close()
			log.Println("Redis connected successfully!")
		}
	}

	// Initialize S3 service
	s3Service, err := services.NewS3Service(cfg.S3Config)
	if err != nil {
		log.Fatal("Failed to initialize S3 service:", err)
	}

	// Setup router
	router := api.SetupRouter(db, s3Service, cfg)

	// Start server
	address := cfg.Host + ":" + cfg.Port
	log.Printf("Server starting on http://%s", address)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
