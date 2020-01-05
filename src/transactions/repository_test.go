package transactions

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lungria/mono"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var txn1 = Transaction{
	ID:        primitive.NewObjectID(),
	AccountID: "fCfTvZgMihvz_URzZSqVEf",
	StatementItem: mono.StatementItem{
		ID:              "qwe",
		Time:            123123123,
		Description:     "Some desc",
		MCC:             456798,
		Hold:            false,
		Amount:          -10,
		OperationAmount: -10,
		CurrencyCode:    826,
		CommissionRate:  0,
		CashbackAmount:  0,
		Balance:         90,
	},
}

var txn2 = Transaction{
	ID:        primitive.NewObjectID(),
	Category:  "Shopping",
	AccountID: "fCfTvZgMihvz_URzZSqVEf",
	StatementItem: mono.StatementItem{
		ID:              "qwe",
		Time:            123123123,
		Description:     "Some desc",
		MCC:             456798,
		Hold:            false,
		Amount:          -100,
		OperationAmount: -100,
		CurrencyCode:    826,
		CommissionRate:  0,
		CashbackAmount:  0,
		Balance:         0,
	},
}

func TestFindAll(t *testing.T) {
	expectedResult := []Transaction{txn1, txn2}
	mockRepo := MockRepository{}
	mockRepo.On("FindAll").Return(expectedResult, nil)
	actualResult, err := mockRepo.FindAll()
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResult, actualResult)
	mockRepo.AssertExpectations(t)
}

func TestFindAllCategorized(t *testing.T) {
	expectedResult := []Transaction{txn2}
	mockRepo := MockRepository{}
	mockRepo.On("FindAllCategorized").Return(expectedResult, nil)
	actualResult, err := mockRepo.FindAllCategorized()
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResult, actualResult)
	mockRepo.AssertExpectations(t)
}

func TestFindAllUncategorized(t *testing.T) {
	expectedResult := []Transaction{txn1}
	mockRepo := MockRepository{}
	mockRepo.On("FindAllUncategorized").Return(expectedResult, nil)
	actualResult, err := mockRepo.FindAllUncategorized()
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResult, actualResult)
	mockRepo.AssertExpectations(t)
}

func TestFindAllByCategory(t *testing.T) {
	expectedResult := []Transaction{txn1}
	mockRepo := MockRepository{}
	mockRepo.On("FindAllByCategory", "Shopping").Return(expectedResult, nil)
	actualResult, err := mockRepo.FindAllByCategory("Shopping")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResult, actualResult)
	mockRepo.AssertExpectations(t)
}

func TestFindAllByCategory_Fail(t *testing.T) {
	expectedResult := []Transaction{}
	mockRepo := MockRepository{}
	mockRepo.On("FindAllByCategory", "Entertainment").Return(expectedResult, nil)
	actualResult, err := mockRepo.FindAllByCategory("Entertainment")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedResult, actualResult)
	mockRepo.AssertExpectations(t)
}

func TestUpdateCategory(t *testing.T) {
	mockRepo := MockRepository{}
	mockRepo.On("UpdateCategory", txn2.ID.Hex(), "Entertainment").Return(nil)
	err := mockRepo.UpdateCategory(txn2.ID.Hex(), "Entertainment")
	assert.Equal(t, nil, err)
	mockRepo.AssertExpectations(t)
}
