package storage_test

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func prepare(t *testing.T) (*pgxpool.Pool, func()) {
	mainPool, err := pgxpool.Connect(context.Background(), "postgres://localhost:5432/postgres?sslmode=disable&user=postgres&password=adminpass123")
	assert.Nil(t, err)

	random := rand.Intn(999)
	dbName := fmt.Sprintf("test%v", random)
	sql := fmt.Sprintf("create db %s;", dbName)
	_, err = mainPool.Exec(context.Background(), sql)

	assert.Nil(t, err)
	cleanup := func() {
		defer mainPool.Close()
		_, err = mainPool.Exec(context.Background(), fmt.Sprintf("dropdb %s;", dbName))
		assert.Nil(t, err)
	}

	testDB, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://localhost:5432/%s?sslmode=disable&user=postgres&password=adminpass123", dbName))
	assert.Nil(t, err)

	return testDB, cleanup
}

func Test_Select_1__WithLocalDb__NoErrors(t *testing.T) {
	db, cleanup := prepare(t)
	defer cleanup()

	_, err := db.Query(context.Background(), "select 1;")

	assert.Nil(t, err)
}


//
//func TestSave_WithLocalDb_NoErrorReturned(t *testing.T) {
//	dbpool, err := pgxpool.Connect(context.Background(), dbConnString)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
//		os.Exit(1)
//	}
//	defer dbpool.Close()
//	db := storage.NewPostgreSQLStorage(dbpool)
//
//	err = db.Save(context.Background(), []storage.Transaction{{
//		"id1",
//		time.Now().UTC(),
//		"food",
//		123,
//		true,
//		1110,
//		"acc1",
//		category.Default,
//		time.Now().UTC(),
//		nil,
//	}, {
//		"id1",
//		time.Now().UTC(),
//		"food",
//		123,
//		true,
//		1110,
//		"acc1",
//		category.Default,
//		time.Now().UTC(),
//		nil,
//	}, {
//		"id2",
//		time.Now().UTC(),
//		"car",
//		3121,
//		true,
//		1500,
//		"acc1",
//		category.Default,
//		time.Now().UTC(),
//		nil,
//	}, {
//		"id3",
//		time.Now().UTC(),
//		"home",
//		3,
//		false,
//		2000,
//		"acc1",
//		category.Default,
//		time.Now().UTC(),
//		nil,
//	}})
//
//	assert.NoError(t, err)
//}
//
//func TestGetLastTransactionDate_WithLocalDb_NoErrorReturned(t *testing.T) {
//	dbpool, err := pgxpool.Connect(context.Background(), dbConnString)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
//		os.Exit(1)
//	}
//	defer dbpool.Close()
//	db := storage.NewPostgreSQLStorage(dbpool)
//
//	_, err = db.GetLastTransactionDate(context.Background(), "acc1")
//
//	assert.NoError(t, err)
//}

//
//// todo: this test doesn't work, because storage.Save method ignores lastUpdatedAt, need to add some workaround
//func TestUpdate_WithLocalDb_NoErrorReturned(t *testing.T) {
//	dbpool, err := pgxpool.Connect(context.Background(), dbConnString)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
//		os.Exit(1)
//	}
//	defer dbpool.Close()
//	db := storage.NewPostgreSQLStorage(dbpool)
//	lastUpdatedAt := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
//	err = db.Save(context.Background(), []storage.Transaction{{
//		"id4",
//		time.Now().UTC(),
//		"food",
//		123,
//		true,
//		1110,
//		"acc1",
//		category.Default,
//		lastUpdatedAt,
//		nil,
//	}})
//
//	// update comment without category
//	comment := "comment"
//	_, err = db.UpdateTransaction(context.Background(), storage.UpdateTransactionCommand{
//		Comment:       &comment,
//		ID:            "id4",
//		LastUpdatedAt: lastUpdatedAt,
//	})
//	assert.NoError(t, err)
//	updatedTransaction, err := db.GetByID(context.Background(), "id4")
//	assert.NoError(t, err)
//	assert.Equal(t, "comment", *updatedTransaction.Comment)
//	assert.Equal(t, category.Default, updatedTransaction.CategoryID)
//	// update category without comment
//	category := int32(10)
//	_, err = db.UpdateTransaction(context.Background(), storage.UpdateTransactionCommand{
//		CategoryID:    &category,
//		ID:            "id4",
//		LastUpdatedAt: updatedTransaction.LastUpdatedAt,
//	})
//	assert.NoError(t, err)
//	updatedTransaction, err = db.GetByID(context.Background(), "id4")
//	assert.NoError(t, err)
//	assert.Equal(t, "comment", *updatedTransaction.Comment)
//	assert.Equal(t, 10, updatedTransaction.CategoryID)
//}
