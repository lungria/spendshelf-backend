package transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/lungria/spendshelf-backend/account"
	"github.com/lungria/spendshelf-backend/importer/mono"
	"github.com/lungria/spendshelf-backend/transaction/category"
)

// BankAPI abstracts bank API related to transactions information.
type BankAPI interface {
	// GetTransactions allows to load list of transactions based on specified query parameters.
	GetTransactions(ctx context.Context, query mono.GetTransactionsQuery) ([]mono.Transaction, error)
}

// Storage abstracts persistent storage for transactions.
type Storage interface {
	// Save transactions to database.
	Save(ctx context.Context, transactions []Transaction) error
}

// ImportIntervalGenerator generates time interval for transaction import.
type ImportIntervalGenerator interface {
	// GetInterval generates time interval for transaction import. Only transactions in this time interval
	// will be loaded.
	GetInterval(ctx context.Context, accountID string) (from, to time.Time, err error)
}

// Importer loads transactions in specified interval for specified accountID and saves
// them to storage.
type Importer struct {
	api          BankAPI
	transactions Storage
	accounts     account.Importer
	intervalGen  ImportIntervalGenerator
}

// NewImporter create new instance of Importer.
func NewImporter(
	api BankAPI,
	storage Storage,
	gen ImportIntervalGenerator) *Importer {
	return &Importer{
		api:          api,
		transactions: storage,
		intervalGen:  gen,
	}
}

// Import loads transactions in specified interval for specified accountID and saves them to storage.
func (i *Importer) Import(ctx context.Context, accountID string) error {
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

func (i *Importer) mapTransactions(accountID string, src []mono.Transaction) []Transaction {
	dst := make([]Transaction, len(src))
	for i, v := range src {
		dst[i] = Transaction{
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
