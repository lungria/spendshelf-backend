package importer

import (
	"context"
	"errors"
	"time"
)

const maxAllowedIntervalDuration = 2682000

type TransactionsStorage interface {
	GetLastTransactionDate(ctx context.Context, accountID string) (time.Time, error)
}

type SimpleIntervalGenerator struct {
	storage TransactionsStorage
}

func (gen *SimpleIntervalGenerator) GetInterval(ctx context.Context, accountID string) (from time.Time, to time.Time, err error) {
	date, err := gen.storage.GetLastTransactionDate(ctx, accountID)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	nowUtc := time.Now().UTC()
	diffSecs := nowUtc.Sub(date.UTC()).Seconds()
	if diffSecs > maxAllowedIntervalDuration {
		return time.Time{}, time.Time{}, errors.New("interval too long")
	}

	return date.UTC(), nowUtc, nil
}
