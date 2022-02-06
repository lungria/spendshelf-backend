package storage

import "github.com/jackc/pgx/v4/pgxpool"

type Worker struct {
	pool *pgxpool.Pool
}

func NewWorker(pool *pgxpool.Pool) *Worker {
	return &Worker{pool: pool}
}

func (w *Worker) Start() {}

func (w *Worker) Close() error {
	w.pool.Close()
	return nil
}
