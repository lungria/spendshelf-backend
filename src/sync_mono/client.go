package sync_mono

import (
	"errors"

	"github.com/lungria/spendshelf-backend/src/config"

	"github.com/lungria/spendshelf-backend/src/transactions"

	"go.uber.org/zap"

	"github.com/gorilla/websocket"
	"github.com/lungria/spendshelf-backend/src/models"
)

type socketError struct {
	Error string `json:"error"`
}

type SyncSocket struct {
	send     chan []models.Transaction
	SendErr  chan error
	Socket   *websocket.Conn
	MonoSync *monoSync
	logger   *zap.SugaredLogger
}

func NewClient(logger *zap.SugaredLogger, cfg *config.EnvironmentConfiguration, txnRepo transactions.Repository) (*SyncSocket, error) {
	m, err := newMonoSync(cfg, logger, txnRepo)
	if err != nil {
		return nil, err
	}

	client := SyncSocket{
		send:     make(chan []models.Transaction),
		SendErr:  make(chan error),
		logger:   logger,
		MonoSync: m,
	}
	go client.run()

	return &client, nil
}

func (c *SyncSocket) Write() {
	for {
		select {
		case txn := <-c.send:
			if err := c.Socket.WriteJSON(txn); err != nil {
				errMsg := "unable to Write transactions"
				c.logger.Errorw(errMsg, "Error", err.Error())
				c.SendErr <- err
			}
		case errc := <-c.SendErr:
			if err := c.Socket.WriteJSON(socketError{Error: errc.Error()}); err != nil {
				c.logger.Errorw("unable to Write error", "Error", err.Error())
			}
		}
	}
}

func (c SyncSocket) run() {
	for {
		select {
		case txns := <-c.MonoSync.transactions:
			toInsert := c.MonoSync.trimDuplicate(txns)
			if len(toInsert) == 0 {
				c.logger.Info("no transactions to save for this period")
				c.SendErr <- errors.New("no transactions to save for this period")
				continue
			}

			c.send <- toInsert

			err := c.MonoSync.txnRepo.InsertManyTransactions(toInsert)
			if err != nil {
				c.MonoSync.errChan <- err
			}
			c.logger.Info("Transactions were saved.")
		case err := <-c.MonoSync.errChan:
			c.logger.Error(err.Error())
			c.SendErr <- err
		}
	}
}
