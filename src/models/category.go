package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Category is general struct for spendshelf categories
type Category struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	NormalizedName string             `json:"normalizedName" bson:"normalizedName"`
}
