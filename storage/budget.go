package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Budget describes single month budget.
type Budget struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	StartsAt  time.Time `json:"startsAt"`
	EndsAt    time.Time `json:"endsAt"`
	Limits    []Limit   `json:"limits"`
}

// Limit describe per-category limit inside monthly budget.
type Limit struct {
	CategoryID int32 `json:"categoryID"`
	Amount     int64 `json:"amount"`
}

// BudgetsStorage implements persistent storage layer for budgets and limits in PostgreSQL.
type BudgetsStorage struct {
	pool *pgxpool.Pool
}

// NewBudgetsStorage creates new instance of BudgetsStorage.
func NewBudgetsStorage(pool *pgxpool.Pool) *BudgetsStorage {
	return &BudgetsStorage{
		pool: pool,
	}
}

// GetLast returns last budget from storage.
func (s *BudgetsStorage) GetLast(ctx context.Context) (Budget, error) {
	row := s.pool.QueryRow(
		ctx,
		`select "ID", "startsAt", "endsAt", "createdAt" from budget
			order by "startsAt" desc
			limit 1`)

	b := Budget{}

	err := row.Scan(
		&b.ID,
		&b.StartsAt,
		&b.EndsAt,
		&b.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Budget{}, ErrNotFound
		}

		return Budget{}, err
	}

	const limit = 50

	rows, err := s.pool.Query(
		ctx,
		`select "categoryID", "amount" from "limit"
			where "budgetID" = $1
			order by "amount" desc
			limit $2`,
		b.ID, limit)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Budget{
				ID:        b.ID,
				CreatedAt: b.CreatedAt,
				StartsAt:  b.StartsAt,
				EndsAt:    b.EndsAt,
			}, nil
		}

		return Budget{}, err
	}

	defer rows.Close()

	limits, err := scanLimits(limit, rows)
	if err != nil {
		return Budget{}, err
	}

	b.Limits = limits

	return b, nil
}

func scanLimits(buffSize int, rows pgx.Rows) ([]Limit, error) {
	buffer := make([]Limit, buffSize)
	i := 0

	for rows.Next() {
		l := Limit{}

		err := rows.Scan(
			&l.CategoryID,
			&l.Amount)
		if err != nil {
			return nil, err
		}

		buffer[i] = l

		i++
	}

	result := make([]Limit, i)
	copy(result, buffer)

	return result, nil
}
