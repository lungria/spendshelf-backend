package mono

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetTransactionsQuery describes parameters for GetTransactions monobank request.
type GetTransactionsQuery struct {
	Account string
	From    time.Time
	To      time.Time
}

func (q *GetTransactionsQuery) asRoute() string {
	var sb strings.Builder

	sb.WriteString("/")
	sb.WriteString(q.Account)
	sb.WriteString("/")
	sb.WriteString(strconv.FormatInt(q.From.Unix(), 10))
	sb.WriteString("/")
	sb.WriteString(strconv.FormatInt(q.To.Unix(), 10))

	return sb.String()
}

// GetTransactions loads transactions from monobank with specified query parameters.
func (c *Client) GetTransactions(ctx context.Context, query GetTransactionsQuery) ([]Transaction, error) {
	uri := fmt.Sprintf("%s/personal/statement%s", c.baseURL, query.asRoute())

	response, err := c.performRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get transactions: %w", err)
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("unable to get transactions: empty response without error received")
	}

	transactions := make([]Transaction, 0)
	if err = json.Unmarshal(response, &transactions); err != nil {
		return nil, fmt.Errorf("failed unmarshal transactions form json body: %w", err)
	}

	return transactions, nil
}

// Transaction describes monobank transaction.
type Transaction struct {
	ID          string `json:"ID"`
	Time        Time   `json:"time"`
	Description string `json:"description"`
	MCC         int32  `json:"mcc"`
	Hold        bool   `json:"hold"`
	Amount      int64  `json:"amount"`
}
