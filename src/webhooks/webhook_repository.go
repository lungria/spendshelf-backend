package webhooks

import (
	"context"
	"errors"
	"time"

	"github.com/lungria/spendshelf-backend/src/transactions"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/lungria/spendshelf-backend/src/models"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo"
)

// MonoBankName name of mono bank, which using in transactions
const MonoBankName = "Mono Bank"

// Repository defines method which inserts the transaction from monoAPI
type Repository interface {
	InsertOneHook(ctx context.Context, transaction *models.WebHook) error
}

// WebHookRepository implements by methods which save the transaction to transactions collection
type WebHookRepository struct {
	collection *mongo.Collection
	logger     *zap.SugaredLogger
}

// NewWebHookRepository create a new repository
func NewWebHookRepository(db *mongo.Database, logger *zap.SugaredLogger) *WebHookRepository {
	return &WebHookRepository{
		collection: db.Collection(transactions.TransactionsCollection),
		logger:     logger,
	}
}

// InsertOneHook insert one Transaction to MongoDB
func (repo *WebHookRepository) InsertOneHook(ctx context.Context, webhook *models.WebHook) error {
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
