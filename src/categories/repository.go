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

type Repository interface {
	Provider
	Insert(ctx context.Context, name string) (primitive.ObjectID, error)
}

type CachedRepository struct {
	provider   Provider
	collection *mongo.Collection
	updates    chan<- Category
}

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
	provider, err := NewProvider(ctx, seed, updates)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to seed categories collection")
	}
	return &CachedRepository{
		provider:   provider,
		updates:    updates,
		collection: collection,
	}, nil
}

func (repo *CachedRepository) GetAll() []Category {
	return repo.provider.GetAll()
}

func (repo *CachedRepository) Find(name string) (Category, bool) {
	return repo.provider.Find(name)
}

func (repo *CachedRepository) Insert(ctx context.Context, name string) (primitive.ObjectID, error) {
	// todo add unique index for normalized name in db
	trimmed := strings.TrimSpace(name)
	normalized := norm.NFC.String(strings.ToUpper(trimmed))
	c, exists := repo.provider.Find(normalized)
	if exists {
		return c.Id, nil
	}
	c = Category{
		NormalizedName: normalized,
		Name:           trimmed,
	}
	result, err := repo.collection.InsertOne(ctx, c)
	if err != nil {
		return primitive.NilObjectID, errors.Wrap(err, "Unable to insert category")
	}

	c.Id = result.InsertedID.(primitive.ObjectID)
	go func() {
		repo.updates <- c
	}()
	return c.Id, nil
}
