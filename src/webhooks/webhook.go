package webhooks

import "github.com/lungria/mono"

// WebHook struct using in response from mono API and model in DB
type WebHook struct {
	AccountID     string             `json:"account" bson:"account_id"`
	StatementItem mono.StatementItem `json:"statementItem" bson:"statement_item"`
}
