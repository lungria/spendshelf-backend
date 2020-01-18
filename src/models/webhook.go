package models

import (
	shalmono "github.com/shal/mono"
)

// WebHook struct using in response from mono API and model in DB
type WebHook struct {
	AccountID     string               `json:"accountId" bson:"accountId"`
	StatementItem shalmono.Transaction `json:"statementItem" bson:"statementItem"`
}
