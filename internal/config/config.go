package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretsManagerService handles AWS Secrets Manager operations
type SecretsManagerService struct {
	client *secretsmanager.Client
	region string
}

// NewSecretsManagerService creates a new Secrets Manager service
func NewSecretsManagerService(region string) (*SecretsManagerService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := secretsmanager.NewFromConfig(cfg)

	return &SecretsManagerService{
		client: client,
		region: region,
	}, nil
}

// GetSecret retrieves a secret from AWS Secrets Manager
func (s *SecretsManagerService) GetSecret(secretName string) (*SecretData, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := s.client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve secret %s: %w", secretName, err)
	}

	var secretData SecretData
	err = json.Unmarshal([]byte(*result.SecretString), &secretData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal secret data: %w", err)
	}

	log.Printf("✓ Successfully retrieved secret: %s", secretName)
	return &secretData, nil
}

func Load() *Config {
	// Check if we should use Secrets Manager
	useSecrets := getEnv("USE_SECRETS_MANAGER", "false") == "true"
	secretName := getEnv("SECRET_NAME", "portfolio-secrets")
	region := getEnv("AWS_REGION", "us-east-1")

	var secretData *SecretData

	if useSecrets {
		log.Printf("Loading configuration from AWS Secrets Manager...")
		secretsService, err := NewSecretsManagerService(region)
		if err != nil {
			log.Printf("Failed to initialize Secrets Manager service: %v", err)
			log.Printf("Falling back to environment variables...")
		} else {
			secretData, err = secretsService.GetSecret(secretName)
			if err != nil {
				log.Printf("Failed to retrieve secrets from Secrets Manager: %v", err)
				log.Printf("Falling back to environment variables...")
			} else {
				log.Printf("Successfully loaded configuration from Secrets Manager")
			}
		}
	}

	// Load config with fallback to environment variables
	config := &Config{
		Port:        getEnv("PORT", "5303"),
		Host:        getEnv("HOST", "localhost"),
		DatabaseURL: getSecretOrEnv(secretData, "database_url", "DATABASE_URL", "portfolio.db"),
		S3Config: S3Config{
			Endpoint:        getSecretOrEnv(secretData, "s3_endpoint", "S3_ENDPOINT", "change-in-production"),
			Region:          getSecretOrEnv(secretData, "s3_region", "S3_REGION", "us-east-1"),
			Bucket:          getSecretOrEnv(secretData, "s3_bucket", "S3_BUCKET", "change-in-production"),
			AccessKeyID:     getSecretOrEnv(secretData, "s3_access_key_id", "S3_ACCESS_KEY_ID", "change-in-production"),
			SecretAccessKey: getSecretOrEnv(secretData, "s3_secret_access_key", "S3_SECRET_ACCESS_KEY", "change-in-production"),
			ForcePathStyle:  true,
		},
		JWTConfig: JWTConfig{
			SecretKey: getSecretOrEnv(secretData, "jwt_secret_key", "JWT_SECRET_KEY", "your-super-secret-jwt-key-change-in-production"),
			Issuer:    getEnv("JWT_ISSUER", "portfolio-api"),
		},
		SecretsManagerConfig: SecretsManagerConfig{
			SecretName: secretName,
			Region:     region,
			UseSecrets: useSecrets,
		},
	}

	// Validate critical S3 configuration
	if err := validateS3Config(config.S3Config); err != nil {
		log.Fatalf("Invalid S3 configuration: %v", err)
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getSecretOrEnv(secretData *SecretData, secretKey, envKey, defaultValue string) string {
	if secretData != nil {
		switch secretKey {
		case "database_url":
			if secretData.DatabaseURL != "" {
				return secretData.DatabaseURL
			}
		case "s3_endpoint":
			if secretData.S3Endpoint != "" {
				return secretData.S3Endpoint
			}
		case "s3_region":
			if secretData.S3Region != "" {
				return secretData.S3Region
			}
		case "s3_bucket":
			if secretData.S3Bucket != "" {
				return secretData.S3Bucket
			}
		case "s3_access_key_id":
			if secretData.S3AccessKeyID != "" {
				return secretData.S3AccessKeyID
			}
		case "s3_secret_access_key":
			if secretData.S3SecretAccessKey != "" {
				return secretData.S3SecretAccessKey
			}
		case "jwt_secret_key":
			if secretData.JWTSecretKey != "" {
				return secretData.JWTSecretKey
			}
		}
	}
	// Fallback to environment variable or default
	return getEnv(envKey, defaultValue)
}

func GetEnv(key, defaultValue string) string {
	return getEnv(key, defaultValue)
}

// validateS3Config validates that S3 configuration doesn't contain placeholder values
func validateS3Config(s3Config S3Config) error {
	placeholderValues := []string{"change-in-production", "your-super-secret-jwt-key-change-in-production"}

	for _, placeholder := range placeholderValues {
		if s3Config.Endpoint == placeholder {
			return fmt.Errorf("S3 Endpoint contains placeholder value: %s. Please set S3_ENDPOINT environment variable or configure AWS Secrets Manager", placeholder)
		}
		if s3Config.Bucket == placeholder {
			return fmt.Errorf("S3 Bucket contains placeholder value: %s. Please set S3_BUCKET environment variable or configure AWS Secrets Manager", placeholder)
		}
		if s3Config.AccessKeyID == placeholder {
			return fmt.Errorf("S3 AccessKeyID contains placeholder value: %s. Please set S3_ACCESS_KEY_ID environment variable or configure AWS Secrets Manager", placeholder)
		}
		if s3Config.SecretAccessKey == placeholder {
			return fmt.Errorf("S3 SecretAccessKey contains placeholder value: %s. Please set S3_SECRET_ACCESS_KEY environment variable or configure AWS Secrets Manager", placeholder)
		}
	}

	return nil
}
