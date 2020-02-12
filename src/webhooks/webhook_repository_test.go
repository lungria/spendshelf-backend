package webhooks

import (
	"testing"

	"github.com/lungria/spendshelf-backend/src/models"

	"github.com/stretchr/testify/assert"

	"github.com/shal/mono"
)

func TestInsertOneHook(t *testing.T) {
	webhook := models.WebHook{
		AccountID: "test_id",
		StatementItem: mono.Transaction{
			ID:              "test_t_id",
			Time:            1577011328,
			Description:     "Some desc",
			MCC:             4900,
			Hold:            false,
			Amount:          100,
			OperationAmount: 100,
			CurrencyCode:    826,
			CommissionRate:  0,
			Balance:         45000,
		},
	}
	repoMock := &MockRepository{}
	repoMock.On("InsertOneHook", &webhook).Return(nil).Once()
	err := repoMock.InsertOneHook(&webhook)
	assert.Equal(t, nil, err)
	repoMock.AssertExpectations(t)
}
