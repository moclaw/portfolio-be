package config

import (
	"os"
)

type Config struct {
	Port        string
	Host        string
	DatabaseURL string
	S3Config    S3Config
	JWTConfig   JWTConfig
}

type S3Config struct {
	Endpoint        string
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	ForcePathStyle  bool
}

type JWTConfig struct {
	SecretKey string
	Issuer    string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "5303"),
		Host:        getEnv("HOST", "localhost"),
		DatabaseURL: getEnv("DATABASE_URL", "portfolio.db"),
		S3Config: S3Config{
			Endpoint:        getEnv("S3_ENDPOINT", "http://localhost:4566"),
			Region:          getEnv("S3_REGION", "us-east-1"),
			Bucket:          getEnv("S3_BUCKET", "portfolio-bucket"),
			AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", "test"),
			SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", "test"),
			ForcePathStyle:  true,
		},
		JWTConfig: JWTConfig{
			SecretKey: getEnv("JWT_SECRET_KEY", "your-super-secret-jwt-key-change-in-production"),
			Issuer:    getEnv("JWT_ISSUER", "portfolio-api"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnv(key, defaultValue string) string {
	return getEnv(key, defaultValue)
}
