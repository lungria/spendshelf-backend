package handlers

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/lungria/spendshelf-backend/src/categories"

	"github.com/gin-gonic/gin"
)

type createCategoryRequest struct {
	Name string `json:"name"`
}

type getAllCategoriesResponse struct {
	Categories []categories.Category `json:"categories"`
}

type insertCategoryResponse struct {
	ID primitive.ObjectID `json:"id"`
}

// CategoriesHandler is a struct which implemented by categories handler
type CategoriesHandler struct {
	repo   categories.Repository
	logger *zap.SugaredLogger
}

// NewCategoriesHandler create a new instance of CategoriesHandler
func NewCategoriesHandler(repo categories.Repository, logger *zap.SugaredLogger) (*CategoriesHandler, error) {
	if repo == nil {
		return nil, errors.New("Repo must not be nil")
	}
	if logger == nil {
		return nil, errors.New("Logger must not be nil")
	}
	return &CategoriesHandler{
		repo:   repo,
		logger: logger}, nil
}

// HandleGet return all existing categories
func (handler *CategoriesHandler) HandleGet(c *gin.Context) {
	c.Header("content-type", "application/json") // todo mb move to middleware?
	c.JSON(http.StatusOK, getAllCategoriesResponse{handler.repo.GetAll()})
	return
}

// HandlePost create a new category
func (handler *CategoriesHandler) HandlePost(c *gin.Context) {
	c.Header("content-type", "application/json") // todo mb move to middleware?
	var req createCategoryRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Unable to parse body as JSON")
		return
	}
	id, err := handler.repo.Insert(c, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Unable to save category to DB")
		handler.logger.Error(err)
		return
	}
	c.JSON(http.StatusOK, insertCategoryResponse{id})
}
