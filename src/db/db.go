package db

import (
	"context"
	"time"

	"github.com/lungria/spendshelf-backend/src/transactions"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

type Config interface {
	GetDBName() string
	GetMongoURI() string
}

func NewDbConnection(cfg Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	m, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.GetMongoURI()))
	if err != nil {
		return nil, err
	}
	db := m.Database(cfg.GetDBName())

	err = buildIndexes(ctx, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type index struct {
	Name string `bson:"name"`
}

func buildIndexes(ctx context.Context, db *mongo.Database) error {
	col := db.Collection(transactions.CollectionName)

	indexes, err := getExisting(ctx, col)
	if err != nil {
		return err
	}

	timeAmountUnique := "time_amount_unique"
	_, ok := indexes[timeAmountUnique]
	if ok {
		return nil
	}

	// index for transaction deduplication
	opts := &options.IndexOptions{}
	opts.SetUnique(true)
	opts.SetName(timeAmountUnique)
	index := mongo.IndexModel{
		Keys:    bson.M{"time": -1, "amount": 1},
		Options: opts,
	}

	_, err = col.Indexes().CreateOne(ctx, index)
	if err != nil {
		return err
	}

	return nil
}

func getExisting(ctx context.Context, col *mongo.Collection) (map[string]struct{}, error) {
	cur, err := col.Indexes().List(ctx)
	if err != nil {
		return nil, err
	}

	var indexes []index
	err = cur.All(ctx, &indexes)
	if err != nil {
		return nil, err
	}
	var namesSet = make(map[string]struct{}, len(indexes))
	for _, v := range indexes {
		namesSet[v.Name] = struct{}{}
	}
	return namesSet, err
}
