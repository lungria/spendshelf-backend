package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lungria/spendshelf-backend/src/transactions"
	"go.uber.org/zap"
)

type patchCategoryRequest struct {
	Category string `json:"category" binding:"required"`
}

type transactionResponse struct {
	ID          string `json:"id"`
	Time        int32  `json:"time"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Amount      int64  `json:"amount"`
	Balance     int64  `json:"balance"`
}

type getTransactionsResponse struct {
	Transactions []transactionResponse `json:"transactions"`
}

// TransactionsHandler is a struct which implemented by transactions handler
type TransactionsHandler struct {
	repo   transactions.Repository
	logger *zap.SugaredLogger
}

// NewTransactionsHandler create a new instance of TransactionsHandler
func NewTransactionsHandler(repo transactions.Repository, logger *zap.SugaredLogger) (*TransactionsHandler, error) {
	if repo == nil {
		return nil, errors.New("repository must not be nil")
	}

	if logger == nil {
		return nil, errors.New("logger must not be nil (Transactions)")
	}

	return &TransactionsHandler{
		repo:   repo,
		logger: logger,
	}, nil
}

// HandleGet can return all transactions, only categorized transactions, only uncategorized transactions and transactions interrelated with one category.
// /transactions or /transactions?category= returned all transactions
// /transactions?category=with returned all categorized transactions
// /transactions?category=without returned all uncategorized transactions
// /transactions?category=SomeCategory returned all transactions which related with one specify category
func (handler *TransactionsHandler) HandleGet(c *gin.Context) {
	category := c.Request.URL.Query().Get("category")
	handler.logger.Info(category)

	switch category {
	case "":
		handler.allTransactions(c)
		return
	case "without":
		handler.onlyUncategorizedTransactions(c)
		return
	case "with":
		handler.onlyCategorizedTransactions(c)
		return
	default:
		handler.oneCategoryTransactions(c, category)
		return
	}
}

// HandlePatch is setting or changing a category for specify transactionResponse
func (handler *TransactionsHandler) HandlePatch(c *gin.Context) {
	var req patchCategoryRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "Bad request", Error: err.Error()})
		return
	}
	transactionID := c.Param("transactionID")
	err = handler.repo.UpdateCategory(transactionID, req.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "Update failed", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (handler *TransactionsHandler) allTransactions(c *gin.Context) {
	t, err := handler.repo.FindAll()
	handler.logger.Info(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "Unable to received all transactions", Error: err.Error()})
		return
	}
	tArr := handler.prepareResponse(t)
	c.JSON(http.StatusOK, getTransactionsResponse{Transactions: tArr})
	return
}

func (handler *TransactionsHandler) onlyCategorizedTransactions(c *gin.Context) {
	t, err := handler.repo.FindAllCategorized()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "Unable to received categorized transactions", Error: err.Error()})
		return
	}
	tArr := handler.prepareResponse(t)
	c.JSON(http.StatusOK, getTransactionsResponse{Transactions: tArr})
	return
}

func (handler *TransactionsHandler) onlyUncategorizedTransactions(c *gin.Context) {
	t, err := handler.repo.FindAllUncategorized()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "Unable to received uncategorized transactions", Error: err.Error()})
		return
	}
	tArr := handler.prepareResponse(t)
	c.JSON(http.StatusOK, getTransactionsResponse{Transactions: tArr})
	return
}

func (handler *TransactionsHandler) oneCategoryTransactions(c *gin.Context, category string) {
	t, err := handler.repo.FindAllByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "Unable to received transactions for specify category " + category, Error: err.Error()})
		return
	}
	tArr := handler.prepareResponse(t)
	c.JSON(http.StatusOK, getTransactionsResponse{Transactions: tArr})
	return
}

func (handler *TransactionsHandler) prepareResponse(rawTxn []transactions.Transaction) []transactionResponse {
	processedTxn := make([]transactionResponse, len(rawTxn))

	for i := 0; i < len(rawTxn); i++ {
		processedTxn[i].ID = rawTxn[i].ID.Hex()
		processedTxn[i].Time = rawTxn[i].StatementItem.Time
		processedTxn[i].Description = rawTxn[i].StatementItem.Description
		processedTxn[i].Category = rawTxn[i].Category
		processedTxn[i].Amount = rawTxn[i].StatementItem.Amount
		processedTxn[i].Balance = rawTxn[i].StatementItem.Balance
	}
	return processedTxn
}
