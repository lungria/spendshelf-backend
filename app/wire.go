// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/lungria/spendshelf-backend/api"
	"github.com/lungria/spendshelf-backend/api/handler"
	"github.com/lungria/spendshelf-backend/app/config"
	"github.com/lungria/spendshelf-backend/importer"
	"github.com/lungria/spendshelf-backend/importer/account"
	"github.com/lungria/spendshelf-backend/importer/interval"
	"github.com/lungria/spendshelf-backend/importer/mono"
	"github.com/lungria/spendshelf-backend/importer/transactions"
	"github.com/lungria/spendshelf-backend/storage"
)

func InitializeApp() (*State, error) {
	wire.Build(
		config.FromEnv,
		NewDbPoolProvider,
		NewMonoAPIProvider,
		storage.NewPostgreSQLStorage,
		storage.NewAccountsStorage,
		interval.NewGenerator,
		NewSchedulerProvider,
		NewRoutesProvider,
		api.NewServer,
		NewAppStateProvider,

		importer.NewImporter,
		transactions.NewImporter,
		account.NewDefaultImporter,

		wire.Bind(new(transactions.TransactionsStorage), new(*storage.PostgreSQLStorage)),
		wire.Bind(new(interval.TransactionsStorage), new(*storage.PostgreSQLStorage)),
		wire.Bind(new(account.Storage), new(*storage.AccountsStorage)),
		wire.Bind(new(importer.AccountImporter), new(*account.DefaultImporter)),
		wire.Bind(new(importer.TransactionsImporter), new(*transactions.DefaultImporter)),
		wire.Bind(new(transactions.TransactionsBankAPI), new(*mono.Client)),
		wire.Bind(new(account.BankAPI), new(*mono.Client)),
		wire.Bind(new(transactions.ImportIntervalGenerator), new(*interval.Generator)),

		handler.NewTransactionHandler,
		handler.NewAccountHandler,

		wire.Bind(new(handler.TransactionStorage), new(*storage.PostgreSQLStorage)),
		wire.Bind(new(handler.AccountsStorage), new(*storage.AccountsStorage)),
	)
	return &State{}, nil
}
