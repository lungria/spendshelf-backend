package main

import (
	"log"
	"net/http"

	"github.com/lungria/spendshelf-backend/src/config"
	"github.com/lungria/spendshelf-backend/src/db"
	"github.com/lungria/spendshelf-backend/src/transactions"
	"go.uber.org/zap"

	"github.com/lungria/spendshelf-backend/src/sync_mono"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln(err)
	}
	database, err := db.NewDatabase(conf.DBName, conf.MongoURI)
	if err != nil {
		log.Fatalln(err)
	}
	repo, err := transactions.NewTransactionRepository(database, logger.Sugar())
	if err != nil {
		log.Fatalln(err)
	}
	m, err := sync_mono.NewMonoSync("", repo, logger.Sugar())
	if err != nil {
		log.Fatalln(err)
	}
	c := sync_mono.NewClient(m)

	http.Handle("/sync", c)

	log.Println("listening...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}

}
