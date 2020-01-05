package transactions

import (
	"github.com/lungria/mono"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Transaction represents a model of transactions in database
type Transaction struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Category      string             `json:"category" bson:"category,omitempty"`
	AccountID     string             `json:"account" bson:"account_id"`
	StatementItem mono.StatementItem `json:"statementItem" bson:"statement_item"`
}
