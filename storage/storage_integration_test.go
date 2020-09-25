package storage

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/transaction"
	"github.com/stretchr/testify/assert"
)

const dbConnString = "postgres://localhost:5432/postgres?sslmode=disable"

func TestSave_WithLocalDb_NoErrorReturned(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), dbConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	storage := &PostgreSQLStorage{pool: dbpool}
	err = storage.Save(context.Background(), []transaction.Transaction{
		{"id1", time.Now().UTC(), "food", 123, true, 1110},
		{"id1", time.Now().UTC(), "food", 123, true, 1110},
		{"id2", time.Now().UTC(), "car", 3121, true, 1500},
		{"id3", time.Now().UTC(), "home", 3, false, 2000},
	})

	assert.NoError(t, err)
}
