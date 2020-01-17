//+build wireinject

package main

import (
	"time"

	"github.com/lungria/spendshelf-backend/src/report"

	"github.com/lungria/spendshelf-backend/src/transactions"

	"github.com/gin-contrib/cors"
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

func mongoDbProvider(cfg *config.EnvironmentConfiguration) (*mongo.Database, error) {
	return db.NewDatabase(cfg.DBName, cfg.MongoURI)
}

func sugarProvider(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}

func zapProvider() (*zap.Logger, error) {
	return zap.NewProduction()
}

func defaultHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("content-type", "application/json")
		c.Next()
	}
}

func routerProvider(logger *zap.Logger, hookHandler *handlers.WebHookHandler, ctgHandler *handlers.CategoriesHandler, txHandler *handlers.TransactionsHandler, rpHandler *handlers.ReportsHandler) *gin.Engine {
	router := gin.New()
	router.Use(gzap.Ginzap(logger, time.RFC3339, true))
	router.Use(gzap.RecoveryWithZap(logger, true))
	router.Use(defaultHeaders())
	router.Use(cors.Default())
	router.GET("/webhook", hookHandler.HandleGet)
	router.POST("/webhook", hookHandler.HandlePost)
	router.POST("/categories", ctgHandler.HandlePost)
	router.GET("/categories", ctgHandler.HandleGet)
	router.GET("/transactions", txHandler.HandleGet)
	router.PATCH("/transactions/:transactionID", txHandler.HandlePatch)
	router.GET("/reports", rpHandler.HandleGet)
	return router
}

func InitializeServer() (*config.Dependencies, error) {
	wire.Build(config.NewConfig,
		mongoDbProvider,
		categories.NewCachedRepository,
		handlers.NewCategoriesHandler,
		transactions.NewTransactionRepository,
		report.NewSequentialReportGenerator,
		handlers.NewTransactionsHandler,
		webhooks.NewWebHookRepository,
		handlers.NewWebHookHandler,
		handlers.NewReportsHandler,
		zapProvider,
		sugarProvider,
		routerProvider,
		api.NewAPI,
		wire.Bind(new(transactions.Repository), new(*transactions.TransactionRepository)),
		wire.Bind(new(webhooks.Repository), new(*webhooks.WebHookRepository)),
		wire.Bind(new(categories.Repository), new(*categories.CachedRepository)),
		wire.Bind(new(report.Generator), new(*report.SequentialReportGenerator)),
		wire.Struct(new(config.Dependencies), "Logger", "Server"),
	)
	return &config.Dependencies{}, nil
}
