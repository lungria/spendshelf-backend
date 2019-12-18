package webhook

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	transactionsCollection = "Transactions"
)

// GetTransactionByID fetch one Transaction by transactionId from MongoDB
func (s *Server) GetTransactionByID(transactionID string) (Transaction, error) {
	var t Transaction
	collection := s.MongoDB.Collection(transactionsCollection)
	err := collection.FindOne(context.Background(), bson.M{"id": transactionID}).Decode(t)
	if err != nil {
		s.Logger.Errorw("GetTransactionByID failed", "Database", s.MongoDB.Name(), "Collection", transactionsCollection, "Transaction ID", transactionID, "Error", err)
		return t, errors.New("retrieve transaction failed")
	}
	return t, err
}

// GetAllTransactions fetch all Transaction by accountId from MongoDB
func (s *Server) GetAllTransactions(accountID string) ([]Transaction, error) {
	var transactions []Transaction
	collection := s.MongoDB.Collection(transactionsCollection)
	cur, err := collection.Find(context.Background(), bson.M{"account_id": accountID})
	if err != nil {
		s.Logger.Errorw("GetAllTransactions failed", "Database", s.MongoDB.Name(), "Collection", transactionsCollection, "Account ID",accountID, "Error", err)
		return nil, errors.New("retrieve transactions failed")
	}
	for cur.Next(context.Background()) {
		var t Transaction
		cur.Decode(&t)
		transactions = append(transactions, t)
	}
	return transactions, nil
}

// SaveOneTransaction save one Transaction to MongoDB
func (s *Server) SaveOneTransaction(transaction *Transaction) error {
	collection := s.MongoDB.Collection(transactionsCollection)
	_, err := collection.InsertOne(context.Background(), transaction)
	if err != nil {
		s.Logger.Errorw("SaveOneTransaction failed", "Database", s.MongoDB.Name(), "Collection", transactionsCollection, "transaction",transaction, "Error", err)
		return errors.New("save transaction to database failed")
	}
	return nil
}
