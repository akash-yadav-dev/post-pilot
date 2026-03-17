package redis

import (
	"context"
	"errors"
	"testing"
	"time"
)

type fakeRepository struct {
	getValue string
	getErr   error

	setKey   string
	setValue string
	setTTL   time.Duration
	setErr   error

	deleteKey string
	deleteErr error

	existsKey string
	existsVal bool
	existsErr error

	incrementKey string
	incrementTTL time.Duration
	incrementVal int64
	incrementErr error

	healthErr error
}

type fakeMetrics struct {
	hits   int
	misses int
	errors int
	lastOp string
}

func (m *fakeMetrics) IncHit() {
	m.hits++
}

func (m *fakeMetrics) IncMiss() {
	m.misses++
}

func (m *fakeMetrics) IncError(operation string) {
	m.errors++
	m.lastOp = operation
}

func (f *fakeRepository) Get(_ context.Context, _ string) (string, error) {
	return f.getValue, f.getErr
}

func (f *fakeRepository) Set(_ context.Context, key string, value string, ttl time.Duration) error {
	f.setKey = key
	f.setValue = value
	f.setTTL = ttl
	return f.setErr
}

func (f *fakeRepository) Delete(_ context.Context, key string) error {
	f.deleteKey = key
	return f.deleteErr
}

func (f *fakeRepository) Exists(_ context.Context, key string) (bool, error) {
	f.existsKey = key
	return f.existsVal, f.existsErr
}

func (f *fakeRepository) Increment(_ context.Context, key string, ttl time.Duration) (int64, error) {
	f.incrementKey = key
	f.incrementTTL = ttl
	return f.incrementVal, f.incrementErr
}

func (f *fakeRepository) HealthCheck(_ context.Context) error {
	return f.healthErr
}

func TestCacheServiceSetUsesDefaultTTL(t *testing.T) {
	repo := &fakeRepository{}
	svc := NewCacheService(repo, CacheOptions{DefaultTTL: 30 * time.Second, KeyPrefix: "app"})

	if err := svc.Set(context.Background(), "jobs:1", "value", 0); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	if repo.setKey != "app:jobs:1" {
		t.Fatalf("set key = %q, want %q", repo.setKey, "app:jobs:1")
	}

	if repo.setTTL != 30*time.Second {
		t.Fatalf("set ttl = %s, want %s", repo.setTTL, 30*time.Second)
	}
}

func TestCacheServiceSetPreservesSubsecondTTL(t *testing.T) {
	repo := &fakeRepository{}
	svc := NewCacheService(repo, CacheOptions{})

	ttl := 250 * time.Millisecond
	if err := svc.Set(context.Background(), "k", "v", ttl); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	if repo.setTTL != ttl {
		t.Fatalf("set ttl = %s, want %s", repo.setTTL, ttl)
	}
}

func TestCacheServiceSetRejectsNegativeTTL(t *testing.T) {
	repo := &fakeRepository{}
	svc := NewCacheService(repo, CacheOptions{})

	err := svc.Set(context.Background(), "k", "v", -1*time.Second)
	if !errors.Is(err, ErrInvalidTTL) {
		t.Fatalf("Set() error = %v, want ErrInvalidTTL", err)
	}
}

func TestCacheServiceSetEncodesStructAsJSON(t *testing.T) {
	repo := &fakeRepository{}
	svc := NewCacheService(repo, CacheOptions{})

	payload := struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{ID: 7, Name: "post"}

	if err := svc.Set(context.Background(), "post", payload, time.Second); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	if repo.setValue != `{"id":7,"name":"post"}` {
		t.Fatalf("set value = %q, want JSON payload", repo.setValue)
	}
}

func TestCacheServiceSetRejectsEmptyKey(t *testing.T) {
	repo := &fakeRepository{}
	svc := NewCacheService(repo, CacheOptions{})

	err := svc.Set(context.Background(), "   ", "v", time.Second)
	if !errors.Is(err, ErrInvalidKey) {
		t.Fatalf("Set() error = %v, want ErrInvalidKey", err)
	}
}

