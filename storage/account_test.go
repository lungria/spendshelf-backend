package storage_test

import (
	"context"
	"testing"

	"github.com/lungria/spendshelf-backend/storage"
	"github.com/lungria/spendshelf-backend/storage/pgtest"
	"github.com/stretchr/testify/assert"
)

func TestAccountStorageSave_WithProductionSchema_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	db := storage.NewAccountsStorage(pool)

	// test insert
	err := db.Save(context.Background(), storage.Account{
		ID:       "acc1",
		Balance:  10000,
		Currency: "UAH",
	})

	assert.NoError(t, err)

	// test on conflict update
	err = db.Save(context.Background(), storage.Account{
		ID:       "acc1",
		Balance:  20000,
		Currency: "UAH",
	})

	assert.NoError(t, err)

	row := pool.QueryRow(
		context.Background(),
		`select "balance" from account
		where "ID" = 'acc1'`)

	var balance int64

	err = row.Scan(&balance)

	assert.Equal(t, int64(20000), balance)
}

func TestAccountStorageGetAll_WithProductionSchema_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	accountID := prepareTestAccount(t, pool)
	db := storage.NewAccountsStorage(pool)

	accounts, err := db.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, accounts, 1)
	assert.Equal(t, accounts[0].ID, accountID)
}
