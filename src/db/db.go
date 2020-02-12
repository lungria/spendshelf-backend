package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewDatabase is create a new database connection
func NewDatabase(ctx context.Context, dbname, URI string) (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI), options.Client().SetMaxPoolSize(50))
	if err != nil {
		return nil, err
	}
	databaseCtx, cancel := context.WithTimeout(ctx, time.Millisecond*4999)
	defer cancel()
	if err := client.Connect(databaseCtx); err != nil {
		return nil, err
	}
	database := client.Database(dbname)
	return database, nil
}
