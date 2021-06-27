package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lungria/spendshelf-backend/api"
	"github.com/lungria/spendshelf-backend/storage"
	"github.com/rs/zerolog/log"
)

// UpdateTransactionBody describes request body for update transaction request.
type UpdateTransactionBody struct {
	CategoryID *int32  `json:"categoryId"`
	Comment    *string `json:"comment"`
}

// UpdateTransactionQuery describes request query for update transaction request.
type UpdateTransactionQuery struct {
	LastUpdatedAt time.Time `form:"lastUpdatedAt" binding:"required"`
}

// GetReportQuery describes request query for transaction report request.
type GetReportQuery struct {
	From time.Time `form:"from" binding:"required"`
	To   time.Time `form:"to" binding:"required"`
}

// TransactionStorage abstracts transactions storage implementation.
type TransactionStorage interface {
	Get(ctx context.Context, query storage.Query, page storage.Page) ([]storage.Transaction, error)
	UpdateTransaction(ctx context.Context, cmd storage.UpdateTransactionCommand) (storage.Transaction, error)
	GetReport(ctx context.Context, from, to time.Time) (map[int32]int64, error)
}

// CategoryStorage abstracts categories storage implementation.
type CategoryStorage interface {
	GetAll(ctx context.Context) ([]storage.Category, error)
}

// TransactionHandler handles /vN/transaction routes.
type TransactionHandler struct {
	transactions TransactionStorage
	categories   CategoryStorage
}

// NewTransactionHandler returns new instance of TransactionHandler.
func NewTransactionHandler(transactions TransactionStorage, categories CategoryStorage) *TransactionHandler {
	return &TransactionHandler{transactions: transactions, categories: categories}
}

// GetTransactionsQuery describes request query for transaction list.
type GetTransactionsQuery struct {
	From       time.Time `form:"from" binding:"required"`
	To         time.Time `form:"to" binding:"required"`
	CategoryID int32     `form:"categoryId" binding:"required"`
}

// GetTransactions returns transactions (without category).
func (t *TransactionHandler) GetTransactions(c *gin.Context) {
	query := GetTransactionsQuery{}

	if err := c.BindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, api.Error{Message: fmt.Sprintf("unable to parse query: %s", err)})
		return
	}

	result, err := t.transactions.Get(c, storage.Query{CategoryID: query.CategoryID}, storage.Page{})
	if err != nil {
		log.Error().Err(err).Msg("unable to load transactions from storage")
		c.JSON(api.InternalServerError())

		return
	}

	c.JSON(http.StatusOK, &result)
}

// PatchTransaction allows to update single transaction.
// Transaction is being selected by id route parameter and filtered by lastUpdatedAt query parameter.
// lastUpdatedAt filtering protects us from concurrent updates issues (simplest implementation of optimistic
// concurrency). Body can contain optional parameters (see UpdateTransactionBody fields). Patch will update
// only not nil fields.
func (t *TransactionHandler) PatchTransaction(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, api.Error{Message: "id required"})
		return
	}

	var query UpdateTransactionQuery
	if err := c.BindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, api.Error{Message: "lastUpdatedAt must be valid time"})
		return
	}

	var req UpdateTransactionBody
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.Error{Message: "unable to parse json"})
		return
	}

	result, err := t.transactions.UpdateTransaction(c, storage.UpdateTransactionCommand{
		Query: storage.Query{
			ID:            id,
			LastUpdatedAt: query.LastUpdatedAt,
		},
		UpdatedFields: storage.UpdatedFields{
			CategoryID: req.CategoryID,
			Comment:    req.Comment,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to update transaction in storage")
		c.JSON(api.InternalServerError())

		return
	}

	c.JSON(http.StatusOK, &result)
}

// GetReport returns monthly spendings report.
func (t *TransactionHandler) GetReport(c *gin.Context) {
	var query GetReportQuery
	if err := c.BindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, api.Error{Message: "from/to query parameters are required"})
		return
	}

	result, err := t.transactions.GetReport(c, query.From, query.To)
	if err != nil {
		log.Error().Err(err).Msg("unable to load transactions from storage")
		c.JSON(api.InternalServerError())

		return
	}

	c.JSON(http.StatusOK, &result)
}

// GetCategories returns list of existing categories.
func (t *TransactionHandler) GetCategories(c *gin.Context) {
	result, err := t.categories.GetAll(c)
	if err != nil {
		log.Error().Err(err).Msg("unable to load categories from storage")
		c.JSON(api.InternalServerError())

		return
	}

	c.JSON(http.StatusOK, &result)
}

// BindRoutes bind gin routes to handler methods.
func (t *TransactionHandler) BindRoutes(router *gin.Engine) {
	router.GET("/v1/transactions", t.GetTransactions)
	router.PATCH("/v1/transactions/:id", t.PatchTransaction)
	router.GET("/v1/transactions/report", t.GetReport)
	router.GET("/v1/transactions/categories", t.GetCategories)
}
