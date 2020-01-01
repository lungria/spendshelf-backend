package api

import (
	"net/http"
	"time"

	gzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/lungria/spendshelf-backend/src/db"
	"go.uber.org/zap"
)

// WebHookAPI is API instance with DB, logger and router
type WebHookAPI struct {
	Database   *db.Database
	HTTPServer *http.Server
	Logger     *zap.SugaredLogger
}

// NewAPI create a new WebHookAPI with DB, logger and router
func NewAPI(addr, dbname, mongoURI string) (*WebHookAPI, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()
	database, err := db.NewDatabase(dbname, mongoURI, sugar)
	if err != nil {
		return nil, err
	}

	a := WebHookAPI{
		Database:   database,
		HTTPServer: nil,
		Logger:     sugar,
	}
	a.InitRouter(addr, logger)

	return &a, nil
}

// InitRouter is initiate a new router also using in tests
func (a *WebHookAPI) InitRouter(addr string, logger *zap.Logger) {
	router := gin.New()

	router.Use(gzap.Ginzap(logger, time.RFC3339, true))
	router.Use(gzap.RecoveryWithZap(logger, true))

	router.Any("/webhook", a.WebHookHandler)

	a.HTTPServer = &http.Server{
		Addr:    addr,
		Handler: router,
	}
}
