package main

import (
	"context"
	"flag"
	"github.com/lungria/spendshelf-backend/src/pkg/webhook"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var err error
	addr := flag.String("addr", ":8080", "HTTP address of server")
	flag.Parse()

	s := webhook.NewServer(*addr, "", "")


	go func() {
		if err = s.HTTPServer.ListenAndServe(); err != nil {
			log.Fatalf("ListenAndServe failed %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = s.HTTPServer.Shutdown(ctx); err != nil {
		log.Fatalln("Server shutdown: ", err)
	}

	log.Println("Server exiting")
}
