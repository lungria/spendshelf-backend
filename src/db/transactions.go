package db

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/lungria/spendshelf-backend/src/models"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	transactionsCollection = "transactions"
)

// TransactionsRepository define all methods which do some work with database
type TransactionsRepository interface {
	GetTransactionByID(transactionID string) (models.Transaction, error)
	GetAllTransactions(accountID string) ([]models.Transaction, error)
	SaveOneTransaction(transaction *models.Transaction) error
}

type TransactionsMongoDbRepository struct {
	collection *mongo.Collection
	logger     *zap.SugaredLogger
}

func NewTransactionsMongoDbRepository(db *mongo.Database, logger *zap.SugaredLogger) (*TransactionsMongoDbRepository, error) {
	if db == nil {
		return nil, errors.New("Db must not be nil")
	}
	if logger == nil {
		return nil, errors.New("Logger must not be nil")
	}
	return &TransactionsMongoDbRepository{
		collection: db.Collection(transactionsCollection),
		logger:     logger,
	}, nil
}

// GetTransactionByID fetch one Transaction by transactionId from MongoDB
func (repo *TransactionsMongoDbRepository) GetTransactionByID(transactionID string) (models.Transaction, error) {
	var t models.Transaction
	err := repo.collection.FindOne(context.Background(), bson.M{"id": transactionID}).Decode(t)
	if err != nil {
		repo.logger.Errorw("GetTransactionByID failed", "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Transaction ID", transactionID, "Error", err)
		return t, errors.New("retrieve transaction failed")
	}
	return t, err
}

// GetAllTransactions fetch all Transaction by accountId from MongoDB
func (repo *TransactionsMongoDbRepository) GetAllTransactions(accountID string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	cur, err := repo.collection.Find(context.Background(), bson.M{"account_id": accountID})
	if err != nil {
		repo.logger.Errorw("GetAllTransactions failed", "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Account ID", accountID, "Error", err)
		return nil, errors.New("retrieve transactions failed")
	}
	for cur.Next(context.Background()) {
		var t models.Transaction
		cur.Decode(&t)
		transactions = append(transactions, t)
	}
	return transactions, nil
}

// SaveOneTransaction save one Transaction to MongoDB
func (repo *TransactionsMongoDbRepository) SaveOneTransaction(transaction *models.Transaction) error {
	_, err := repo.collection.InsertOne(context.Background(), transaction)
	if err != nil {
		repo.logger.Errorw("SaveOneTransaction failed", "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "transaction", transaction, "Error", err)
		return errors.New("save transaction to database failed")
	}
	return nil
}
