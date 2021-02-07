package app

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/api"
	"github.com/lungria/spendshelf-backend/api/handler"
	"github.com/lungria/spendshelf-backend/app/config"
	"github.com/lungria/spendshelf-backend/app/job"
	"github.com/lungria/spendshelf-backend/importer"
	"github.com/lungria/spendshelf-backend/importer/mono"
)

func NewDbPoolProvider(cfg config.Config) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.Connect(context.Background(), cfg.DBConnString)
	return dbpool, err
}

func NewMonoAPIProvider(cfg config.Config) *mono.Client {
	return mono.NewClient(cfg.MonoBaseURL, cfg.MonoAPIKey)
}

func NewSchedulerProvider() *job.Scheduler {
	scheduler := job.NewScheduler()
	return scheduler
}

func NewRoutesProvider(t *handler.TransactionHandler, a *handler.AccountHandler) []api.RouteBinder {
	return []api.RouteBinder{t, a}
}

func NewAppStateProvider(
	s *job.Scheduler,
	i *importer.Importer,
	a *api.Server,
	pool *pgxpool.Pool,
	cfg config.Config) *State {
	return &State{
		API:       a,
		Scheduler: s,
		DB:        pool,
		Importer:  i,
		Config:    cfg,
	}
}
