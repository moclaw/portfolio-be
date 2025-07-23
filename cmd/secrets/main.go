package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretData represents the structure of secrets stored in AWS Secrets Manager
type SecretData struct {
	DatabaseURL       string `json:"database_url"`
	JWTSecretKey      string `json:"jwt_secret_key"`
	S3Endpoint        string `json:"s3_endpoint"`
	S3Region          string `json:"s3_region"`
	S3Bucket          string `json:"s3_bucket"`
	S3AccessKeyID     string `json:"s3_access_key_id"`
	S3SecretAccessKey string `json:"s3_secret_access_key"`
}

func main() {
	var (
		region     = flag.String("region", "us-east-1", "AWS region")
		secretName = flag.String("secret-name", "portfolio-backend-secrets", "Name of the secret in AWS Secrets Manager")
		action     = flag.String("action", "create", "Action to perform: create, update, get, delete, list")
	)
	flag.Parse()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(*region))
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}

	client := secretsmanager.NewFromConfig(cfg)

	switch *action {
	case "create":
		createSecret(client, *secretName)
	case "update":
		updateSecret(client, *secretName)
	case "get":
		getSecret(client, *secretName)
	case "delete":
		deleteSecret(client, *secretName)
	case "list":
		listSecrets(client)
	default:
		fmt.Printf("Unknown action: %s\n", *action)
		fmt.Println("Available actions: create, update, get, delete, list")
		os.Exit(1)
	}
}

func createSecret(client *secretsmanager.Client, secretName string) {
	// Create default secret data from environment variables or prompts
	secretData := SecretData{
		DatabaseURL:       getEnvOrDefault("DATABASE_URL", "portfolio.db"),
		JWTSecretKey:      getEnvOrDefault("JWT_SECRET_KEY", "your-super-secret-jwt-key-change-in-production"),
		S3Endpoint:        getEnvOrDefault("S3_ENDPOINT", "https://media.moclawr.com"),
		S3Region:          getEnvOrDefault("S3_REGION", "us-east-1"),
		S3Bucket:          getEnvOrDefault("S3_BUCKET", "portfolio-bucket"),
		S3AccessKeyID:     getEnvOrDefault("S3_ACCESS_KEY_ID", "test"),
		S3SecretAccessKey: getEnvOrDefault("S3_SECRET_ACCESS_KEY", "test"),
	}

	secretString, err := json.Marshal(secretData)
	if err != nil {
		log.Fatalf("Failed to marshal secret data: %v", err)
	}

	description := "Portfolio Backend API configuration secrets"
	secretStr := string(secretString)
	input := &secretsmanager.CreateSecretInput{
		Name:         &secretName,
		SecretString: &secretStr,
		Description:  &description,
	}

	result, err := client.CreateSecret(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to create secret: %v", err)
	}

	fmt.Printf("✓ Successfully created secret: %s\n", *result.Name)
	fmt.Printf("  ARN: %s\n", *result.ARN)
	fmt.Println("\nSecret contains the following configuration:")
	printSecretData(secretData)
}

func updateSecret(client *secretsmanager.Client, secretName string) {
	// First, get the existing secret to preserve values
	existing, err := getSecretData(client, secretName)
	if err != nil {
		log.Fatalf("Failed to get existing secret: %v", err)
	}

	fmt.Println("Current secret values:")
	printSecretData(*existing)
	fmt.Println("\nUpdating with new values from environment variables...")

	// Update with new values from environment variables
	secretData := SecretData{
		DatabaseURL:       getEnvOrDefault("DATABASE_URL", existing.DatabaseURL),
		JWTSecretKey:      getEnvOrDefault("JWT_SECRET_KEY", existing.JWTSecretKey),
		S3Endpoint:        getEnvOrDefault("S3_ENDPOINT", existing.S3Endpoint),
		S3Region:          getEnvOrDefault("S3_REGION", existing.S3Region),
		S3Bucket:          getEnvOrDefault("S3_BUCKET", existing.S3Bucket),
		S3AccessKeyID:     getEnvOrDefault("S3_ACCESS_KEY_ID", existing.S3AccessKeyID),
		S3SecretAccessKey: getEnvOrDefault("S3_SECRET_ACCESS_KEY", existing.S3SecretAccessKey),
	}

	secretString, err := json.Marshal(secretData)
	if err != nil {
		log.Fatalf("Failed to marshal secret data: %v", err)
	}

	secretStr := string(secretString)
	input := &secretsmanager.UpdateSecretInput{
		SecretId:     &secretName,
		SecretString: &secretStr,
	}

	_, err = client.UpdateSecret(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to update secret: %v", err)
	}

	fmt.Printf("✓ Successfully updated secret: %s\n", secretName)
	fmt.Println("\nNew secret values:")
	printSecretData(secretData)
}

func getSecret(client *secretsmanager.Client, secretName string) {
	secretData, err := getSecretData(client, secretName)
	if err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}

	fmt.Printf("Secret: %s\n", secretName)
	printSecretData(*secretData)
}

func getSecretData(client *secretsmanager.Client, secretName string) (*SecretData, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	var secretData SecretData
	err = json.Unmarshal([]byte(*result.SecretString), &secretData)
	if err != nil {
		return nil, err
	}

	return &secretData, nil
}

func deleteSecret(client *secretsmanager.Client, secretName string) {
	fmt.Printf("Are you sure you want to delete secret '%s'? This action cannot be undone.\n", secretName)
	fmt.Print("Type 'yes' to confirm: ")

	var confirmation string
	fmt.Scanln(&confirmation)

	if confirmation != "yes" {
		fmt.Println("Operation cancelled.")
		return
	}

	forceDelete := true
	input := &secretsmanager.DeleteSecretInput{
		SecretId:                   &secretName,
		ForceDeleteWithoutRecovery: &forceDelete,
	}

	_, err := client.DeleteSecret(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to delete secret: %v", err)
	}

	fmt.Printf("✓ Successfully deleted secret: %s\n", secretName)
}

func listSecrets(client *secretsmanager.Client) {
	input := &secretsmanager.ListSecretsInput{}

	result, err := client.ListSecrets(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to list secrets: %v", err)
	}

	if len(result.SecretList) == 0 {
		fmt.Println("No secrets found.")
		return
	}

	fmt.Println("Available secrets:")
	for _, secret := range result.SecretList {
		if secret.Name != nil {
			fmt.Printf("  - %s", *secret.Name)
			if secret.Description != nil {
				fmt.Printf(" (%s)", *secret.Description)
			}
			fmt.Println()
		}
	}
}

func printSecretData(data SecretData) {
	fmt.Printf("  Database URL: %s\n", maskSensitive(data.DatabaseURL))
	fmt.Printf("  JWT Secret Key: %s\n", maskSensitive(data.JWTSecretKey))
	fmt.Printf("  S3 Endpoint: %s\n", data.S3Endpoint)
	fmt.Printf("  S3 Region: %s\n", data.S3Region)
	fmt.Printf("  S3 Bucket: %s\n", data.S3Bucket)
	fmt.Printf("  S3 Access Key ID: %s\n", maskSensitive(data.S3AccessKeyID))
	fmt.Printf("  S3 Secret Access Key: %s\n", maskSensitive(data.S3SecretAccessKey))
}

func maskSensitive(value string) string {
	if len(value) <= 8 {
		return "****"
	}
	return value[:4] + "****" + value[len(value)-4:]
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
