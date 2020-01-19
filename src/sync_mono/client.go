package sync_mono

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/lungria/spendshelf-backend/src/models"
)

type client struct {
	send   chan []models.Transaction
	socket *websocket.Conn
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			log.Fatalln(err)
		}
	}
}
