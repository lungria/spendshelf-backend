package storage_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/category"
	"github.com/lungria/spendshelf-backend/storage"
	"github.com/stretchr/testify/assert"
)

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
