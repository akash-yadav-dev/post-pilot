package auth

import (
	"context"
	"sync"
	"time"
)

type InMemoryTokenStore struct {
	mu      sync.RWMutex
	tokens  map[string]time.Time
	ticker  *time.Ticker
	stopped chan struct{}
}

func NewInMemoryTokenStore() *InMemoryTokenStore {
	store := &InMemoryTokenStore{
		tokens:  make(map[string]time.Time),
		ticker:  time.NewTicker(1 * time.Minute),
		stopped: make(chan struct{}),
	}

	go store.cleanupLoop()

	return store
}

func (s *InMemoryTokenStore) SaveRefreshToken(_ context.Context, jti string, ttl time.Duration) error {
	expiry := time.Now().Add(ttl)

	s.mu.Lock()
	s.tokens[jti] = expiry
	s.mu.Unlock()

	return nil
}

func (s *InMemoryTokenStore) ExistsRefreshToken(_ context.Context, jti string) (bool, error) {
	now := time.Now()

	s.mu.RLock()
	expiry, ok := s.tokens[jti]
	s.mu.RUnlock()

	if !ok {
		return false, nil
	}

	if now.After(expiry) {
		s.mu.Lock()
		delete(s.tokens, jti)
		s.mu.Unlock()
		return false, nil
	}

	return true, nil
}

func (s *InMemoryTokenStore) DeleteRefreshToken(_ context.Context, jti string) error {
	s.mu.Lock()
	delete(s.tokens, jti)
	s.mu.Unlock()

	return nil
}

func (s *InMemoryTokenStore) cleanupLoop() {
	for {
		select {
		case <-s.ticker.C:
			s.cleanupExpired()
		case <-s.stopped:
			s.ticker.Stop()
			return
		}
	}
}

func (s *InMemoryTokenStore) cleanupExpired() {
	now := time.Now()

	s.mu.Lock()
	for jti, expiry := range s.tokens {
		if now.After(expiry) {
			delete(s.tokens, jti)
		}
	}
	s.mu.Unlock()
}

func (s *InMemoryTokenStore) Close() {
	close(s.stopped)
}
