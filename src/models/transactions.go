package models

import "github.com/lungria/mono"

// Transaction struct using in response from mono API and model in DB
type Transaction struct {
	AccountId     string             `json:"account" bson:"account_id"`
	StatementItem mono.StatementItem `json:"statementItem" bson:"statement_item"`
}
