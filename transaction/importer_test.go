package transaction_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lungria/spendshelf-backend/importer/mono"
	"github.com/lungria/spendshelf-backend/transaction"
	"github.com/lungria/spendshelf-backend/transaction/mock"
	"github.com/stretchr/testify/assert"
)

func TestImport_WhenGetIntervalFails_ReturnsError(t *testing.T) {
	testError := errors.New("something failed")
	api := &mock.BankAPIMock{}
	storage := &mock.StorageMock{}
	gen := &mock.ImportIntervalGeneratorMock{}
	gen.GetIntervalFunc = func(ctx context.Context, accountID string) (time.Time, time.Time, error) {
		return time.Time{}, time.Time{}, testError
	}
	svc := transaction.NewImporter(api, storage, gen)

	err := svc.Import(context.Background(), "acc")

	assert.True(t, errors.Is(err, testError))
	assert.Zero(t, len(storage.SaveCalls()))
	assert.Zero(t, len(api.GetTransactionsCalls()))
	assert.NotZero(t, len(gen.GetIntervalCalls()))
}

func TestImport_WhenApiGetTransactionsFails_ReturnsError(t *testing.T) {
	testError := errors.New("something failed")
	api := &mock.BankAPIMock{}
	api.GetTransactionsFunc = func(ctx context.Context, query mono.GetTransactionsQuery) ([]mono.Transaction, error) {
		return nil, testError
	}
	storage := &mock.StorageMock{}
	gen := &mock.ImportIntervalGeneratorMock{}
	gen.GetIntervalFunc = func(ctx context.Context, accountID string) (time.Time, time.Time, error) {
		return time.Now(), time.Now(), nil
	}
	svc := transaction.NewImporter(api, storage, gen)

	err := svc.Import(context.Background(), "acc")

	assert.True(t, errors.Is(err, testError))
	assert.Zero(t, len(storage.SaveCalls()))
	assert.NotZero(t, len(api.GetTransactionsCalls()))
	assert.NotZero(t, len(gen.GetIntervalCalls()))
}

func TestImport_WhenApiGetTransactionsReturnsNothing_StorageNotCalled(t *testing.T) {
	api := &mock.BankAPIMock{}
	api.GetTransactionsFunc = func(ctx context.Context, query mono.GetTransactionsQuery) ([]mono.Transaction, error) {
		return nil, nil
	}
	storage := &mock.StorageMock{}
	gen := &mock.ImportIntervalGeneratorMock{}
	gen.GetIntervalFunc = func(ctx context.Context, accountID string) (time.Time, time.Time, error) {
		return time.Now(), time.Now(), nil
	}
	svc := transaction.NewImporter(api, storage, gen)

	err := svc.Import(context.Background(), "acc")

	assert.Nil(t, err)
	assert.Equal(t, 0, len(storage.SaveCalls()))
	assert.Zero(t, len(storage.SaveCalls()))
	assert.NotZero(t, len(api.GetTransactionsCalls()))
	assert.NotZero(t, len(gen.GetIntervalCalls()))
}

func TestImport_WhenStorageSaveReturnsError_ReturnsError(t *testing.T) {
	transactionID := "trID"
	testError := errors.New("something failed")
	api := &mock.BankAPIMock{}
	api.GetTransactionsFunc = func(ctx context.Context, query mono.GetTransactionsQuery) ([]mono.Transaction, error) {
		return []mono.Transaction{
			{
				ID: transactionID,
			},
		}, nil
	}
	db := &mock.StorageMock{}
	db.SaveFunc = func(ctx context.Context, transactions []transaction.Transaction) error {
		return testError
	}
	gen := &mock.ImportIntervalGeneratorMock{}
	gen.GetIntervalFunc = func(ctx context.Context, accountID string) (time.Time, time.Time, error) {
		return time.Now(), time.Now(), nil
	}
	svc := transaction.NewImporter(api, db, gen)

	err := svc.Import(context.Background(), "acc")
	saveCalls := db.SaveCalls()

	assert.Error(t, err)
	assert.Equal(t, 1, len(saveCalls))
	assert.Equal(t, transactionID, saveCalls[0].Transactions[0].ID)
	assert.NotZero(t, len(db.SaveCalls()))
	assert.NotZero(t, len(api.GetTransactionsCalls()))
	assert.NotZero(t, len(gen.GetIntervalCalls()))
}

func TestImport_WhenDataIsSaved_ReturnsNil(t *testing.T) {
	transactionID := "trID"
	api := &mock.BankAPIMock{}
	api.GetTransactionsFunc = func(ctx context.Context, query mono.GetTransactionsQuery) ([]mono.Transaction, error) {
		return []mono.Transaction{
			{
				ID: transactionID,
			},
		}, nil
	}
	db := &mock.StorageMock{}
	db.SaveFunc = func(ctx context.Context, transactions []transaction.Transaction) error {
		return nil
	}
	gen := &mock.ImportIntervalGeneratorMock{}
	gen.GetIntervalFunc = func(ctx context.Context, accountID string) (time.Time, time.Time, error) {
		return time.Now(), time.Now(), nil
	}
	svc := transaction.NewImporter(api, db, gen)

	err := svc.Import(context.Background(), "acc")
	saveCalls := db.SaveCalls()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(saveCalls))
	assert.Equal(t, transactionID, saveCalls[0].Transactions[0].ID)
	assert.NotZero(t, len(db.SaveCalls()))
	assert.NotZero(t, len(api.GetTransactionsCalls()))
	assert.NotZero(t, len(gen.GetIntervalCalls()))
}
