package mono

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lungria/spendshelf-backend/transaction"

	"strconv"
	"strings"
	"time"
)

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

func (c *Client) GetTransactions(ctx context.Context, query GetTransactionsQuery) ([]transaction.Transaction, error) {
	uri := fmt.Sprintf("%s/personal/statement%s", c.baseURL, query.asRoute())

	response, err := c.performRequest(ctx, uri, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get transactions: %w", err)
	}
	if len(response) == 0 {
		return nil, fmt.Errorf("unable to get transactions: empty response without error received")
	}

	transactions := make(responseTransactions, 0)
	if err = json.Unmarshal(response, &transactions); err != nil {
		return nil, fmt.Errorf("failed unmarshal transactions form json body: %w", err)
	}
	return transactions.AsPublicAPIModel(), nil
}

type responseTransaction struct {
	ID          string `json:"id"`
	Time        Time   `json:"time"`
	Description string `json:"description"`
	MCC         int32  `json:"mcc"`
	Hold        bool   `json:"hold"`
	Amount      int64  `json:"amount"`
}

type responseTransactions []responseTransaction

func (t responseTransactions) AsPublicAPIModel() []transaction.Transaction {
	transactions := make([]transaction.Transaction, len(t))
	for i, v := range t {
		transactions[i] = transaction.Transaction{
			BankID:      v.ID,
			Time:        time.Time(v.Time),
			Description: v.Description,
			MCC:         v.MCC,
			Hold:        v.Hold,
			Amount:      v.Amount,
		}
	}
}
