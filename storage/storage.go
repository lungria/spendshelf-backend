package storage

import (
	"context"

	"github.com/lungria/spendshelf-backend/transaction"
)

type Storage struct {
}

func (s *Storage) Save(ctx context.Context, transactions []transaction.Transaction) error {
	// sql insert
	// on conflict - ignore
	// todo : Using Prepared Statements
	panic("implement me")
}
