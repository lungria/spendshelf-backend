package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Account describes single user's bank account.
type Account struct {
	ID            string    `json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	Description   string    `json:"description"`
	Balance       int64     `json:"balance"`
	Currency      string    `json:"currency"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}

// AccountsStorage implements persistent storage layer for accounts in PostgreSQL.
type AccountsStorage struct {
	pool *pgxpool.Pool
}

// NewAccountsStorage creates new instance of AccountsStorage.
func NewAccountsStorage(pool *pgxpool.Pool) *AccountsStorage {
	return &AccountsStorage{
		pool: pool,
	}
}

// Save account to db. If conflict (on ID) occurs - only "lastUpdatedAt" and "balance" fields would be updated.
func (s *AccountsStorage) Save(ctx context.Context, account Account) error {
	cmd, err := s.pool.Exec(
		ctx,
		`insert into "account"
			 ("ID", "createdAt", "description", "balance", "currency", "lastUpdatedAt")
			 values ($1, current_timestamp(0), $1, $2, $3, current_timestamp(0))
			 on conflict ("ID") do update
			 set "balance" = "excluded"."balance", "lastUpdatedAt" = current_timestamp(0)`,
		account.ID, account.Balance, account.Currency)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("failed to upsert account: %v", account.ID)
	}

	return nil
}

// GetAll accounts from database.
func (s *AccountsStorage) GetAll(ctx context.Context) ([]Account, error) {
	panic("implement me")
}
