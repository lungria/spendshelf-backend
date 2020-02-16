package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/lungria/spendshelf-backend/src/db"
	"github.com/lungria/spendshelf-backend/src/mqtt"
)

type App struct {
	web      *Server
	queue    *mqtt.Listener
	database *db.Connection
	logger   *zap.SugaredLogger
}

func NewApp(web *Server, queue *mqtt.Listener, database *db.Connection, logger *zap.SugaredLogger) *App {
	return &App{web: web, queue: queue, database: database, logger: logger}
}

func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err := a.database.KeepConnected(ctx)
		if err != nil {
			a.logger.Error("db keepconnected returned error", zap.Error(err))
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		err := a.queue.Listen(ctx)
		if err != nil {
			a.logger.Error("mqtt listener returned error", zap.Error(err))
		}
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		err := a.web.Listen(ctx)
		if err != nil {
			a.logger.Error("webserver returned error", zap.Error(err))
		}
		wg.Done()
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	<-sigChan
	waitCh := make(chan struct{})
	go func() {
		cancel()
		wg.Wait()
		waitCh <- struct{}{}
	}()
	select {
	case <-waitCh:
		a.logger.Info("shut down. see ya")
	case <-time.After(4 * time.Second):
		a.logger.Warn("unable to shut down in time")
	}

}
