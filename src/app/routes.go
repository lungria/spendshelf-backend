package app

import (
	"github.com/lungria/spendshelf-backend/src/transactions"
)

func RoutesProvider(t *transactions.Handler) []RouterBinder {
	return []RouterBinder{
		t,
	}
}
