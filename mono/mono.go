package mono

import (
	"context"

	"github.com/lungria/spendshelf-backend/transaction"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (c *Client) GetTransactions(ctx context.Context) ([]transaction.Transaction, error) {
	panic("implement me")
}
