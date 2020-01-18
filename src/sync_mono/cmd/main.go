package main

import (
	"log"
	"time"

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
	m, err := sync_mono.NewSync("uS9e4JRPxCykO7yz53cFUOQsIQ1BwzX5Est0TizsCNQI", repo)
	if err != nil {
		log.Fatalln(err)
	}
	tm := time.Unix(1574246567, 0).UTC()
	m.Transactions(tm)

}
