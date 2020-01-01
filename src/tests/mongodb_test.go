package tests

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lungria/mono"
	"github.com/lungria/spendshelf-backend/src/models"
	mock_db "github.com/lungria/spendshelf-backend/src/tests/mocks"
)

func TestSaveOneTransaction(t *testing.T) {
	transaction := models.Transaction{
		AccountID: "test_id",
		StatementItem: mono.StatementItem{
			ID:              "test_t_id",
			Time:            1577011328,
			Description:     "Some desc",
			MCC:             4900,
			Hold:            false,
			Amount:          100,
			OperationAmount: 100,
			CurrencyCode:    826,
			CommissionRate:  0,
			CashbackAmount:  0,
			Balance:         45000,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mock_db.NewMockWebHookDB(ctrl)

	db.EXPECT().SaveOneTransaction(&transaction)
	if err := db.SaveOneTransaction(&transaction); err != nil {
		t.Errorf("Save transaction failed. Error: %v", err)
	}
}

func TestGetAllTransactions(t *testing.T) {
	item := mono.StatementItem{
		ID:              "test",
		Time:            1577014638,
		Description:     "test desc",
		MCC:             1212,
		Hold:            false,
		Amount:          100,
		OperationAmount: 100,
		CurrencyCode:    826,
		CommissionRate:  0,
		CashbackAmount:  0,
		Balance:         45648,
	}
	expectedTransactions := []models.Transaction{{
		AccountID:     "test1",
		StatementItem: item,
	},
		{
			AccountID:     "test1",
			StatementItem: item,
		}}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mock_db.NewMockWebHookDB(ctrl)
	db.EXPECT().GetAllTransactions("test1").Return(expectedTransactions, nil)

	actualTransactions, err := db.GetAllTransactions("test1")
	if err != nil {
		t.Errorf("Get all transactions failed. Error: %v", err)
	}
	for i := 0; i > len(actualTransactions); i++ {
		if actualTransactions[i] != expectedTransactions[i] {
			t.Errorf("Transactions aren't same. Expected %v, got %v", expectedTransactions[i], actualTransactions[i])
		}
	}
}

func TestGetOneTransactions(t *testing.T) {
	expectedTransaction := models.Transaction{
		AccountID: "test_id",
		StatementItem: mono.StatementItem{
			ID:              "test_t_id",
			Time:            1577011328,
			Description:     "Some desc",
			MCC:             4900,
			Hold:            false,
			Amount:          100,
			OperationAmount: 100,
			CurrencyCode:    826,
			CommissionRate:  0,
			CashbackAmount:  0,
			Balance:         45000,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := mock_db.NewMockWebHookDB(ctrl)
	db.EXPECT().GetTransactionByID("test_t_id").Return(expectedTransaction, nil)

	actualTransaction, err := db.GetTransactionByID("test_t_id")
	if err != nil {
		t.Errorf("Get transaction by transactions id failed. Error: %v", err)
	}
	if actualTransaction != expectedTransaction {
		t.Errorf("Transactions aren't same. Expected %v, got %v", expectedTransaction, actualTransaction)
	}
}
