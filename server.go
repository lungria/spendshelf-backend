package main

import (
	"github.com/lungria/spendshelf-backend/pkg/webhook"
	"net/http"
	"time"

	gzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewServer is create a new server
func NewServer(addr string) *http.Server {
	router := gin.New()

	logger, _ := zap.NewProduction()

	router.Use(gzap.Ginzap(logger, time.RFC3339, true))
	router.Use(gzap.RecoveryWithZap(logger, true))

	router.GET("/webhook", webhook.WebHookHandlerGet)
	router.POST("/webhook", webhook.WebHookHandlerPost)

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
