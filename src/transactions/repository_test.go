package transactions

import (
	"testing"

	"github.com/lungria/spendshelf-backend/src/models"

	"github.com/stretchr/testify/assert"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var txn1 = models.Transaction{
	ID:          primitive.NewObjectID(),
	Time:        1578391653,
	Description: "Some Desc",
	Amount:      -10,
	Balance:     90,
}

var txn2 = models.Transaction{
	ID:          primitive.NewObjectID(),
	Time:        1578391654,
	Description: "Test description",
	Category:    "Shopping",
	Amount:      -90,
	Balance:     0,
}

func TestFindAll(t *testing.T) {
	expectedResult := []models.Transaction{txn1, txn2}
	mockRepo := MockRepository{}
	mockRepo.On("FindAll").Return(expectedResult, nil)
	actualResult, err := mockRepo.FindAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResult, actualResult)
	mockRepo.AssertExpectations(t)
}

func TestFindAllCategorized(t *testing.T) {
	expectedResult := []models.Transaction{txn2}
	mockRepo := MockRepository{}
	mockRepo.On("FindAllCategorized").Return(expectedResult, nil)
	actualResult, err := mockRepo.FindAllCategorized()
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResult, actualResult)
	mockRepo.AssertExpectations(t)
}

func TestFindAllUncategorized(t *testing.T) {
	expectedResult := []models.Transaction{txn1}
	mockRepo := MockRepository{}
	mockRepo.On("FindAllUncategorized").Return(expectedResult, nil)
	actualResult, err := mockRepo.FindAllUncategorized()
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResult, actualResult)
	mockRepo.AssertExpectations(t)
}

func TestFindAllByCategory(t *testing.T) {
	expectedResult := []models.Transaction{txn1}
	mockRepo := MockRepository{}
	mockRepo.On("FindAllByCategory", "Shopping").Return(expectedResult, nil)
	actualResult, err := mockRepo.FindAllByCategory("Shopping")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResult, actualResult)
	mockRepo.AssertExpectations(t)
}

func TestFindAllByCategory_Fail(t *testing.T) {
	expectedResult := []models.Transaction{}
	mockRepo := MockRepository{}
	mockRepo.On("FindAllByCategory", "Entertainment").Return(expectedResult, nil)
	actualResult, err := mockRepo.FindAllByCategory("Entertainment")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResult, actualResult)
	mockRepo.AssertExpectations(t)
}

func TestUpdateCategory(t *testing.T) {
	mockRepo := MockRepository{}
	mockRepo.On("UpdateCategory", txn2.ID, "Entertainment").Return(nil)
	err := mockRepo.UpdateCategory(txn2.ID, "Entertainment")
	assert.Equal(t, nil, err)
	mockRepo.AssertExpectations(t)
}
