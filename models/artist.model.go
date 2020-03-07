package models

import (
	"context"
	"fmt"
	"time"

	"github.com/brocococonut/waveline-server-go/wavespotify"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Artist - an artist of music stored in the DB
type Artist struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Picture   string             `json:"picture" bson:"picture"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

// FindOrCreate - Find or create an artist document
func (*Artist) FindOrCreate(
	db *mongo.Database,
	names []string,
	spotifyClient, spotifySecret string,
) (artists []Artist, err error) {
	col := db.Collection("artists")
	spot := wavespotify.Spotify{}
	spot.Authorize(spotifyClient, spotifySecret)

	// Loop over names to find artists
	for _, name := range names {
		var artist Artist

		// Check to see if the artist already exists
		if err := col.FindOne(context.TODO(), bson.M{
			"name": name,
		}).Decode(&artist); err != nil && err.Error() == mongo.ErrNoDocuments.Error() {
			// Artist not found, construct one
			artist = Artist{
				ID:        primitive.NewObjectID(),
				Name:      name,
				CreatedAt: time.Now(),
				Picture:   spot.ArtistPicture(fmt.Sprintf("artist:%s", name)),
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
