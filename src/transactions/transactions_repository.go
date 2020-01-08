package transactions

import (
	"context"
	"errors"

	"github.com/lungria/spendshelf-backend/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	transactionsCollection = "transactions"
)

// Repository defines methods which find the transactions and update the category
type Repository interface {
	FindAll() ([]models.Transaction, error)
	FindAllCategorized() ([]models.Transaction, error)
	FindAllUncategorized() ([]models.Transaction, error)
	FindAllByCategoryID(categoryID primitive.ObjectID) ([]models.Transaction, error)
	UpdateCategory(transactionID primitive.ObjectID, categoryID primitive.ObjectID) (int64, error)
}

// TransactionRepository implements by methods which define in Repository interface
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
func (repo *TransactionRepository) FindAllUncategorized() ([]models.Transaction, error) {
	var transactions []models.Transaction
	ctx := context.Background()
	cur, err := repo.collection.Find(ctx, bson.M{"$or": bson.A{bson.M{"category_id": bson.M{"$exists": false}}, bson.M{"category_id": nil}}})
	if err != nil {
		errMsg := "unable to received transactions without category"
		repo.logger.Errorw(errMsg, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return nil, errors.New(errMsg)
	}

	return transactionsDecoder(ctx, cur, transactions), nil
}

// FindAll returns all transactions
func (repo *TransactionRepository) FindAll() ([]models.Transaction, error) {
	var transactions []models.Transaction
	ctx := context.Background()
	cur, err := repo.collection.Find(ctx, bson.M{})
	if err != nil {
		errMsg := "unable to received transactions with all categories"
		repo.logger.Errorw(errMsg, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return nil, errors.New(errMsg)
	}

	return transactionsDecoder(ctx, cur, transactions), nil
}

// FindAllByCategoryID returns all transactions which relate with specify category
func (repo *TransactionRepository) FindAllByCategoryID(categoryID primitive.ObjectID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	ctx := context.Background()
	cur, err := repo.collection.Find(ctx, bson.M{"category_id": categoryID})
	if err != nil {
		errMsg := "unable to received transactions with category"
		repo.logger.Errorw(errMsg, "CategoryID", categoryID, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return nil, errors.New(errMsg)
	}

	return transactionsDecoder(ctx, cur, transactions), nil
}

// FindAllCategorized returns all categorized transactions
func (repo *TransactionRepository) FindAllCategorized() ([]models.Transaction, error) {
	var transactions []models.Transaction
	ctx := context.Background()
	cur, err := repo.collection.Find(ctx, bson.M{"$and": bson.A{bson.M{"category_id": bson.M{"$exists": true}}, bson.M{"category_id": bson.M{"$ne": nil}}}})
	if err != nil {
		errMsg := "unable to received transactions with category"
		repo.logger.Errorw(errMsg, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return nil, errors.New(errMsg)
	}

	return transactionsDecoder(ctx, cur, transactions), nil
}

// UpdateCategory changes the category for appropriate transaction
func (repo *TransactionRepository) UpdateCategory(transactionID primitive.ObjectID, categoryID primitive.ObjectID) (int64, error) {
	txn, err := repo.collection.UpdateOne(context.Background(), bson.M{"_id": transactionID}, bson.M{"$set": bson.M{"category_id": categoryID}})
	if err != nil {
		errMsg := "unable to update transaction"
		repo.logger.Errorw(errMsg, "TransactionID", transactionID, "CategoryID", categoryID, "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Error", err)
		return txn.ModifiedCount, errors.New(errMsg)
	}
	return txn.ModifiedCount, nil
}

func transactionsDecoder(ctx context.Context, cursor *mongo.Cursor, transactions []models.Transaction) []models.Transaction {
	for cursor.Next(ctx) {
		var t models.Transaction
		cursor.Decode(&t)
		transactions = append(transactions, t)
	}
	return transactions
}
