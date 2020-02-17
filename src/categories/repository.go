package categories

import (
	"context"
	"fmt"
	"strings"

	"github.com/lungria/spendshelf-backend/src/models"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/text/unicode/norm"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	categoriesCollection = "categories"
)

type Repository struct {
	collection *mongo.Collection
}

// NewRepository is creating a new Repository
func NewRepository(db *mongo.Database) *Repository {
	collection := db.Collection(categoriesCollection)
	return &Repository{
		collection: collection,
	}
}

// GetAll return all categories
func (repo *Repository) GetAll(ctx context.Context) ([]models.Category, error) {
	result, err := repo.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all transactions: %w", err)
	}
	var categories []models.Category
	err = result.All(ctx, &categories)
	if err != nil {
		return nil, fmt.Errorf("failed to get all transactions: %w", err)
	}
	return categories, nil
}

// Find is returning preferred category
func (repo *Repository) FindByNormalizedName(ctx context.Context, name string) *models.Category {
	query := bson.M{"normalizedName": name}
	result := repo.collection.FindOne(ctx, query)
	err := result.Err()
	if err != nil {
		return nil
	}
	model := &models.Category{}
	err = result.Decode(model)
	if err != nil {
		return nil
	}

	return model
}

// FindByID returns the category which was found by ObjectID
func (repo *Repository) FindByID(ctx context.Context, id primitive.ObjectID) *models.Category {
	query := bson.M{"_id": id}
	result := repo.collection.FindOne(ctx, query)
	err := result.Err()
	if err != nil {
		return nil
	}
	model := &models.Category{}
	err = result.Decode(model)
	if err != nil {
		return nil
	}

	return model
}

// Insert is writing category to database and in memory cache
func (repo *Repository) Insert(ctx context.Context, name string) (primitive.ObjectID, error) {
	trimmed := strings.TrimSpace(name)
	normalized := norm.NFC.String(strings.ToUpper(trimmed))
	ctg := repo.FindByNormalizedName(ctx, normalized)
	if ctg != nil {
		return ctg.ID, nil
	}

	newCtg := models.Category{
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
