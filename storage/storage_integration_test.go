package storage_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/storage"
	"github.com/lungria/spendshelf-backend/storage/category"
	"github.com/stretchr/testify/assert"
)

// todo: setup before test, clear after test
const dbConnString = "postgres://localhost:5432/spendshelf-test?sslmode=disable"

func TestSave_WithLocalDb_NoErrorReturned(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), dbConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	db := storage.NewPostgreSQLStorage(dbpool)

	err = db.Save(context.Background(), []storage.Transaction{{
		"id1",
		time.Now().UTC(),
		"food",
		123,
		true,
		1110,
		"acc1",
		category.Default,
		time.Now().UTC(),
		nil,
	}, {
		"id1",
		time.Now().UTC(),
		"food",
		123,
		true,
		1110,
		"acc1",
		category.Default,
		time.Now().UTC(),
		nil,
	}, {
		"id2",
		time.Now().UTC(),
		"car",
		3121,
		true,
		1500,
		"acc1",
		category.Default,
		time.Now().UTC(),
		nil,
	}, {
		"id3",
		time.Now().UTC(),
		"home",
		3,
		false,
		2000,
		"acc1",
		category.Default,
		time.Now().UTC(),
		nil,
	}})

	assert.NoError(t, err)
}

func TestGetLastTransactionDate_WithLocalDb_NoErrorReturned(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), dbConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	db := storage.NewPostgreSQLStorage(dbpool)

	_, err = db.GetLastTransactionDate(context.Background(), "acc1")

	assert.NoError(t, err)
}

// todo: this test doesn't work, because storage.Save method ignores lastUpdatedAt, need to add some workaround
func TestUpdate_WithLocalDb_NoErrorReturned(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), dbConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	db := storage.NewPostgreSQLStorage(dbpool)
	lastUpdatedAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	err = db.Save(context.Background(), []storage.Transaction{{
		"id4",
		time.Now().UTC(),
		"food",
		123,
		true,
		1110,
		"acc1",
		category.Default,
		lastUpdatedAt,
		nil,
	}})

	// update comment without category
	comment := "comment"
	_, err = db.UpdateTransaction(context.Background(), storage.UpdateTransactionCommand{
		Comment:       &comment,
		ID:            "id4",
		LastUpdatedAt: lastUpdatedAt,
	})
	assert.NoError(t, err)
	updatedTransaction, err := db.GetByID(context.Background(), "id4")
	assert.NoError(t, err)
	assert.Equal(t, "comment", *updatedTransaction.Comment)
	assert.Equal(t, category.Default, updatedTransaction.CategoryID)
	// update category without comment
	category := int32(10)
	_, err = db.UpdateTransaction(context.Background(), storage.UpdateTransactionCommand{
		CategoryID:    &category,
		ID:            "id4",
		LastUpdatedAt: updatedTransaction.LastUpdatedAt,
	})
	assert.NoError(t, err)
	updatedTransaction, err = db.GetByID(context.Background(), "id4")
	assert.NoError(t, err)
	assert.Equal(t, "comment", *updatedTransaction.Comment)
	assert.Equal(t, 10, updatedTransaction.CategoryID)
}
