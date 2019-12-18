package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func Connect(dbName, URI string) (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI), options.Client().SetMaxPoolSize(50))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*4999)
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}
	database := client.Database(dbName)
	return database, nil
}
