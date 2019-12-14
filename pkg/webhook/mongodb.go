package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type mongoDB struct {
	Database              *mongo.Database
	Cancel                context.CancelFunc
	Context               context.Context
	SpendShelfDB          string
	TransactionCollection string
}

// NewConnection create a new connection to mongoDB
func NewConnection() (Repository, error) {
	mongoURI := os.Getenv("MONGO_URI")
	spendShelfDB := os.Getenv("SPEND_SHELF_DB")
	transactionsColl := os.Getenv("TRANSACTIONS_COLLECTION")

	db, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	if err := db.Connect(ctx); err != nil {
		return nil, err
	}
	database := db.Database(spendShelfDB)
	return &mongoDB{
		Database: database,
		Cancel:   cancel,
		Context:  ctx,
		TransactionCollection: transactionsColl,
	}, nil
}

// GetTransactionByID fetch one transaction by transactionId from MongoDB
func (m *mongoDB) GetTransactionByID(transactionID string) (Transaction, error) {
	defer m.Cancel()
	var t Transaction
	collection := m.Database.Collection(m.TransactionCollection)
	err := collection.FindOne(m.Context, bson.M{"id": transactionID}).Decode(t)
	if err != nil {
		return t, err
	}
	return t, err
}

// GetAllTransactions fetch all transaction by accountId from MongoDB
func (m *mongoDB) GetAllTransactions(accountID string) ([]Transaction, error) {
	defer m.Cancel()
	var transactions []Transaction
	collection := m.Database.Collection(m.TransactionCollection)
	cur, err := collection.Find(m.Context, bson.M{"account_id": accountID})
	if err != nil {
		return nil, nil
	}
	for cur.Next(m.Context) {
		var t Transaction
		cur.Decode(&t)
		transactions = append(transactions, t)
	}
	return transactions, nil
}

// SaveOneTransaction save one transaction to MongoDB
func (m *mongoDB) SaveOneTransaction(transaction *Transaction) error {
	defer m.Cancel()
	collection := m.Database.Collection(m.TransactionCollection)
	_, err := collection.InsertOne(m.Context, transaction)
	if err != nil {
		return err
	}
	return nil
}
