package transactions

import (
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

type patchRequest struct {
	CategoryID string `json:"categoryId" binding:"required"`
}

type getResponse struct {
	Transactions []Transaction `json:"transactions"`
}

// Handler is a struct which implemented by transactions handler
type Handler struct {
	txnRepo *Repository
	logger  *zap.SugaredLogger
}

// NewHandler create a new instance of Handler
func NewHandler(txnRepo *Repository, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		txnRepo: txnRepo,
		logger:  logger,
	}
}

// Get can return all transactions, only categorized transactions, only uncategorized transactions and transactions interrelated with one category.
// /transactions or /transactions returned all transactions
// /transactions?category=without returned all uncategorized transactions
// /transactions?category=categoryID returned all transactions which related with one specify category
func (handler *Handler) Get(c *gin.Context) {
	category := c.Request.URL.Query().Get("category")

	switch category {
	case "":
		handler.allTransactions(c)
		return
	case "without":
		handler.onlyUncategorizedTransactions(c)
		return
	default:
		handler.oneCategoryTransactions(c, category)
		return
	}
}

// Patch is setting or changing a category for specify transactionResponse
func (handler *Handler) Patch(c *gin.Context) {
	var req patchRequest
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

	id, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("invalid id: %w", err))
		return
	}

	countModifiedDocs, err := handler.txnRepo.UpdateCategory(c, tObjID, id)
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

func (handler *Handler) allTransactions(c *gin.Context) {
	t, err := handler.txnRepo.FindAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("unable to received all transactions: %w", err))
		return
	}

	c.JSON(http.StatusOK, getResponse{Transactions: t})
}

func (handler *Handler) onlyUncategorizedTransactions(c *gin.Context) {
	t, err := handler.txnRepo.FindAllUncategorized(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("unable to received uncategorized transactions: %w", err))
		return
	}

	c.JSON(http.StatusOK, getResponse{Transactions: t})
}

func (handler *Handler) oneCategoryTransactions(c *gin.Context, categoryID string) {
	id, err := primitive.ObjectIDFromHex(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("invalid id: %w", err))
		return
	}
	t, err := handler.txnRepo.FindByCategoryID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("unable to find transactions for specify category: %w", err))
		return
	}

	c.JSON(http.StatusOK, getResponse{Transactions: t})
}

// Get allows to save transaction.
func (handler *Handler) Post(c *gin.Context) {
	var req Transaction
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = handler.txnRepo.Insert(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("unable to save transactions: %w", err))
		return
	}

	c.Status(http.StatusCreated)
}

func (handler *Handler) BindRoutes(router *gin.Engine) {
	router.GET("/transactions", handler.Get)
	router.PATCH("/transactions/:transactionID", handler.Patch)
	router.POST("/transactions", handler.Post)
}
