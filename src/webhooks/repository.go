package webhooks

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	transactionsCollection = "transactions"
)

// Repository define all methods which do some work with database
type Repository interface {
	GetHookByID(transactionID string) (WebHook, error)
	GetAllHooks(accountID string) ([]WebHook, error)
	SaveOneHook(transaction *WebHook) error
}

// WebHookRepository is contain the connection and logger
type WebHookRepository struct {
	collection *mongo.Collection
	logger     *zap.SugaredLogger
}

// NewWebHookRepository create a new repository
func NewWebHookRepository(db *mongo.Database, logger *zap.SugaredLogger) (*WebHookRepository, error) {
	if db == nil {
		return nil, errors.New("DB must not be nil")
	}
	if logger == nil {
		return nil, errors.New("logger must not be nil")
	}
	return &WebHookRepository{
		collection: db.Collection(transactionsCollection),
		logger:     logger,
	}, nil
}

// GetHookByID fetch one Transaction by transactionId from MongoDB
func (repo *WebHookRepository) GetHookByID(transactionID string) (WebHook, error) {
	var w WebHook
	err := repo.collection.FindOne(context.Background(), bson.M{"id": transactionID}).Decode(w)
	if err != nil {
		repo.logger.Errorw("GetHookByID failed", "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Transaction ID", transactionID, "Error", err)
		return w, errors.New("retrieve transaction failed")
	}
	return w, err
}

// GetAllHooks fetch all Transaction by accountId from MongoDB
func (repo *WebHookRepository) GetAllHooks(accountID string) ([]WebHook, error) {
	var webhooks []WebHook
	cur, err := repo.collection.Find(context.Background(), bson.M{"account_id": accountID})
	if err != nil {
		repo.logger.Errorw("GetAllHooks failed", "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "Account ID", accountID, "Error", err)
		return nil, errors.New("retrieve webHooks failed")
	}
	for cur.Next(context.Background()) {
		var w WebHook
		cur.Decode(&w)
		webhooks = append(webhooks, w)
	}
	return webhooks, nil
}

// SaveOneHook save one Transaction to MongoDB
func (repo *WebHookRepository) SaveOneHook(webhook *WebHook) error {
	_, err := repo.collection.InsertOne(context.Background(), webhook)
	if err != nil {
		repo.logger.Errorw("SaveOneHook failed", "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "webHook", webhook, "Error", err)
		return errors.New("save webHook to database failed")
	}
	return nil
}
