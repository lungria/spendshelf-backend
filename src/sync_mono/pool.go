package sync_mono

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lungria/spendshelf-backend/src/models"
)

type pool struct {
	clients  map[*client]bool
	leave    chan *client
	join     chan *client
	monoSync *MonoSync
}

func NewPool(monoSync *MonoSync) *pool {
	clients := make(map[*client]bool)
	leave := make(chan *client)
	join := make(chan *client)

	p := pool{
		clients:  clients,
		leave:    leave,
		join:     join,
		monoSync: monoSync,
	}
	go p.run()
	return &p
}

func (p pool) run() {
	for {
		select {
		case client := <-p.leave:
			delete(p.clients, client)
		case client := <-p.join:
			p.clients[client] = true
		case txns := <-p.monoSync.transactions:
			toInsert := p.monoSync.trimDuplicate(txns)
			if len(toInsert) == 0 {
				continue
			}

			for client := range p.clients {
				client.send <- toInsert
			}

			p.monoSync.Lock()
			err := p.monoSync.txnRepo.InsertManyTransactions(toInsert)
			p.monoSync.errChan <- err
			p.monoSync.Unlock()
		}
	}
}

func (p *pool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{WriteBufferSize: 1024, ReadBufferSize: 1024}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Socket serving failed", err)
		return
	}
	client := &client{
		send:   make(chan []models.Transaction),
		socket: socket,
	}

	p.join <- client
	defer func() { p.leave <- client }()
	go client.write()

	p.monoSync.Transactions(time.Unix(1574158956, 0))
}
