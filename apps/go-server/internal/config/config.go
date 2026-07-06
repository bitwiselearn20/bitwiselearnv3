// Package config loads 12-factor environment configuration.
// Mirrors the settings exposed by the legacy Python config.py so the
// Go service is a drop-in replacement during the strangler migration.
package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds all runtime settings. Defaults match config.py.
type Config struct {
	// Server
	Port        string
	FrontendURL string
	BackendURL  string

	// Database
	DatabaseURL                 string
	MongoServerSelectionTimeout int // milliseconds

	// JWT secrets
	JWTAccessSecret  string
	JWTRefreshSecret string
	ResetTokenSecret string

	// Object storage: Azure Blob is the storage backend for the Go rewrite.
	AzureStorageConnectionString string
	AzureStorageContainer        string

	// S3 fields kept only so the strangler period can read pre-migration URLs.
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSS3Region        string
	AWSS3Bucket        string

	// Cloudinary
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string

	// Email
	EmailUser string
	EmailPass string

	// Async messaging + code execution
	MQClient            string
	CodeExecutionServer string

	// Cache
	RedisURL string
}

func getEnv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			return n
		}
	}
	return def
}

// Load reads configuration from the environment.
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		BackendURL:  getEnv("BACKEND_PUBLIC_URL", "http://localhost:8080"),

		DatabaseURL:                 getEnv("DATABASE_URL", ""),
		MongoServerSelectionTimeout: getEnvInt("MONGO_SERVER_SELECTION_TIMEOUT_MS", 8000),

		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", "access-secret-change-me"),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", "refresh-secret-change-me"),
		ResetTokenSecret: getEnv("RESET_TOKEN_SECRET", "reset-token-secret-change-me"),

		AzureStorageConnectionString: getEnv("AZURE_STORAGE_CONNECTION_STRING", ""),
		AzureStorageContainer:        getEnv("AZURE_STORAGE_CONTAINER", "bitwise-learn"),

		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSS3Region:        getEnv("AWS_S3_REGION", "ap-south-1"),
		AWSS3Bucket:        getEnv("AWS_S3_BUCKET", "bitwise-learn"),

		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:    getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret: getEnv("CLOUDINARY_API_SECRET", ""),

		EmailUser: getEnv("EMAIL_USER", ""),
		EmailPass: getEnv("EMAIL_PASS", ""),

		MQClient:            getEnv("MQ_CLIENT", "amqp://guest:guest@localhost/"),
		CodeExecutionServer: getEnv("CODE_EXECUTION_SERVER", "http://localhost:2000/"),

		RedisURL: getEnv("REDIS_URL", ""),
	}
}

// DBName derives the Mongo database name from the connection string,
// matching the Python logic: url.rsplit("/",1)[-1].split("?")[0] or "bitwiselearn".
func (c *Config) DBName() string {
	url := c.DatabaseURL
	if i := strings.LastIndex(url, "/"); i >= 0 {
		url = url[i+1:]
	}
	if i := strings.Index(url, "?"); i >= 0 {
		url = url[:i]
	}
	if url == "" {
		return "bitwiselearn"
	}
	return url
}
