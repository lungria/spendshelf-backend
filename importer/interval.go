package importer

import (
	"context"
	"time"
)

type InMemoryPersistantIntervalGenerator struct {
}

func (gen *InMemoryPersistantIntervalGenerator) GetInterval(ctx context.Context) (from time.Time, to time.Time, err error) {
	// 1. get last successfuly synced transaction date - lastKnownTransactionDate
	// 2. get current date/time unix
	// 3. if diff (currDate, lastKnownTransactionDate)  > 2682000 secs - log err, edge case, do not care now

}
