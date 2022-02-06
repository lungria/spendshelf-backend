package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lungria/spendshelf-backend/api"
	"github.com/rs/zerolog/log"
)

// Handler handles /vN/account* and /vN/accounts* routes.
type Handler struct {
	storage *Repository
}

// NewHandler returns new instance of Handler.
func NewHandler(storage *Repository) *Handler {
	return &Handler{storage: storage}
}

// GetAccounts returns accounts list.
func (t *Handler) GetAccounts(c *gin.Context) {
	result, err := t.storage.GetAll(c)
	if err != nil {
		log.Error().Err(err).Msg("unable to load accounts from storage")
		c.JSON(api.InternalServerError())

		return
	}

	c.JSON(http.StatusOK, &result)
}

// BindRoutes bind gin routes to handler methods.
func (t *Handler) BindRoutes(router *gin.Engine) {
	router.GET("/v1/accounts", t.GetAccounts)
}
