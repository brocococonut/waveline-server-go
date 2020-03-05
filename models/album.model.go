package model

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/brocococonut/go-waveline-server/wavelineutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// Album - Album of songs
	Album struct {
		ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Name      string             `json:"name,omitempty" bson:"name,omitempty"`
		Artist    *Artist            `json:"artist,omitempty" bson:"artist,omitempty"`
		Artists   []*Artist          `json:"artists,omitempty" bson:"artists,omitempty"`
		Picture   string             `json:"picture,omitempty" bson:"picture,omitempty"`
		Year      int                `json:"year,omitempty" bson:"year,omitempty"`
		CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
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

func (*Album) findOrCreate(db *mongo.Database, data AlbumSearchData) (album Album, err error) {
	// Find the artist
	if err = db.Collection("albums").FindOne(context.TODO(), bson.M{
		"name": data.Album,
	}).Decode(&album); err != nil {
		return album, err
	}

	// Create a new album if
	if album.Name == "" {
		var artPath, _ = os.Getwd()
		var picURL string

		if providedPath := os.Getenv("ART_PATH"); providedPath != "" {
			artPath = providedPath
		}

		if len(data.Picture) == 0 {
			fileName := fmt.Sprintf("%s-%s", data.Artist.ID.Hex(), data.Album)
			hash := md5.Sum([]byte(fileName))
			err = ioutil.WriteFile(fmt.Sprintf("%s/%s", artPath, hash), data.Picture, 0604)

			picURL = fmt.Sprintf("/albums/art/%s", hash)
		} else {
			spot := wavelineutils.Spotify{}

			spot.Authorize()
			picURL = spot.AlbumPicture(fmt.Sprintf("album:%s artist:%s", data.Album, data.Artist.Name))
		}

		album = Album{
			Name:      data.Album,
			Year:      data.Year,
			Artist:    data.Artist,
			Artists:   data.Artists,
			CreatedAt: time.Now(),
			Picture:   picURL,
		}

		_, err = db.Collection("albums").InsertOne(context.TODO(), album)
	}

	return album, err
}
