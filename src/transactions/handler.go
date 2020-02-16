package transactions

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

type patchRequest struct {
	CategoryID    uint8     `json:"categoryId" binding:"required"`
	TransactionID time.Time `json:"transactionId" binding:"required"`
}

type getResponse struct {
	Transactions []Transaction `json:"transactions"`
}

// Handler is a struct which implemented by transactions handler
type Handler struct {
	store  *Store
	logger *zap.SugaredLogger
}

// NewHandler create a new instance of Handler
func NewHandler(store *Store, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		store:  store,
		logger: logger,
	}
}

// GetUncategorized can return uncategorized transactions.
func (handler *Handler) GetUncategorized(c *gin.Context) {
	transactions, err := handler.store.ReadUncategorized()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, getResponse{Transactions: transactions})
}

// GetUncategorized can return uncategorized transactions.
func (handler *Handler) GetCategorized(c *gin.Context) {
	transactions, err := handler.store.ReadUncategorized()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, getResponse{Transactions: transactions})
}

// Patch is setting or changing a category for specify transactionResponse
func (handler *Handler) Patch(c *gin.Context) {
	var req patchRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = handler.store.SetCategory(req.TransactionID, req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

// Get allows to save transaction.
func (handler *Handler) Post(c *gin.Context) {
	var req Transaction
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = handler.store.Insert(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusCreated)
}

func (handler *Handler) BindRoutes(router *gin.Engine) {
	router.GET("/transactions/uncategorized", handler.GetUncategorized)
	router.PATCH("/transactions/:transactionID", handler.Patch)
	router.POST("/transactions", handler.Post)
	router.GET("/transactions/categorized", handler.GetCategorized)
}
