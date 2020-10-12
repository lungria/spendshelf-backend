// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lungria/spendshelf-backend/api"
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
		NewAppStateProvider,

		wire.Bind(new(interval.TransactionsStorage), new(*storage.PostgreSQLStorage)),
		wire.Bind(new(importer.BankAPI), new(*mono.Client)),
		wire.Bind(new(importer.TransactionsStorage), new(*storage.PostgreSQLStorage)),
		wire.Bind(new(importer.ImportIntervalGenerator), new(*interval.Generator)),
	)
	return &State{}, nil
}
