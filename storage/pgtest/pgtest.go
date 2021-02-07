package pgtest

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"

	"github.com/caarlos0/env/v6"

	"github.com/jackc/pgx/v4/pgxpool"
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
	return fmt.Sprintf("postgres://%s/%s?sslmode=disable&user=%s&password=%s", conStr.Host, conStr.DB, conStr.Username, conStr.Password)
}

func prepare(t *testing.T) (*pgxpool.Pool, func()) {
	conStr := getConnectionString(t)

	// mainPool would be used to create temporary testing database
	mainPool, err := pgxpool.Connect(context.Background(), conStr.String())
	if err != nil {
		t.Fatalf("expected nil, found: %v", err)
	}

	// create temporary testing database
	random := rand.Intn(9999)
	dbName := fmt.Sprintf("pgtest%v", random)
	_, err = mainPool.Exec(context.Background(), fmt.Sprintf("create database %s", dbName))
	if err != nil {
		mainPool.Close()
		t.Fatalf("expected nil, found: %v", err)
	}
	_, err = mainPool.Exec(context.Background(), fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", dbName, conStr.Username))
	if err != nil {
		mainPool.Close()
		t.Fatalf("expected nil, found: %v", err)
	}

	// connect to temporary testing database
	testDB, err := pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://localhost:5432/%s?sslmode=disable&user=postgres&password=adminpass123", dbName))
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
