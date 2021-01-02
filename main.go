package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lungria/spendshelf-backend/job"

	"github.com/rs/zerolog/log"

	"github.com/lungria/spendshelf-backend/app"
)

func main() {
	state, err := app.InitializeApp()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize app")
	}

	defer state.Close()

	state.API.Start()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if state.Config.EnableImportJob {
		state.Scheduler.Schedule(ctx, job.Job{
			Run:            state.Importer.Import(state.Config.MonoAccountID),
			WaitBeforeRuns: 1 * time.Minute,
			Timeout:        30 * time.Second,
		})
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
