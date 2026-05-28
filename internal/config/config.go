package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config holds all configuration for the application.
// It is loaded once at startup and passed explicitly to all layers.
type Config struct {
	// HTTP server settings
	HTTPPort string `env:"HTTP_PORT" env-default:"8080"`

	// Database connection settings
	DBHost     string `env:"DB_HOST"     env-required:"true"`
	DBPort     string `env:"DB_PORT"     env-default:"5432"`
	DBName     string `env:"DB_NAME"     env-required:"true"`
	DBUser     string `env:"DB_USER"     env-required:"true"`
	DBPassword string `env:"DB_PASSWORD" env-required:"true"`
	DBSSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`

	// JWT configuration
	JWTSecret      string `env:"JWT_SECRET"       env-required:"true"`
	JWTExpireHours int    `env:"JWT_EXPIRE_HOURS" env-default:"24"`

	// Security parameters
	BCryptCost int `env:"BCRYPT_COST" env-default:"12"`

	// Application behaviour
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
}

// Load reads configuration from a .env file (if present) and then from
// environment variables. It validates the resulting values and returns
// a parsed Config or an error.
func Load() (*Config, error) {
	var cfg Config

	// Optionally load .env file if it exists
	if _, err := os.Stat(".env"); err == nil {
		if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
			return nil, fmt.Errorf("reading .env file: %w", err)
		}
	}

	// Overwrite with actual environment variables
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("reading env: %w", err)
	}

	// Validation
	if cfg.BCryptCost < 12 {
		return nil, fmt.Errorf("BCRYPT_COST must be >= 12, got %d", cfg.BCryptCost)
	}

	if cfg.JWTExpireHours <= 0 {
		return nil, fmt.Errorf("JWT_EXPIRE_HOURS must be > 0, got %d", cfg.JWTExpireHours)
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[cfg.LogLevel] {
		cfg.LogLevel = "info"
	}

	return &cfg, nil
}

// DSN returns a PostgreSQL connection string suitable for pgx.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}
