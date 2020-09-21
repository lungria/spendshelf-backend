package mono

import (
	"context"

	"github.com/lungria/spendshelf-backend/transaction"
)

type Client struct {
}

func (c *Client) GetTransactions(ctx context.Context) ([]transaction.Transaction, error) {
	panic("implement me")
}
