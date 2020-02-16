package db

import (
	"context"

	"go.uber.org/zap"

	"go.etcd.io/bbolt"
)

type Config interface {
	GetDBName() string
}

// bucket names
var (
	UncategorizedTransactionsBucket = "uncategorized"
	TransactionsBucket              = "transactions"
	CategoriesBucket                = "categories"
)

type Connection struct {
	*bbolt.DB
	logger *zap.SugaredLogger
}

func OpenConnection(cfg Config, logger *zap.SugaredLogger) (*Connection, error) {
	db, err := bbolt.Open(cfg.GetDBName(), 0666, nil)
	if err != nil {
		return nil, err
	}
	connection := &Connection{db, logger}

	err = connection.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(TransactionsBucket))
		if err != nil {
			connection.logger.Error("create bucket: %s", zap.Error(err))
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(UncategorizedTransactionsBucket))
		if err != nil {
			connection.logger.Error("create bucket: %s", zap.Error(err))
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(CategoriesBucket))
		if err != nil {
			connection.logger.Error("create bucket: %s", zap.Error(err))
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return connection, err
}

func (db *Connection) CloseWithTimeout(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		err := db.Close()
		if err != nil {
			db.logger.Error("unable to close db connection", zap.Error(err))
		}
		done <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		db.logger.Error("unable to gracefully close db connection")
		return
	case <-done:
		db.logger.Info("db connection closed successfully")
		return
	}
}
