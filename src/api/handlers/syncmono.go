package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/lungria/spendshelf-backend/src/syncmono"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type SyncMonoHandler struct {
	logger *zap.SugaredLogger
	sync   *syncmono.SyncSocket
}

func NewSyncMonoHandler(logger *zap.SugaredLogger, syncClient *syncmono.SyncSocket) *SyncMonoHandler {
	return &SyncMonoHandler{
		logger: logger,
		sync:   syncClient,
	}
}

func (handler *SyncMonoHandler) HandleSocket(c *gin.Context) {
	from, exist := c.GetQuery("from")
	if !exist {
		c.JSON(http.StatusBadRequest, messageResponse{Message: "Query from is required"})
	}
	fromInt, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, messageResponse{Message: err.Error()})
	}
	upgrader := websocket.Upgrader{WriteBufferSize: 1024, ReadBufferSize: 1024}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		errMsg := "Connection serving failed"
		handler.logger.Errorw(errMsg, "Error", err.Error())
		handler.sync.SendErr <- err
		return
	}

	handler.sync.Conn = socket

	go handler.sync.Write()

	handler.sync.Transactions(time.Unix(fromInt, 0))
}
