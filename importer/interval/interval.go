package interval

import (
	"context"
	"fmt"
	"time"
)

const maxAllowedIntervalDuration = 2682000

type TransactionsStorage interface {
	GetLastTransactionDate(ctx context.Context, accountID string) (time.Time, error)
}

type SimpleIntervalGenerator struct {
	storage TransactionsStorage
}

func NewSimpleIntervalGenerator(storage TransactionsStorage) *SimpleIntervalGenerator {
	return &SimpleIntervalGenerator{storage: storage}
}

func (gen *SimpleIntervalGenerator) GetInterval(ctx context.Context, accountID string) (from time.Time, to time.Time, err error) {
	lastKnownTransactionDate, err := gen.storage.GetLastTransactionDate(ctx, accountID)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	nowUtc := time.Now().UTC()
	diffSecs := nowUtc.Sub(lastKnownTransactionDate.UTC()).Seconds()
	if diffSecs > maxAllowedIntervalDuration {
		return time.Time{}, time.Time{}, fmt.Errorf("interval too long, lastKnownTransactionDate: %v", lastKnownTransactionDate)
	}

	return lastKnownTransactionDate.UTC(), nowUtc, nil
}
