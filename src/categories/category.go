package categories

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/text/unicode/norm"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	categoriesCollection = "categories"
)

// Category is general struct for spendshelf categories
type Category struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	NormalizedName string             `json:"normalizedName" bson:"normalizedName"`
}

type Repository struct {
	collection *mongo.Collection
	logger     *zap.SugaredLogger
}

// NewRepository is creating a new Repository
func NewRepository(db *mongo.Database, logger *zap.SugaredLogger) *Repository {
	collection := db.Collection(categoriesCollection)
	return &Repository{
		collection: collection,
		logger:     logger,
	}
}

// GetAll return all categories
func (repo *Repository) GetAll(ctx context.Context) ([]Category, error) {
	result, err := repo.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all transactions: %w", err)
	}
	var categories []Category
	err = result.All(ctx, &categories)
	if err != nil {
		return nil, fmt.Errorf("failed to get all transactions: %w", err)
	}
	return categories, nil
}

// Find is returning preferred category
func (repo *Repository) findByNormalizedName(ctx context.Context, name string) *Category {
	query := bson.M{"normalizedName": name}
	result := repo.collection.FindOne(ctx, query)
	err := result.Err()
	if err != nil {
		return nil
	}
	model := &Category{}
	err = result.Decode(model)
	if err != nil {
		return nil
	}

	return model
}

// Any checks if category with specified id exists
func (repo *Repository) Any(ctx context.Context, id primitive.ObjectID) (bool, error) {
	query := bson.M{"_id": id}
	opts := &options.CountOptions{}
	opts.SetLimit(1)
	count, err := repo.collection.CountDocuments(ctx, query)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Insert is writing category to database and in memory cache
func (repo *Repository) Insert(ctx context.Context, name string) (primitive.ObjectID, error) {
	trimmed := strings.TrimSpace(name)
	normalized := norm.NFC.String(strings.ToUpper(trimmed))
	ctg := repo.findByNormalizedName(ctx, normalized)
	if ctg != nil {
		return ctg.ID, nil
	}

	newCtg := Category{
		NormalizedName: normalized,
		Name:           trimmed,
	}
	result, err := repo.collection.InsertOne(ctx, newCtg)
	if err != nil {
		return primitive.NilObjectID, fmt.Errorf("unable to insert category: %w", err)
	}

	newCtg.ID = result.InsertedID.(primitive.ObjectID)
	return newCtg.ID, nil
}
