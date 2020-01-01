package api

import (
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/lungria/spendshelf-backend/src/api/handlers"

	"github.com/lungria/spendshelf-backend/src/categories"
	"github.com/pkg/errors"

	gzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/lungria/spendshelf-backend/src/db"
)

// NewAPI create a new WebHookAPI with DB, logger and router
func NewAPI(addr, dbname, mongoURI string, logger *zap.Logger, sugar *zap.SugaredLogger) (*http.Server, error) {
	database, err := db.NewDatabase(dbname, mongoURI)
	if err != nil {
		return nil, err
	}
	ctgRepo, err := categories.NewCachedRepository(database)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create cached repository")
	}
	ctgHandler, err := handlers.NewCategoriesHandler(ctgRepo)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create categories handler")
	}
	trRepo, err := db.NewTransactionsMongoDbRepository(database, sugar)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create transactions repository")
	}
	hookHandler, err := handlers.NewWebHookHandler(trRepo, sugar)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create transactions handler")
	}
	router := gin.New()
	router.Use(gzap.Ginzap(logger, time.RFC3339, true))
	router.Use(gzap.RecoveryWithZap(logger, true))

	router.Any("/webhook", hookHandler.Handle)
	router.POST("/categories", ctgHandler.Handle)

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	return server, nil
}
