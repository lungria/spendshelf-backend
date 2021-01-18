package account_test

import (
	"context"
	"errors"
	"testing"

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

func TestImport_WhenApiDoesntReturnAccountInfo_ReturnsError(t *testing.T) {
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

	var parsedError *account.NotFoundInAPIError
	errors.As(err, &parsedError)

	assert.NotNil(t, parsedError)
	assert.Equal(t, "unknown_account_id", parsedError.GetAccountID())
}
