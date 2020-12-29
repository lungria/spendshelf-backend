package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/lungria/spendshelf-backend/api"
	"github.com/lungria/spendshelf-backend/category"
	"github.com/lungria/spendshelf-backend/storage"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

// UpdateTransactionBody describes request body for update transaction request.
type UpdateTransactionBody struct {
	CategoryID *int32  `json:"categoryID"`
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

// TransactionStorage abstracts storage implementation
type TransactionStorage interface {
	GetByCategory(ctx context.Context, categoryID int32) ([]storage.Transaction, error)
	UpdateTransaction(ctx context.Context, params storage.UpdateTransactionCommand) (storage.Transaction, error)
	GetReport(ctx context.Context, from, to time.Time) (map[int32]int64, error)
	GetCategories(ctx context.Context) ([]storage.Category, error)
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
	// todo: bind category ID from query
	result, err := t.storage.GetByCategory(c, category.Default)
	if err != nil {
		log.Error().Err(err).Msg("failed to query transactions")
		c.JSON(http.StatusInternalServerError, api.Error{Message: "unable to load transactions from database"})

		return
	}

	c.JSON(http.StatusOK, &result)
}

// PatchTransaction allows to update
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

	result, err := t.storage.UpdateTransaction(c, storage.UpdateTransactionCommand{
		ID:            id,
		LastUpdatedAt: query.LastUpdatedAt,
		CategoryID:    req.CategoryID,
		Comment:       req.Comment,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to update transaction")
		c.JSON(http.StatusInternalServerError, api.Error{Message: "failed to update transaction in database"})

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

	result, err := t.storage.GetReport(c, query.From, query.To)
	if err != nil {
		log.Error().Err(err).Msg("failed to query transactions")
		c.JSON(http.StatusInternalServerError, api.Error{Message: "unable to load transactions from database"})

		return
	}

	c.JSON(http.StatusOK, &result)
}

// GetCategories returns list of existing categories.
func (t *TransactionHandler) GetCategories(c *gin.Context) {
	result, err := t.storage.GetCategories(c)
	if err != nil {
		log.Error().Err(err).Msg("failed to query categories")
		c.JSON(http.StatusInternalServerError, api.Error{Message: "unable to load categories from database"})

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
