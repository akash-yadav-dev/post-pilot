package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")
var ErrNilRedisClient = errors.New("redis client is nil")

type Repository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Increment(ctx context.Context, key string, ttl time.Duration) (int64, error)
	HealthCheck(ctx context.Context) error
}

type RedisRepository struct {
	client *Client
}

func NewRepository(client *Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	if err := r.ensureClient(); err != nil {
		return "", err
	}

	value, err := r.client.Get(ctx, key)
	if err == goredis.Nil {
		return "", ErrCacheMiss
	}
	if err != nil {
		return "", fmt.Errorf("redis get key %q: %w", key, err)
	}
	return value, err
}

func (r *RedisRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	if err := r.ensureClient(); err != nil {
		return err
	}

	if err := r.client.Set(ctx, key, value, ttl); err != nil {
		return fmt.Errorf("redis set key %q: %w", key, err)
	}
	return nil
}

func (r *RedisRepository) Delete(ctx context.Context, key string) error {
	if err := r.ensureClient(); err != nil {
		return err
	}

	if err := r.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("redis delete key %q: %w", key, err)
	}
	return nil
}

func (r *RedisRepository) Exists(ctx context.Context, key string) (bool, error) {
	if err := r.ensureClient(); err != nil {
		return false, err
	}

	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists key %q: %w", key, err)
	}
	return count > 0, nil
}

func (r *RedisRepository) Increment(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	if err := r.ensureClient(); err != nil {
		return 0, err
	}

	pipe := r.client.TxPipeline()
	incrCmd := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)

	if _, err := pipe.Exec(ctx); err != nil {
		return 0, fmt.Errorf("redis increment key %q: %w", key, err)
	}

	return incrCmd.Val(), nil
}

func (r *RedisRepository) HealthCheck(ctx context.Context) error {
	if err := r.ensureClient(); err != nil {
		return err
	}

	if err := r.client.HealthCheck(ctx); err != nil {
		return fmt.Errorf("redis repository health check failed: %w", err)
	}

	return nil
}

func (r *RedisRepository) ensureClient() error {
	if r == nil || r.client == nil {
		return ErrNilRedisClient
	}

	return nil
}
