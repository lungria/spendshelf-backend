package storage

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/transaction"
)

type Storage struct {
	pool *pgxpool.Pool
}

func (s *Storage) Save(ctx context.Context, transactions []transaction.Transaction) error {
	// sql insert
	// on conflict - ignore
	// todo : Using Prepared Statements
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	// todo: tx.Prepare()
	// todo: foreach in transactions - batch.Queue
	batch := pgx.Batch{}
	batch.Queue(`insert into transactions values ($1, $2, $3, $4, $5, $6) on conflict do nothing`, "test", "1/8/1999", "desc", 123, true, 10)
	result := tx.SendBatch(ctx, &batch)
	err = result.Close()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
