package api

import (
	"github.com/lungria/spendshelf-backend/src/categories"
	"github.com/lungria/spendshelf-backend/src/transactions"
)

func RoutesProvider(c *categories.Handler, t *transactions.Handler) []RouterBinder {
	return []RouterBinder{
		c, t,
	}
}
