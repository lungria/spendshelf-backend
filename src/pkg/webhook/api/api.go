package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lungria/spendshelf-backend/src/pkg/webhook/db"
	gzap "github.com/gin-contrib/zap"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type WebHookAPI struct {
	Database  *db.Database
	HTTPServer *http.Server
	Logger 		*zap.SugaredLogger
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
		Database: database,
		HTTPServer:   nil,
		Logger: logger.Sugar(),
	}
	a.initRouter(addr)


}

func (a *WebHookAPI) initRouter(addr string) {
	router := gin.New()

	logger, _ := zap.NewProduction()

	router.Use(gzap.Ginzap(logger, time.RFC3339, true))
	router.Use(gzap.RecoveryWithZap(logger, true))

	router.Any("/webhook", a.WebHookHandler)

	a.HTTPServer.Addr = addr
	a.HTTPServer.Handler = router
}
