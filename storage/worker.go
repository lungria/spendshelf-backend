package storage

import "github.com/jackc/pgx/v4/pgxpool"

// Worker is a dummy app/Worker interface implementation, so DB pool can be managed with other parts of the
// app lifecycle.
type Worker struct {
	pool *pgxpool.Pool
}

// NewWorker creates new instance of a Worker.
func NewWorker(pool *pgxpool.Pool) *Worker {
	return &Worker{pool: pool}
}

// Start is a stub implementation for app/Worker interface. Does nothing.
func (w *Worker) Start() {}

// Close pgx connection pool.
func (w *Worker) Close() error {
	w.pool.Close()
	return nil
}
