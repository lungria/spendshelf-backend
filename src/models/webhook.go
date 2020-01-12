package models

import "github.com/lungria/mono"

// WebHook struct using in response from mono API and model in DB
type WebHook struct {
	AccountID     string             `json:"accountId" bson:"accountId"`
	StatementItem mono.StatementItem `json:"statementItem" bson:"statementItem"`
}
