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
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Time time.Time          `json:"time" bson:"time"`
	// LocalDate describes transaction date in local timezone of the user.
	// Used for reporting purposes, so we can aggregate using it.
	LocalDate   ShortDate          `json:"-" bson:"localDate"`
	Description string             `json:"description" bson:"description"`
	CategoryID  primitive.ObjectID `json:"categoryId,omitempty" bson:"categoryId,omitempty"`
	Amount      int32              `json:"amount" json:"amount"`
	BankId      Bank               `json:"bankId" json:"bankId"`
}

type ShortDate struct {
	Day   int8  `bson:"day"`
	Month int8  `bson:"month"`
	Year  int16 `bson:"year"`
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
	t.LocalDate = ShortDate{
		Day:   int8(d),
		Month: int8(m),
		Year:  int16(y),
	}
	_, err := s.db.InsertOne(ctx, t)
	return err
}

// Insert multiple transaction to database.
func (s *Repository) InsertMany(ctx context.Context, t []Transaction) error {
	data := make([]interface{}, len(t))
	for i, v := range t {
		y, m, d := v.Time.Date()
		v.LocalDate = ShortDate{
			Day:   int8(d),
			Month: int8(m),
			Year:  int16(y),
		}
		data[i] = v
	}
	_, err := s.db.InsertMany(ctx, data)
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

type ReportEntry struct {
	Date          string `json:"date" bson:"_id"`
	Sum           int32  `json:"sum" bson:"sum"`
	EndDayBalance int32  `json:"endDayBalance" bson:"-"`
}

// BuildDailyReport returns report with spendings per day.
func (s *Repository) BuildDailyReport(ctx context.Context, from time.Time, to time.Time, balance int32) ([]ReportEntry, error) {
	var list []ReportEntry
	filter := bson.A{

		bson.M{"$match": bson.M{"time": bson.M{
			"$lte": to,
			"$gte": from}},
		},
		bson.M{"$group": bson.M{
			"_id": "$localDate",
			"sum": bson.M{"$sum": "$amount"},
		}},
		bson.M{"$sort": bson.M{"_id.year": 1, "_id.month": 1, "_id.day": 1}},
		bson.M{"$project": bson.M{"_id": bson.M{"$concat": bson.A{bson.M{"$toString": "$_id.day"}, ".", bson.M{"$toString": "$_id.month"}, ".", bson.M{"$toString": "$_id.year"}}}, "sum": "$sum"}},
	}
	cursor, err := s.db.Aggregate(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &list)
	if err != nil {
		return nil, err
	}

	for i, v := range list {
		balance = balance + v.Sum
		list[i].EndDayBalance = balance
	}
	return list, err
}
