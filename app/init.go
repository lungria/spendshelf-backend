package app

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/account"
	"github.com/lungria/spendshelf-backend/api"
	"github.com/lungria/spendshelf-backend/app/config"
	"github.com/lungria/spendshelf-backend/budget"
	"github.com/lungria/spendshelf-backend/importer"
	"github.com/lungria/spendshelf-backend/importer/mono"
	"github.com/lungria/spendshelf-backend/storage"
	"github.com/lungria/spendshelf-backend/transaction"
	"github.com/lungria/spendshelf-backend/transaction/category"
	"github.com/lungria/spendshelf-backend/transaction/interval"
	"github.com/rs/zerolog/log"
)

func Initialize() (*App, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, err
	}

	client := mono.NewClient(cfg.MonoBaseURL, cfg.MonoAPIKey)
	// todo: ping API?

	pool, err := pgxpool.Connect(context.Background(), cfg.DBConnString)
	if err != nil {
		return nil, err
	}

	api, acRepo, txRepo := initAPI(cfg, pool)
	im := initImporter(cfg, client, acRepo, txRepo)
	db := storage.NewWorker(pool)

	state := NewApp(&log.Logger)
	state.RegisterWorkers(api, im, db)

	return state, nil
}

func initAPI(cfg config.Config, pool *pgxpool.Pool) (Worker, *account.Repository, *transaction.Repository) {
	acRepo := account.NewRepository(pool)
	acHandler := account.NewHandler(acRepo)

	bgRepo := budget.NewRepository(pool)
	bgHandler := budget.NewHandler(bgRepo)

	ctRepo := category.NewRepository(pool)
	txRepo := transaction.NewRepository(pool)
	txHandler := transaction.NewHandler(txRepo, ctRepo)

	return api.NewServer(cfg, acHandler, bgHandler, txHandler), acRepo, txRepo
}

func initImporter(cfg config.Config, client *mono.Client, acRepo *account.Repository, txRepo *transaction.Repository) Worker {
	gen := interval.NewGenerator(txRepo)
	acIm := account.NewImporter(client, acRepo)
	txIm := transaction.NewImporter(client, txRepo, gen)
	globalIm := importer.NewImporter(acIm, txIm)

	return importer.NewWorker(globalIm, cfg.MonoAccountID)
}
