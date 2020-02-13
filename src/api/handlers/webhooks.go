package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/lungria/spendshelf-backend/src/models"

	"github.com/lungria/spendshelf-backend/src/webhooks"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type webHookRequest struct {
	Type string          `json:"type"`
	Data *models.WebHook `json:"data"`
}

// WebHookHandler is a struct which implemented by webhooks handlers
type WebHookHandler struct {
	repo   webhooks.Repository
	Logger *zap.SugaredLogger
}

// NewWebHookHandler create a new instance of WebHookHandler
func NewWebHookHandler(repo webhooks.Repository, logger *zap.SugaredLogger) (*WebHookHandler, error) {
	if repo == nil {
		return nil, errors.New("repo must not be nil")
	}
	if logger == nil {
		return nil, errors.New("logger must not be nil")
	}
	return &WebHookHandler{
		repo:   repo,
		Logger: logger,
	}, nil
}

// HandlePost catch the request from monoAPI and save to DB
func (handler *WebHookHandler) HandlePost(c *gin.Context) {
	var err error
	var req *webHookRequest

	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = handler.repo.InsertOneHook(req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("unable save transaction: %w", err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// HandleGet respond 200 to monoAPI when WebHook was set
func (handler *WebHookHandler) HandleGet(c *gin.Context) {
	c.String(http.StatusOK, "")
}
