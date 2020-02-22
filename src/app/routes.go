package app

import (
	"github.com/lungria/spendshelf-backend/src/categories"
	"github.com/lungria/spendshelf-backend/src/transactions"
)

func RoutesProvider(t *transactions.Handler, c *categories.Handler) []RouterBinder {
	return []RouterBinder{
		t, c,
	}
}
