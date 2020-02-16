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
	cfg Config
	*bbolt.DB
	logger *zap.SugaredLogger
}

func NewConnection(cfg Config, logger *zap.SugaredLogger) *Connection {
	return &Connection{cfg: cfg, logger: logger}
}

// KeepConnected keeps connection to db until ctx is cancelled. This method is blocking.
func (db *Connection) KeepConnected(ctx context.Context) error {
	bolt, err := bbolt.Open(db.cfg.GetDBName(), 0666, nil)
	if err != nil {
		db.logger.Fatal("unable connect to db", zap.Error(err))
	}
	db.DB = bolt
	err = db.ensureBucketsCreated()
	if err != nil {
		db.logger.Fatal("unable connect to db", zap.Error(err))
	}
	<-ctx.Done()
	return db.Close()
}

func (db *Connection) ensureBucketsCreated() error {
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(TransactionsBucket))
		if err != nil {
			db.logger.Error("create bucket: %s", zap.Error(err))
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(UncategorizedTransactionsBucket))
		if err != nil {
			db.logger.Error("create bucket: %s", zap.Error(err))
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(CategoriesBucket))
		if err != nil {
			db.logger.Error("create bucket: %s", zap.Error(err))
			return err
		}

		return nil
	})
}
