package categories

import (
	"context"
	"strings"

	"github.com/lungria/spendshelf-backend/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/text/unicode/norm"

	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	categoriesCollection = "categories"
)

// Repository define methods which Insert category to database and inherit methods of  memory cache
type Repository interface {
	provider
	Insert(ctx context.Context, name string) (primitive.ObjectID, error)
}

// CachedRepository implements by methods which define in Repository
type CachedRepository struct {
	provider   provider
	collection *mongo.Collection
	updates    chan<- models.Category
}

// NewCachedRepository is creating a new CachedRepository
func NewCachedRepository(db *mongo.Database) (*CachedRepository, error) {
	//todo get shutdown context
	ctx := context.Background()
	collection := db.Collection(categoriesCollection)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "Unable to seed categories cache")
	}
	seed := make([]models.Category, 0)
	err = cursor.All(ctx, &seed)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to seed categories cache")
	}
	updates := make(chan models.Category)
	provider, err := newProvider(ctx, seed, updates)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to seed categories collection")
	}
	return &CachedRepository{
		provider:   provider,
		updates:    updates,
		collection: collection,
	}, nil
}

// GetAll return all categories
func (repo *CachedRepository) GetAll() []models.Category {
	return repo.provider.GetAll()
}

// Find is returning preferred category
func (repo *CachedRepository) Find(name string) (models.Category, bool) {
	return repo.provider.Find(name)
}

// FindByID returns the category which was found by ObjectID
func (repo *CachedRepository) FindByID(id primitive.ObjectID) (models.Category, bool) {
	return repo.provider.FindByID(id)
}

// Insert is writing category to database and in memory cache
func (repo *CachedRepository) Insert(ctx context.Context, name string) (primitive.ObjectID, error) {
	// todo add unique index for normalized name in db
	trimmed := strings.TrimSpace(name)
	normalized := norm.NFC.String(strings.ToUpper(trimmed))
	c, exists := repo.provider.Find(normalized)
	if exists {
		return c.ID, nil
	}
	c = models.Category{
		NormalizedName: normalized,
		Name:           trimmed,
	}
	result, err := repo.collection.InsertOne(ctx, c)
	if err != nil {
		return primitive.NilObjectID, errors.Wrap(err, "Unable to insert category")
	}

	c.ID = result.InsertedID.(primitive.ObjectID)
	go func() {
		repo.updates <- c
	}()
	return c.ID, nil
}
