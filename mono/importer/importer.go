package importer

import (
	"context"
	"time"

	"github.com/lungria/spendshelf-backend/mono"
	"github.com/lungria/spendshelf-backend/transaction"
	"github.com/rs/zerolog/log"
)

type BankAPI interface {
	GetTransactions(ctx context.Context, query mono.GetTransactionsQuery) ([]mono.Transaction, error)
}

type TransactionsStorage interface {
	Save(ctx context.Context, transactions []transaction.Transaction) error
}

type ImportIntervalGenerator interface {
	GetInterval(ctx context.Context, accountID string) (from, to time.Time, err error)
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

// todo test
func (i *Importer) Import(accountID string) func(context.Context) {
	return func(ctx context.Context) {
		monoTransactions, err := i.getMonoTransactions(ctx, accountID)
		if err != nil {
			log.Err(err).Msg("failed import")
			return
		}

		if len(monoTransactions) == 0 {
			return
		}

		transactions := mapTransactions(monoTransactions)

		err = i.storage.Save(ctx, transactions)
		if err != nil {
			log.Err(err).Msg("failed import")
		}
	}
}

func (i *Importer) getMonoTransactions(ctx context.Context, accountID string) ([]mono.Transaction, error) {
	from, to, err := i.intervalGen.GetInterval(ctx, accountID)
	if err != nil {
		return nil, err
	}

	query := mono.GetTransactionsQuery{
		Account: accountID,
		From:    from,
		To:      to,
	}

	return i.api.GetTransactions(ctx, query)
}

func mapTransactions(monoTransactions []mono.Transaction) []transaction.Transaction {
	transactions := make([]transaction.Transaction, len(monoTransactions))
	for i, v := range monoTransactions {
		transactions[i] = transaction.Transaction{
			BankID:      v.ID,
			Time:        time.Time(v.Time),
			Description: v.Description,
			MCC:         v.MCC,
			Hold:        v.Hold,
			Amount:      v.Amount,
		}
	}

	return transactions
}
