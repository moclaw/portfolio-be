package config

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
