package storage

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/storage/category"
)

var (
	// ErrNotFound is being returned, if no data was found in database.
	ErrNotFound = errors.New("data not found")
	// ErrAtLeastOneUpdateFieldRequired is being returned, if all update fields are nil.
	ErrAtLeastOneUpdateFieldRequired = errors.New("at least single update field must not be nil")
)

// Transaction describes single user's transaction.
type Transaction struct {
	ID            string    `json:"id"`
	Time          time.Time `json:"time"`
	Description   string    `json:"description"`
	MCC           int32     `json:"mcc"`
	Hold          bool      `json:"hold"`
	Amount        int64     `json:"amount"`
	AccountID     string    `json:"accountId"`
	CategoryID    int32     `json:"categoryId"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
	Comment       *string   `json:"comment"`
}

// UpdateTransactionCommand describes transaction update parameters.
type UpdateTransactionCommand struct {
	Query         Query
	UpdatedFields UpdatedFields
}

// Page controls pagination settings. Both fields are optional and would be set to defaults if not passed.
type Page struct {
	Limit  int
	Offset int
}

func (p Page) withDefaults() Page {
	limit := p.Limit
	offset := p.Offset

	if limit <= 0 {
		limit = 50
	}

	if offset < 0 {
		offset = 0
	}

	return Page{
		Limit:  limit,
		Offset: offset,
	}
}

// UpdatedFields for UpdateTransactionCommand. All fields are optional, but at least one field must be non-nil.
type UpdatedFields struct {
	CategoryID *int32
	Comment    *string
}

// Valid checks if updated fields are valid. Checks that at least single field is not nil.
func (f UpdatedFields) Valid() bool {
	if f.CategoryID == nil && f.Comment == nil {
		return false
	}

	return true
}

// TransactionStorage for transactions.
type TransactionStorage struct {
	pool *pgxpool.Pool
}

// NewTransactionStorage creates new instance of TransactionStorage.
func NewTransactionStorage(pool *pgxpool.Pool) *TransactionStorage {
	return &TransactionStorage{pool: pool}
}

const insertPrepStatementName = "insert_transactions"

// Save transactions to db with deduplication using transaction ID.
func (s *TransactionStorage) Save(ctx context.Context, transactions []Transaction) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if _, err = tx.Prepare(ctx, insertPrepStatementName,
		`insert into transaction 
		("ID", "time", "description", "mcc", "hold", "amount", "accountID", "categoryID", "lastUpdatedAt") 
		 values ($1, $2, $3, $4, $5, $6, $7, $8, current_timestamp(0)) on conflict do nothing`); err != nil {
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
// Returns storage.ErrNotFound if transaction not found by query.
func (s *TransactionStorage) GetLastTransactionDate(ctx context.Context, accountID string) (time.Time, error) {
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

// GetOne tries to find one transaction by query.
func (s *TransactionStorage) GetOne(ctx context.Context, query Query) (Transaction, error) {
	transactions, err := s.Get(ctx, query, Page{Limit: 1})
	if err != nil {
		return Transaction{}, err
	}

	return transactions[0], nil
}

// Get returns transactions by filter.
// Returns storage.ErrNotFound if transaction not found by query.
func (s *TransactionStorage) Get(ctx context.Context, query Query, page Page) ([]Transaction, error) {
	page = page.withDefaults()
	sqlBuilder := strings.Builder{}

	sqlBuilder.WriteString("select * from transaction ")

	filter, filterParams := query.AsSQL()
	sqlBuilder.WriteString(filter)

	sqlBuilder.WriteString(`order by "time" desc `)

	pageParams := make([]interface{}, 0)
	pageParams = append(pageParams, page.Limit)
	sqlBuilder.WriteString(fmt.Sprintf(`limit $%v `, len(pageParams)))

	pageParams = append(pageParams, page.Offset)
	sqlBuilder.WriteString(fmt.Sprintf(`offset $%v `, len(pageParams)))

	sqls := sqlBuilder.String()

	rows, err := s.pool.Query(
		ctx,
		sqls,
		append(filterParams, pageParams)...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, err
	}

	defer rows.Close()

	return scanTransactions(page.Limit, rows)
}

// UpdateTransaction allows to partially update transaction.
func (s *TransactionStorage) UpdateTransaction(
	ctx context.Context,
	cmd UpdateTransactionCommand) (Transaction, error) {
	sqlBuilder := strings.Builder{}

	if !cmd.UpdatedFields.Valid() {
		return Transaction{}, ErrAtLeastOneUpdateFieldRequired
	}

	sqlBuilder.WriteString(`update "transaction" set `)

	setParams := make([]interface{}, 0)

	if cmd.UpdatedFields.CategoryID != nil {
		setParams = append(setParams, *cmd.UpdatedFields.CategoryID)
		sqlBuilder.WriteString(fmt.Sprintf(`"categoryID" = $%v, `, len(setParams)))
	}

	if cmd.UpdatedFields.Comment != nil {
		setParams = append(setParams, *cmd.UpdatedFields.Comment)
		sqlBuilder.WriteString(fmt.Sprintf(`"comment" = $%v, `, len(setParams)))
	}

	sqlBuilder.WriteString(`"lastUpdatedAt" = current_timestamp(0) `)

	filter, filterParams := cmd.Query.AsSQL()
	sqlBuilder.WriteString(filter)

	cmdResult, err := s.pool.Exec(
		ctx,
		sqlBuilder.String(),
		append(setParams, filterParams)...)
	if err != nil {
		return Transaction{}, err
	}

	if cmdResult.RowsAffected() == 0 {
		return Transaction{}, errors.New("failed to update transaction")
	}

	row, err := s.pool.Query(
		ctx,
		`select * from transaction
		where "ID" = $1`,
		cmd.Query.ID)
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
func (s *TransactionStorage) GetReport(ctx context.Context, from, to time.Time) (map[int32]int64, error) {
	rows, err := s.pool.Query(
		ctx,
		`select "categoryID", sum("amount") as "amount" from transaction
		 where "time" > $1 AND "time" <= $2 AND "categoryID" != $3
		 group by "categoryID"`,
		from, to, category.Ignored)
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

func scanTransactions(buffSize int, rows pgx.Rows) ([]Transaction, error) {
	buffer := make([]Transaction, buffSize)
	i := 0

	// todo: what if only one row?!
	// todo2: check similar code in mongo client
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
			&t.LastUpdatedAt,
			&t.Comment)
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
