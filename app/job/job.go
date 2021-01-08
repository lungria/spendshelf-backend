package job

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// Job describes functions, that can be run by JobRunner.
type Job struct {
	// Run is the method which would be executed as job.
	Run func(ctx context.Context)
	// WaitBeforeRuns is a duration of time between each job execution.
	// Timeout can be used to set single job run timeout.
	WaitBeforeRuns, Timeout time.Duration
}

// Scheduler allows to schedule background jobs with graceful cancellation.
type Scheduler struct {
	wg *sync.WaitGroup
}

// NewScheduler creates new instance of Scheduler.
func NewScheduler() *Scheduler {
	return &Scheduler{wg: &sync.WaitGroup{}}
}

// Schedule background job. It will be run with Job.WaitBeforeRuns period and canceled, when ctx is canceled.
// Job.Timeout can be used to set single job run timeout.
func (s *Scheduler) Schedule(ctx context.Context, job Job) {
	s.wg.Add(1)

	ticker := time.NewTicker(job.WaitBeforeRuns)
	go func(t *time.Ticker, j Job) {
		defer s.wg.Done()
		defer t.Stop()

		for {
			select {
			case _ = <-t.C:
				{
					executeWithTimeout(ctx, j)
					log.Debug().Msg("job finished")
				}
			case _ = <-ctx.Done():
				{
					return
				}
			}
		}
	}(ticker, job)
}

// Wait blocks until the all scheduled jobs exist. They will be canceled when ctx parameter of the Schedule method is
// canceled, but if you need to know if they were actually canceled - you should wait on this method.
func (s *Scheduler) Wait() {
	s.wg.Wait()
}

func executeWithTimeout(ctx context.Context, job Job) {
	timeoutCtx, cancel := context.WithTimeout(ctx, job.Timeout)
	defer cancel()
	job.Run(timeoutCtx)
}
