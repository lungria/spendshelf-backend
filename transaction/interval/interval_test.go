package interval_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lungria/spendshelf-backend/transaction/interval"
	"github.com/lungria/spendshelf-backend/transaction/interval/moq"

	"github.com/stretchr/testify/assert"
)

func TestGetInterval_WhenDateDiffIsBelowLimit_IntervalReturned(t *testing.T) {
	mockStorage := moq.TransactionsStorageMock{}
	timeFrom := time.Now().UTC().Add(-1 * time.Minute)
	mockStorage.GetLastTransactionDateFunc = func(ctx context.Context, accountID string) (time.Time, error) {
		return timeFrom, nil
	}
	svc := interval.NewGenerator(&mockStorage)

	from, to, err := svc.GetInterval(context.Background(), "accID")

	assert.NoError(t, err)
	assert.Equal(t, timeFrom, from)
	assert.GreaterOrEqual(t, time.Now().UTC().Unix(), to.Unix())
}

func TestGetInterval_WhenDateDiffIsAboveLimit_TimeMoved(t *testing.T) {
	mockStorage := moq.TransactionsStorageMock{}
	timeFrom := time.Now().UTC().Add(-2682001 * time.Second)
	mockStorage.GetLastTransactionDateFunc = func(ctx context.Context, accountID string) (time.Time, error) {
		return timeFrom, nil
	}
	svc := interval.NewGenerator(&mockStorage)

	from, to, err := svc.GetInterval(context.Background(), "accID")

	assert.NoError(t, err)
	assert.Equal(t, timeFrom, from)
	assert.GreaterOrEqual(t, to.Unix(), time.Now().UTC().Unix()/2)
}

func TestGetInterval_WhenStorageReturnsError_ErrorReturned(t *testing.T) {
	expectedErr := errors.New("something went wrong")
	mockStorage := moq.TransactionsStorageMock{}
	mockStorage.GetLastTransactionDateFunc = func(ctx context.Context, accountID string) (time.Time, error) {
		return time.Time{}, expectedErr
	}
	svc := interval.NewGenerator(&mockStorage)

	_, _, err := svc.GetInterval(context.Background(), "accID")

	assert.Equal(t, expectedErr, err)
}
