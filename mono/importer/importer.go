package importer

import (
	"context"
	"time"

	"github.com/lungria/spendshelf-backend/mono"
	"github.com/lungria/spendshelf-backend/storage/transaction"
	"github.com/rs/zerolog/log"
)

// BankAPI abstracts bank API.
type BankAPI interface {
	// GetTransactions allows to load list of transactions based on specified query parameters.
	GetTransactions(ctx context.Context, query mono.GetTransactionsQuery) ([]mono.Transaction, error)
}

// TransactionsStorage abstracts persistent storage for transactions.
type TransactionsStorage interface {
	// Save transactions to database.
	Save(ctx context.Context, transactions []transaction.Transaction) error
}

// ImportIntervalGenerator generates interval for transaction import.
type ImportIntervalGenerator interface {
	// GetInterval generates interval for transaction import.
	GetInterval(ctx context.Context, accountID string) (from, to time.Time, err error)
}

// Importer handles one-time import of transactions list for selected interval.
type Importer struct {
	api         BankAPI
	storage     TransactionsStorage
	intervalGen ImportIntervalGenerator
}

// NewImporter create new instance of Importer.
func NewImporter(api BankAPI, storage TransactionsStorage, gen ImportIntervalGenerator) *Importer {
	return &Importer{
		api:         api,
		storage:     storage,
		intervalGen: gen,
	}
}

// Import loads latest transactions from mono API and stores them to DB.
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

		transactions := mapTransactions(accountID, monoTransactions)

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

func mapTransactions(accountID string, monoTransactions []mono.Transaction) []transaction.Transaction {
	transactions := make([]transaction.Transaction, len(monoTransactions))
	for i, v := range monoTransactions {
		transactions[i] = transaction.Transaction{
			ID:          v.ID,
			Time:        time.Time(v.Time),
			Description: v.Description,
			MCC:         v.MCC,
			Hold:        v.Hold,
			Amount:      v.Amount,
			AccountID:   accountID,
			CategoryID:  transaction.DefaultCategoryID,
		}
	}

	return transactions
}
