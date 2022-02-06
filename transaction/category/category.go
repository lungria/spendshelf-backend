package category

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	// common categories

	// Default is the ID of category, that must be used for all new imported transactions.
	Default = int32(1)
	// Ignored is the ID of category, that would be ignored in report.
	Ignored = int32(127)
)

// Category describes transaction category.
type Category struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

// Repository for categories.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates new instance of Repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// GetAll returns existing categories.
func (s *Repository) GetAll(ctx context.Context) ([]Category, error) {
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
