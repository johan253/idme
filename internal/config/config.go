package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/joho/godotenv/autoload" // Load .env file automatically
)

type Config struct {
	Port               int
	JwtSecret          string
	JwtRefreshSecret   string
	JwtTTLSeconds      int
	DatabaseURL        string
	Kek                []byte
	JwkRefreshInterval time.Duration
	JwkGraceDuration   time.Duration
}

func Load() (*Config, error) {
	// Server port
	port := envInt("PORT", 8080)

	// JWT secret and TTL
	jwtSecret := envOr("JWT_SECRET", "")
	jwtRefreshSecret := envOr("JWT_REFRESH_SECRET", "")
	jwtTTLSeconds := envInt("JWT_TTL_SECONDS", 600) // Default to 10 minutes

	// Database configuration
	postgresHost := envOr("POSTGRES_HOST", "localhost")
	postgresPort := envInt("POSTGRES_PORT", 5432)
	postgresDatabase := envOr("POSTGRES_DB", "identity")
	postgresUser := envOr("POSTGRES_USER", "postgres")
	postgresPassword := envOr("POSTGRES_PASSWORD", "password")
	postgresSSLMode := envOr("POSTGRES_SSL_MODE", "disable")
	postgresURL := envOr("POSTGRES_URL", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", postgresUser, postgresPassword, postgresHost, postgresPort, postgresDatabase, postgresSSLMode))

	// JWK params
	kekRaw := envOr("JWK_KEK", "")
	jwkRefreshInterval := envInt("JWK_REFRESH_SECONDS", 60)
	jwkGraceSeconds := envInt("JWT_GRACE_SECONDS", 60)

	// Validate required environment variables
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}
	if jwtRefreshSecret == "" {
		return nil, fmt.Errorf("JWT_REFRESH_SECRET environment variable is required")
	}
	if err := validPostgresURL(postgresURL); err != nil {
		return nil, fmt.Errorf("invalid Postgres URL: %s", postgresURL)
	}
	if kekRaw == "" {
		return nil, fmt.Errorf("JWK_KEK environment variable is required")
	}
	kek, err := base64.StdEncoding.DecodeString(kekRaw)
	if err != nil {
		return nil, fmt.Errorf("JWK_KEK is not valid base64: %w", err)
	}
	if len(kek) != 32 {
		return nil, fmt.Errorf("JWR_KEK must decode to 32 bytes, got %d", len(kek))
	}

	// Return the configuration
	cfg := &Config{
		Port:               port,
		JwtSecret:          jwtSecret,
		JwtRefreshSecret:   jwtRefreshSecret,
		JwtTTLSeconds:      jwtTTLSeconds,
		DatabaseURL:        postgresURL,
		Kek:                kek,
		JwkRefreshInterval: time.Duration(jwkRefreshInterval) * time.Second,
		JwkGraceDuration:   time.Duration(jwkGraceSeconds) * time.Second,
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
