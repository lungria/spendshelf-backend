// todo: rename package to transaction
package transactions

import (
	"context"
	"fmt"
	"time"

	"github.com/lungria/spendshelf-backend/importer/account"

	"github.com/lungria/spendshelf-backend/importer/mono"
	"github.com/lungria/spendshelf-backend/storage"
	"github.com/lungria/spendshelf-backend/storage/category"
)

// BankAPI abstracts bank API related to transactions information.
type BankAPI interface {
	// GetTransactions allows to load list of transactions based on specified query parameters.
	GetTransactions(ctx context.Context, query mono.GetTransactionsQuery) ([]mono.Transaction, error)
}

// Storage abstracts persistent storage for transactions.
type Storage interface {
	// Save transactions to database.
	Save(ctx context.Context, transactions []storage.Transaction) error
}

// ImportIntervalGenerator generates time interval for transaction import.
type ImportIntervalGenerator interface {
	// GetInterval generates time interval for transaction import. Only transactions in this time interval
	// will be loaded.
	GetInterval(ctx context.Context, accountID string) (from, to time.Time, err error)
}

// DefaultImporter loads transactions in specified interval for specified accountID and saves
// them to storage.
type DefaultImporter struct {
	api          BankAPI
	transactions Storage
	accounts     account.DefaultImporter
	intervalGen  ImportIntervalGenerator
}

// NewImporter create new instance of DefaultImporter.
func NewImporter(
	api BankAPI,
	storage Storage,
	gen ImportIntervalGenerator) *DefaultImporter {
	return &DefaultImporter{
		api:          api,
		transactions: storage,
		intervalGen:  gen,
	}
}

// Import loads transactions in specified interval for specified accountID and saves them to storage.
// todo: tests.
func (i *DefaultImporter) Import(ctx context.Context, accountID string) error {
	from, to, err := i.intervalGen.GetInterval(ctx, accountID)
	if err != nil {
		return fmt.Errorf("failed import transaction for account '%s': %w", accountID, err)
	}

	query := mono.GetTransactionsQuery{
		Account: accountID,
		From:    from,
		To:      to,
	}

	monoTransactions, err := i.api.GetTransactions(ctx, query)
	if err != nil {
		return fmt.Errorf("failed import transaction for account '%s': %w", accountID, err)
	}

	if len(monoTransactions) == 0 {
		return nil
	}

	transactions := i.mapTransactions(accountID, monoTransactions)

	err = i.transactions.Save(ctx, transactions)
	if err != nil {
		return fmt.Errorf("failed import transaction for account '%s': %w", accountID, err)
	}

	return nil
}

func (i *DefaultImporter) mapTransactions(accountID string, src []mono.Transaction) []storage.Transaction {
	dst := make([]storage.Transaction, len(src))
	for i, v := range src {
		dst[i] = storage.Transaction{
			ID:          v.ID,
			Time:        time.Time(v.Time),
			Description: v.Description,
			MCC:         v.MCC,
			Hold:        v.Hold,
			Amount:      v.Amount,
			AccountID:   accountID,
			CategoryID:  category.Default,
		}
	}

	return dst
}
