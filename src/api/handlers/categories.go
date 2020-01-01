package handlers

import (
	"github.com/pkg/errors"

	"github.com/lungria/spendshelf-backend/src/categories"

	"github.com/gin-gonic/gin"
)

type CategoriesHandler struct {
	repo categories.Repository
}

func NewCategoriesHandler(repo categories.Repository) (*CategoriesHandler, error) {
	if repo == nil {
		return nil, errors.New("Repo must not be nil")
	}
	return &CategoriesHandler{repo: repo}, nil
}

func (handler *CategoriesHandler) Handle(c *gin.Context) {
	c.Header("content-type", "application/json") // todo mb move to middleware?

}
