package sync_mono

//
//import (
//	"log"
//	"net/http"
//	"time"
//
//	"github.com/gorilla/websocket"
//	"github.com/lungria/spendshelf-backend/src/models"
//)
//
//type pool struct {
//	clients  map[*SyncSocket]bool
//	leave    chan *SyncSocket
//	join     chan *SyncSocket
//	MonoSync *MonoSync
//}
//
//func NewPool(MonoSync *MonoSync) *pool {
//	clients := make(map[*SyncSocket]bool)
//	leave := make(chan *SyncSocket)
//	join := make(chan *SyncSocket)
//
//	p := pool{
//		clients:  clients,
//		leave:    leave,
//		join:     join,
//		MonoSync: MonoSync,
//	}
//	go p.run()
//	return &p
//}
//
//func (p pool) run() {
//	for {
//		select {
//		case SyncSocket := <-p.leave:
//			delete(p.clients, SyncSocket)
//		case SyncSocket := <-p.join:
//			p.clients[SyncSocket] = true
//		case txns := <-p.MonoSync.transactions:
//			toInsert := p.MonoSync.trimDuplicate(txns)
//			if len(toInsert) == 0 {
//				continue
//			}
//
//			for SyncSocket := range p.clients {
//				SyncSocket.send <- toInsert
//			}
//
//			p.MonoSync.Lock()
//			err := p.MonoSync.txnRepo.InsertManyTransactions(toInsert)
//			p.MonoSync.errChan <- err
//			p.MonoSync.Unlock()
//		}
//	}
//}
//
//func (p *pool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	upgrader := websocket.Upgrader{WriteBufferSize: 1024, ReadBufferSize: 1024}
//	upgrader.CheckOrigin = func(r *http.Request) bool {
//		return true
//	}
//
//	Socket, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Println("Socket serving failed", err)
//		return
//	}
//	SyncSocket := &SyncSocket{
//		send:   make(chan []models.Transaction),
//		Socket: Socket,
//	}
//
//	p.join <- SyncSocket
//	defer func() { p.leave <- SyncSocket }()
//	go SyncSocket.Write()
//
//	p.MonoSync.Transactions(time.Unix(1574158956, 0))
//}
