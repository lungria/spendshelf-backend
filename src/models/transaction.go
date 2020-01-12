package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Transaction represents a model of transactions in database
type Transaction struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Time            time.Time          `json:"dateTime" bson:"time"`
	Description     string             `json:"description" bson:"description"`
	Category        *Category          `json:"category,omitempty" bson:"category,omitempty"`
	Amount          int64              `json:"amount" bson:"amount"`
	Balance         int64              `json:"balance" bson:"balance"`
	Bank            string             `json:"bank" bson:"bank"`
	BankTransaction WebHook            `json:"-" bson:"bankTransaction"`
}
