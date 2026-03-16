package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServePort string

	// Database configuration
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// JWT configuration
	JWTSecretKey string
	JWTExpiry    int

	AppBaseURL string
}

func LoadConfig() (*Config, error) {
	dbPort, err := getEnvAsInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	jwtExpiry, err := getEnvAsInt("JWT_EXPIRY", 3600)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ServePort: getEnv("APP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "post-pilot"),

		JWTSecretKey: getEnv("JWT_SECRET_KEY", "mysecretkey"),
		JWTExpiry:    jwtExpiry,

		AppBaseURL: getEnv("APP_BASE_URL", "http://localhost:8080"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.ServePort == "" {
		return fmt.Errorf("APP_PORT is required")
	}
	if c.DBHost == "" || c.DBUser == "" || c.DBName == "" {
		return fmt.Errorf("DB_HOST, DB_USER and DB_NAME are required")
	}
	if c.DBPort <= 0 || c.DBPort > 65535 {
		return fmt.Errorf("DB_PORT must be between 1 and 65535")
	}
	if c.JWTSecretKey == "" {
		return fmt.Errorf("JWT_SECRET_KEY is required for production use")
	}
	if c.JWTExpiry <= 0 {
		return fmt.Errorf("JWT_EXPIRY must be greater than 0")
	}

	return nil
}

func getEnv(key string, fallback string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) (int, error) {
	raw := getEnv(key, "")
	if raw == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}

	return value, nil
}
