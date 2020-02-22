package transactions

import (
	"context"
	"errors"
	"time"

	"github.com/lungria/spendshelf-backend/src/categories"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"

	"go.uber.org/zap"
)

// Transaction represents a model of transactions in database
type Transaction struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Time time.Time          `json:"time" bson:"time"`
	// LocalDate describes transaction date in local timezone of the user.
	// Used for reporting purposes, so we can aggregate using it.
	// Timezone is stored as UTC, but must not be considered or used for any purpose.
	LocalDate   time.Time `json:"-" bson:"localDate"`
	Description string    `json:"description" bson:"description"`
	CategoryID  uint8     `json:"categoryId,omitempty" bson:"categoryId,omitempty"`
	Amount      int32     `json:"amount" json:"amount"`
	BankId      Bank      `json:"bankId" json:"bankId"`
}

type Bank uint8

const (
	CollectionName = "transactions"

	Mono Bank = iota + 1
	Privat
)

// Repository implements by methods which define in Repository interface
type Repository struct {
	logger     *zap.SugaredLogger
	db         *mongo.Collection
	categories *categories.Repository
}

// NewRepository creates a new instance of Repository
func NewRepository(mongo *mongo.Database, logger *zap.SugaredLogger, categories *categories.Repository) *Repository {
	return &Repository{
		logger:     logger,
		db:         mongo.Collection(CollectionName),
		categories: categories,
	}
}

// Insert transaction to database.
func (s *Repository) Insert(ctx context.Context, t *Transaction) error {
	y, m, d := t.Time.Date()
	t.LocalDate = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	_, err := s.db.InsertOne(ctx, t)
	return err
}

// Insert multiple transaction to database.
func (s *Repository) InsertMany(ctx context.Context, t []Transaction) error {
	d := make([]interface{}, len(t))
	for i := range t {
		d[i] = t[i]
	}
	_, err := s.db.InsertMany(ctx, d)
	return err
}

// SetCategory changes the category for uncategorized transaction
func (s *Repository) SetCategory(ctx context.Context, trId primitive.ObjectID, catId primitive.ObjectID) error {
	exists, err := s.categories.Any(ctx, catId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("category doesn't exist")
	}

	filter := bson.M{"_id": trId}
	update := bson.M{"$set": bson.M{"categoryId": catId}}
	_, err = s.db.UpdateOne(ctx, filter, update)
	return err
}

// ReadUncategorized returns all uncategorized transactions
func (s *Repository) ReadUncategorized(ctx context.Context) ([]Transaction, error) {
	var list []Transaction
	filter := bson.M{"categoryId": primitive.Null{}}
	cursor, err := s.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &list)
	if err != nil {
		return nil, err
	}
	return list, err
}
