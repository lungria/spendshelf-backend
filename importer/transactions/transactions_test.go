package transactions_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lungria/spendshelf-backend/storage"

	"github.com/lungria/spendshelf-backend/importer/mono"
	"github.com/lungria/spendshelf-backend/importer/transactions"
	"github.com/lungria/spendshelf-backend/importer/transactions/mock"
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
	svc := transactions.NewImporter(api, storage, gen)

	err := svc.Import(context.Background(), "acc")

	assert.True(t, errors.Is(err, testError))
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
	svc := transactions.NewImporter(api, storage, gen)

	err := svc.Import(context.Background(), "acc")

	assert.True(t, errors.Is(err, testError))
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
	svc := transactions.NewImporter(api, storage, gen)

	err := svc.Import(context.Background(), "acc")

	assert.Nil(t, err)
	calls := storage.SaveCalls()
	assert.Equal(t, 0, calls)
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
	db.SaveFunc = func(ctx context.Context, transactions []storage.Transaction) error {
		return testError
	}
	gen := &mock.ImportIntervalGeneratorMock{}
	gen.GetIntervalFunc = func(ctx context.Context, accountID string) (time.Time, time.Time, error) {
		return time.Now(), time.Now(), nil
	}
	svc := transactions.NewImporter(api, db, gen)

	err := svc.Import(context.Background(), "acc")

	assert.Error(t, err)
	saveCalls := db.SaveCalls()
	assert.Equal(t, 1, len(saveCalls))
	assert.Equal(t, transactionID, saveCalls[0].Transactions[0].ID)
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
	db.SaveFunc = func(ctx context.Context, transactions []storage.Transaction) error {
		return nil
	}
	gen := &mock.ImportIntervalGeneratorMock{}
	gen.GetIntervalFunc = func(ctx context.Context, accountID string) (time.Time, time.Time, error) {
		return time.Now(), time.Now(), nil
	}
	svc := transactions.NewImporter(api, db, gen)

	err := svc.Import(context.Background(), "acc")

	assert.Nil(t, err)
	saveCalls := db.SaveCalls()
	assert.Equal(t, 1, len(saveCalls))
	assert.Equal(t, transactionID, saveCalls[0].Transactions[0].ID)
}
