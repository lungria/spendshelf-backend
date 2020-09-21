package importer

import (
	"context"

	"github.com/lungria/spendshelf-backend/transaction"
	"github.com/sirupsen/logrus"
)

type BankAPI interface {
	GetTransactions(ctx context.Context) ([]transaction.Transaction, error)
}

type TransactionsStorage interface {
	Save(ctx context.Context, transactions []transaction.Transaction) error
}

type Importer struct {
	api     BankAPI
	log     *logrus.Logger
	storage TransactionsStorage
}

func NewImporeter(api BankAPI, storage TransactionsStorage, log *logrus.Logger) *Importer {
	return &Importer{
		api:     api,
		storage: storage,
		log:     log,
	}
}

func (i *Importer) Import(ctx context.Context) {
	transactions, err := i.api.GetTransactions(ctx)
	if err != nil {
		i.log.WithError(err).Error("failed import")
		return
	}

	err = i.storage.Save(ctx, transactions)
	if err != nil {
		i.log.WithError(err).Error("failed import") // todo zerolog?
	}
}
