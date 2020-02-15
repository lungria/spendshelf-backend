package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config interface {
	GetDBName() string
	GetMongoURI() string
}

// NewDatabase is create a new database connection
func NewDatabase(config Config) (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.GetMongoURI()), options.Client().SetMaxPoolSize(50))
	if err != nil {
		return nil, err
	}
	databaseCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := client.Connect(databaseCtx); err != nil {
		return nil, err
	}
	database := client.Database(config.GetDBName())
	return database, nil
}
