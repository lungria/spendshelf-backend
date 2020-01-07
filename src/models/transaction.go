package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Transaction represents a model of transactions in database
type Transaction struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Time        int32              `json:"time" bson:"time"`
	Description string             `json:"description" bson:"description"`
	Category    string             `json:"category" bson:"category,omitempty"`
	Amount      int64              `json:"amount" bson:"amount"`
	Balance     int64              `json:"balance" bson:"balance"`
}
