package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	MigrationsPath string

	AllowDestructive bool
}

func Load() (*Config, error) {
	dbPort, err := getEnvAsInt("DB_PORT", 5432)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "post_pilot"),
		DBSSLMode:  getEnv("DB_SSLMODE", "require"),

		MigrationsPath: getEnv("MIGRATIONS_PATH", filepath.ToSlash("packages/database/migrations")),

		AllowDestructive: getEnvAsBool("ALLOW_DESTRUCTIVE_MIGRATIONS", false),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.DBHost) == "" || strings.TrimSpace(c.DBUser) == "" || strings.TrimSpace(c.DBName) == "" {
		return fmt.Errorf("DB_HOST, DB_USER and DB_NAME are required")
	}
	if c.DBPort <= 0 || c.DBPort > 65535 {
		return fmt.Errorf("DB_PORT must be between 1 and 65535")
	}
	if strings.TrimSpace(c.DBSSLMode) == "" {
		return fmt.Errorf("DB_SSLMODE is required")
	}
	if strings.TrimSpace(c.MigrationsPath) == "" {
		return fmt.Errorf("MIGRATIONS_PATH is required")
	}

	return nil
}

func (c *Config) SourceURL() string {
	path := filepath.ToSlash(c.MigrationsPath)
	if strings.HasPrefix(path, "file://") {
		return path
	}
	return "file://" + path
}

func (c *Config) DatabaseURL() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUser, c.DBPassword),
		Host:   fmt.Sprintf("%s:%d", c.DBHost, c.DBPort),
		Path:   c.DBName,
	}

	q := url.Values{}
	q.Set("sslmode", c.DBSSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) (int, error) {
	raw := strings.TrimSpace(getEnv(key, ""))
	if raw == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}

	return value, nil
}

func getEnvAsBool(key string, fallback bool) bool {
	raw := strings.TrimSpace(getEnv(key, ""))
	if raw == "" {
		return fallback
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}

	return value
}
