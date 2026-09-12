package config

import (
	"fmt"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload" // Load .env file automatically
)

type Config struct {
	Port      int
	JwtSecret string
}

func Load() (*Config, error) {
	port := envInt("PORT", 8080)
	jwtSecret := envOr("JWT_SECRET", "")

	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	cfg := &Config{
		Port:      port,
		JwtSecret: jwtSecret,
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
