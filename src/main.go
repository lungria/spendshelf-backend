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
	s, err := InitializeServer()
	if err != nil {
		log.Fatal("Unable to initialize server")
	}
	done := make(chan bool, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	go func() {
		<-sigChan
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		s.SetKeepAlivesEnabled(false)
		if err = s.Shutdown(ctx); err != nil {
			log.Fatalf("Couldn't gracefully shutdown the server: %+v\n", err)
		}
		close(done)
	}()

	if err = s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Couldn't listen: %+v\n", err)
	}

	<-done
}
