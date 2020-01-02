package handlers

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/lungria/spendshelf-backend/src/categories"

	"github.com/gin-gonic/gin"
)

type CreateCategoryRequest struct {
	Name string
}

type CategoriesHandler struct {
	repo categories.Repository
}

func NewCategoriesHandler(repo categories.Repository) (*CategoriesHandler, error) {
	if repo == nil {
		return nil, errors.New("Repo must not be nil")
	}
	return &CategoriesHandler{repo: repo}, nil
}

func (handler *CategoriesHandler) HandleGet(c *gin.Context) {
	c.Header("content-type", "application/json") // todo mb move to middleware?
	c.JSON(http.StatusOK, handler.repo.GetAll())
	return
}

func (handler *CategoriesHandler) HandlePost(c *gin.Context) {
	c.Header("content-type", "application/json") // todo mb move to middleware?
	var req CreateCategoryRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Unable to parse body as JSON")
		return
	}
	id, err := handler.repo.Insert(c, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Unable to save category to DB")
		return
	}
	c.JSON(http.StatusOK, id)
}
