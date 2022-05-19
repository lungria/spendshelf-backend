package budget

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lungria/spendshelf-backend/transaction"
)

// Budget describes single month budget.
type Budget struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	StartsAt  time.Time `json:"startsAt"`
	EndsAt    time.Time `json:"endsAt"`
	Limits    []Limit   `json:"limits"`
}

// CreateBudgetRequest describes single month budget creation request.
type CreateBudgetRequest struct {
	Days   int     `json:"days"`
	Limits []Limit `json:"limits"`
}

// Limit describe per-category limit inside monthly budget.
type Limit struct {
	CategoryID int32 `json:"categoryId"`
	Amount     int64 `json:"amount"`
}

// Repository implements persistent storage layer for budgets and limits in PostgreSQL.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates new instance of Repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

// Create returns last budget from repository.
func (s *Repository) Create(ctx context.Context, budget CreateBudgetRequest) (Budget, error) {
	// todo
	return Budget{}, nil
}

// GetLast returns last budget from repository.
func (s *Repository) GetLast(ctx context.Context) (Budget, error) {
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
			return Budget{}, transaction.ErrNotFound
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
