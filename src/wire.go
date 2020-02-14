/*//+build wireinject
 */
package main

import (
	"github.com/lungria/spendshelf-backend/src/report"

	"github.com/lungria/spendshelf-backend/src/transactions"

	"github.com/lungria/spendshelf-backend/src/webhooks"

	"go.uber.org/zap"

	"github.com/lungria/spendshelf-backend/src/api"

	"github.com/lungria/spendshelf-backend/src/config"

	"github.com/lungria/spendshelf-backend/src/api/handlers"
	"github.com/lungria/spendshelf-backend/src/categories"

	"github.com/lungria/spendshelf-backend/src/db"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/wire"
)

func mongoDbProvider(cfg *config.EnvironmentConfiguration) (*mongo.Database, error) {
	return db.NewDatabase(cfg.DBName, cfg.MongoURI)
}

func sugarProvider(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}

func zapProvider() (*zap.Logger, error) {
	return zap.NewProduction()
}

func InitializeServer() (*api.Server, error) {
	wire.Build(config.NewConfig,
		mongoDbProvider,
		categories.NewRepository,
		handlers.NewCategoriesHandler,
		transactions.NewTransactionRepository,
		report.NewSequentialReportGenerator,
		handlers.NewTransactionsHandler,
		webhooks.NewWebHookRepository,
		handlers.NewWebHookHandler,
		handlers.NewReportsHandler,
		zapProvider,
		sugarProvider,
		wire.Bind(new(api.ServerConfig), new(*config.EnvironmentConfiguration)),
		wire.Bind(new(transactions.Repository), new(*transactions.TransactionRepository)),
		wire.Bind(new(webhooks.Repository), new(*webhooks.WebHookRepository)),
		wire.Bind(new(report.Generator), new(*report.SequentialReportGenerator)),
		api.RoutesProvider,
		api.NewPipelineBuilder,
		api.NewServer,
	)
	return &api.Server{}, nil
}
