package webhooks

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/lungria/spendshelf-backend/src/models"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo"
)

const transactionsCollection = "transactions"

const MonoBankName = "Mono Bank"

// Repository defines method which inserts the transaction from monoAPI
type Repository interface {
	InsertOneHook(transaction *models.WebHook) error
}

// WebHookRepository implements by methods which save the transaction to transactions collection
type WebHookRepository struct {
	collection *mongo.Collection
	logger     *zap.SugaredLogger
	context    context.Context
}

// NewWebHookRepository create a new repository
func NewWebHookRepository(ctx context.Context, db *mongo.Database, logger *zap.SugaredLogger) (*WebHookRepository, error) {
	if db == nil {
		return nil, errors.New("DB must not be nil")
	}
	if logger == nil {
		return nil, errors.New("logger must not be nil")
	}
	return &WebHookRepository{
		collection: db.Collection(transactionsCollection),
		logger:     logger,
		context:    ctx,
	}, nil
}

// InsertOneHook insert one Transaction to MongoDB
func (repo *WebHookRepository) InsertOneHook(webhook *models.WebHook) error {
	ctx, cancel := context.WithCancel(repo.context)
	defer cancel()

	_, err := repo.collection.InsertOne(ctx, repo.txnFromHook(webhook))
	if err != nil {
		repo.logger.Errorw("InsertOneHook failed", "Database", repo.collection.Database().Name(), "Collection", repo.collection.Name(), "webHook", webhook, "Error", err)
		return errors.New("save webHook to database failed")
	}
	return nil
}

func (repo *WebHookRepository) txnFromHook(webhook *models.WebHook) models.Transaction {
	dest := models.Transaction{}

	dateTime := time.Unix(int64(webhook.StatementItem.Time), 0).UTC()

	dest.ID = primitive.NewObjectID()
	dest.Amount = webhook.StatementItem.Amount
	dest.Balance = webhook.StatementItem.Balance
	dest.Description = webhook.StatementItem.Description
	dest.Time = dateTime
	dest.Bank = MonoBankName
	dest.BankTransaction = *webhook
	return dest
}
