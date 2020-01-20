package syncmono

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
//	monoSync *monoSync
//}
//
//func NewPool(monoSync *monoSync) *pool {
//	clients := make(map[*SyncSocket]bool)
//	leave := make(chan *SyncSocket)
//	join := make(chan *SyncSocket)
//
//	p := pool{
//		clients:  clients,
//		leave:    leave,
//		join:     join,
//		monoSync: monoSync,
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
//		case txns := <-p.monoSync.transactions:
//			toInsert := p.monoSync.trimDuplicate(txns)
//			if len(toInsert) == 0 {
//				continue
//			}
//
//			for SyncSocket := range p.clients {
//				SyncSocket.send <- toInsert
//			}
//
//			p.monoSync.Lock()
//			err := p.monoSync.txnRepo.InsertManyTransactions(toInsert)
//			p.monoSync.errChan <- err
//			p.monoSync.Unlock()
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
//	Conn, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Println("Conn serving failed", err)
//		return
//	}
//	SyncSocket := &SyncSocket{
//		send:   make(chan []models.Transaction),
//		Conn: Conn,
//	}
//
//	p.join <- SyncSocket
//	defer func() { p.leave <- SyncSocket }()
//	go SyncSocket.Write()
//
//	p.monoSync.Transactions(time.Unix(1574158956, 0))
//}
