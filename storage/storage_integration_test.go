package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/lungria/spendshelf-backend/storage/category"

	"github.com/stretchr/testify/require"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/lungria/spendshelf-backend/storage/pgtest"

	"github.com/lungria/spendshelf-backend/storage"

	"github.com/stretchr/testify/assert"
)

var defaultCategory = storage.Category{
	ID:   1,
	Name: "Unknown",
	Logo: "creditcard",
}

func TestSave_OnDuplicateInsert_DoesNothing(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	accountID := prepareTestAccount(t, pool)
	prepareTestCategory(t, pool, defaultCategory)

	db := storage.NewPostgreSQLStorage(pool)
	// try insert
	err := db.Save(context.Background(), []storage.Transaction{{
		"id1",
		time.Now().UTC(),
		"ORIGINAL DESCRIPTION",
		123,
		true,
		1110,
		accountID,
		defaultCategory.ID,
		time.Now().UTC(),
		nil,
	}, {
		"id2",
		time.Now().UTC(),
		"car",
		3121,
		true,
		1500,
		accountID,
		defaultCategory.ID,
		time.Now().UTC(),
		nil,
	}})

	assert.NoError(t, err)

	// try insert with same ID
	err = db.Save(context.Background(), []storage.Transaction{{
		"id1",
		time.Now().UTC(),
		"UPDATED DESCRIPTION",
		123,
		true,
		1110,
		accountID,
		defaultCategory.ID,
		time.Now().UTC(),
		nil,
	}})

	assert.NoError(t, err)

	// check that description was not changed
	row := pool.QueryRow(
		context.Background(),
		`select "description" from transaction
		where "ID" = 'id1'`)

	var description string

	err = row.Scan(&description)

	assert.NoError(t, err)
	assert.Equal(t, "ORIGINAL DESCRIPTION", description)
}

func TestGetLastTransactionDate_WithLocalDb_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	accountID := prepareTestAccount(t, pool)
	prepareTestCategory(t, pool, defaultCategory)
	db := storage.NewPostgreSQLStorage(pool)
	mockTransactions := []storage.Transaction{
		{
			ID:          "old-tr",
			Time:        time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC),
			Description: "desc",
			MCC:         10,
			Hold:        false,
			Amount:      100,
			AccountID:   accountID,
			CategoryID:  defaultCategory.ID,
		},
		{
			ID:          "new-tr",
			Time:        time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC),
			Description: "desc",
			MCC:         10,
			Hold:        false,
			Amount:      100,
			AccountID:   accountID,
			CategoryID:  defaultCategory.ID,
		},
	}

	err := db.Save(context.Background(), mockTransactions)
	assert.NoError(t, err)

	lastTransactionDate, err := db.GetLastTransactionDate(context.Background(), accountID)

	assert.NoError(t, err)
	assert.Equal(t, mockTransactions[1].Time, lastTransactionDate)
}

func TestGetByID_WithLocalDb_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	accountID := prepareTestAccount(t, pool)
	prepareTestCategory(t, pool, defaultCategory)
	db := storage.NewPostgreSQLStorage(pool)
	mockTransactions := []storage.Transaction{
		{
			ID:          "1",
			Time:        time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC),
			Description: "desc",
			MCC:         10,
			Hold:        false,
			Amount:      100,
			AccountID:   accountID,
			CategoryID:  defaultCategory.ID,
		},
		{
			ID:          "2",
			Time:        time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC),
			Description: "desc",
			MCC:         10,
			Hold:        false,
			Amount:      100,
			AccountID:   accountID,
			CategoryID:  defaultCategory.ID,
		},
	}

	err := db.Save(context.Background(), mockTransactions)
	assert.NoError(t, err)

	transaction, err := db.GetByID(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, "1", transaction.ID)
}

func TestGetByCategory_WithLocalDb_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	newCategory := storage.Category{
		ID:   22,
		Name: "test_category",
		Logo: "no_logo",
	}
	accountID := prepareTestAccount(t, pool)
	prepareTestCategory(t, pool, defaultCategory)
	prepareTestCategory(t, pool, newCategory)
	db := storage.NewPostgreSQLStorage(pool)
	mockTransactions := []storage.Transaction{
		{
			ID:          "1",
			Time:        time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC),
			Description: "desc",
			MCC:         10,
			Hold:        false,
			Amount:      100,
			AccountID:   accountID,
			CategoryID:  newCategory.ID,
		},
		{
			ID:          "2",
			Time:        time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC),
			Description: "desc",
			MCC:         10,
			Hold:        false,
			Amount:      100,
			AccountID:   accountID,
			CategoryID:  defaultCategory.ID,
		},
	}

	err := db.Save(context.Background(), mockTransactions)
	assert.NoError(t, err)

	transaction, err := db.GetByCategory(context.Background(), newCategory.ID)

	assert.NoError(t, err)
	assert.Len(t, transaction, 1)
	assert.Equal(t, mockTransactions[0].ID, transaction[0].ID)
}

