package app

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
)

// Worker describes something that starts working when Start is called and stops when Close is called.
type Worker interface {
	Start()
	Close() error
}

// App handles application lifecycle.
type App struct {
	log     *zerolog.Logger
	workers []Worker
}

// NewApp creates and returns new App.
func NewApp(log *zerolog.Logger) *App {
	app := &App{
		log:     log,
		workers: []Worker{},
	}

	return app
}

// Run application, block until interruption received, and handle graceful shutdown.
func (a *App) Run() {
	a.log.Info().Msg("app started")

	wg := &sync.WaitGroup{}

	for _, worker := range a.workers {
		wg.Add(1)

		a.log.Trace().Func(func(e *zerolog.Event) {
			typeOf := fmt.Sprintf("%T", worker)
			e.Str("typeOf", typeOf).Msg("starting worker")
		})

		go worker.Start()
	}

	a.handleGracefulShutdown(wg)
}

// RegisterWorkers allows to add workers into app lifecycle.
// App will keep them alive until shutdown signal is received.
func (a *App) RegisterWorkers(workers ...Worker) {
	a.workers = append(a.workers, workers...)
}

// handleGracefulShutdown listens to sigterm and allows to gracefully shutdown whole app.
func (a *App) handleGracefulShutdown(wg *sync.WaitGroup) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	for i := range a.workers {
		worker := a.workers[i]

		a.log.Trace().Func(func(e *zerolog.Event) {
			typeOf := fmt.Sprintf("%T", worker)
			e.Str("typeOf", typeOf).Msg("closing worker")
		})
		wg.Done()

		if err := worker.Close(); err != nil {
			a.log.Error().Err(err).Msg("error when closing worker")
		}
	}

	wg.Wait()

	a.log.Info().Msg("app shutdown gracefully")
}
