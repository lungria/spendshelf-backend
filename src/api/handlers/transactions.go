package handlers

import (
	"fmt"
	"net/http"

	"github.com/lungria/spendshelf-backend/src/categories"

	"github.com/lungria/spendshelf-backend/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"

	"github.com/lungria/spendshelf-backend/src/transactions"
	"go.uber.org/zap"
)

type patchCategoryRequest struct {
	CategoryID string `json:"categoryId" binding:"required"`
}

type getTransactionsResponse struct {
	Transactions []models.Transaction `json:"transactions"`
}

// TransactionsHandler is a struct which implemented by transactions handler
type TransactionsHandler struct {
	txnRepo transactions.Repository
	ctgRepo *categories.Repository
	logger  *zap.SugaredLogger
}

// NewTransactionsHandler create a new instance of TransactionsHandler
func NewTransactionsHandler(txnRepo transactions.Repository, ctgRepo *categories.Repository, logger *zap.SugaredLogger) *TransactionsHandler {
	return &TransactionsHandler{
		txnRepo: txnRepo,
		ctgRepo: ctgRepo,
		logger:  logger,
	}
}

// HandleGet can return all transactions, only categorized transactions, only uncategorized transactions and transactions interrelated with one category.
// /transactions or /transactions?category= returned all transactions
// /transactions?category=with returned all categorized transactions
// /transactions?category=without returned all uncategorized transactions
// /transactions?category=categoryID returned all transactions which related with one specify category
func (handler *TransactionsHandler) HandleGet(c *gin.Context) {
	category := c.Request.URL.Query().Get("category")

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
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	transactionID := c.Param("transactionID")
	tObjID, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		handler.logger.Errorw("Transaction ID wrong or invalid", "TransactionID", tObjID, "Error", err)
		c.JSON(http.StatusBadRequest, fmt.Errorf("invalid id: %w", err).Error())
		return
	}

	var ctg models.Category
	ok := handler.findCategoryByID(c, req.CategoryID, &ctg)
	if !ok {
		return
	}

	countModifiedDocs, err := handler.txnRepo.UpdateCategory(tObjID, ctg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("update failed: %w", err).Error())
		return
	}
	if countModifiedDocs == 0 {
		c.JSON(http.StatusNotFound, fmt.Errorf("transaction not found: %w", err).Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (handler *TransactionsHandler) allTransactions(c *gin.Context) {
	t, err := handler.txnRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("unable to received all transactions: %w", err))
		return
	}

	c.JSON(http.StatusOK, getTransactionsResponse{Transactions: t})
	return
}

func (handler *TransactionsHandler) onlyCategorizedTransactions(c *gin.Context) {
	t, err := handler.txnRepo.FindAllCategorized()
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("unable to received categorized transactions: %w", err))
		return
	}

	c.JSON(http.StatusOK, getTransactionsResponse{Transactions: t})
	return
}

func (handler *TransactionsHandler) onlyUncategorizedTransactions(c *gin.Context) {
	t, err := handler.txnRepo.FindAllUncategorized()
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("unable to received uncategorized transactions: %w", err))
		return
	}

	c.JSON(http.StatusOK, getTransactionsResponse{Transactions: t})
	return
}

func (handler *TransactionsHandler) oneCategoryTransactions(c *gin.Context, categoryID string) {
	ctg := models.Category{}
	ok := handler.findCategoryByID(c, categoryID, &ctg)
	if !ok {
		return
	}
	t, err := handler.txnRepo.FindAllByCategoryID(ctg.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("unable to received transactions for specify category: %w", err))
		return
	}

	c.JSON(http.StatusOK, getTransactionsResponse{Transactions: t})
	return
}

func (handler *TransactionsHandler) findCategoryByID(c *gin.Context, categoryID string, category *models.Category) bool {
	ctgObjID, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		handler.logger.Errorw("CategoryID wrong or invalid", "CategoryID", ctgObjID, "Error", err)
		c.JSON(http.StatusBadRequest, fmt.Errorf("invalid id: %w", err))
		return false
	}
	category = handler.ctgRepo.FindByID(c, ctgObjID)
	if category == nil {
		handler.logger.Infow("Category not found by ID", "CategoryID", ctgObjID)
		c.JSON(http.StatusNotFound, fmt.Errorf("unable to find category: %w", err))
		return false
	}

	return true
}

func (handler *TransactionsHandler) BindRoutes(router *gin.Engine) {
	router.GET("/transactions", handler.HandleGet)
	router.PATCH("/transactions/:transactionID", handler.HandlePatch)
}
