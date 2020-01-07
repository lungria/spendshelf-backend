package webhooks

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/lungria/spendshelf-backend/src/models"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	transactionsCollection = "transactions"
)

// Repository define all methods which do some work with database
type Repository interface {
	SaveOneHook(transaction *WebHook) error
}

// WebHookRepository implements by methods which save the transaction to transactions collection
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

// SaveOneHook save one Transaction to MongoDB
func (repo *WebHookRepository) SaveOneHook(webhook *WebHook) error {
	_, err := repo.collection.InsertOne(context.Background(), repo.convert(webhook))
	if err != nil {
		repo.logger.Errorw("SaveOneHook failed", "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "webHook", webhook, "Error", err)
		return errors.New("save webHook to database failed")
	}
	return nil
}

func (repo *WebHookRepository) convert(webhook *WebHook) models.Transaction {
	dest := models.Transaction{}

	dest.ID = primitive.NewObjectID()
	dest.Amount = webhook.StatementItem.Amount
	dest.Balance = webhook.StatementItem.Balance
	dest.Description = webhook.StatementItem.Description
	dest.Time = webhook.StatementItem.Time
	return dest
}
