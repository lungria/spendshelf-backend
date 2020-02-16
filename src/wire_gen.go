// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/lungria/spendshelf-backend/src/app"
	"github.com/lungria/spendshelf-backend/src/config"
	"github.com/lungria/spendshelf-backend/src/db"
	"github.com/lungria/spendshelf-backend/src/transactions"
	"go.uber.org/zap"
)

// Injectors from wire.go:

func InitializeServer() (*app.Server, error) {
	environmentConfiguration, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	logger, err := zapProvider()
	if err != nil {
		return nil, err
	}
	sugaredLogger := sugarProvider(logger)
	connection, err := db.OpenConnection(environmentConfiguration, sugaredLogger)
	if err != nil {
		return nil, err
	}
	store := transactions.NewStore(connection, sugaredLogger)
	handler := transactions.NewHandler(store, sugaredLogger)
	v := app.RoutesProvider(handler)
	pipelineBuilder := app.NewPipelineBuilder(logger, v)
	server := app.NewServer(environmentConfiguration, logger, pipelineBuilder, connection)
	return server, nil
}

// wire.go:

func sugarProvider(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}

func zapProvider() (*zap.Logger, error) {
	return zap.NewProduction()
}
