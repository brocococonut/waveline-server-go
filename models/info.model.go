package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Info - information
type Info struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Start    time.Time          `json:"start,omitempty" bson:"start,omitempty"`
	End      time.Time          `json:"end,omitempty" bson:"end,omitempty"`
	LastScan time.Time          `json:"last_scan,omitempty" bson:"last_scan,omitempty"`
	Seconds  int                `json:"seconds,omitempty" bson:"seconds,omitempty"`
	Tracks   int                `json:"tracks,omitempty" bson:"tracks,omitempty"`
	Albums   int                `json:"albums,omitempty" bson:"albums,omitempty"`
	Artists  int                `json:"artists,omitempty" bson:"artists,omitempty"`
	Size     int                `json:"size,omitempty" bson:"size,omitempty"`
	Mount    string             `json:"mount,omitempty" bson:"mount,omitempty"`
}
