// Package pgtest provides simple utility functions for creating temporary databases in PostgreSQL.
// It was designed to be used in integration tests in the following way:
//
//     func TestSomethingThatUsesPostgres(t *testing.T) {
//
//         // create pool and schedule cleanup function
//         pool, cleanup := pgtest.PrepareWithSchema(t, "schema/schema.sql")
//         defer cleanup()
//         // use pool in your code
//         accountID := prepareTestAccount(t, pool)
//         ...
//     }
//
// It's not production ready and was specifically tailored for spendshelf-backend requirements.
package pgtest

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	mutex      = &sync.Mutex{}
)

type config struct {
	Username string `env:"PG_TEST_USERNAME,required"`
	Password string `env:"PG_TEST_PASSWORD,required"`
	Host     string `env:"PG_TEST_HOST,required"`
}

type connectionString struct {
	Username string
	Password string
	Host     string
	DB       string
}

func (conStr connectionString) String() string {
	return fmt.Sprintf(
		"postgres://%s/%s?sslmode=disable&user=%s&password=%s",
		conStr.Host,
		conStr.DB,
		conStr.Username,
		conStr.Password)
}

func prepare(t *testing.T) (*pgxpool.Pool, func()) {
	conStr := getConnectionString(t)

	// mainPool would be used to create temporary testing database
	mainPool, err := pgxpool.Connect(context.Background(), conStr.String())
	if err != nil {
		t.Fatalf("expected nil, found: %v", err)
	}

	// create temporary testing database
	mutex.Lock()
	dbName := fmt.Sprintf("pgtest%v", seededRand.Uint64())
	mutex.Unlock()

	_, err = mainPool.Exec(context.Background(), fmt.Sprintf("create database %s", dbName))
	if err != nil {
		mainPool.Close()
		t.Fatalf("expected nil, found: %v", err)
	}

	conStr.DB = dbName

	// connect to temporary testing database
	testDB, err := pgxpool.Connect(context.Background(), conStr.String())
	if err != nil {
		mainPool.Close()
		t.Fatalf("expected nil, found: %v", err)
	}

	// prepare cleanup function that
	cleanup := func() {
		defer mainPool.Close()
		testDB.Close()

		_, err = mainPool.Exec(context.Background(), fmt.Sprintf("drop database %s;", dbName))
		if err != nil {
			t.Fatalf("expected nil, found: %v", err)
		}
	}

	return testDB, cleanup
}

func getConnectionString(t *testing.T) connectionString {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		t.Fatalf("expected nil, found: %v", err)
	}

	return connectionString{
		Username: cfg.Username,
		Password: cfg.Password,
		Host:     cfg.Host,
		DB:       "postgres",
	}
}

// PrepareWithSchema returns pgxpool.Pool instance and cleanup function, that can be used in integration tests.
func PrepareWithSchema(t *testing.T, schemaFilePath string) (*pgxpool.Pool, func()) {
	db, cleanup := prepare(t)

	content, err := ioutil.ReadFile(schemaFilePath)
	if err != nil {
		cleanup()
		t.Fatalf("expected nil, found: %v", err)
	}

	_, err = db.Exec(context.Background(), string(content))
	if err != nil {
		cleanup()
		t.Fatalf("expected nil, found: %v", err)
	}

	return db, cleanup
}
