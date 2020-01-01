package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	"github.com/lungria/spendshelf-backend/src/api"
)

func main() {
	var err error

	config, err := NewConfig()
	if err != nil {
		log.Fatalln("Couldn't parse environment variables", err)
	}
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Couldn't initialize logger %+v", err)
	}
	sugar := logger.Sugar()
	s, err := api.NewAPI(config.HTTPAddr, config.DBName, config.MongoURI, logger, sugar)
	if err != nil {
		log.Fatalln("Couldn't create a new server")
	}

	done := make(chan bool, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	go func() {
		<-sigChan
		sugar.Info("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s.SetKeepAlivesEnabled(false)
		if err = s.Shutdown(ctx); err != nil {
			sugar.Fatalf("Couldn't gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	if err = s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		sugar.Fatalf("Couldn't listen on %v: %v\n", config.HTTPAddr, err)
	}

	<-done
	sugar.Info("Server stopped")
}