func TestGetCategories_WithLocalDb_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	newCategory := storage.Category{
		ID:   22,
		Name: "test_category",
		Logo: "no_logo",
	}

	prepareTestCategory(t, pool, defaultCategory)
	prepareTestCategory(t, pool, newCategory)

	db := storage.NewPostgreSQLStorage(pool)

	categories, err := db.GetCategories(context.Background())
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.Len(t, categories, 2)

	for _, v := range categories {
		assert.True(t, v.ID == 22 || v.ID == 1)
		assert.True(t, v.Name == "test_category" || v.Name == "Unknown")
	}
}

func TestUpdate_WithLocalDb_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	accountID := prepareTestAccount(t, pool)
	prepareTestCategory(t, pool, defaultCategory)
	db := storage.NewPostgreSQLStorage(pool)
	// prepare transaction
	err := db.Save(context.Background(), []storage.Transaction{{
		"id4",
		time.Now().UTC(),
		"food",
		123,
		true,
		1110,
		accountID,
		defaultCategory.ID,
		time.Now(),
		nil,
	}})
	assert.NoError(t, err)
	transaction, err := db.GetByID(context.Background(), "id4")
	assert.NoError(t, err)

	// update comment without category
	comment := "comment"
	_, err = db.UpdateTransaction(context.Background(), storage.UpdateTransactionCommand{
		Query: storage.Query{
			ID:            "id4",
			LastUpdatedAt: transaction.LastUpdatedAt,
		},
		UpdatedFields: storage.UpdatedFields{
			Comment: &comment,
		},
	})
	assert.NoError(t, err)
	updatedTransaction, err := db.GetByID(context.Background(), "id4")
	assert.NoError(t, err)
	assert.Equal(t, "comment", *updatedTransaction.Comment)
	assert.Equal(t, category.Default, updatedTransaction.CategoryID)

	transaction, err = db.GetByID(context.Background(), "id4")
	require.NoError(t, err)

	// update category without comment
	newCategory := storage.Category{
		ID:   99,
		Name: "Food",
		Logo: "food",
	}

	prepareTestCategory(t, pool, newCategory)

	_, err = db.UpdateTransaction(context.Background(), storage.UpdateTransactionCommand{
		Query: storage.Query{
			ID:            "id4",
			LastUpdatedAt: transaction.LastUpdatedAt,
		},
		UpdatedFields: storage.UpdatedFields{
			CategoryID: &newCategory.ID,
		},
	})
	assert.NoError(t, err)
	updatedTransaction, err = db.GetByID(context.Background(), "id4")
	assert.NoError(t, err)
	assert.Equal(t, "comment", *updatedTransaction.Comment)
	assert.Equal(t, newCategory.ID, updatedTransaction.CategoryID)
}

func TestGetReport_WithLocalDb_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	newCategory := storage.Category{
		ID:   22,
		Name: "test_category",
		Logo: "no_logo",
	}
	accountID := prepareTestAccount(t, pool)
	prepareTestCategory(t, pool, defaultCategory)
	prepareTestCategory(t, pool, newCategory)
	db := storage.NewPostgreSQLStorage(pool)
	mockTransactions := []storage.Transaction{
		{
			ID:          "1",
			Time:        time.Date(2020, 10, 10, 0, 0, 0, 0, time.UTC),
			Description: "desc",
			MCC:         10,
			Hold:        false,
			Amount:      100,
			AccountID:   accountID,
			CategoryID:  newCategory.ID,
		},
		{
			ID:          "2",
			Time:        time.Date(2020, 10, 11, 0, 0, 0, 0, time.UTC),
			Description: "desc",
			MCC:         10,
			Hold:        false,
			Amount:      100,
			AccountID:   accountID,
			CategoryID:  defaultCategory.ID,
		},

		{
			ID:          "2",
			Time:        time.Date(2020, 10, 13, 0, 0, 0, 0, time.UTC),
			Description: "desc",
			MCC:         10,
			Hold:        false,
			Amount:      100,
			AccountID:   accountID,
			CategoryID:  defaultCategory.ID,
		},
	}

	err := db.Save(context.Background(), mockTransactions)
	assert.NoError(t, err)

	report, err := db.GetReport(
		context.Background(),
		time.Date(2020, 10, 0o1, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 10, 12, 0, 0, 0, 0, time.UTC))

	assert.NoError(t, err)
	assert.Len(t, report, 2)
	assert.Equal(t, mockTransactions[0].Amount, report[22])
	assert.Equal(t, mockTransactions[1].Amount, report[1])
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

func prepareTestCategory(t *testing.T, db *pgxpool.Pool, category storage.Category) {
	_, err := db.Exec(context.Background(), `
				insert into category ("ID", "name", "logo", "createdAt")
				values ($1, $2, $3, current_timestamp(0))`, category.ID, category.Name, category.Logo)

	require.NoError(t, err)
}
