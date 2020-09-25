package importer

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/lungria/spendshelf-backend/transaction"
 )

type BankAPI interface {
	GetTransactions(ctx context.Context) ([]transaction.Transaction, error)
}

type TransactionsStorage interface {
	Save(ctx context.Context, transactions []transaction.Transaction) error
}

type Importer struct {
	api     BankAPI
 	storage TransactionsStorage
}

func NewImporeter(api BankAPI, storage TransactionsStorage) *Importer {
	return &Importer{
		api:     api,
		storage: storage,
 	}
}

func (i *Importer) Import(ctx context.Context) {
	transactions, err := i.api.GetTransactions(ctx)
	if err != nil {
		log.Err(err).Msg("failed import")
		return
	}

	err = i.storage.Save(ctx, transactions)
	if err != nil {
		log.Err(err).Msg("failed import")
	}
}
