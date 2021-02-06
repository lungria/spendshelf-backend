package pgtest_test

import (
	"testing"

	"github.com/lungria/spendshelf-backend/storage/pgtest"
)

func TestCreateDb(t *testing.T) {
	_, clear := pgtest.PrepareWithSchema(t, "../schema/schema.sql")
	clear()
}
