package scheduler

import (
	"context"
	"time"

	"post-pilot/apps/api/internal/posts/model"
	"post-pilot/apps/api/internal/posts/service"
)

// Scheduler polls for due posts and enqueues them for publishing.
// The actual publish work is performed by the worker app via the queue.
type Scheduler struct {
	svc      service.PostService
	interval time.Duration
	stop     chan struct{}
}

func NewScheduler(svc service.PostService, interval time.Duration) *Scheduler {
	return &Scheduler{
		svc:      svc,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start runs the scheduler in the background until Stop is called.
func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.tick(ctx)
		case <-s.stop:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stop)
}

func (s *Scheduler) tick(ctx context.Context) {
	// TODO: query posts with status=scheduled and scheduled_at <= now(),
	//       then enqueue each via the queue package.
	_ = model.PostStatusScheduled
}
