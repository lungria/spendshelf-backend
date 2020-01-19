package sync_mono

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lungria/spendshelf-backend/src/models"
)

type client struct {
	send     chan []models.Transaction
	socket   *websocket.Conn
	monoSync *MonoSync
}

func NewClient(m *MonoSync) *client {
	client := client{
		send:     make(chan []models.Transaction),
		monoSync: m,
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

func (c client) run() {
	for {
		select {
		case txns := <-c.monoSync.transactions:
			toInsert := c.monoSync.trimDuplicate(txns)
			if len(toInsert) == 0 {
				continue
			}

			c.send <- toInsert

			err := c.monoSync.txnRepo.InsertManyTransactions(toInsert)
			c.monoSync.errChan <- err
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

	c.monoSync.Transactions(time.Unix(1574158956, 0))
}
