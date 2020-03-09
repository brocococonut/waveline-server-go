package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Track - a single track able to be played
type Track struct {
	ID         primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Album      primitive.ObjectID   `json:"album,omitempty" bson:"album,omitempty"`
	Artist     string               `json:"artist" bson:"artist"`
	Artists    []primitive.ObjectID `json:"artists,omitempty" bson:"artists,omitempty"`
	CreatedAt  time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Duration   int                  `json:"duration" bson:"duration"`
	Favourited bool                 `json:"favourited" bson:"favourited"`
	Genre      primitive.ObjectID   `json:"genre,omitempty" bson:"genre,omitempty"`
	LastPlay   time.Time            `json:"last_play,omitempty" bson:"last_play,omitempty"`
	Lossless   bool                 `json:"lossless" bson:"lossless"`
	Name       string               `json:"name" bson:"name"`
	Path       string               `json:"path,omitempty" bson:"path,omitempty"`
	Plays      int                  `json:"plays" bson:"plays"`
	UpdatedAt  time.Time            `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Year       int                  `json:"year,omitempty" bson:"year,omitempty"`
}

type TrackExtended struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Album      Album              `json:"album,omitempty" bson:"album,omitempty"`
	Artist     string             `json:"artist" bson:"artist"`
	Artists    []Artist           `json:"artists,omitempty" bson:"artists,omitempty"`
	CreatedAt  time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Duration   int                `json:"duration" bson:"duration"`
	Favourited bool               `json:"favourited" bson:"favourited"`
	Genre      Genre              `json:"genre,omitempty" bson:"genre,omitempty"`
	LastPlay   time.Time          `json:"last_play,omitempty" bson:"last_play,omitempty"`
	Lossless   bool               `json:"lossless" bson:"lossless"`
	Name       string             `json:"name" bson:"name"`
	Path       string             `json:"path,omitempty" bson:"path,omitempty"`
	Plays      int                `json:"plays" bson:"plays"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	Year       int                `json:"year,omitempty" bson:"year,omitempty"`
}
