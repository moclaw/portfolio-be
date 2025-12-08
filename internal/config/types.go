package config

// Config represents the application configuration
type Config struct {
	Port                 string
	Host                 string
	DatabaseType         string // "postgres" or "sqlite"
	DatabaseURL          string
	RedisConfig          RedisConfig
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
	PublicURL       string // Public URL for serving files (e.g., https://media.moclawr.com)
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey string
	Issuer    string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Enabled  bool
}
