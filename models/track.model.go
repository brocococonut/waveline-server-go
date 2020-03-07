package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Track - a single track able to be played
type Track struct {
	ID         primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name" bson:"name"`
	Artists    []primitive.ObjectID `json:"artists,omitempty" bson:"artists,omitempty"`
	Artist     string               `json:"artist" bson:"artist"`
	Album      primitive.ObjectID   `json:"album,omitempty" bson:"album,omitempty"`
	Genre      primitive.ObjectID   `json:"genre,omitempty" bson:"genre,omitempty"`
	Duration   int                  `json:"duration" bson:"duration"`
	Plays      int                  `json:"plays" bson:"plays"`
	Path       string               `json:"path,omitempty" bson:"path,omitempty"`
	Favourited bool                 `json:"favourited" bson:"favourited"`
	LastPlay   time.Time            `json:"last_play,omitempty" bson:"last_play,omitempty"`
	Year       int                  `json:"year,omitempty" bson:"year,omitempty"`
	Lossless   bool                 `json:"lossless" bson:"lossless"`
	UpdatedAt  time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt  time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type TrackExtended struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Artist     string             `json:"artist" bson:"artist"`
	Duration   int                `json:"duration" bson:"duration"`
	Plays      int                `json:"plays" bson:"plays"`
	Path       string             `json:"path,omitempty" bson:"path,omitempty"`
	Favourited bool               `json:"favourited" bson:"favourited"`
	LastPlay   time.Time          `json:"last_play,omitempty" bson:"last_play,omitempty"`
	Year       int                `json:"year,omitempty" bson:"year,omitempty"`
	Lossless   bool               `json:"lossless" bson:"lossless"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt  time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Album      Album              `json:"album,omitempty" bson:"album,omitempty"`
	Artists    []Artist           `json:"artists,omitempty" bson:"artists,omitempty"`
	Genre      []Genre            `json:"genre,omitempty" bson:"genre,omitempty"`
}
