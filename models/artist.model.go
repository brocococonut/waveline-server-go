package models

import (
	"context"
	"fmt"
	"time"

	"github.com/brocococonut/waveline-server-go/wavelineutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Artist - an artist of music stored in the DB
type Artist struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Picture   string             `json:"picture,omitempty" bson:"picture,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

func (*Artist) findOrCreate(
	col *mongo.Collection,
	names []string,
	spotifyClient, spotifySecret string,
) (artists []Artist, err error) {
	spot := wavelineutils.Spotify{}
	spot.Authorize(spotifyClient, spotifySecret)

	// Loop over names to find artists
	for _, name := range names {
		var artist Artist

		// Find the artist
		if err = col.FindOne(context.TODO(), bson.M{
			name: name,
		}).Decode(&artist); err != nil {
			continue
		}

		// Create a new artist if that one didn't exist
		if artist.Name == "" {
			artist = Artist{
				ID:        primitive.NewObjectID(),
				Name:      name,
				CreatedAt: time.Now(),
				Picture:   spot.ArtistPicture(fmt.Sprintf("album:%s artist:%s", "", name)),
			}

			// Insert the new artist to the db
			if _, err = col.InsertOne(context.TODO(), artist); err != nil {
				continue
			}
		}

		artists = append(artists, artist)
	}

	return artists, err
}
