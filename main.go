package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/config"
	"github.com/lungria/spendshelf-backend/importer"
	"github.com/lungria/spendshelf-backend/importer/interval"
	"github.com/lungria/spendshelf-backend/job"
	"github.com/lungria/spendshelf-backend/mono"
	"github.com/lungria/spendshelf-backend/storage"
)

const dbConnString = "postgres://localhost:5432/postgres?sslmode=disable"

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	bckgCtx := context.Background()
	ctx, cancel := context.WithCancel(bckgCtx)
	defer cancel()

	apiClient := mono.NewClient(cfg.MonoBaseURL, cfg.MonoAPIKey)
	dbpool, err := pgxpool.Connect(context.Background(), cfg.DBConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	s := storage.NewPostgreSQLStorage(dbpool)
	intervalGen := interval.NewSimpleIntervalGenerator(s)
	i := importer.NewImporeter(apiClient, s, intervalGen)

	scheduler := job.Scheduler{}
	scheduler.Schedule(ctx, i.Import(cfg.MonoAccountID), 1*time.Minute)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	cancel()

	scheduler.Wait()
}
