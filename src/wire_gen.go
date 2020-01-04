// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/lungria/spendshelf-backend/src/api"
	"github.com/lungria/spendshelf-backend/src/api/handlers"
	"github.com/lungria/spendshelf-backend/src/categories"
	"github.com/lungria/spendshelf-backend/src/config"
	"github.com/lungria/spendshelf-backend/src/db"
	"github.com/lungria/spendshelf-backend/src/webhooks"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

// Injectors from wire.go:

func InitializeServer() (*config.Dependencies, error) {
	logger, err := zapProvider()
	if err != nil {
		return nil, err
	}
	sugaredLogger := sugarProvider(logger)
	environmentConfiguration, err := config.NewConfig()
	if err != nil {
		return nil, err
	}
	database, err := mongoDbProvider(environmentConfiguration)
	if err != nil {
		return nil, err
	}
	webHookRepository, err := webhooks.NewWebHookRepository(database, sugaredLogger)
	if err != nil {
		return nil, err
	}
	webHookHandler, err := handlers.NewWebHookHandler(webHookRepository, sugaredLogger)
	if err != nil {
		return nil, err
	}
	cachedRepository, err := categories.NewCachedRepository(database)
	if err != nil {
		return nil, err
	}
	categoriesHandler, err := handlers.NewCategoriesHandler(cachedRepository, sugaredLogger)
	if err != nil {
		return nil, err
	}
	engine := routerProvider(logger, webHookHandler, categoriesHandler)
	server, err := api.NewAPI(environmentConfiguration, engine)
	if err != nil {
		return nil, err
	}
	dependencies := &config.Dependencies{
		Logger: sugaredLogger,
		Server: server,
	}
	return dependencies, nil
}

// wire.go:

func mongoDbProvider(cfg *config.EnvironmentConfiguration) (*mongo.Database, error) {
	return db.NewDatabase(cfg.DBName, cfg.MongoURI)
}

func sugarProvider(logger *zap.Logger) *zap.SugaredLogger {
	return logger.Sugar()
}

func zapProvider() (*zap.Logger, error) {
	return zap.NewProduction()
}

func routerProvider(logger *zap.Logger, hookHandler *handlers.WebHookHandler, ctgHandler *handlers.CategoriesHandler) *gin.Engine {
	router := gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	router.GET("/webhook", hookHandler.HandleGet)
	router.POST("/webhook", hookHandler.HandlePost)
	router.POST("/categories", ctgHandler.HandlePost)
	router.GET("/categories", ctgHandler.HandleGet)
	return router
}
