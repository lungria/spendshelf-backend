package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/lungria/spendshelf-backend/src/api"
)

func main() {
	var err error

	config, err := NewConfig()
	if err != nil {
		log.Fatalln("Couldn't parse environment variables", err)
	}

	s, err := api.NewAPI(config.HTTPAddr, config.DBName, config.MongoURI)
	if err != nil {
		log.Fatalln("Couldn't create a new server")
	}

	done := make(chan bool, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	go func() {
		<-sigChan
		s.Logger.Info("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s.HTTPServer.SetKeepAlivesEnabled(false)
		if err = s.HTTPServer.Shutdown(ctx); err != nil {
			s.Logger.Fatalf("Couldn't gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	if err = s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.Logger.Fatalf("Couldn't listen on %v: %v\n", config.HTTPAddr, err)
	}

	<-done
	s.Logger.Info("Server stopped")

}
