/*//+build wireinject
 */
package main

import (
	"github.com/lungria/spendshelf-backend/src/transactions"

	"go.uber.org/zap"

	"github.com/lungria/spendshelf-backend/src/api"

	"github.com/lungria/spendshelf-backend/src/config"

	"github.com/lungria/spendshelf-backend/src/categories"

	"github.com/google/wire"
	"github.com/lungria/spendshelf-backend/src/db"
)

func sugarProvider(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}

func zapProvider() (*zap.Logger, error) {
	return zap.NewProduction()
}

func InitializeServer() (*api.Server, error) {
	wire.Build(
		config.NewConfig,
		wire.Bind(new(api.ServerConfig), new(*config.EnvironmentConfiguration)),
		wire.Bind(new(db.Config), new(*config.EnvironmentConfiguration)),
		db.NewDatabase,
		zapProvider,
		sugarProvider,
		categories.NewRepository,
		categories.NewHandler,
		transactions.NewRepository,
		transactions.NewHandler,
		api.RoutesProvider,
		api.NewPipelineBuilder,
		api.NewServer,
	)
	return &api.Server{}, nil
}
