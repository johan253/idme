package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload" // Load .env file automatically
)

type Config struct {
	Port          int
	JwtSecret     string
	JwtTtlSeconds int
	DatabaseURL   string
}

func Load() (*Config, error) {
	// Server port
	port := envInt("PORT", 8080)

	// JWT secret and TTL
	jwtSecret := envOr("JWT_SECRET", "")
	jwtTtlSeconds := envInt("JWT_TTL_SECONDS", 600) // Default to 10 minutes

	// Database configuration
	postgresHost := envOr("POSTGRES_HOST", "localhost")
	postgresPort := envInt("POSTGRES_PORT", 5432)
	postgresDatabase := envOr("POSTGRES_DB", "identity")
	postgresUser := envOr("POSTGRES_USER", "postgres")
	postgresPassword := envOr("POSTGRES_PASSWORD", "password")
	postgresSSLMode := envOr("POSTGRES_SSL_MODE", "disable")
	postgresURL := envOr("POSTGRES_URL", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", postgresUser, postgresPassword, postgresHost, postgresPort, postgresDatabase, postgresSSLMode))

	// Validate required environment variables
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	if err := validPostgresURL(postgresURL); err != nil {
		return nil, fmt.Errorf("invalid Postgres URL: %s", postgresURL)
	}

	// Return the configuration
	cfg := &Config{
		Port:          port,
		JwtSecret:     jwtSecret,
		JwtTtlSeconds: jwtTtlSeconds,
	}

	return cfg, nil
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func envOr(key string, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func validPostgresURL(url string) error {
	config, err := pgx.ParseConfig(url)
	if err != nil {
		return fmt.Errorf("invalid Postgres URL: %w", err)
	}
	if config.Host == "" || config.Port == 0 || config.Database == "" || config.User == "" {
		return fmt.Errorf("incomplete Postgres URL: %s", url)
	}
	return nil
}
