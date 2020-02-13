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

// SyncMonoHandler is struct for sync socket handler
type SyncMonoHandler struct {
	logger *zap.SugaredLogger
	sync   *syncmono.SyncSocket
}

// NewSyncMonoHandler creates a new SyncMonoHandler
func NewSyncMonoHandler(logger *zap.SugaredLogger, syncClient *syncmono.SyncSocket) *SyncMonoHandler {
	return &SyncMonoHandler{
		logger: logger,
		sync:   syncClient,
	}
}

// HandleSocket triggers the sync transactions from monoAPI and returns the result via websocket. ws://base_url/sync?from=1574153172
func (handler *SyncMonoHandler) HandleSocket(c *gin.Context) {
	from, exist := c.GetQuery("from")
	if !exist {
		c.JSON(http.StatusBadRequest, "'from' query parameter required: %w")
	}
	fromInt, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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

	handler.sync.MonoSync.Transactions(time.Unix(fromInt, 0))
}
