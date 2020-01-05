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
	err := repoMock.SaveOneHook(&webhook)
	assert.Equal(t, nil, err)
	repoMock.AssertExpectations(t)
}
