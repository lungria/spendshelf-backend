package api

import (
	"net/http"

	"github.com/lungria/spendshelf-backend/src/models"

	"github.com/gin-gonic/gin"
)

type webHookRequest struct {
	Type string              `json:"type"`
	Data *models.Transaction `json:"data"`
}

// WebHookHandler is routing for different HTTP methods
func (a *WebHookAPI) WebHookHandler(c *gin.Context) {
	c.Header("content-type", "application/json")
	switch c.Request.Method {
	case http.MethodGet:
		a.WebHookHandlerGet(c)
	case http.MethodPost:
		a.WebHookHandlerPost(c)
	default:
		c.JSON(http.StatusMethodNotAllowed, "")
	}
}

// WebHookHandlerPost catch the request from monoAPI and save to DB
func (a *WebHookAPI) WebHookHandlerPost(c *gin.Context) {
	var err error
	var req *webHookRequest

	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "Bad request", Error: err.Error()})
		return
	}
	err = a.Database.SaveOneTransaction(req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "Saving Transaction failed", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// WebHookHandlerGet respond 200 to monoAPI when WebHook was set
func (a *WebHookAPI) WebHookHandlerGet(c *gin.Context) {
	c.String(http.StatusOK, "")
}
