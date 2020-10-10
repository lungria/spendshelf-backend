package app

import (
	"context"

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
	// scheduler.Schedule(ctx, i.Import(cfg.MonoAccountID), 1*time.Minute, 30*time.Second)
	return scheduler
}

func NewRoutesProvider() []api.RouteBinder {
	return nil
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
