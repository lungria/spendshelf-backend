//+build wireinject

package main

import (
	"net/http"

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

func InitializeServer() (*http.Server, error) {
	wire.Build(config.NewConfig,
		mongoDbProvider,
		categories.NewCachedRepository,
		handlers.NewCategoriesHandler,
		db.NewTransactionsMongoDbRepository,
		handlers.NewWebHookHandler,
		zapProvider,
		sugarProvider,
		api.NewAPI,
		wire.Bind(new(db.TransactionsRepository), new(*db.TransactionsMongoDbRepository)),
		wire.Bind(new(categories.Repository), new(*categories.CachedRepository)),
	)
	return &http.Server{}, nil
}
