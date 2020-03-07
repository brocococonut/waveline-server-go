package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Genre - Genre model
type Genre struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}
