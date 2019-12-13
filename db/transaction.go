package db

import "github.com/lungria/mono"

// Transaction ...
type Transaction struct {
	AccountId     string             `json:"account" bson:"account_id"`
	StatementItem mono.StatementItem `json:"statementItem" bson:"statement_item"`
}
