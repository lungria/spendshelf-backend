package importer

import (
	"context"
	"time"

	"github.com/lungria/spendshelf-backend/mono"
	"github.com/lungria/spendshelf-backend/transaction"
	"github.com/rs/zerolog/log"
)

type BankAPI interface {
	GetTransactions(ctx context.Context, query mono.GetTransactionsQuery) ([]transaction.Transaction, error)
}

type TransactionsStorage interface {
	Save(ctx context.Context, transactions []transaction.Transaction) error
}

type ImportIntervalGenerator interface {
	GetInterval(ctx context.Context) (from time.Time, to time.Time, err error)
}

type Importer struct {
	api         BankAPI
	storage     TransactionsStorage
	intervalGen ImportIntervalGenerator
}

func NewImporeter(api BankAPI, storage TransactionsStorage, gen ImportIntervalGenerator) *Importer {
	return &Importer{
		api:         api,
		storage:     storage,
		intervalGen: gen,
	}
}

func (i *Importer) Import(accountID string) func(context.Context) {
	return func(ctx context.Context) {
		from, to, err := i.intervalGen.GetInterval(ctx)
		if err != nil {
			log.Err(err).Msg("failed import")
			return
		}
		query := mono.GetTransactionsQuery{
			Account: accountID,
			From:    from,
			To:      to,
		}
		transactions, err := i.api.GetTransactions(ctx, query)
		if err != nil {
			log.Err(err).Msg("failed import")
			return
		}

		err = i.storage.Save(ctx, transactions)
		if err != nil {
			log.Err(err).Msg("failed import")
		}
	}
}
