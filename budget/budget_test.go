package budget_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/budget"
	"github.com/lungria/spendshelf-backend/storage/pgtest"
	"github.com/lungria/spendshelf-backend/transaction/category"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBudgetsStorageGetLast_WhenNoLimitsExist_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	budgetID := prepareTestBudget(t, pool)
	db := budget.NewRepository(pool)

	budget, err := db.GetLast(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, budget.ID, budgetID)
}

func TestBudgetsStorageGetLast_WhenLimitsExist_NoErrorReturned(t *testing.T) {
	pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
	defer cleanup()

	budgetID := prepareTestBudget(t, pool)
	prepareTestCategory(t, pool, category.Category{
		ID:   1,
		Name: "test1",
		Logo: "test1",
	})
	prepareTestCategory(t, pool, category.Category{
		ID:   2,
		Name: "test2",
		Logo: "test2",
	})
	prepareTestLimit(t, pool, budgetID, budget.Limit{
		CategoryID: 1,
		Amount:     100,
	})
	prepareTestLimit(t, pool, budgetID, budget.Limit{
		CategoryID: 2,
		Amount:     200,
	})

	db := budget.NewRepository(pool)

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

func prepareTestLimit(t *testing.T, db *pgxpool.Pool, budgetID int, limit budget.Limit) {
	_, err := db.Exec(context.Background(), `
				insert into "limit"
					("budgetID", "categoryID", amount) 
					VALUES ($1,$2,$3)`, budgetID, limit.CategoryID, limit.Amount)

	require.NoError(t, err)
}

func prepareTestCategory(t *testing.T, db *pgxpool.Pool, category category.Category) {
	_, err := db.Exec(context.Background(), `
				insert into category ("ID", "name", "logo", "createdAt")
				values ($1, $2, $3, current_timestamp(0))`, category.ID, category.Name, category.Logo)

	require.NoError(t, err)
}