func TestCacheServiceSetRejectsNilValue(t *testing.T) {
	repo := &fakeRepository{}
	svc := NewCacheService(repo, CacheOptions{})

	err := svc.Set(context.Background(), "key", nil, time.Second)
	if !errors.Is(err, ErrNilValue) {
		t.Fatalf("Set() error = %v, want ErrNilValue", err)
	}
}

func TestCacheServiceSetRejectsNilRepository(t *testing.T) {
	svc := NewCacheService(nil, CacheOptions{})

	err := svc.Set(context.Background(), "key", "v", time.Second)
	if !errors.Is(err, ErrNilRepository) {
		t.Fatalf("Set() error = %v, want ErrNilRepository", err)
	}
}

func TestCacheServiceGetRecordsHitOnSuccess(t *testing.T) {
	repo := &fakeRepository{getValue: "value"}
	metrics := &fakeMetrics{}
	svc := NewCacheService(repo, CacheOptions{Metrics: metrics})

	got, err := svc.Get(context.Background(), "key")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got != "value" {
		t.Fatalf("Get() value = %q, want %q", got, "value")
	}
	if metrics.hits != 1 || metrics.misses != 0 || metrics.errors != 0 {
		t.Fatalf("metrics = %+v, want hit=1 miss=0 error=0", metrics)
	}
}

func TestCacheServiceGetRecordsMissOnCacheMiss(t *testing.T) {
	repo := &fakeRepository{getErr: ErrCacheMiss}
	metrics := &fakeMetrics{}
	svc := NewCacheService(repo, CacheOptions{Metrics: metrics})

	_, err := svc.Get(context.Background(), "key")
	if !errors.Is(err, ErrCacheMiss) {
		t.Fatalf("Get() error = %v, want ErrCacheMiss", err)
	}
	if metrics.hits != 0 || metrics.misses != 1 || metrics.errors != 0 {
		t.Fatalf("metrics = %+v, want hit=0 miss=1 error=0", metrics)
	}
}

func TestCacheServiceGetRecordsErrorOnFailure(t *testing.T) {
	repo := &fakeRepository{getErr: errors.New("redis down")}
	metrics := &fakeMetrics{}
	svc := NewCacheService(repo, CacheOptions{Metrics: metrics})

	_, err := svc.Get(context.Background(), "key")
	if err == nil {
		t.Fatal("Get() error = nil, want non-nil")
	}
	if metrics.errors != 1 || metrics.lastOp != "get" {
		t.Fatalf("metrics = %+v, want errors=1 op=get", metrics)
	}
}

func TestCacheServiceExistsRecordsHitAndMiss(t *testing.T) {
	metrics := &fakeMetrics{}
	repo := &fakeRepository{existsVal: true}
	svc := NewCacheService(repo, CacheOptions{Metrics: metrics})

	exists, err := svc.Exists(context.Background(), "key")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if !exists {
		t.Fatal("Exists() = false, want true")
	}
	if metrics.hits != 1 || metrics.misses != 0 {
		t.Fatalf("metrics after hit = %+v, want hit=1 miss=0", metrics)
	}

	repo.existsVal = false
	exists, err = svc.Exists(context.Background(), "key")
	if err != nil {
		t.Fatalf("Exists() error = %v", err)
	}
	if exists {
		t.Fatal("Exists() = true, want false")
	}
	if metrics.hits != 1 || metrics.misses != 1 {
		t.Fatalf("metrics after miss = %+v, want hit=1 miss=1", metrics)
	}
}

func TestCacheServiceHealthCheck(t *testing.T) {
	repo := &fakeRepository{}
	metrics := &fakeMetrics{}
	svc := NewCacheService(repo, CacheOptions{Metrics: metrics})

	if err := svc.HealthCheck(context.Background()); err != nil {
		t.Fatalf("HealthCheck() error = %v", err)
	}
	if metrics.errors != 0 {
		t.Fatalf("metrics.errors = %d, want 0", metrics.errors)
	}

	repo.healthErr = errors.New("ping failed")
	err := svc.HealthCheck(context.Background())
	if err == nil {
		t.Fatal("HealthCheck() error = nil, want non-nil")
	}
	if metrics.errors != 1 || metrics.lastOp != "health_check" {
		t.Fatalf("metrics = %+v, want errors=1 op=health_check", metrics)
	}
}
