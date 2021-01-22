// todo: verify calls count on mocks!
package account_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lungria/spendshelf-backend/storage"

	"github.com/lungria/spendshelf-backend/importer/account"

	"github.com/lungria/spendshelf-backend/importer/account/mock"
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
	svc := account.NewDefaultImporter(api, storage)

	err := svc.Import(context.Background(), "acc")

	assert.True(t, errors.Is(err, testError))
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
	svc := account.NewDefaultImporter(api, storage)

	err := svc.Import(context.Background(), "unknown_account_id")

	assert.Error(t, err)
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
	db.SaveFunc = func(ctx context.Context, account storage.Account) error {
		return testError
	}
	svc := account.NewDefaultImporter(api, db)

	err := svc.Import(context.Background(), "id1")

	assert.True(t, errors.Is(err, testError))
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
	db.SaveFunc = func(ctx context.Context, account storage.Account) error {
		return nil
	}
	svc := account.NewDefaultImporter(api, db)

	err := svc.Import(context.Background(), "id1")

	assert.Nil(t, err)
}
