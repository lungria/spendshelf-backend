package interval

import (
	"context"
	"fmt"
	"time"

	"github.com/lungria/spendshelf-backend/transaction"
)

// maxAllowedIntervalDuration limits max length of interval (in seconds) that mono API allows us to query.
const maxAllowedIntervalDuration = 2682000

// TransactionsStorage abstracts data access layer for already imported transactions.
// todo: untie implementation from DB: cache latest transaction date in memory.
type TransactionsStorage interface {
	// GetLastTransactionDate returns date property of latest transaction (sorted by date desc).
	GetLastTransactionDate(ctx context.Context, accountID string) (time.Time, error)
}

// Generator creates interval based on latest stored transaction.
type Generator struct {
	storage TransactionsStorage
}

// NewGenerator creates new instance of Generator.
func NewGenerator(storage TransactionsStorage) *Generator {
	return &Generator{storage: storage}
}

// GetInterval creates interval based on latest stored transaction. It will return error if there latest transaction was
// created more than maxAllowedIntervalDuration seconds ago.
func (gen *Generator) GetInterval(ctx context.Context, accountID string) (from, to time.Time, err error) {
	nowUtc := time.Now().UTC()

	lastKnownTransactionDate, err := gen.storage.GetLastTransactionDate(ctx, accountID)
	switch err {
	case transaction.ErrNotFound:
		{
			lastKnownTransactionDate = nowUtc.Add(-maxAllowedIntervalDuration * time.Second)
		}
	case nil:
		{
		}
	default:
		{
			return time.Time{}, time.Time{}, err
		}
	}

	diffSecs := nowUtc.Sub(lastKnownTransactionDate.UTC()).Seconds()
	if diffSecs > maxAllowedIntervalDuration {
		return time.Time{}, time.Time{}, intervalTooLongErr(lastKnownTransactionDate)
	}

	return lastKnownTransactionDate.UTC(), nowUtc, nil
}

func intervalTooLongErr(lastKnownTransactionDate time.Time) error {
	return fmt.Errorf("interval too long, lastKnownTransactionDate: %v", lastKnownTransactionDate)
}
