package api

import "github.com/lungria/spendshelf-backend/src/api/handlers"

func RoutesProvider(w *handlers.WebHookHandler, c *handlers.CategoriesHandler, t *handlers.TransactionsHandler, r *handlers.ReportsHandler) []RouterBinder {
	return []RouterBinder{
		w, c, t, r,
	}
}
