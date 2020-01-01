package api

import (
	"net/http"

	"github.com/lungria/spendshelf-backend/src/config"

	"github.com/gin-gonic/gin"
)

// NewAPI create a new WebHookAPI with DB, logger and router
func NewAPI(cfg *config.EnvironmentConfiguration, router *gin.Engine) (*http.Server, error) {
	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}
	return server, nil
}
