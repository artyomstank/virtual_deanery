// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	HTTPPort           int
	GRPCPort           int
	DatabaseURL        string
	LogLevel           string
	JWTSecretKey       string
	JWTAlgo            string
	JWTAccessTTL       time.Duration
	JWTRefreshTTL      time.Duration
	BCryptCost         int
	DBMaxConns         int
	CORSAllowedOrigins []string
}

// LoadFromEnv loads configuration from environment variables.
func LoadFromEnv() (*Config, error) {
	// TODO: Implement environment variable loading
	// Parse HTTP_PORT, GRPC_PORT, DATABASE_URL, LOG_LEVEL, etc.
	// Handle parsing errors and type conversions
	// Return Config or error
	return nil, fmt.Errorf("not implemented")
}

// Validate checks if configuration is valid.
func (c *Config) Validate() error {
	// TODO: Validate required fields
	// TODO: Validate port ranges
	// TODO: Validate database URL format
	// TODO: Return error if validation fails
	return nil
}

// getEnvInt gets integer environment variable with default.
func getEnvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

// getEnvString gets string environment variable with default.
func getEnvString(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// getEnvDuration parses environment variable as duration with default.
func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	valStr := os.Getenv(key)
	if val, err := time.ParseDuration(valStr); err == nil {
		return val
	}
	return defaultVal
}
