package db

import (
	"context"
	"errors"

	"github.com/lungria/spendshelf-backend/src/models"
	"go.mongodb.org/mongo-driver/bson"
)

// GetTransactionByID fetch one Transaction by transactionId from MongoDB
func (d *Database) GetTransactionByID(transactionID string) (models.Transaction, error) {
	var t models.Transaction
	collection := d.MongoDB.Collection(transactionsCollection)
	err := collection.FindOne(context.Background(), bson.M{"id": transactionID}).Decode(t)
	if err != nil {
		d.logger.Errorw("GetTransactionByID failed", "Database", d.MongoDB.Name(), "Collection", transactionsCollection, "Transaction ID", transactionID, "Error", err)
		return t, errors.New("retrieve transaction failed")
	}
	return t, err
}

// GetAllTransactions fetch all Transaction by accountId from MongoDB
func (d *Database) GetAllTransactions(accountID string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	collection := d.MongoDB.Collection(transactionsCollection)
	cur, err := collection.Find(context.Background(), bson.M{"account_id": accountID})
	if err != nil {
		d.logger.Errorw("GetAllTransactions failed", "Database", d.MongoDB.Name(), "Collection", transactionsCollection, "Account ID", accountID, "Error", err)
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
func (d *Database) SaveOneTransaction(transaction *models.Transaction) error {
	collection := d.MongoDB.Collection(transactionsCollection)
	_, err := collection.InsertOne(context.Background(), transaction)
	if err != nil {
		d.logger.Errorw("SaveOneTransaction failed", "Database", d.MongoDB.Name(), "Collection", transactionsCollection, "transaction", transaction, "Error", err)
		return errors.New("save transaction to database failed")
	}
	return nil
}
