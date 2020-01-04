package categories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"golang.org/x/text/unicode/norm"
)

func Test_GetAll_ForTwoSeededCategories_ReturnsTwoCategories(t *testing.T) {
	seed := []Category{
		{NormalizedName: "test"},
		{NormalizedName: "test2"},
	}
	provider := getProvider(seed, make(chan Category))

	categories := provider.GetAll()

	if len(categories) != 2 {
		t.Error(fmt.Sprintf("Expected 2 categories, received %v", len(categories)))
	}
}

func Test_Find_ForExistingCategory_ReturnsCategory(t *testing.T) {
	name := norm.NFC.String("test")
	seedID := primitive.NewObjectID()
	seed := []Category{
		{NormalizedName: name, ID: seedID},
		{NormalizedName: "other_category", ID: primitive.NewObjectID()},
	}
	provider := getProvider(seed, make(chan Category))

	category, _ := provider.Find(name)

	if category.ID != seedID {
		t.Error(fmt.Sprintf("Expected category with id %v, received %v", seedID, category.ID))
	}
}

func Test_Find_ForNewlyInsertedCategory_ReturnsCategory(t *testing.T) {
	seed := []Category{
		{NormalizedName: norm.NFC.String("test"), ID: primitive.NewObjectID()},
	}
	newCategory := Category{primitive.NewObjectID(), norm.NFC.String("test2"), norm.NFC.String("test2")}
	updates := make(chan Category)
	provider := getProvider(seed, updates)

	_, exists := provider.Find(newCategory.NormalizedName)
	if exists {
		t.Error("Found category that wasn't expected to be in collection")
	}
	updates <- newCategory
	<-time.After(time.Second * 1)

	_, existsAfterUpdate := provider.Find(newCategory.NormalizedName)
	if !existsAfterUpdate {
		t.Error("Category not found after insert")
	}
}

func getProvider(categories []Category, updates chan Category) *inMemoryProvider {
	provider, _ := newProvider(context.Background(), categories, updates)
	return provider
}
