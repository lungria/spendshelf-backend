package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
	"github.com/lungria/spendshelf-backend/src/sync_mono"
	"go.uber.org/zap"
)

type SyncMonoHandler struct {
	logger *zap.SugaredLogger
	socket *sync_mono.SyncSocket
}

func NewSyncMonoHandler(logger *zap.SugaredLogger, syncClient *sync_mono.SyncSocket) *SyncMonoHandler {
	return &SyncMonoHandler{
		logger: logger,
		socket: syncClient,
	}
}

func (handler *SyncMonoHandler) HandleSocket(c *gin.Context) {
	from, exist := c.GetQuery("from")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
	fromInt, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
	upgrader := websocket.Upgrader{WriteBufferSize: 1024, ReadBufferSize: 1024}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		errMsg := "Socket serving failed"
		handler.logger.Errorw(errMsg, "Error", err.Error())
		handler.socket.SendErr <- err
		return
	}

	handler.socket.Socket = socket

	go handler.socket.Write()

	handler.socket.MonoSync.Transactions(time.Unix(fromInt, 0))
}
