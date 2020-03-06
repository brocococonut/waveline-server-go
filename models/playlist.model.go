package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Playlist - A list of songs grouped together
type Playlist struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Tracks    []Track            `json:"tracks,omitempty" bson:"tracks,omitempty"`
}
