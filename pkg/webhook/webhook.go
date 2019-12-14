package api

import (
	"github.com/lungria/spendshelf-backend/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type webHookRequest struct {
	Type string          `json:"type"`
	Data *db.Transaction `json:"data"`
}

// WebHookHandlerPost catch the request from monoAPI and save to DB
func WebHookHandlerPost(c *gin.Context) {
	var err error
	var req *webHookRequest
	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "Bad request", Error: err.Error()})
	}
	database, err := db.NewConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "Connection to database failed", Error: err.Error()})
	}
	err = database.SaveOneTransaction(req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "Saving transaction failed", Error: err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// WebHookHandlerGet respond 200 to monoAPI when WebHook was set
func WebHookHandlerGet(c *gin.Context) {
	c.String(http.StatusOK, "")
}
