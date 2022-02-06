package account

import (
	"context"
	"fmt"

	"github.com/lungria/spendshelf-backend/importer/mono"
)

// BankAPI abstracts bank API related to accounts information.
type BankAPI interface {
	// GetUserInfo loads user accounts list.
	GetUserInfo(ctx context.Context) ([]mono.Account, error)
}

// Storage abstracts persistent storage for accounts.
type Storage interface {
	// Save information about account.
	Save(ctx context.Context, account Account) error
}

// Importer loads the latest account related data from bank for specified accountID and saves it to storage.
type Importer struct {
	api      BankAPI
	accounts Storage
}

// NewImporter create new instance of Importer.
func NewImporter(api BankAPI, accounts Storage) *Importer {
	return &Importer{api: api, accounts: accounts}
}

// Import latest account related data from bank for specified accountID and save it to storage.
func (i *Importer) Import(ctx context.Context, accountID string) error {
	accounts, err := i.api.GetUserInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed import account '%s' data: %w", accountID, err)
	}

	monoAccount, found := i.findByID(accounts, accountID)
	if !found {
		return fmt.Errorf("API response doesn't contain required information for account '%s'", accountID)
	}

	err = i.accounts.Save(ctx, Account{
		ID:      monoAccount.ID,
		Balance: monoAccount.Balance,
	})
	if err != nil {
		return fmt.Errorf("failed import account '%s' data: %w", accountID, err)
	}

	return nil
}

func (i *Importer) findByID(accounts []mono.Account, accountID string) (mono.Account, bool) {
	for _, v := range accounts {
		if v.ID == accountID {
			return v, true
		}
	}

	return mono.Account{}, false
}
