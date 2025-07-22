package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretsManagerService handles AWS Secrets Manager operations
type SecretsManagerService struct {
	client *secretsmanager.Client
	region string
}

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

// CreateSecret creates a new secret in AWS Secrets Manager
func (s *SecretsManagerService) CreateSecret(secretName string, secretData *SecretData, description string) error {
	secretString, err := json.Marshal(secretData)
	if err != nil {
		return fmt.Errorf("failed to marshal secret data: %w", err)
	}

	secretStr := string(secretString)
	input := &secretsmanager.CreateSecretInput{
		Name:         &secretName,
		SecretString: &secretStr,
		Description:  &description,
	}

	_, err = s.client.CreateSecret(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to create secret %s: %w", secretName, err)
	}

	log.Printf("✓ Successfully created secret: %s", secretName)
	return nil
}

// UpdateSecret updates an existing secret in AWS Secrets Manager
func (s *SecretsManagerService) UpdateSecret(secretName string, secretData *SecretData) error {
	secretString, err := json.Marshal(secretData)
	if err != nil {
		return fmt.Errorf("failed to marshal secret data: %w", err)
	}

	secretStr := string(secretString)
	input := &secretsmanager.UpdateSecretInput{
		SecretId:     &secretName,
		SecretString: &secretStr,
	}

	_, err = s.client.UpdateSecret(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to update secret %s: %w", secretName, err)
	}

	log.Printf("✓ Successfully updated secret: %s", secretName)
	return nil
}

// DeleteSecret deletes a secret from AWS Secrets Manager
func (s *SecretsManagerService) DeleteSecret(secretName string, forceDelete bool) error {
	input := &secretsmanager.DeleteSecretInput{
		SecretId: &secretName,
	}

	if forceDelete {
		// Delete immediately without recovery window
		forceDeleteWithoutRecovery := true
		input.ForceDeleteWithoutRecovery = &forceDeleteWithoutRecovery
	}

	_, err := s.client.DeleteSecret(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete secret %s: %w", secretName, err)
	}

	log.Printf("✓ Successfully deleted secret: %s", secretName)
	return nil
}

// ListSecrets lists all secrets in the account
func (s *SecretsManagerService) ListSecrets() ([]string, error) {
	input := &secretsmanager.ListSecretsInput{}

	result, err := s.client.ListSecrets(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	var secretNames []string
	for _, secret := range result.SecretList {
		if secret.Name != nil {
			secretNames = append(secretNames, *secret.Name)
		}
	}

	return secretNames, nil
}
