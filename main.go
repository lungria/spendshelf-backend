package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/importer"
	"github.com/lungria/spendshelf-backend/job"
	"github.com/lungria/spendshelf-backend/mono"
	"github.com/lungria/spendshelf-backend/storage"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const dbConnString = "postgres://localhost:5432/postgres?sslmode=disable"

func main() {
	bckgCtx := context.Background()
	ctx, cancel := context.WithCancel(bckgCtx)
	defer cancel()

	apiClient := mono.Client{}
	dbpool, err := pgxpool.Connect(context.Background(), dbConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	s := storage.NewPostgreSQLStorage(dbpool)
 	i :=	importer.NewImporeter(&apiClient, s)

	scheduler := job.Scheduler{}
	scheduler.Schedule(ctx, i.Import, 1*time.Minute)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	cancel()

	scheduler.Wait()
}
