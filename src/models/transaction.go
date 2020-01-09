package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Transaction represents a model of transactions in database
type Transaction struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Time        primitive.DateTime `json:"dateTime" bson:"time"`
	Description string             `json:"description" bson:"description"`
	CategoryID  primitive.ObjectID `json:"categoryId,omitempty" bson:"category_id,omitempty"`
	Amount      int64              `json:"amount" bson:"amount"`
	Balance     int64              `json:"balance" bson:"balance"`
}
