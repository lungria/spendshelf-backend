package categories

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	Id             CategoryId `json:"account" bson:"_id"`
	Name           string     `json:"name" bson:"name"`
	NormalizedName string     `json:"normalizedName" bson:"normalizedName"`
}

type CategoryId primitive.ObjectID
