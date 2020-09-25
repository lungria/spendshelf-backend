package mono_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lungria/spendshelf-backend/mono"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetTransactions_NoErrorsReturned(t *testing.T) {
	key := os.Getenv("MONO_API_KEY")
	if key == "" {
		t.Fatal("MONO_API_KEY is required for test")
		return
	}
	client := mono.NewClient("https://api.monobank.ua", key)

	transactions, err := client.GetTransactions(context.Background(), mono.GetTransactionsQuery{
		Account: "0",
		From:    time.Now().Add(-48 * time.Hour),
		To:      time.Now(),
	})

	assert.NoError(t, err)
	fmt.Printf("%v\n", transactions)
}
