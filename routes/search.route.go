package routes

import (
	"context"

	"github.com/brocococonut/waveline-server-go/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// SearchIndex - Search
func (r *Router) SearchIndex(c echo.Context) (err error) {
	type queryType struct {
		q string
	}
	type searchResults struct {
		Albums  []models.Album         `json:"albums"`
		Artists []models.Artist        `json:"artists"`
		Tracks  []models.TrackExtended `json:"tracks"`
	}

	var (
		query = c.QueryParam("q")
		// results = searchResults{}
	)

	trackPipe := []bson.M{
		bson.M{"$match": bson.M{
			"$or": []bson.M{
				bson.M{
					"name": bson.M{
						"$regex":   query,
						"$options": "i",
					},
				},
				bson.M{
					"artist": bson.M{
						"$regex":   query,
						"$options": "i",
					},
				},
			},
		}},
		bson.M{"$lookup": bson.M{
			"from":         "albums",
			"localField":   "album",
			"foreignField": "_id",
			"as":           "album",
		}},
		bson.M{"$unwind": "$album"},

		bson.M{"$lookup": bson.M{
			"from":         "artists",
			"localField":   "artists",
			"foreignField": "_id",
			"as":           "artists",
		}},

		bson.M{"$lookup": bson.M{
			"from":         "genres",
			"localField":   "genre",
			"foreignField": "_id",
			"as":           "genre",
		}},
		bson.M{"$unwind": "$genre"},
	}

	albumPipe := []bson.M{
		bson.M{"$match": bson.M{
			"name": bson.M{
				"$regex":   query,
				"$options": "i",
			},
		}},

		bson.M{"$lookup": bson.M{
			"from":         "artists",
			"localField":   "artist",
			"foreignField": "_id",
			"as":           "artist",
		}},
		bson.M{"$unwind": "$artist"},
	}

	var trackCur, albumCur, artistCur *mongo.Cursor

	if trackCur, err = r.Client.
		Database(r.Env.DB).
		Collection("tracks").
		Aggregate(context.TODO(), trackPipe); err != nil {
		return c.JSON(500, err)
	}
	if albumCur, err = r.Client.
		Database(r.Env.DB).
		Collection("albums").
		Aggregate(context.TODO(), albumPipe); err != nil {
		return c.JSON(500, err)
	}

	if artistCur, err = r.Client.Database(r.Env.DB).Collection("artists").Find(context.TODO(), bson.M{
		"name": bson.M{
			"$regex": query,
			// "$options": "i",
		},
	}); err != nil {
		return c.JSON(500, err)
	}

	var (
		tracks  = []models.TrackExtended{}
		albums  = []models.Album{}
		artists = []models.Artist{}
	)

	trackCur.All(context.TODO(), &tracks)
	albumCur.All(context.TODO(), &albums)
	artistCur.All(context.TODO(), &artists)

	results := searchResults{
		Tracks:  tracks,
		Albums:  albums,
		Artists: artists,
	}

	return c.JSON(200, results)
}
