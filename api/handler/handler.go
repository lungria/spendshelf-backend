package handler

import (
	"context"
	"net/http"

	"github.com/lungria/spendshelf-backend/api"

	"github.com/rs/zerolog/log"

	"github.com/lungria/spendshelf-backend/storage/transaction"

	"github.com/gin-gonic/gin"
)

// TransactionStorage abstracts storage implementation
type TransactionStorage interface {
	GetByCategory(ctx context.Context, categoryID int32) ([]transaction.Transaction, error)
}

// TransactionHandler handles /vN/transaction routes.
type TransactionHandler struct {
	storage TransactionStorage
}

// NewTransactionHandler returns new instance of TransactionHandler.
func NewTransactionHandler(storage TransactionStorage) *TransactionHandler {
	return &TransactionHandler{storage: storage}
}

// GetTransactions returns transactions (without category).
func (t *TransactionHandler) GetTransactions(c *gin.Context) {
	result, err := t.storage.GetByCategory(c, transaction.DefaultCategoryID)
	if err != nil {
		log.Error().Err(err).Msg("failed to query transactions")
		c.JSON(http.StatusInternalServerError, api.Error{Message: "unable to load transactions from database"})
		return
	}

	c.JSON(http.StatusOK, &result)
}

// BindRoutes bind gin routes to handler methods.
func (t *TransactionHandler) BindRoutes(router *gin.Engine) {
	router.GET("/v1/transactions/", t.GetTransactions)
}
