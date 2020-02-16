package main

import (
	"log"
)

func main() {
	server, err := InitializeServer()
	if err != nil {
		log.Fatalf("Unable to initialize app.go: %+v", err)
	}
	server.Run()
}
