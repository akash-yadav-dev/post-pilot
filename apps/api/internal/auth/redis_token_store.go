package auth

import (
	"context"
	"time"

	rediscache "post-pilot/packages/cache/redis"
)

type RedisTokenStore struct {
	cache *rediscache.CacheService
}

func NewRedisTokenStore(cache *rediscache.CacheService) *RedisTokenStore {
	return &RedisTokenStore{cache: cache}
}

func (s *RedisTokenStore) SaveRefreshToken(ctx context.Context, jti string, ttl time.Duration) error {
	return s.cache.Set(ctx, refreshTokenStoreKey(jti), "1", ttl)
}

func (s *RedisTokenStore) ExistsRefreshToken(ctx context.Context, jti string) (bool, error) {
	return s.cache.Exists(ctx, refreshTokenStoreKey(jti))
}

func (s *RedisTokenStore) DeleteRefreshToken(ctx context.Context, jti string) error {
	return s.cache.Delete(ctx, refreshTokenStoreKey(jti))
}

func refreshTokenStoreKey(jti string) string {
	return "refresh_token:" + jti
}
