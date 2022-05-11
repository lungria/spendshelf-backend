package importer

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	waitBeforeRuns = 1 * time.Minute
	timeout        = 30 * time.Second
)

// Worker is the app/Worker interface implementation, so importer lifecycle can be managed with other parts of the
// app lifecycle.
type Worker struct {
	im        *Importer
	accountID string
	ctx       context.Context
	cancel    context.CancelFunc
	wg        *sync.WaitGroup
}

// NewWorker creates new instance of Worker.
func NewWorker(im *Importer, accountID string) *Worker {
	ctx, cancel := context.WithCancel(context.Background())

	return &Worker{im: im, accountID: accountID, ctx: ctx, cancel: cancel, wg: &sync.WaitGroup{}}
}

// Start importer worker. Blocks until Close() is called.
func (w *Worker) Start() {
	ticker := time.NewTicker(waitBeforeRuns)

	w.wg.Add(1)

	for {
		select {
		case _ = <-ticker.C:
			w.executeWithTimeout(w.ctx)
			log.Debug().Msg("import finished")
		case _ = <-w.ctx.Done():
			return
		}
	}
}

func (w *Worker) executeWithTimeout(ctx context.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	w.im.Import(timeoutCtx, w.accountID)
}

// Close importer worker, blocks until import is fully stopped.
func (w Worker) Close() error {
	w.ctx.Done()
	w.wg.Wait()

	return nil
}
