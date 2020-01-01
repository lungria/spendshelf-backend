package handlers

import (
	"errors"
	"net/http"

	"github.com/lungria/spendshelf-backend/src/db"
	"go.uber.org/zap"

	"github.com/lungria/spendshelf-backend/src/models"

	"github.com/gin-gonic/gin"
)

type webHookRequest struct {
	Type string              `json:"type"`
	Data *models.Transaction `json:"data"`
}

type WebHookHandler struct {
	repo   db.TransactionsRepository
	Logger *zap.SugaredLogger
}

func NewWebHookHandler(repo db.TransactionsRepository, logger *zap.SugaredLogger) (*WebHookHandler, error) {
	if repo == nil {
		return nil, errors.New("Repo must not be nil")
	}
	if logger == nil {
		return nil, errors.New("Logger must not be nil")
	}
	return &WebHookHandler{
		repo:   repo,
		Logger: logger,
	}, nil
}

// WebHookHandler is routing for different HTTP methods
func (handler *WebHookHandler) Handle(c *gin.Context) {
	c.Header("content-type", "application/json")
	switch c.Request.Method {
	case http.MethodGet:
		handler.webHookHandlerGet(c)
	case http.MethodPost:
		handler.webHookHandlerPost(c)
	default:
		c.JSON(http.StatusMethodNotAllowed, "")
	}
}

// webHookHandlerPost catch the request from monoAPI and save to DB
func (handler *WebHookHandler) webHookHandlerPost(c *gin.Context) {
	var err error
	var req *webHookRequest

	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "Bad request", Error: err.Error()})
		return
	}
	err = handler.repo.SaveOneTransaction(req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "Saving Transaction failed", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// webHookHandlerGet respond 200 to monoAPI when WebHook was set
func (handler *WebHookHandler) webHookHandlerGet(c *gin.Context) {
	c.String(http.StatusOK, "")
}
