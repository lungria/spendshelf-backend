package transactions

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

type patchRequest struct {
	CategoryID primitive.ObjectID `json:"categoryId" binding:"required"`
}

type getResponse struct {
	Transactions []Transaction `json:"transactions"`
}

// Handler is a struct which implemented by transactions handler
type Handler struct {
	repo   *Repository
	logger *zap.SugaredLogger
}

// NewHandler create a new instance of Handler
func NewHandler(repo *Repository, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		repo:   repo,
		logger: logger,
	}
}

// GetUncategorized can return uncategorized transactions.
func (handler *Handler) GetUncategorized(c *gin.Context) {
	transactions, err := handler.repo.ReadUncategorized(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	if transactions == nil {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, getResponse{Transactions: transactions})
}

// Patch is setting or changing a category for specify transactionResponse
func (handler *Handler) Patch(c *gin.Context) {
	param := c.Param("transactionID")
	transactionID, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var req patchRequest
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = handler.repo.SetCategory(c, transactionID, req.CategoryID)
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

	err = handler.repo.Insert(c, &req)
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
}
