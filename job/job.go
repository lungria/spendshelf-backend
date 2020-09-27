package job

import (
	"context"
	"sync"
	"time"
)

// Scheduler allows to schedule background jobs with graceful cancellation.
type Scheduler struct {
	wg *sync.WaitGroup
}

// Job describes functions, that can be runned by JobRunner. Job must implement cancellation via context.
type Job func(ctx context.Context)

// Schedule background job. It will be runned with waitBeforeRuns period (it will wait for this period even for first)
// and cancelled, when ctx is cancelled.
func (s *Scheduler) Schedule(ctx context.Context, job Job, waitBeforeRuns, jobTimeout time.Duration) {
	s.wg.Add(1)
	ticker := time.NewTicker(waitBeforeRuns)
	go func(t *time.Ticker, j Job, timeout time.Duration) {
		defer s.wg.Done()
		defer t.Stop()
		for {
			select {
			case _ = <-t.C:
				{
					executeWithTimeout(ctx, j, timeout)
				}
			case _ = <-ctx.Done():
				{
					return
				}

			}
		}
	}(ticker, job, jobTimeout)
}

// Wait blocks until the all scheduled jobs exist.
func (s *Scheduler) Wait() {
	s.wg.Wait()
}

func executeWithTimeout(ctx context.Context, job Job, jobTimeout time.Duration) {
	timeoutCtx, cancel := context.WithTimeout(ctx, jobTimeout)
	defer cancel()
	job(timeoutCtx)
}
