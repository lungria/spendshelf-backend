package mono_test

import (
	"context"
	"fmt"
	"github.com/lungria/spendshelf-backend/mono"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestClient_GetTransactions_NoErrorsReturned(t *testing.T) {
	key := os.Getenv("MONO_API_KEY")
	if key == "" {
		t.Fatal("MONO_API_KEY is required for test")
		return
	}
	client := mono.NewClient(key)

	transactions, err := client.GetTransactions(context.Background())

	assert.NoError(t, err)
	fmt.Printf("%v\n", transactions)
}
