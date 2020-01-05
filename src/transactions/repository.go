package transactions

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	transactionsCollection = "transactions"
)

// Repository defines methods which do some work with database
type Repository interface {
	FindAll() ([]Transaction, error)
	FindAllCategorized() ([]Transaction, error)
	FindAllUncategorized() ([]Transaction, error)
	FindAllByCategory(category string) ([]Transaction, error)
	UpdateCategory(id string, category string) error
}

// TransactionRepository implements by methods which do some work with transaction collection
type TransactionRepository struct {
	logger     *zap.SugaredLogger
	collection *mongo.Collection
}

// NewTransactionRepository creates a new instance of TransactionRepository
func NewTransactionRepository(db *mongo.Database, logger *zap.SugaredLogger) (*TransactionRepository, error) {
	if db == nil {
		return nil, errors.New("database must not be nil")
	}
	if logger == nil {
		return nil, errors.New("logger must not be nil")
	}

	return &TransactionRepository{
		logger:     logger,
		collection: db.Collection(transactionsCollection),
	}, nil
}

// FindAllUncategorized returns all uncategorized transactions
func (repo *TransactionRepository) FindAllUncategorized() ([]Transaction, error) {
	var transactions []Transaction
	ctx := context.Background()
	cur, err := repo.collection.Find(ctx, bson.M{"$or": bson.A{bson.M{"category": bson.M{"$exists": false}}, bson.M{"category": nil}}})
	if err != nil {
		errMsg := "unable to received transactions without category"
		repo.logger.Errorw(errMsg, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return nil, errors.New(errMsg)
	}

	return transactionsDecoder(ctx, cur, transactions), nil
}

// FindAll returns all transactions
func (repo *TransactionRepository) FindAll() ([]Transaction, error) {
	var transactions []Transaction
	ctx := context.Background()
	cur, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		errMsg := "unable to received transactions with all categories"
		repo.logger.Errorw(errMsg, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return nil, errors.New(errMsg)
	}

	return transactionsDecoder(ctx, cur, transactions), nil
}

// FindAllByCategory returns all transactions which relate with specify category
func (repo *TransactionRepository) FindAllByCategory(category string) ([]Transaction, error) {
	var transactions []Transaction
	ctx := context.Background()
	cur, err := repo.collection.Find(ctx, bson.M{"category": category})
	if err != nil {
		errMsg := "unable to received transactions with category"
		repo.logger.Errorw(errMsg, "Category", category, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return nil, errors.New(errMsg)
	}

	return transactionsDecoder(ctx, cur, transactions), nil
}

// FindAllCategorized returns all categorized transactions
func (repo *TransactionRepository) FindAllCategorized() ([]Transaction, error) {
	var transactions []Transaction
	ctx := context.Background()
	cur, err := repo.collection.Find(ctx, bson.M{"$and": bson.A{bson.M{"category": bson.M{"$exists": true}}, bson.M{"category": bson.M{"$ne": nil}}}})
	if err != nil {
		errMsg := "unable to received transactions with category"
		repo.logger.Errorw(errMsg, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return nil, errors.New(errMsg)
	}

	return transactionsDecoder(ctx, cur, transactions), nil
}

// UpdateCategory changes the category for appropriate transaction
func (repo *TransactionRepository) UpdateCategory(id string, category string) error {
	transactionID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		errMsg := "unable to convert transaction id to ObjectID"
		repo.logger.Errorw(errMsg, "TransactionID", id, "Category", category, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return errors.New(errMsg)
	}

	_, err = repo.collection.UpdateOne(context.Background(), bson.M{"_id": transactionID}, bson.M{"$set": bson.M{"category": category}})
	if err != nil {
		errMsg := "unable to update transaction"
		repo.logger.Errorw(errMsg, "TransactionID", id, "Category", category, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
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
