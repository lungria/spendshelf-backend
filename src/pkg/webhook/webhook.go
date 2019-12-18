package webhook

import (
	"github.com/lungria/mono"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebHookAPI struct {
	HTTPServer *http.Server
}

type webHookRequest struct {
	Type string       `json:"type"`
	Data *Transaction `json:"data"`
}

// Transaction ...
type Transaction struct {
	AccountId     string             `json:"account" bson:"account_id"`
	StatementItem mono.StatementItem `json:"statementItem" bson:"statement_item"`
}

// WebHookHandlerPost catch the request from monoAPI and save to DB
func (s *Server) WebHookHandlerPost(c *gin.Context) {
	var err error
	var req *webHookRequest

	err = c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: "Bad request", Error: err.Error()})
		return
	}
	err = s.MongoDB.SaveOneTransaction(req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: "Saving Transaction failed", Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

// WebHookHandlerGet respond 200 to monoAPI when WebHook was set
func (s *Server) WebHookHandlerGet(c *gin.Context) {
	c.String(http.StatusOK, "")
}
