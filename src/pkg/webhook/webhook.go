package webhook

import (
	"github.com/lungria/mono"
	"net/http"

	"github.com/gin-gonic/gin"
)

type webHookRequest struct {
	Type string       `json:"type"`
	Data *Transaction `json:"data"`
}

// Transaction ...
type Transaction struct {
	AccountId     string             `json:"account" bson:"account_id"`
	StatementItem mono.StatementItem `json:"statementItem" bson:"statement_item"`
}

func (s *Server) WebHookHandler(c *gin.Context) {
	c.Header("content-type", "application/json")
	switch c.Request.Method {
	case http.MethodGet:
		s.WebHookHandlerGet(c)
	case http.MethodPost:
		s.WebHookHandlerPost(c)
	default:
		c.JSON(http.StatusMethodNotAllowed, "")
	}
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
	err = s.SaveOneTransaction(req.Data)
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
