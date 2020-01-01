package api

import (
	"net/http"
	"time"

	"github.com/lungria/spendshelf-backend/src/config"

	"github.com/lungria/spendshelf-backend/src/api/handlers"
	"go.uber.org/zap"

	gzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

// NewAPI create a new WebHookAPI with DB, logger and router
func NewAPI(cfg *config.EnvironmentConfiguration, logger *zap.Logger, hookHandler *handlers.WebHookHandler, ctgHandler *handlers.CategoriesHandler) (*http.Server, error) {

	router := gin.New()
	router.Use(gzap.Ginzap(logger, time.RFC3339, true))
	router.Use(gzap.RecoveryWithZap(logger, true))

	router.Any("/webhook", hookHandler.Handle)
	router.POST("/categories", ctgHandler.Handle)

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}
	return server, nil
}
