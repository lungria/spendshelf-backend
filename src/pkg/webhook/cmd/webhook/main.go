package main

import (
	"context"
	"flag"
	"github.com/lungria/spendshelf-backend/src/pkg/webhook"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var err error

	addr := flag.String("addr", ":8080", "HTTP address of server")
	flag.Parse()

	s := webhook.NewServer(*addr, "SpendShelf", "mongodb://root:toor@localhost:27017")

	done := make(chan bool, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

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
		s.Logger.Fatalf("Couldn't listen on %v: %v\n", &addr, err)
	}

	<- done
	s.Logger.Info("Server stopped")

}
