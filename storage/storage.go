package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// DefaultCategoryID is the ID of category, that must be used for all new imported transactions.
const DefaultCategoryID = 1

// ErrNotFound is being returned, if no data was found in database.
var ErrNotFound = errors.New("data not found")

// Transaction describes single user's transaction.
type Transaction struct {
	ID            string    `json:"id"`
	Time          time.Time `json:"time"`
	Description   string    `json:"description"`
	MCC           int32     `json:"mcc"`
	Hold          bool      `json:"hold"`
	Amount        int64     `json:"amount"`
	AccountID     string    `json:"accountID"`
	CategoryID    int32     `json:"categoryID"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}

// Category describes transaction category.
type Category struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

// UpdateTransactionCommand describes transaction update parameters.
type UpdateTransactionCommand struct {
	// filter
	ID            string
	LastUpdatedAt time.Time
	// would be updated
	CategoryID int32
}

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
func (s *PostgreSQLStorage) Save(ctx context.Context, transactions []Transaction) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if _, err = tx.Prepare(ctx, insertPrepStatementName,
		`insert into transaction ("ID", "time", "description", "mcc", "hold", "amount", "accountID", "categoryID", "lastUpdatedAt") 
		 values ($1, $2, $3, $4, $5, $6, $7, $8, current_timestamp) on conflict do nothing`); err != nil {
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
		where "accountID" = $1
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
func (s *PostgreSQLStorage) GetByCategory(ctx context.Context, categoryID int32) ([]Transaction, error) {
	const limit = 50

	rows, err := s.pool.Query(
		ctx,
		`select * from transaction
			where "categoryID" = $1
			order by "time" desc
			limit $2`,
		categoryID, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return scanTransactions(limit, rows)
}

func scanTransactions(buffSize int, rows pgx.Rows) ([]Transaction, error) {
	buffer := make([]Transaction, buffSize)
	i := 0

	for rows.Next() {
		t := Transaction{}

		err := rows.Scan(
			&t.ID,
			&t.Time,
			&t.Description,
			&t.MCC,
			&t.Hold,
			&t.Amount,
			&t.AccountID,
			&t.CategoryID,
			&t.LastUpdatedAt)
		if err != nil {
			return nil, err
		}

		buffer[i] = t

		i++
	}

	result := make([]Transaction, i)
	copy(result, buffer)

	return result, nil
}

// UpdateTransaction allows to partially update transaction.
func (s *PostgreSQLStorage) UpdateTransaction(
	ctx context.Context,
	params UpdateTransactionCommand) (Transaction, error) {
	cmd, err := s.pool.Exec(
		ctx,
		`update "transaction"
			set "categoryID" = $1,
			"lastUpdatedAt" = current_timestamp
		 where "ID" = $2 AND "lastUpdatedAt" = $3`,
		params.CategoryID, params.ID, params.LastUpdatedAt)
	if err != nil {
		return Transaction{}, err
	}

	if cmd.RowsAffected() == 0 {
		return Transaction{}, errors.New("failed to update transaction")
	}

	row, err := s.pool.Query(
		ctx,
		`select * from transaction
		where "ID" = $1`,
		params.ID)
	if err != nil {
		return Transaction{}, err
	}

	result, err := scanTransactions(1, row)
	if err != nil {
		return Transaction{}, err
	}

	return result[0], nil
}

// GetReport generates spending report for set date/time interval.
func (s *PostgreSQLStorage) GetReport(ctx context.Context, from, to time.Time) (map[int32]int64, error) {
	rows, err := s.pool.Query(
		ctx,
		`select "categoryID", sum("amount") as "amount" from transaction
		 where "time" > $1 AND "time" <= $2
		 group by "categoryID"`,
		from, to)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	result := make(map[int32]int64)

	for rows.Next() {
		var categoryID int32
		var amount int64

		err := rows.Scan(
			&categoryID,
			&amount)
		if err != nil {
			return nil, err
		}

		result[categoryID] = amount
	}

	return result, nil
}

// GetCategories returns existing categories.
func (s *PostgreSQLStorage) GetCategories(ctx context.Context) ([]Category, error) {
	const limit = 20

	rows, err := s.pool.Query(
		ctx,
		`select "ID", "name", "logo" from category
		 limit $1`,
		limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	buffer := make([]Category, limit)
	i := 0

	for rows.Next() {
		var c Category

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Logo)
		if err != nil {
			return nil, err
		}

		buffer[i] = c

		i++
	}

	result := make([]Category, i)
	copy(result, buffer)

	return result, nil
}
