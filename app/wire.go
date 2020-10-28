// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lungria/spendshelf-backend/api"
	"github.com/lungria/spendshelf-backend/api/handler"
	"github.com/lungria/spendshelf-backend/config"
	"github.com/lungria/spendshelf-backend/mono"
	"github.com/lungria/spendshelf-backend/mono/importer"
	"github.com/lungria/spendshelf-backend/mono/importer/interval"
	"github.com/lungria/spendshelf-backend/storage"
)

func InitializeApp() (*State, error) {
	wire.Build(
		config.FromEnv,
		NewCtxProvider,
		NewDbPoolProvider,
		NewMonoAPIProvider,
		storage.NewPostgreSQLStorage,
		interval.NewIntervalGenerator,
		importer.NewImporter,
		NewSchedulerProvider,
		NewRoutesProvider,
		api.NewServer,
		handler.NewTransactionHandler,
		NewAppStateProvider,

		wire.Bind(new(importer.BankAPI), new(*mono.Client)),
		wire.Bind(new(importer.ImportIntervalGenerator), new(*interval.Generator)),

		wire.Bind(new(importer.TransactionsStorage), new(*storage.PostgreSQLStorage)),
		wire.Bind(new(interval.TransactionsStorage), new(*storage.PostgreSQLStorage)),
		wire.Bind(new(handler.TransactionStorage), new(*storage.PostgreSQLStorage)),
	)
	return &State{}, nil
}
