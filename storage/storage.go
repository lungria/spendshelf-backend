package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/transaction"
)

// PostgreSQLStorage for transactions.
type PostgreSQLStorage struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLStorage(pool *pgxpool.Pool) *PostgreSQLStorage {
	return &PostgreSQLStorage{pool: pool}
}

const insertPreparedStatementName = "insert_transactions"

// Save transactions to db with dedublication using transaction ID.
func (s *PostgreSQLStorage) Save(ctx context.Context, transactions []transaction.Transaction) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Prepare(
		ctx,
		insertPreparedStatementName,
		`insert into transactions (bankID, time, description, mcc, hold, amount) 
		 values ($1, $2, $3, $4, $5, $6) on conflict do nothing`)
	if err != nil {
		return err
	}

	batch := pgx.Batch{}
	for _, t := range transactions {
		batch.Queue(insertPreparedStatementName, t.BankID, t.Time, t.Description, t.MCC, t.Hold, t.Amount)
	}
	result := tx.SendBatch(ctx, &batch)
	err = result.Close()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *PostgreSQLStorage) GetLastTransactionDate(ctx context.Context, accountID string) (time.Time, error) {
	// todo
	panic("implement me")
}
