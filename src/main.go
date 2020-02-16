package main

import (
	"log"
)

func main() {
	a, err := InitializeServer()
	if err != nil {
		log.Fatalf("Unable to initialize app.go: %+v", err)
	}
	a.Run()
}
