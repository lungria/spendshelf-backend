package categories

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/text/unicode/norm"

	"github.com/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	categoriesCollection = "categories"
)

// Repository define methods which do some work with database
type Repository interface {
	provider
	Insert(ctx context.Context, name string) (primitive.ObjectID, error)
}

// CachedRepository database repository
type CachedRepository struct {
	provider   provider
	collection *mongo.Collection
	updates    chan<- Category
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
	seed := make([]Category, 0)
	err = cursor.All(ctx, &seed)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to seed categories cache")
	}
	updates := make(chan Category)
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
func (repo *CachedRepository) GetAll() []Category {
	return repo.provider.GetAll()
}

// Find is returning preferred category
func (repo *CachedRepository) Find(name string) (Category, bool) {
	return repo.provider.Find(name)
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
	c = Category{
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
