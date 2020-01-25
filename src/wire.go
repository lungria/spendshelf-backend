//+build wireinject

package main

import (
	"context"
	"time"

	"github.com/lungria/spendshelf-backend/src/transactions"

	"github.com/gin-contrib/cors"
	"github.com/lungria/spendshelf-backend/src/syncmono"
	"github.com/lungria/spendshelf-backend/src/webhooks"

	gzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"

	"github.com/lungria/spendshelf-backend/src/api"

	"github.com/lungria/spendshelf-backend/src/config"

	"github.com/lungria/spendshelf-backend/src/api/handlers"
	"github.com/lungria/spendshelf-backend/src/categories"

	"github.com/lungria/spendshelf-backend/src/db"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/wire"
)

func mongoDbProvider(ctx context.Context, cfg *config.EnvironmentConfiguration) (*mongo.Database, error) {
	return db.NewDatabase(ctx, cfg.DBName, cfg.MongoURI)
}

func sugarProvider(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}

func zapProvider() (*zap.Logger, error) {
	return zap.NewProduction()
}

func ctxProvider() context.Context {
	return context.Background()
}

func defaultHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("content-type", "application/json")
		c.Next()
	}
}

func routerProvider(logger *zap.Logger, hookHandler *handlers.WebHookHandler, ctgHandler *handlers.CategoriesHandler, txnHandler *handlers.TransactionsHandler, syncHandler *handlers.SyncMonoHandler) *gin.Engine {
	router := gin.New()
	router.Use(gzap.Ginzap(logger, time.RFC3339, true))
	router.Use(gzap.RecoveryWithZap(logger, true))
	router.Use(defaultHeaders())
	router.Use(cors.Default())
	router.GET("/webhook", hookHandler.HandleGet)
	router.POST("/webhook", hookHandler.HandlePost)
	router.POST("/categories", ctgHandler.HandlePost)
	router.GET("/categories", ctgHandler.HandleGet)
	router.GET("/transactions", txnHandler.HandleGet)
	router.PATCH("/transactions/:transactionID", txnHandler.HandlePatch)
	router.GET("/sync", syncHandler.HandleSocket)
	return router
}

func InitializeServer() (*config.Dependencies, error) {
	wire.Build(config.NewConfig,
		mongoDbProvider,
		categories.NewCachedRepository,
		handlers.NewCategoriesHandler,
		transactions.NewTransactionRepository,
		handlers.NewTransactionsHandler,
		webhooks.NewWebHookRepository,
		handlers.NewWebHookHandler,
		syncmono.NewSyncSocket,
		handlers.NewSyncMonoHandler,
		zapProvider,
		sugarProvider,
		routerProvider,
		ctxProvider,
		api.NewAPI,
		wire.Bind(new(transactions.Repository), new(*transactions.TransactionRepository)),
		wire.Bind(new(webhooks.Repository), new(*webhooks.WebHookRepository)),
		wire.Bind(new(categories.Repository), new(*categories.CachedRepository)),
		wire.Struct(new(config.Dependencies), "Logger", "Server", "Context"),
	)
	return &config.Dependencies{}, nil
}
