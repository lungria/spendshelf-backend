package account_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/account"
	"github.com/lungria/spendshelf-backend/storage/pgtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountStorageSave_WithProductionSchema_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	db := account.NewRepository(pool)

	// test insert
	err := db.Save(context.Background(), account.Account{
		ID:       "acc1",
		Balance:  10000,
		Currency: "UAH",
	})

	assert.NoError(t, err)

	// test on conflict update
	err = db.Save(context.Background(), account.Account{
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

	require.NoError(t, err)
	assert.Equal(t, int64(20000), balance)
}

func TestAccountStorageGetAll_WithProductionSchema_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	accountID := prepareTestAccount(t, pool)
	db := account.NewRepository(pool)

	accounts, err := db.GetAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, accounts, 1)
	assert.Equal(t, accounts[0].ID, accountID)
}

func prepareTestAccount(t *testing.T, db *pgxpool.Pool) string {
	accountID := "test-acc-id"
	_, err := db.Exec(context.Background(), `
				insert into "account"
							 ("ID", "createdAt", "description", "balance", "currency", "lastUpdatedAt")
							 values ($1, current_timestamp(0), 'desc', 0, 'UAH', current_timestamp(0))
				`, accountID)

	require.NoError(t, err)

	return accountID
}
