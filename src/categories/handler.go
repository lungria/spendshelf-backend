package categories

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type getResponse struct {
	Categories []Category `json:"categories"`
}

type createRequest struct {
	Name string `json:"name"`
}

type createResponse struct {
	ID primitive.ObjectID `json:"id"`
}

// Handler is a struct which implemented by categories handler
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

// Get return all existing categories
func (handler *Handler) Get(c *gin.Context) {
	ctg, err := handler.repo.GetAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getResponse{ctg})
	return
}

// Post create a new category
func (handler *Handler) Post(c *gin.Context) {
	var req createRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	id, err := handler.repo.Insert(c, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		handler.logger.Error(err)
		return
	}
	c.JSON(http.StatusOK, createResponse{id})
}

func (handler *Handler) BindRoutes(router *gin.Engine) {
	router.POST("/categories", handler.Post)
	router.GET("/categories", handler.Get)
}
