package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lungria/spendshelf-backend/api"
	"github.com/lungria/spendshelf-backend/storage"
	"github.com/rs/zerolog/log"
)

// BudgetsStorage abstracts storage implementation.
type BudgetsStorage interface {
	// GetLast returns last budget from storage.
	GetLast(ctx context.Context) (storage.Budget, error)
}

// BudgetHandler handles /vN/budget* and /vN/budgets* routes.
type BudgetHandler struct {
	storage BudgetsStorage
}

// NewBudgetHandler returns new instance of BudgetHandler.
func NewBudgetHandler(storage BudgetsStorage) *BudgetHandler {
	return &BudgetHandler{storage: storage}
}

// GetCurrentBudget returns current active budget.
func (h *BudgetHandler) GetCurrentBudget(c *gin.Context) {
	result, err := h.storage.GetLast(c)
	if err != nil {
		log.Error().Err(err).Msg("unable to load current budget from storage")
		c.JSON(api.InternalServerError())
		return
	}

	c.JSON(http.StatusOK, &result)
}

// BindRoutes bind gin routes to handler methods.
func (t *BudgetHandler) BindRoutes(router *gin.Engine) {
	router.GET("/v1/budget?month=current", t.GetCurrentBudget)
}
