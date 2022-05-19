package budget

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lungria/spendshelf-backend/api"
	"github.com/rs/zerolog/log"
)

// Handler handles /vN/budget* and /vN/budgets* routes.
type Handler struct {
	storage *Repository
}

// NewHandler returns new instance of Handler.
func NewHandler(storage *Repository) *Handler {
	return &Handler{storage: storage}
}

// GetCurrentBudget returns current active budget.
func (h *Handler) GetCurrentBudget(c *gin.Context) {
	result, err := h.storage.GetLast(c)
	if err != nil {
		log.Error().Err(err).Msg("unable to load current budget from storage")
		c.JSON(api.InternalServerError())

		return
	}

	c.JSON(http.StatusOK, &result)
}

// CreateNewBudget ends current active budget (if present) and creates new budget.
func (h *Handler) CreateNewBudget(c *gin.Context) {
	var req CreateBudgetRequest

	if err := c.BindJSON(&req); err != nil {
		log.Error().Err(err).Msg("unable to create new budget from api")
		c.JSON(api.InternalServerError())

		return
	}

	b, err := h.storage.Create(c, req)
	if err != nil {
		log.Error().Err(err).Msg("unable to create new budget from api")
		c.JSON(api.InternalServerError())

		return
	}

	c.JSON(http.StatusCreated, &b)
}

// BindRoutes bind gin routes to handler methods.
func (h *Handler) BindRoutes(router *gin.Engine) {
	router.GET("/v1/budget", h.GetCurrentBudget)
	router.POST("/v1/budget", h.GetCurrentBudget)
}
