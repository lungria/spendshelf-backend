package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/lungria/spendshelf-backend/src/transactions"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

type Config interface {
	GetDBName() string
	GetMongoURI() string
}

func NewDbConnection(cfg Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	m, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.GetMongoURI()))
	if err != nil {
		return nil, err
	}
	db := m.Database(cfg.GetDBName())

	// index for transaction deduplication
	opts := &options.IndexOptions{}
	opts.SetUnique(true)
	index := mongo.IndexModel{
		Keys:    bson.M{"time": -1, "amount": 1},
		Options: opts,
	}
	_, err = db.Collection(transactions.CollectionName).Indexes().CreateOne(context.TODO(), index)
	if err != nil {
		return nil, err
	}

	return db, nil
}
