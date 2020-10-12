package app

import (
	"context"

	"github.com/lungria/spendshelf-backend/api/handler"

	"github.com/lungria/spendshelf-backend/mono/importer"

	"github.com/lungria/spendshelf-backend/api"

	"github.com/lungria/spendshelf-backend/job"

	"github.com/lungria/spendshelf-backend/mono"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/config"
)

func NewDbPoolProvider(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.Connect(ctx, cfg.DBConnString)
	return dbpool, err
}

func NewMonoAPIProvider(cfg config.Config) *mono.Client {
	return mono.NewClient(cfg.MonoBaseURL, cfg.MonoAPIKey)
}

func NewCtxProvider() context.Context {
	return context.Background()
}

func NewSchedulerProvider() *job.Scheduler {
	scheduler := job.NewScheduler()
	return scheduler
}

func NewRoutesProvider(t *handler.TransactionHandler) []api.RouteBinder {
	return []api.RouteBinder{t}
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
