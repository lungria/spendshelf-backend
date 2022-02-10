package main

import (
	"github.com/lungria/spendshelf-backend/app"
	"github.com/rs/zerolog/log"
)

func main() {
	app, err := app.Initialize()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize app")
	}

	app.Run()
}
