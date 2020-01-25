package models

import (
	"github.com/shal/mono"
)

// WebHook struct using in response from mono API and model in DB
type WebHook struct {
	AccountID     string               `json:"accountId" bson:"accountId"`
	StatementItem mono.Transaction `json:"statementItem" bson:"statementItem"`
}
