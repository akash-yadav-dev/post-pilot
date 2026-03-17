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
	JWTAccessSecretKey  string
	JWTRefreshSecretKey string
	JWTExpiry           int
	JWTRefreshExpiry    int

	PasswordHashCost int

	AuthRateLimitWindowSeconds int
	AuthLoginRateLimitMax      int
	AuthRegisterRateLimitMax   int
	AuthRefreshRateLimitMax    int

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

	jwtRefreshExpiry, err := getEnvAsInt("JWT_REFRESH_EXPIRY", 604800)
	if err != nil {
		return nil, err
	}

	passwordHashCost, err := getEnvAsInt("BCRYPT_COST", 12)
	if err != nil {
		return nil, err
	}

	rateLimitWindow, err := getEnvAsInt("AUTH_RATE_LIMIT_WINDOW_SECONDS", 60)
	if err != nil {
		return nil, err
	}

	loginRateLimitMax, err := getEnvAsInt("AUTH_LOGIN_RATE_LIMIT_MAX", 10)
	if err != nil {
		return nil, err
	}

	registerRateLimitMax, err := getEnvAsInt("AUTH_REGISTER_RATE_LIMIT_MAX", 5)
	if err != nil {
		return nil, err
	}

	refreshRateLimitMax, err := getEnvAsInt("AUTH_REFRESH_RATE_LIMIT_MAX", 20)
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

		JWTAccessSecretKey:  getEnv("JWT_ACCESS_SECRET_KEY", "dev-access-secret-change-me-at-least-32-bytes"),
		JWTRefreshSecretKey: getEnv("JWT_REFRESH_SECRET_KEY", "dev-refresh-secret-change-me-at-least-32-bytes"),
		JWTExpiry:           jwtExpiry,
		JWTRefreshExpiry:    jwtRefreshExpiry,

		PasswordHashCost: passwordHashCost,

		AuthRateLimitWindowSeconds: rateLimitWindow,
		AuthLoginRateLimitMax:      loginRateLimitMax,
		AuthRegisterRateLimitMax:   registerRateLimitMax,
		AuthRefreshRateLimitMax:    refreshRateLimitMax,

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
	if c.JWTAccessSecretKey == "" {
		return fmt.Errorf("JWT_ACCESS_SECRET_KEY is required for production use")
	}
	if c.JWTRefreshSecretKey == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET_KEY is required for production use")
	}
	if len(c.JWTAccessSecretKey) < 32 {
		return fmt.Errorf("JWT_ACCESS_SECRET_KEY must be at least 32 bytes")
	}
	if len(c.JWTRefreshSecretKey) < 32 {
		return fmt.Errorf("JWT_REFRESH_SECRET_KEY must be at least 32 bytes")
	}
	if c.JWTAccessSecretKey == c.JWTRefreshSecretKey {
		return fmt.Errorf("JWT access and refresh secrets must be different")
	}
	if c.JWTExpiry <= 0 {
		return fmt.Errorf("JWT_EXPIRY must be greater than 0")
	}
	if c.JWTRefreshExpiry <= 0 {
		return fmt.Errorf("JWT_REFRESH_EXPIRY must be greater than 0")
	}
	if c.PasswordHashCost < 12 {
		return fmt.Errorf("BCRYPT_COST must be >= 12")
	}
	if c.AuthRateLimitWindowSeconds <= 0 {
		return fmt.Errorf("AUTH_RATE_LIMIT_WINDOW_SECONDS must be greater than 0")
	}
	if c.AuthLoginRateLimitMax <= 0 {
		return fmt.Errorf("AUTH_LOGIN_RATE_LIMIT_MAX must be greater than 0")
	}
	if c.AuthRegisterRateLimitMax <= 0 {
		return fmt.Errorf("AUTH_REGISTER_RATE_LIMIT_MAX must be greater than 0")
	}
	if c.AuthRefreshRateLimitMax <= 0 {
		return fmt.Errorf("AUTH_REFRESH_RATE_LIMIT_MAX must be greater than 0")
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
