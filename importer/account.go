package importer

import (
	"context"
	"fmt"

	"github.com/lungria/spendshelf-backend/importer/mono"
	"github.com/lungria/spendshelf-backend/storage"
)

// UserInfoBankAPI abstracts bank API related to user information.
type UserInfoBankAPI interface {
	// GetUserInfo loads user accounts list.
	GetUserInfo(ctx context.Context) ([]mono.Account, error)
}

// AccountsStorage abstracts persistent storage for accounts.
type AccountsStorage interface {
	// Save information about account.
	Save(ctx context.Context, account storage.Account) error
}

// DefaultAccountImporter loads latest account related data from bank for specified accountID and saves it to storage.
type DefaultAccountImporter struct {
	api      UserInfoBankAPI
	accounts AccountsStorage
}

// NewDefaultAccountImporter create new instance of DefaultAccountImporter.
func NewDefaultAccountImporter(api UserInfoBankAPI, accounts AccountsStorage) *DefaultAccountImporter {
	return &DefaultAccountImporter{api: api, accounts: accounts}
}

// Import latest account related data from bank for specified accountID and save it to storage.
// todo: tests.
func (i *DefaultAccountImporter) Import(ctx context.Context, accountID string) error {
	accounts, err := i.api.GetUserInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed import account '%s' data: %w", accountID, err)
	}

	monoAccount, found := i.findByID(accounts, accountID)
	if !found {
		return fmt.Errorf("failed import account '%s' data: %w", accountID, err)
	}

	err = i.accounts.Save(ctx, storage.Account{
		ID:      monoAccount.ID,
		Balance: monoAccount.Balance,
	})
	if err != nil {
		return fmt.Errorf("failed import account '%s' data: %w", accountID, err)
	}

	return nil
}

func (i *DefaultAccountImporter) findByID(accounts []mono.Account, accountID string) (mono.Account, bool) {
	for _, v := range accounts {
		if v.ID == accountID {
			return v, true
		}
	}

	return mono.Account{}, false
}
