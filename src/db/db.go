package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewDatabase is create a new database connection
func NewDatabase(dbname, URI string) (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(URI), options.Client().SetMaxPoolSize(50))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*4999)
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}
	database := client.Database(dbname)
	err = initCappedCollection(database, transactionsCollection)
	if err != nil {
		return nil, err
	}
	return database, nil
}

func initCappedCollection(database *mongo.Database, collName string) error {
	var cursor []bson.M
	result, _ := database.ListCollections(context.Background(), bson.D{})
	if err := result.All(context.Background(), &cursor); err != nil {
		return err
	}
	cmd := bson.D{{"convertToCapped", collName}, {"size", 5000000}}
	if len(cursor) == 0 {
		database.RunCommand(context.Background(), cmd)
		return nil
	}
	for coll := 0; coll < len(cursor); coll++ {
		if cursor[coll]["name"] == collName {
			database.RunCommand(context.Background(), cmd)
			break
		}
	}
	return nil
}
