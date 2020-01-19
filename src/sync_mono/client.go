package sync_mono

import (
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
	"github.com/lungria/spendshelf-backend/src/models"
)

type socketError struct {
	Error string `json:"error"`
}

type client struct {
	send     chan []models.Transaction
	sendErr  chan error
	socket   *websocket.Conn
	monoSync *MonoSync
	logger   *zap.SugaredLogger
}

func NewClient(m *MonoSync, l *zap.SugaredLogger) *client {
	client := client{
		send:     make(chan []models.Transaction),
		sendErr:  make(chan error),
		monoSync: m,
		logger:   l,
	}
	go client.run()

	return &client
}

func (c *client) write() {
	for {
		select {
		case txn := <-c.send:
			if err := c.socket.WriteJSON(txn); err != nil {
				errMsg := "unable to write transactions"
				c.logger.Errorw(errMsg, "Error", err.Error())
				c.sendErr <- err
			}
		case errc := <-c.sendErr:
			if err := c.socket.WriteJSON(socketError{Error: errc.Error()}); err != nil {
				c.logger.Errorw("unable to write error", "Error", err.Error())
			}
		}
	}
}

func (c client) run() {
	for {
		select {
		case txns := <-c.monoSync.transactions:
			toInsert := c.monoSync.trimDuplicate(txns)
			if len(toInsert) == 0 {
				c.logger.Info("no transactions to save for this period")
				c.sendErr <- errors.New("no transactions to save for this period")
				continue
			}

			c.send <- toInsert

			err := c.monoSync.txnRepo.InsertManyTransactions(toInsert)
			if err != nil {
				c.monoSync.errChan <- err
			}
			c.logger.Info("Transactions were saved.")
		case err := <-c.monoSync.errChan:
			c.logger.Error(err.Error())
			c.sendErr <- err
		}
	}
}

func (c *client) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{WriteBufferSize: 1024, ReadBufferSize: 1024}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errMsg := "Socket serving failed"
		c.logger.Errorw(errMsg, "Error", err.Error())
		c.sendErr <- err
		return
	}
	c.socket = socket

	go c.write()

	c.monoSync.Transactions(time.Unix(1574158956, 0))
}
