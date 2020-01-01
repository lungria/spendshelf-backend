package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	services, err := InitializeServer()
	if err != nil {
		log.Fatal("Unable to initialize server")
	}
	done := make(chan bool, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	go func() {
		<-sigChan
		services.Logger.Info("Shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		services.Server.SetKeepAlivesEnabled(false)
		if err = services.Server.Shutdown(ctx); err != nil {
			services.Logger.Fatalf("Couldn't gracefully shutdown the server: %+v\n", err)
		}
		close(done)
	}()

	if err = services.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		services.Logger.Fatalf("Couldn't listen: %+v\n", err)
	}
	services.Logger.Info("Server shutdown")
	<-done
}
