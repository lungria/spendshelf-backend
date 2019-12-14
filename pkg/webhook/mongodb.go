package webhook

import (
	"github.com/lungria/spendshelf-backend/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

const (
	transactionsCollection = "Transactions"
)

// SpendShelf is struct for SpendShelf database in MongoDB
type SpendShelf struct {
	*db.MongoDB
}

// NewConnection create a new connection to specified database
func NewConnection() (*SpendShelf, error) {
	dbName := os.Getenv("SPEND_SHELF_DB")
	mongoClient, err := db.Connect(dbName)
	if err != nil {
		return nil, err
	}
	return &SpendShelf{mongoClient}, nil
}

// GetTransactionByID fetch one Transaction by transactionId from MongoDB
func (s SpendShelf) GetTransactionByID(transactionID string) (Transaction, error) {
	defer s.MongoDB.Cancel()
	var t Transaction
	collection := s.Database.Collection(transactionsCollection)
	err := collection.FindOne(s.Context, bson.M{"id": transactionID}).Decode(t)
	if err != nil {
		return t, err
	}
	return t, err
}

// GetAllTransactions fetch all Transaction by accountId from MongoDB
func (s SpendShelf) GetAllTransactions(accountID string) ([]Transaction, error) {
	defer s.MongoDB.Cancel()
	var transactions []Transaction
	collection := s.Database.Collection(transactionsCollection)
	cur, err := collection.Find(s.MongoDB.Context, bson.M{"account_id": accountID})
	if err != nil {
		return nil, nil
	}
	for cur.Next(s.MongoDB.Context) {
		var t Transaction
		cur.Decode(&t)
		transactions = append(transactions, t)
	}
	return transactions, nil
}

// SaveOneTransaction save one Transaction to MongoDB
func (s SpendShelf) SaveOneTransaction(transaction *Transaction) error {
	defer s.MongoDB.Cancel()
	collection := s.Database.Collection(transactionsCollection)
	_, err := collection.InsertOne(s.MongoDB.Context, transaction)
	if err != nil {
		return err
	}
	return nil
}
