package mono_test

import (
	"context"
	"os"
	"testing"

	"github.com/lungria/spendshelf-backend/importer/mono"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetUserInfo_NoErrorsReturned(t *testing.T) {
	key := os.Getenv("MONO_API_KEY")
	if key == "" {
		t.Fatal("MONO_API_KEY is required for test")
		return
	}

	client := mono.NewClient("https://api.monobank.ua", key)

	_, err := client.GetUserInfo(context.Background())

	assert.NoError(t, err)
}
