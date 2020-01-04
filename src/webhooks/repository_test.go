package webhooks

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lungria/mono"
)

func TestSaveOneHook(t *testing.T) {
	webhook := WebHook{
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
	repoMock := &MockRepository{}
	repoMock.On("SaveOneHook", &webhook).Return(nil).Once()
	_ = repoMock.SaveOneHook(&webhook)
	repoMock.AssertExpectations(t)
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
	expectedTransactions := []WebHook{{
		AccountID:     "test1",
		StatementItem: item,
	},
		{
			AccountID:     "test1",
			StatementItem: item,
		},
	}
	repoMock := &MockRepository{}
	repoMock.On("GetAllHooks", "test1").Return(expectedTransactions, nil).Once()
	actualTransactions, err := repoMock.GetAllHooks("test1")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedTransactions, actualTransactions)
	repoMock.AssertExpectations(t)

}

func TestGetOneTransactions(t *testing.T) {
	expectedTransaction := WebHook{
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
	repoMock := &MockRepository{}
	repoMock.On("GetHookByID", "test_t_id").Return(expectedTransaction, nil).Once()
	actualTransaction, err := repoMock.GetHookByID("test_t_id")
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedTransaction, actualTransaction)
	repoMock.AssertExpectations(t)
}
