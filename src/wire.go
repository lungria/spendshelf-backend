//+build wireinject

package main

import (
	"github.com/lungria/spendshelf-backend/src/transactions"

	"go.uber.org/zap"

	"github.com/lungria/spendshelf-backend/src/app"

	"github.com/lungria/spendshelf-backend/src/config"

	"github.com/google/wire"
	"github.com/lungria/spendshelf-backend/src/db"
)

func sugarProvider(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}

func zapProvider() (*zap.Logger, error) {
	return zap.NewProduction()
}

func InitializeServer() (*app.Server, error) {
	wire.Build(
		config.NewConfig,
		wire.Bind(new(app.ServerConfig), new(*config.EnvironmentConfiguration)),
		wire.Bind(new(db.Config), new(*config.EnvironmentConfiguration)),
		db.OpenConnection,
		zapProvider,
		sugarProvider,
		transactions.NewStore,
		transactions.NewHandler,
		app.RoutesProvider,
		app.NewPipelineBuilder,
		app.NewServer,
	)
	return &app.Server{}, nil
}
