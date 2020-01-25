package syncmono

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/lungria/spendshelf-backend/src/models"
	"github.com/lungria/spendshelf-backend/src/webhooks"
	"github.com/stretchr/testify/assert"

	"github.com/shal/mono"
)

func TestTxnFromSync(t *testing.T) {
	s := monoSync{accountUAH: &mono.Account{
		ID:           "acc_id",
		Balance:      0,
		CreditLimit:  0,
		CurrencyCode: 0,
		CashBackType: "",
	}}

	actualTxn := s.txnFromSyncTxn(monoTxn1)

	modelTxn1.BankTransaction.AccountID = s.accountUAH.ID
	modelTxn1.ID = actualTxn.ID

	assert.Equal(t, modelTxn1, actualTxn)

}

func TestTrimDuplicate(t *testing.T) {
	s := monoSync{accountUAH: &mono.Account{ID: "acc_id"}}

	txnFromDB := []models.Transaction{modelTxn1, modelTxn2}
	syncTxn := []mono.Transaction{monoTxn2, monoTxn3}

	actualTxn := s.trimDuplicate(syncTxn, txnFromDB)
	modelTxn3.ID = actualTxn[0].ID
	expectedTxn := []models.Transaction{modelTxn3}

	for i := 0; i < len(expectedTxn); i++ {
		assert.Equal(t, actualTxn[0], expectedTxn[0])
	}
}

var monoTxn1 = mono.Transaction{
	ID:              "test_id_1",
	Time:            1579954948,
	Description:     "Some description",
	MCC:             123123,
	Hold:            false,
	Amount:          100,
	OperationAmount: 100,
	CurrencyCode:    826,
	CommissionRate:  0,
	Balance:         900,
}

var modelTxn1 = models.Transaction{
	ID:          primitive.ObjectID{},
	Time:        time.Unix(int64(monoTxn1.Time), 0),
	Description: monoTxn1.Description,
	Category:    nil,
	Amount:      monoTxn1.Amount,
	Balance:     monoTxn1.Balance,
	Bank:        webhooks.MonoBankName,
	BankTransaction: models.WebHook{
		AccountID:     mono.Account{}.ID,
		StatementItem: monoTxn1,
	},
}

var monoTxn2 = mono.Transaction{
	ID:              "test_id_2",
	Time:            1579954948,
	Description:     "Some description",
	MCC:             123123,
	Hold:            false,
	Amount:          100,
	OperationAmount: 100,
	CurrencyCode:    826,
	CommissionRate:  0,
	Balance:         900,
}

var modelTxn2 = models.Transaction{
	ID:          primitive.ObjectID{},
	Time:        time.Unix(int64(monoTxn2.Time), 0),
	Description: monoTxn1.Description,
	Category:    nil,
	Amount:      monoTxn2.Amount,
	Balance:     monoTxn2.Balance,
	Bank:        webhooks.MonoBankName,
	BankTransaction: models.WebHook{
		AccountID:     mono.Account{}.ID,
		StatementItem: monoTxn2,
	},
}

var monoTxn3 = mono.Transaction{
	ID:              "test_id_3",
	Time:            1579954948,
	Description:     "Some description",
	MCC:             123123,
	Hold:            false,
	Amount:          100,
	OperationAmount: 100,
	CurrencyCode:    826,
	CommissionRate:  0,
	Balance:         900,
}

var modelTxn3 = models.Transaction{
	ID:          primitive.ObjectID{},
	Time:        time.Unix(int64(monoTxn3.Time), 0),
	Description: monoTxn1.Description,
	Category:    nil,
	Amount:      monoTxn3.Amount,
	Balance:     monoTxn3.Balance,
	Bank:        webhooks.MonoBankName,
	BankTransaction: models.WebHook{
		AccountID:     mono.Account{}.ID,
		StatementItem: monoTxn3,
	},
}
