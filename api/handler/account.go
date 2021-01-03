package handler

import (
	"context"
	"net/http"

	"github.com/lungria/spendshelf-backend/api"
	"github.com/lungria/spendshelf-backend/storage"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

// AccountsStorage abstracts storage implementation
type AccountsStorage interface {
	// GetAll accounts from storage.
	GetAll(ctx context.Context) ([]storage.Account, error)
}

// AccountHandler handles /vN/account* and /vN/accounts* routes.
type AccountHandler struct {
	storage AccountsStorage
}

// NewAccountHandler returns new instance of AccountHandler.
func NewAccountHandler(storage AccountsStorage) *AccountHandler {
	return &AccountHandler{storage: storage}
}

// GetAccounts returns accounts list.
func (t *AccountHandler) GetAccounts(c *gin.Context) {
	result, err := t.storage.GetAll(c)
	if err != nil {
		log.Error().Err(err).Msg("unable to load accounts from storage")
		c.JSON(api.InternalServerError())

		return
	}

	c.JSON(http.StatusOK, &result)
}

// BindRoutes bind gin routes to handler methods.
func (t *AccountHandler) BindRoutes(router *gin.Engine) {
	router.GET("/v1/accounts", t.GetAccounts)
}
