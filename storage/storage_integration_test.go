package storage

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
)

const dbConnString = "postgres://localhost:5432/postgres?sslmode=disable"

func TestSave_IntegrationTest(t *testing.T) {
	dbpool, err := pgxpool.Connect(context.Background(), dbConnString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	storage := &Storage{pool: dbpool}
	err = storage.Save(context.Background(), nil)
	fmt.Println(err)
}
