package importer

import (
	"context"
	"fmt"
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
	Save(ctx context.Context, account storage.Account) error
}

// ImportIntervalGenerator generates interval for transaction import.
type ImportIntervalGenerator interface {
	// GetInterval generates interval for transaction import.
	GetInterval(ctx context.Context, accountID string) (from, to time.Time, err error)
}

// Importer loads latest data from bank for specified accountID.
type Importer struct {
	api          BankAPI
	transactions TransactionsStorage
	accounts     AccountsStorage
	intervalGen  ImportIntervalGenerator
}

// NewImporter create new instance of Importer.
func NewImporter(api BankAPI, t TransactionsStorage, a AccountsStorage, gen ImportIntervalGenerator) *Importer {
	return &Importer{
		api:          api,
		transactions: t,
		accounts:     a,
		intervalGen:  gen,
	}
}

// Import latest data from bank for specified accountID.
func (i *Importer) Import(accountID string) func(context.Context) {
	return func(ctx context.Context) {
		account, err := i.fetchAccount(ctx, accountID)
		if err != nil {
			log.Err(err).Msg("failed import")
			return
		}

		err = i.accounts.Save(ctx, account)
		if err != nil {
			log.Err(err).Msg("failed import")
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

		err = i.transactions.Save(ctx, transactions)
		if err != nil {
			log.Err(err).Msg("failed import")
		}
	}
}

func (i *Importer) fetchAccount(ctx context.Context, accountID string) (storage.Account, error) {
	accounts, err := i.api.GetUserInfo(ctx)
	if err != nil {
		return storage.Account{}, err
	}

	monoAccount, found := findByID(accounts, accountID)
	if !found {
		return storage.Account{}, fmt.Errorf("account not found: %v", accountID)
	}

	return storage.Account{
		ID:      monoAccount.ID,
		Balance: monoAccount.Balance,
	}, nil
}

func findByID(accounts []mono.Account, accountID string) (mono.Account, bool) {
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
