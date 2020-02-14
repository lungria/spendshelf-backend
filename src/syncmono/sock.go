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

type responseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type responseSuccess struct {
	Status string               `json:"status"`
	Result []models.Transaction `json:"result"`
}

// SyncSocket contains the channels for transactions and error which read from MonoSyncalso contains WebSocket connection
type SyncSocket struct {
	send     chan []models.Transaction
	SendErr  chan error
	Conn     *websocket.Conn
	MonoSync *monoSync
	logger   *zap.SugaredLogger
	context  context.Context
}

// NewSyncSocket creates a new SyncSocket
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

// Write reads from channels transactions and errors then parsed as JSON and write to websocket
func (c *SyncSocket) Write() {
	for {
		select {
		case txn := <-c.send:
			resp := responseSuccess{
				Status: "Success",
				Result: txn,
			}
			if err := c.Conn.WriteJSON(resp); err != nil {
				errMsg := "unable to Write transactions"
				c.logger.Errorw(errMsg, "Error", err.Error())
				c.SendErr <- err
			}
		case sockErr := <-c.SendErr:
			resp := responseError{
				Status:  "Error",
				Message: sockErr.Error(),
			}
			if err := c.Conn.WriteJSON(resp); err != nil {
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

			err := c.MonoSync.txnRepo.InsertMany(context.TODO(), txns)
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
