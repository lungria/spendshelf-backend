package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Category describes transaction category.
type Category struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Logo    string `json:"logo"`
	Visible bool   `json:"visible"`
}

// CategoryStorage for categories.
type CategoryStorage struct {
	pool *pgxpool.Pool
}

// NewCategoryStorage creates new instance of CategoryStorage.
func NewCategoryStorage(pool *pgxpool.Pool) *CategoryStorage {
	return &CategoryStorage{pool: pool}
}

// GetAll returns existing categories.
func (s *CategoryStorage) GetAll(ctx context.Context) ([]Category, error) {
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
