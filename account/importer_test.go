package account_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lungria/spendshelf-backend/account"
	"github.com/lungria/spendshelf-backend/account/mock"

	"github.com/lungria/spendshelf-backend/importer/mono"
	"github.com/stretchr/testify/assert"
)

func TestImport_WhenGetUserInfoFails_ReturnsError(t *testing.T) {
	testError := errors.New("something failed")
	api := &mock.BankAPIMock{}
	api.GetUserInfoFunc = func(ctx context.Context) ([]mono.Account, error) {
		return nil, testError
	}
	storage := &mock.StorageMock{}
	svc := account.NewImporter(api, storage)

	err := svc.Import(context.Background(), "acc")

	assert.True(t, errors.Is(err, testError))
	assert.Zero(t, len(storage.SaveCalls()))
	assert.NotZero(t, len(api.GetUserInfoCalls()))
}

func TestImport_WhenApiDoNotReturnAccountInfo_ReturnsError(t *testing.T) {
	api := &mock.BankAPIMock{}
	api.GetUserInfoFunc = func(ctx context.Context) ([]mono.Account, error) {
		return []mono.Account{
			{
				ID: "id1",
			},
			{
				ID: "id2",
			},
			{
				ID: "id3",
			},
		}, nil
	}
	storage := &mock.StorageMock{}
	svc := account.NewImporter(api, storage)

	err := svc.Import(context.Background(), "unknown_account_id")

	assert.Error(t, err)
	assert.Zero(t, len(storage.SaveCalls()))
	assert.NotZero(t, len(api.GetUserInfoCalls()))
}

func TestImport_WhenStorageSaveReturnsError_ReturnsError(t *testing.T) {
	api := &mock.BankAPIMock{}
	api.GetUserInfoFunc = func(ctx context.Context) ([]mono.Account, error) {
		return []mono.Account{
			{
				ID: "id1",
			},
		}, nil
	}
	testError := errors.New("something failed")
	db := &mock.StorageMock{}
	db.SaveFunc = func(ctx context.Context, account account.Account) error {
		return testError
	}
	svc := account.NewImporter(api, db)

	err := svc.Import(context.Background(), "id1")

	assert.True(t, errors.Is(err, testError))
	assert.NotZero(t, len(db.SaveCalls()))
	assert.NotZero(t, len(api.GetUserInfoCalls()))
}

func TestImport_WhenDataIsSaved_ReturnsNil(t *testing.T) {
	api := &mock.BankAPIMock{}
	api.GetUserInfoFunc = func(ctx context.Context) ([]mono.Account, error) {
		return []mono.Account{
			{
				ID: "id1",
			},
		}, nil
	}
	db := &mock.StorageMock{}
	db.SaveFunc = func(ctx context.Context, account account.Account) error {
		return nil
	}
	svc := account.NewImporter(api, db)

	err := svc.Import(context.Background(), "id1")

	assert.Nil(t, err)
	assert.NotZero(t, len(db.SaveCalls()))
	assert.NotZero(t, len(api.GetUserInfoCalls()))
}
