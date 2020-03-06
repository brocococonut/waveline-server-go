package models

import "time"

// Track - a single track able to be played
type Track struct {
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Artists    []Artist  `json:"artists,omitempty" bson:"artists,omitempty"`
	Artist     string    `json:"artist,omitempty" bson:"artist,omitempty"`
	Album      Album     `json:"album,omitempty" bson:"album,omitempty"`
	Genre      Genre     `json:"genre,omitempty" bson:"genre,omitempty"`
	Duration   int       `json:"duration,omitempty" bson:"duration,omitempty"`
	Plays      int       `json:"plays,omitempty" bson:"plays,omitempty"`
	Path       string    `json:"path,omitempty" bson:"path,omitempty"`
	Favourited bool      `json:"favourited,omitempty" bson:"favourited,omitempty"`
	LastPlay   time.Time `json:"last_play,omitempty" bson:"last_play,omitempty"`
	Year       int       `json:"year,omitempty" bson:"year,omitempty"`
	Lossless   bool      `json:"lossless,omitempty" bson:"lossless,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
