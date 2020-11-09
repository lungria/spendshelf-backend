package importer

import (
	"context"
	"time"

	"github.com/lungria/spendshelf-backend/category"
	"github.com/lungria/spendshelf-backend/mono"
	"github.com/lungria/spendshelf-backend/storage"
	"github.com/rs/zerolog/log"
)

// BankAPI abstracts bank API.
type BankAPI interface {
	// GetTransactions allows to load list of transactions based on specified query parameters.
	GetTransactions(ctx context.Context, query mono.GetTransactionsQuery) ([]mono.Transaction, error)
	// GetUserInfo loads user accounts list.
	GetUserInfo(ctx context.Context) ([]mono.Account, error)
}

// TransactionsStorage abstracts persistent storage for transactions.
type TransactionsStorage interface {
	// Save transactions to database.
	Save(ctx context.Context, transactions []storage.Transaction) error
}

// AccountsStorage abstracts persistent storage for accounts.
type AccountsStorage interface {
	// todo: update by id, if not found - insert new record
	// todo: do not use mono.Account, introduce custom model
	Save(ctx context.Context, account mono.Account) error
}

// ImportIntervalGenerator generates interval for transaction import.
type ImportIntervalGenerator interface {
	// GetInterval generates interval for transaction import.
	GetInterval(ctx context.Context, accountID string) (from, to time.Time, err error)
}

// Importer loads latest data from bank for specified accountID.
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

// Import latest data from bank for specified accountID.
func (i *Importer) Import(accountID string) func(context.Context) {
	return func(ctx context.Context) {
		accounts, err := i.api.GetUserInfo(ctx)
		if err != nil {
			log.Err(err).Msg("failed import")
			return
		}
		account, found := i.findAccount(accounts, accountID)
		if !found {
			log.Err(err).Str("accountID", accountID).Msg("failed import: account not found")
			return
		}

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

func (i *Importer) findAccount(accounts []mono.Account, accountID string) (mono.Account, bool) {
	for _, v := range accounts {
		if v.ID == accountID {
			return v, true
		}
	}

	return mono.Account{}, false
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

func mapTransactions(accountID string, monoTransactions []mono.Transaction) []storage.Transaction {
	transactions := make([]storage.Transaction, len(monoTransactions))
	for i, v := range monoTransactions {
		transactions[i] = storage.Transaction{
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

	return transactions
}
