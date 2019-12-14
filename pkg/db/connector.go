package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type MongoDB struct {
	Database   *mongo.Database
	Cancel     context.CancelFunc
	Context    context.Context
	Collection string
}

func Connect(dbName string) (*MongoDB, error) {
	mongoURI := os.Getenv("MONGO_URI")
	db, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return &MongoDB{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	if err := db.Connect(ctx); err != nil {
		return &MongoDB{}, err
	}
	database := db.Database(dbName)
	return &MongoDB{
		Database:   database,
		Cancel:     cancel,
		Context:    ctx,
	}, nil
}
