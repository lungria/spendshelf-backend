package job

import (
	"context"
	"time"
)

// Scheduler allows to schedule background jobs with graceful cancellation.
type Scheduler struct {
}

// Job describes functions, thaca can be runned by JobRunner. Job must implement cancellation via context.
type Job func(ctx context.Context)

// Schedule background job. It will be runned with waitBeforeRuns period (it will wait for this period even for first)
// and cancelled, when ctx is cancelled.
func (s Scheduler) Schedule(ctx context.Context, job Job, waitBeforeRuns time.Duration) {
	ticker := time.NewTicker(waitBeforeRuns)
	go func(t *time.Ticker, j Job) {
		defer t.Stop()
		for {
			select {
			case _ = <-t.C:
				{
					j(ctx)
				}
			case _ = <-ctx.Done():
				{
					return
				}

			}
		}
	}(ticker, job)
}
