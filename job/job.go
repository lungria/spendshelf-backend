package job

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Job describes functions, that can be runned by JobRunner. Job must implement cancelation via context.
type Job func(ctx context.Context)

// Scheduler allows to schedule background jobs with graceful cancelation.
type Scheduler struct {
	wg *sync.WaitGroup
}

// NewScheduler creates new instance of Scheduler.
func NewScheduler() *Scheduler {
	return &Scheduler{wg: &sync.WaitGroup{}}
}

// Schedule background job. It will be runned with waitBeforeRuns period and canceled, when ctx is canceled.
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
					log.Debug().Msg("job finished")
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
