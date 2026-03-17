package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrNilRepository = errors.New("redis cache repository is nil")
	ErrInvalidKey    = errors.New("cache key must not be empty")
	ErrInvalidTTL    = errors.New("cache ttl must be >= 0")
	ErrNilValue      = errors.New("cache value must not be nil")
)

type CacheOptions struct {
	DefaultTTL time.Duration
	KeyPrefix  string
	Metrics    MetricsRecorder
}

type MetricsRecorder interface {
	IncHit()
	IncMiss()
	IncError(operation string)
}

type CacheService struct {
	repo       Repository
	defaultTTL time.Duration
	keyPrefix  string
	metrics    MetricsRecorder
}

func NewCacheService(repo Repository, opts CacheOptions) *CacheService {
	ttl := opts.DefaultTTL
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}

	return &CacheService{
		repo:       repo,
		defaultTTL: ttl,
		keyPrefix:  strings.Trim(opts.KeyPrefix, ":"),
		metrics:    opts.Metrics,
	}
}

func (s *CacheService) Get(ctx context.Context, key string) (string, error) {
	builtKey, err := s.buildKey(key)
	if err != nil {
		s.recordError("get")
		return "", err
	}

	value, err := s.repo.Get(ctx, builtKey)
	if err == nil {
		s.recordHit()
		return value, nil
	}

	if errors.Is(err, ErrCacheMiss) {
		s.recordMiss()
		return "", err
	}

	s.recordError("get")
	return "", err
}

func (s *CacheService) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if s.repo == nil {
		s.recordError("set")
		return ErrNilRepository
	}

	builtKey, err := s.buildKey(key)
	if err != nil {
		s.recordError("set")
		return err
	}

	encodedValue, err := normalizeValue(value)
	if err != nil {
		s.recordError("set")
		return err
	}

	if ttl == 0 {
		ttl = s.defaultTTL
	}

	if ttl < 0 {
		s.recordError("set")
		return ErrInvalidTTL
	}

	if err := s.repo.Set(ctx, builtKey, encodedValue, ttl); err != nil {
		s.recordError("set")
		return err
	}

	return nil
}

func (s *CacheService) Delete(ctx context.Context, key string) error {
	if s.repo == nil {
		s.recordError("delete")
		return ErrNilRepository
	}
	builtKey, err := s.buildKey(key)
	if err != nil {
		s.recordError("delete")
		return err
	}
	if err := s.repo.Delete(ctx, builtKey); err != nil {
		s.recordError("delete")
		return err
	}
	return nil
}

func (s *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	if s.repo == nil {
		s.recordError("exists")
		return false, ErrNilRepository
	}
	builtKey, err := s.buildKey(key)
	if err != nil {
		s.recordError("exists")
		return false, err
	}

	exists, err := s.repo.Exists(ctx, builtKey)
	if err != nil {
		s.recordError("exists")
		return false, err
	}

	if exists {
		s.recordHit()
	} else {
		s.recordMiss()
	}

	return exists, nil
}

func (s *CacheService) Increment(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	if s.repo == nil {
		s.recordError("increment")
		return 0, ErrNilRepository
	}

	builtKey, err := s.buildKey(key)
	if err != nil {
		s.recordError("increment")
		return 0, err
	}

	if ttl < 0 {
		s.recordError("increment")
		return 0, ErrInvalidTTL
	}

	if ttl == 0 {
		ttl = s.defaultTTL
	}

	value, err := s.repo.Increment(ctx, builtKey, ttl)
	if err != nil {
		s.recordError("increment")
		return 0, err
	}

	return value, nil
}

func (s *CacheService) HealthCheck(ctx context.Context) error {
	if s.repo == nil {
		s.recordError("health_check")
		return ErrNilRepository
	}

	if err := s.repo.HealthCheck(ctx); err != nil {
		s.recordError("health_check")
		return err
	}

	return nil
}

func (s *CacheService) buildKey(key string) (string, error) {
	if s.repo == nil {
		return "", ErrNilRepository
	}

	key = strings.TrimSpace(key)
	if key == "" {
		return "", ErrInvalidKey
	}

	if s.keyPrefix == "" {
		return key, nil
	}
	return s.keyPrefix + ":" + key, nil
}

func normalizeValue(value any) (string, error) {
	if value == nil {
		return "", ErrNilValue
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		payload, err := json.Marshal(v)
		if err != nil {
			return "", fmt.Errorf("marshal cache value: %w", err)
		}
		return string(payload), nil
	}
}

func (s *CacheService) recordHit() {
	if s.metrics != nil {
		s.metrics.IncHit()
	}
}

func (s *CacheService) recordMiss() {
	if s.metrics != nil {
		s.metrics.IncMiss()
	}
}

func (s *CacheService) recordError(operation string) {
	if s.metrics != nil {
		s.metrics.IncError(operation)
	}
}
