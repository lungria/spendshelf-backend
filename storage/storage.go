package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/storage/transaction"
)

// ErrNotFound is being returned, if no data was found in database.
var ErrNotFound = errors.New("data not found")

// PostgreSQLStorage for transactions.
type PostgreSQLStorage struct {
	pool *pgxpool.Pool
}

// NewPostgreSQLStorage creates new instance of PostgreSQLStorage.
func NewPostgreSQLStorage(pool *pgxpool.Pool) *PostgreSQLStorage {
	return &PostgreSQLStorage{pool: pool}
}

const insertPrepStatementName = "insert_transactions"

// Save transactions to db with deduplication using transaction ID.
func (s *PostgreSQLStorage) Save(ctx context.Context, transactions []transaction.Transaction) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if _, err = tx.Prepare(ctx, insertPrepStatementName,
		`insert into transaction (ID, time, description, mcc, hold, amount, accountID, categoryID) 
		 values ($1, $2, $3, $4, $5, $6, $7, $8) on conflict do nothing`); err != nil {
		return err
	}

	batch := pgx.Batch{}

	for _, t := range transactions {
		batch.Queue(insertPrepStatementName, t.ID, t.Time, t.Description, t.MCC, t.Hold, t.Amount, t.AccountID, t.CategoryID)
	}

	result := tx.SendBatch(ctx, &batch)

	err = result.Close()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetLastTransactionDate returns date property of latest transaction (sorted by date desc).
func (s *PostgreSQLStorage) GetLastTransactionDate(ctx context.Context, accountID string) (time.Time, error) {
	row := s.pool.QueryRow(
		ctx,
		`select "time" from transaction
		where accountID = $1
		order by "time" desc
		limit 1`,
		accountID)

	var lastKnownTransaction time.Time

	err := row.Scan(&lastKnownTransaction)
	if err != nil {
		if err == pgx.ErrNoRows {
			return time.Time{}, ErrNotFound
		}

		return time.Time{}, err
	}

	return lastKnownTransaction, nil
}

// GetByCategory returns transactions by category.
func (s *PostgreSQLStorage) GetByCategory(ctx context.Context, categoryID int32) ([]transaction.Transaction, error) {
	// todo: add index to categoryID

	rows, err := s.pool.Query(
		ctx,
		`select * from transaction
			where categoryID = $1
			order by "time" desc
			limit 50`,
		categoryID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	buffer := make([]transaction.Transaction, 50)
	i := 0

	for rows.Next() {
		var id string
		var time time.Time
		var description string
		var mcc int32
		var hold bool
		var amount int64
		var accountID string
		var categoryID int32

		err := rows.Scan(&id, &time, &description, &mcc, &hold, &amount, &accountID, &categoryID)
		if err != nil {
			return nil, err
		}

		buffer[i] = transaction.Transaction{
			ID:          id,
			Time:        time,
			Description: description,
			MCC:         mcc,
			Hold:        hold,
			Amount:      amount,
			AccountID:   accountID,
			CategoryID:  categoryID,
		}

		i++
	}
	result := make([]transaction.Transaction, i)
	copy(result, buffer)
	return result, nil
}
