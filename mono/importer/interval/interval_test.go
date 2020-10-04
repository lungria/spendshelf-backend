package interval_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lungria/spendshelf-backend/mono/importer/interval"
	"github.com/lungria/spendshelf-backend/mono/importer/interval/moq"
	"github.com/stretchr/testify/assert"
)

func TestGetInterval_WhenDateDiffIsBelowLimit_IntervalReturned(t *testing.T) {
	mockStorage := moq.TransactionsStorageMock{}
	timeFrom := time.Now().UTC().Add(-1 * time.Minute)
	mockStorage.GetLastTransactionDateFunc = func(ctx context.Context, accountID string) (time.Time, error) {
		return timeFrom, nil
	}
	svc := interval.NewIntervalGenerator(&mockStorage)

	from, to, err := svc.GetInterval(context.Background(), "accID")

	assert.NoError(t, err)
	assert.Equal(t, timeFrom, from)
	assert.GreaterOrEqual(t, time.Now().UTC().Unix(), to.Unix())
}

func TestGetInterval_WhenDateDiffIsAboveLimit_ErrorReturned(t *testing.T) {
	mockStorage := moq.TransactionsStorageMock{}
	timeFrom := time.Now().UTC().Add(-2682001 * time.Second)
	mockStorage.GetLastTransactionDateFunc = func(ctx context.Context, accountID string) (time.Time, error) {
		return timeFrom, nil
	}
	svc := interval.NewIntervalGenerator(&mockStorage)

	_, _, err := svc.GetInterval(context.Background(), "accID")

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "interval too long")
}

func TestGetInterval_WhenStorageReturnsError_ErrorReturned(t *testing.T) {
	expectedErr := errors.New("something went wrong")
	mockStorage := moq.TransactionsStorageMock{}
	mockStorage.GetLastTransactionDateFunc = func(ctx context.Context, accountID string) (time.Time, error) {
		return time.Time{}, expectedErr
	}
	svc := interval.NewIntervalGenerator(&mockStorage)

	_, _, err := svc.GetInterval(context.Background(), "accID")

	assert.Equal(t, expectedErr, err)
}
