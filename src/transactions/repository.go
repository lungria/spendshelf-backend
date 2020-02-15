package transactions

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	TransactionsCollection = "transactions"
)

// Repository implements by methods which define in Repository interface
type Repository struct {
	logger     *zap.SugaredLogger
	collection *mongo.Collection
}

// NewRepository creates a new instance of Repository
func NewRepository(db *mongo.Database, logger *zap.SugaredLogger) *Repository {
	return &Repository{
		logger:     logger,
		collection: db.Collection(TransactionsCollection),
	}
}

// FindAllUncategorized returns all uncategorized transactions
func (repo *Repository) FindAllUncategorized(ctx context.Context) ([]Transaction, error) {
	var transactions []Transaction
	optProjections := options.Find().SetProjection(bson.M{"_id": 1, "time": 1, "description": 1, "category": 1, "amount": 1, "balance": 1, "bank": 1})
	filter := bson.M{"$or": bson.A{bson.M{"category": bson.M{"$exists": false}}, bson.M{"category": nil}}}
	cur, err := repo.collection.Find(ctx, filter, optProjections)

	if err != nil {
		errMsg := "unable to received transactions without category"
		repo.logger.Errorw(errMsg, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return nil, errors.New(errMsg)
	}

	return transactionsDecoder(ctx, cur, transactions), nil
}

// FindAll returns all transactions
func (repo *Repository) FindAll(ctx context.Context) ([]Transaction, error) {
	var transactions []Transaction
	cur, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		errMsg := "unable to received transactions with all categories"
		repo.logger.Errorw(errMsg, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return nil, errors.New(errMsg)
	}

	return transactionsDecoder(ctx, cur, transactions), nil
}

// FindByCategoryID returns all transactions which relate with specify category
func (repo *Repository) FindByCategoryID(ctx context.Context, categoryID primitive.ObjectID) ([]Transaction, error) {
	var transactions []Transaction
	cur, err := repo.collection.Find(ctx, bson.M{"category._id": categoryID})
	if err != nil {
		errMsg := "unable to received transactions with category"
		repo.logger.Errorw(errMsg, "Category", categoryID, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return nil, errors.New(errMsg)
	}

	return transactionsDecoder(ctx, cur, transactions), nil
}

// FindAllCategorized returns all categorized transactions
func (repo *Repository) FindAllCategorized(ctx context.Context) ([]Transaction, error) {
	var transactions []Transaction

	cur, err := repo.collection.Find(ctx, bson.M{"$and": bson.A{bson.M{"category": bson.M{"$exists": true}}, bson.M{"category": bson.M{"$ne": nil}}}})
	if err != nil {
		errMsg := "unable to received transactions with category"
		repo.logger.Errorw(errMsg, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return nil, errors.New(errMsg)
	}

	return transactionsDecoder(ctx, cur, transactions), nil
}

// UpdateCategory changes the category for appropriate transaction
func (repo *Repository) UpdateCategory(ctx context.Context, transactionID primitive.ObjectID, category primitive.ObjectID) (int64, error) {
	txn, err := repo.collection.UpdateOne(ctx, bson.M{"_id": transactionID}, bson.M{"$set": bson.M{"categoryID": category}})
	if err != nil {
		errMsg := "unable to update transaction"
		repo.logger.Errorw(errMsg, "TransactionID", transactionID, "categoryID", category, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return txn.ModifiedCount, errors.New(errMsg)
	}
	return txn.ModifiedCount, nil
}

// InsertMany inserts slice of transactions to transactions collection
func (repo *Repository) InsertMany(ctx context.Context, txns []Transaction) error {
	txnInterface := make([]interface{}, len(txns))
	for i := 0; i < len(txns); i++ {
		txnInterface[i] = txns[i]
	}

	_, err := repo.collection.InsertMany(ctx, txnInterface)
	if err != nil {
		errMsg := "unable to insert transaction"
		repo.logger.Errorw(errMsg, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return errors.New(errMsg)
	}

	return nil
}

// Insert inserts slice of transactions to transactions collection
func (repo *Repository) Insert(ctx context.Context, t Transaction) error {
	_, err := repo.collection.InsertOne(ctx, t)
	if err != nil {
		errMsg := "unable to insert transaction"
		repo.logger.Errorw(errMsg, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return errors.New(errMsg)
	}

	return nil
}

func transactionsDecoder(ctx context.Context, cursor *mongo.Cursor, transactions []Transaction) []Transaction {
	for cursor.Next(ctx) {
		var t Transaction
		cursor.Decode(&t)
		transactions = append(transactions, t)
	}
	return transactions
}
