package main

import (
	"net/http"
	"time"

	gzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/lungria/spendshelf-backend/api"
	"go.uber.org/zap"
)

// NewServer is create a new server
func NewServer(addr string) *http.Server {
	router := gin.New()

	logger, _ := zap.NewProduction()

	router.Use(gzap.Ginzap(logger, time.RFC3339, true))
	router.Use(gin.Recovery())

	router.GET("/webhook", api.WebHookHandlerGet)
	router.POST("/webhook", api.WebHookHandlerPost)

	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
