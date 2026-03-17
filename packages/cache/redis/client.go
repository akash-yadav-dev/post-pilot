package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

type Config struct {
	Addr            string
	Username        string
	Password        string
	DB              int
	PoolSize        int
	MinIdleConns    int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	PoolTimeout     time.Duration
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
	MaxRetries      int
	MinRetryBackoff time.Duration
	MaxRetryBackoff time.Duration
	TLSConfig       *tls.Config
}

func NewRedisClient() (*Client, error) {
	cfg := Config{
		Addr:            getEnv("REDIS_ADDR", "localhost:6379"),
		Username:        getEnv("REDIS_USERNAME", ""),
		Password:        getEnv("REDIS_PASSWORD", ""),
		DB:              getEnvAsInt("REDIS_DB", 0),
		PoolSize:        getEnvAsInt("REDIS_POOL_SIZE", 10),
		MinIdleConns:    getEnvAsInt("REDIS_MIN_IDLE", 2),
		DialTimeout:     getEnvAsDuration("REDIS_DIAL_TIMEOUT", 3*time.Second),
		ReadTimeout:     getEnvAsDuration("REDIS_READ_TIMEOUT", 2*time.Second),
		WriteTimeout:    getEnvAsDuration("REDIS_WRITE_TIMEOUT", 2*time.Second),
		PoolTimeout:     getEnvAsDuration("REDIS_POOL_TIMEOUT", 4*time.Second),
		ConnMaxIdleTime: getEnvAsDuration("REDIS_CONN_MAX_IDLE", 5*time.Minute),
		ConnMaxLifetime: getEnvAsDuration("REDIS_CONN_MAX_LIFE", 30*time.Minute),
		MaxRetries:      getEnvAsInt("REDIS_MAX_RETRIES", 3),
		MinRetryBackoff: getEnvAsDuration("REDIS_MIN_RETRY_BACKOFF", 100*time.Millisecond),
		MaxRetryBackoff: getEnvAsDuration("REDIS_MAX_RETRY_BACKOFF", 2*time.Second),
		TLSConfig:       redisTLSConfig(),
	}

	return NewClient(cfg)
}

func NewClient(cfg Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:            cfg.Addr,
		Username:        cfg.Username,
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		DialTimeout:     cfg.DialTimeout,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		PoolTimeout:     cfg.PoolTimeout,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
		MaxRetries:      cfg.MaxRetries,
		MinRetryBackoff: cfg.MinRetryBackoff,
		MaxRetryBackoff: cfg.MaxRetryBackoff,
		TLSConfig:       cfg.TLSConfig,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Client{Client: rdb}, nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.Client.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}

func (c *Client) Close() error {
	return c.Client.Close()
}

func (c *Client) HealthCheck(ctx context.Context) error {
	if c == nil || c.Client == nil {
		return ErrNilRedisClient
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
	}

	if err := c.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}

func getEnv(key, fallback string) string {

	val := os.Getenv(key)

	if val == "" {
		return fallback
	}

	return val
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}

func redisTLSConfig() *tls.Config {
	if !getEnvAsBool("REDIS_TLS_ENABLED", false) {
		return nil
	}

	return &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: getEnvAsBool("REDIS_TLS_INSECURE_SKIP_VERIFY", false),
	}
}
