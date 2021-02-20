package storage_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/lungria/spendshelf-backend/storage"
	"github.com/lungria/spendshelf-backend/storage/pgtest"
	"github.com/stretchr/testify/assert"
)

func TestBudgetsStorageGetLast_WhenNoLimitsExist_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	budgetID := prepareTestBudget(t, pool)
	db := storage.NewBudgetsStorage(pool)

	budget, err := db.GetLast(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, budget.ID, budgetID)
}

func TestBudgetsStorageGetLast_WhenLimitsExist_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	budgetID := prepareTestBudget(t, pool)
	prepareTestCategory(t, pool, storage.Category{
		ID:      1,
		Name:    "test1",
		Logo:    "test1",
		Visible: true,
	})
	prepareTestCategory(t, pool, storage.Category{
		ID:      2,
		Name:    "test2",
		Logo:    "test2",
		Visible: true,
	})
	prepareTestLimit(t, pool, budgetID, storage.Limit{
		CategoryID: 1,
		Amount:     100,
	})
	prepareTestLimit(t, pool, budgetID, storage.Limit{
		CategoryID: 2,
		Amount:     200,
	})

	db := storage.NewBudgetsStorage(pool)

	budget, err := db.GetLast(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, budget.ID, budgetID)
	assert.Len(t, budget.Limits, 2)
	assert.Equal(t, int64(200), budget.Limits[0].Amount)
}

func prepareTestBudget(t *testing.T, db *pgxpool.Pool) int {
	budgetID := 1
	_, err := db.Exec(context.Background(), `
				insert into budget
					("ID", "startsAt", "endsAt", "createdAt")
 					VALUES ($1,  current_timestamp(0),  current_timestamp(0),  current_timestamp(0))
				`, budgetID)

	require.NoError(t, err)

	return budgetID
}

func prepareTestLimit(t *testing.T, db *pgxpool.Pool, budgetID int, limit storage.Limit) {
	_, err := db.Exec(context.Background(), `
				insert into "limit"
					("budgetID", "categoryID", amount) 
					VALUES ($1,$2,$3)`, budgetID, limit.CategoryID, limit.Amount)

	require.NoError(t, err)
}
