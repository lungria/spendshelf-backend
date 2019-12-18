package api

import (
	"net/http"
	"time"

	gzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/lungria/spendshelf-backend/src/db"
	"go.uber.org/zap"
)

type WebHookAPI struct {
	Database   *db.Database
	HTTPServer *http.Server
	Logger     *zap.SugaredLogger
}

func NewAPI(addr, dbname, mongoURI string) (*WebHookAPI, error) {
	database, err := db.NewDatabase(dbname, mongoURI)
	if err != nil {
		return nil, err
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	a := WebHookAPI{
		Database:   database,
		HTTPServer: nil,
		Logger:     logger.Sugar(),
	}
	a.initRouter(addr)

	return &a, nil
}

func (a *WebHookAPI) initRouter(addr string) {
	router := gin.New()

	logger, _ := zap.NewProduction()

	router.Use(gzap.Ginzap(logger, time.RFC3339, true))
	router.Use(gzap.RecoveryWithZap(logger, true))

	router.Any("/webhook", a.WebHookHandler)

	a.HTTPServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
