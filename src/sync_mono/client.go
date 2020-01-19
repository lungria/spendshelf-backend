package sync_mono

import (
	"errors"
	"log"
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
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			log.Fatalln(err)
		}
	}
}

func (c *client) writeErr() {
	defer c.socket.Close()
	for e := range c.sendErr {
		if err := c.socket.WriteJSON(socketError{Error: e.Error()}); err != nil {
			log.Println(err)
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
		log.Println("Socket serving failed", err)
		return
	}
	c.socket = socket

	go c.write()
	go c.writeErr()

	c.monoSync.Transactions(time.Unix(1574158956, 0))
}
