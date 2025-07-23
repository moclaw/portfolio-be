package config

// Config represents the application configuration
type Config struct {
	Port                 string
	Host                 string
	DatabaseURL          string
	S3Config             S3Config
	JWTConfig            JWTConfig
	SecretsManagerConfig SecretsManagerConfig
}

// SecretsManagerConfig holds AWS Secrets Manager configuration
type SecretsManagerConfig struct {
	SecretName string
	Region     string
	UseSecrets bool
}

// S3Config holds S3 configuration
type S3Config struct {
	Endpoint        string
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	ForcePathStyle  bool
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey string
	Issuer    string
}
