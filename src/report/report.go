package report

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/lungria/spendshelf-backend/src/transactions"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/lungria/spendshelf-backend/src/categories"
	"go.mongodb.org/mongo-driver/mongo"
)

type Element struct {
	CategoryId  primitive.ObjectID `json:"categoryId" bson:"_id"`
	TotalAmount int                `json:"totalAmount" bson:"totalAmount"`
}

type Generator interface {
	GetReport(ctx context.Context, from time.Time, to time.Time) ([]Element, error)
}

type SequentialReportGenerator struct {
	transactions *mongo.Collection
	categories   categories.Repository
	logger       *zap.SugaredLogger
}

func NewSequentialReportGenerator(db *mongo.Database, categories categories.Repository, logger *zap.SugaredLogger) *SequentialReportGenerator {
	return &SequentialReportGenerator{transactions: db.Collection(transactions.TransactionsCollection), categories: categories, logger: logger}
}

func (s *SequentialReportGenerator) GetReport(ctx context.Context, from time.Time, to time.Time) ([]Element, error) {
	s.logger.Infow("Generating report", "from", from.String(), "to", to.String())
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{"time": bson.D{
				{"$gte", from},
				{"$lt", to},
			}},
		},
		bson.M{
			"$match": bson.M{"category._id": bson.D{
				{"$exists", true},
				{"$ne", primitive.Null{}},
			}},
		},
		bson.M{
			"$group": bson.M{
				"_id":         "$category._id",
				"totalAmount": bson.M{"$sum": "$amount"}},
		},
	}
	result, err := s.transactions.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var report []Element
	err = result.All(ctx, &report)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return []Element{}, nil
	}
	return report, nil
}
