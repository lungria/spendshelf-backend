package syncmono

import (
	"context"
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
	Conn     *websocket.Conn
	MonoSync *monoSync
	logger   *zap.SugaredLogger
	context  context.Context
}

func NewSyncSocket(ctx context.Context, logger *zap.SugaredLogger, cfg *config.EnvironmentConfiguration, txnRepo transactions.Repository) (*SyncSocket, error) {
	m, err := newMonoSync(cfg, logger, txnRepo)
	if err != nil {
		return nil, err
	}

	syncSocket := SyncSocket{
		send:     make(chan []models.Transaction),
		SendErr:  make(chan error),
		logger:   logger,
		MonoSync: m,
		context:  ctx,
	}
	go syncSocket.run()

	return &syncSocket, nil
}

func (c *SyncSocket) Write() {
	for {
		select {
		case txn := <-c.send:
			if err := c.Conn.WriteJSON(txn); err != nil {
				errMsg := "unable to Write transactions"
				c.logger.Errorw(errMsg, "Error", err.Error())
				c.SendErr <- err
			}
		case sockErr := <-c.SendErr:
			if err := c.Conn.WriteJSON(socketError{Error: sockErr.Error()}); err != nil {
				c.logger.Errorw("unable to Write error", "Error", err.Error())
			}
		}
	}
}

func (c SyncSocket) run() {
	for {
		select {
		case txns := <-c.MonoSync.transactions:
			if len(txns) == 0 {
				c.logger.Info("no transactions to save for this period")
				c.SendErr <- errors.New("no transactions to save for this period")
				continue
			}

			c.send <- txns

			err := c.MonoSync.txnRepo.InsertManyTransactions(txns)
			if err != nil {
				c.SendErr <- err
			}
			c.logger.Info("Transactions were saved.")
		case err := <-c.MonoSync.errChan:
			c.logger.Error(err.Error())
			c.SendErr <- err
		case <-c.context.Done():
			return
		}
	}
}
