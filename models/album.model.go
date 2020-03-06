package models

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/brocococonut/waveline-server-go/wavelineutils"
	"github.com/davecgh/go-spew/spew"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// Album - Album of songs
	Album struct {
		ID        primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
		Name      string               `json:"name,omitempty" bson:"name,omitempty"`
		Artist    primitive.ObjectID   `json:"-" bson:"artist,omitempty"`
		Artists   []primitive.ObjectID `json:"-" bson:"artists,omitempty"`
		Picture   string               `json:"picture,omitempty" bson:"picture,omitempty"`
		Year      int                  `json:"year,omitempty" bson:"year,omitempty"`
		CreatedAt time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`

		ArtistPop  *Artist   `json:"artist,omitempty" bson:"-"`
		ArtistsPop []*Artist `json:"artists,omitempty" bson:"-"`
	}
	// AlbumSearchData - Search data structure for the findOrCreate function
	AlbumSearchData struct {
		Album   string    `json:"album,omitempty" bson:"album,omitempty"`
		Artist  *Artist   `json:"artist,omitempty" bson:"artist,omitempty"`
		Artists []*Artist `json:"artists,omitempty" bson:"artists,omitempty"`
		Year    int       `json:"year,omitempty" bson:"year,omitempty"`
		Picture []byte    `json:"picture,omitempty" bson:"picture,omitempty"`
	}
)

// FindOrCreate - Find or create an album
func (*Album) FindOrCreate(
	data AlbumSearchData,
	col *mongo.Collection,
	host, artPath string,
	spotifyClient, spotifySecret string,
) (album Album, err error) {
	// Find the artist
	if err = col.FindOne(context.TODO(), bson.M{
		"name": data.Album,
	}).Decode(&album); err != nil {

		if err.Error() != "mongo: no documents in result" {
			return album, err
		}
	}

	// Create a new album if
	if album.Name == "" {
		var picURL string
		// Check to see if a picture was provided
		spew.Dump(data.Picture)
		if len(data.Picture) != 0 {
			// Create a unique file hash and store the file
			fileName := fmt.Sprintf("%s-%s", data.Artist.ID.Hex(), data.Album)
			hash := md5.Sum([]byte(fileName))
			err = ioutil.WriteFile(fmt.Sprintf("%s/%s", artPath, hash), data.Picture, 0604)

			picURL = fmt.Sprintf("%s/albums/art/%s", host, hash)
		} else {
			// Get the url from Spotify
			spot := wavelineutils.Spotify{}

			spot.Authorize(spotifyClient, spotifySecret)
			picURL = spot.AlbumPicture(fmt.Sprintf("album:%s artist:%s", data.Album, data.Artist.Name))
		}

		artistsIds := []primitive.ObjectID{}
		for _, a := range data.Artists {
			artistsIds = append(artistsIds, a.ID)
		}

		album = Album{
			Name:      data.Album,
			Year:      data.Year,
			Artist:    data.Artist.ID,
			Artists:   artistsIds,
			CreatedAt: time.Now(),
			Picture:   picURL,
		}

		_, err = col.InsertOne(context.TODO(), album)
	}

	return album, err
}
