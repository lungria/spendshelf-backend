package importer_test

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/lungria/spendshelf-backend/importer"
	"github.com/lungria/spendshelf-backend/importer/mock"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestImport_WhenCalled_ImportsAccountsAndTransactions(t *testing.T) {
	accounts := &mock.AccountImporterMock{}
	accounts.ImportFunc = func(ctx context.Context, accountID string) error {
		return nil
	}
	transactions := &mock.TransactionsImporterMock{}
	transactions.ImportFunc = func(ctx context.Context, accountID string) error {
		return nil
	}
	log.Logger = zerolog.New(ioutil.Discard)
	svc := importer.NewImporter(transactions, accounts)

	svc.Import(context.Background(), "account")

	assert.Equal(t, 1, len(accounts.ImportCalls()))
	assert.Equal(t, 1, len(transactions.ImportCalls()))
}
