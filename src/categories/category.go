package categories

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	Id             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	NormalizedName string             `json:"normalizedName" bson:"normalizedName"`
}
