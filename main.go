package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	addr := flag.String("addr", ":80", "HTTP address of server")
	flag.Parse()

	s := NewServer(*addr)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatalf("ListenAndServe failed %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutdown server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln("Server shutdown: ", err)
	}

	log.Println("Server exiting")
}
