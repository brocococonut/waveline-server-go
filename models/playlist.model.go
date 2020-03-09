package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Playlist - A list of songs grouped together
type Playlist struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	Tracks    []Track            `json:"tracks" bson:"tracks"`
}

type PlaylistBasic struct {
	ID        primitive.ObjectID   `json:"_id" bson:"_id"`
	Name      string               `json:"name" bson:"name"`
	UpdatedAt time.Time            `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	Tracks    []primitive.ObjectID `json:"tracks" bson:"tracks"`
}

type PlaylistPopulated struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	Tracks    []TrackExtended    `json:"tracks" bson:"tracks"`
}

type PlaylistExtended struct {
	ID        string    `json:"_id" bson:"_id"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Name      string    `json:"name" bson:"name"`
	Tracks    int64     `json:"tracks"`
	Pictures  []string  `json:"pictures"`
	ReadOnly  bool      `json:"readonly"`
}
