package importer

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

// AccountImporter abstracts account import logic implementation.
type AccountImporter interface {
	// Import latest account related data from bank for specified accountID and save it to storage.
	Import(ctx context.Context, accountID string) error
}

// TransactionsImporter abstracts transactions import logic implementation.
type TransactionsImporter interface {
	// Import loads transactions in specified interval for specified accountID and saves them to storage.
	Import(ctx context.Context, accountID string) (time.Time, time.Time, error)
}

// Importer loads latest data from bank for specified accountID.
type Importer struct {
	transactions TransactionsImporter
	accounts     AccountImporter
}

// NewImporter create new instance of Importer.
func NewImporter(transactions TransactionsImporter, accounts AccountImporter) *Importer {
	return &Importer{
		transactions: transactions,
		accounts:     accounts,
	}
}

// Import the latest data from bank for specified accountID.
func (i *Importer) Import(ctx context.Context, accountID string) {
	err := i.accounts.Import(ctx, accountID)
	if err != nil {
		log.Err(err).Msg("failed account import")
		return
	}

	from, to, err := i.transactions.Import(ctx, accountID)
	if err != nil {
		log.Err(err).
			Time("from", from).
			Time("to", to).
			Msg("failed transactions import")
		return
	}

	log.Debug().
		Time("from", from).
		Time("to", to).
		Msg("import finished")
}
