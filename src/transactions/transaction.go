package transactions

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Transaction represents a model of transactions in database
type Transaction struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Time        time.Time          `json:"dateTime" bson:"time"`
	Description string             `json:"description" bson:"description"`
	CategoryID  primitive.ObjectID `json:"categoryID,omitempty" bson:"categoryID,omitempty"`
	Amount      int64              `json:"amount" bson:"amount"`
	Bank        primitive.ObjectID `json:"bankId" bson:"bank"`
}
