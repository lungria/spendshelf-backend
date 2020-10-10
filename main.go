package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/lungria/spendshelf-backend/app"
)

func main() {
	state, err := app.InitializeApp()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize app")
	}

	defer state.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	state.API.Start()
	if state.Config.EnableImportJob {
		state.Scheduler.Schedule(ctx, state.Importer.Import(state.Config.MonoAccountID), 1*time.Minute, 30*time.Second)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
}
