package routes

import (
	"context"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/brocococonut/waveline-server-go/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// TracksIndex - Get track list
// GET /tracks/
func (r *Router) TracksIndex(c echo.Context) (err error) {
	var (
		skipStr       = c.QueryParam("skip")
		limitStr      = c.QueryParam("limit")
		shuffleStr    = c.QueryParam("shuffle")
		genreStr      = c.QueryParam("genre")
		favouritesStr = c.QueryParam("favourites")
		artistStr     = c.QueryParam("artist")
		albumStr      = c.QueryParam("album")
	)
	var (
		skip    = 0
		limit   = 25
		shuffle = false
		lookup  = map[string]interface{}{}
	)
	if skipStr != "" {
		skip, _ = strconv.Atoi(skipStr)
	}
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}
	if shuffleStr != "" && shuffleStr != "false" {
		shuffle = true
	}
	if genreStr != "" {
		lookup["genre"] = genreStr
	}
	if favouritesStr != "" && favouritesStr != "false" {
		lookup["favourites"] = true
	}
	if artistStr != "" {
		lookup["artist"] = artistStr
	}
	if albumStr != "" {
		lookup["album"] = albumStr
	}

	var trackCur, trackCountCur *mongo.Cursor

	pipe := []bson.M{
		bson.M{"$match": lookup},
		bson.M{"$sort": bson.M{
			"created_at": -1,
		}},
	}

	if trackCountCur, err = r.Client.
		Database(r.Env.DB).
		Collection("artists").
		Aggregate(context.TODO(), append(pipe, bson.M{"$count": "total"})); err != nil {
		return c.JSON(500, err)
	}

	tempRes := []map[string]int{}
	trackCountCur.All(context.TODO(), &tempRes)
	total := 0
	if len(tempRes) > 0 {
		total = tempRes[0]["total"]
	}

	if limit > 0 && !shuffle {
		pipe = append(pipe, bson.M{"$limit": limit})
	}
	if skip > 0 && !shuffle {
		pipe = append(pipe, bson.M{"$skip": skip})
	}

	if trackCur, err = r.Client.
		Database(r.Env.DB).
		Collection("tracks").
		Aggregate(context.TODO(), append(
			pipe,
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
			bson.M{"$unwind": "$genre"})); err != nil {
		return c.JSON(500, err)
	}

	tracks := []map[string]interface{}{}
	trackCur.All(context.TODO(), &tracks)

	if shuffle {
		var (
			min     float64
			n       = []map[string]interface{}{}
			iterMax = 10
		)
		if len(tracks) < iterMax {
			iterMax = len(tracks)
		}

		for i := 0; i < 10; i++ {
			var ind = int(math.Floor(rand.Float64()*(float64(len(tracks))-min+1)) + min)
			if ind >= len(tracks) {
				ind = len(tracks) - 1
			}
			n = append(n, tracks[ind])
		}

		tracks = n
	}

	return c.JSON(200, map[string]interface{}{
		"tracks": tracks,
		"total":  total,
		"query": map[string]interface{}{
			"skip":    skip,
			"limit":   limit,
			"shuffle": shuffle,
			"lookup":  lookup,
		},
	})
}

// TracksPlay - Play a specified track
// GET /tracks/play/:id
func (r *Router) TracksPlay(c echo.Context) (err error) {
	var (
		idStr = c.Param("id")
		id    primitive.ObjectID
	)

	id, err = primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}

	var track models.Track

	if err := r.Client.
		Database(r.Env.DB).
		Collection("tracks").
		FindOne(context.TODO(), bson.M{
			"_id": id,
		}).Decode(&track); err != nil {
		return err
	}

	track.Plays = track.Plays + 1
	track.LastPlay = time.Now()

	r.Client.
		Database(r.Env.DB).
		Collection("tracks").
		FindOneAndUpdate(context.TODO(), bson.M{"_id": id}, track)

	return c.File(track.Path)
}

// TracksLike - Set a track to liked/unliked
// GET /tracks/like/:id
func (r *Router) TracksLike(c echo.Context) (err error) {
	var (
		idStr = c.Param("id")
		id    primitive.ObjectID
	)

	id, err = primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}

	var track models.Track

	if err := r.Client.
		Database(r.Env.DB).
		Collection("tracks").
		FindOne(context.TODO(), bson.M{
			"_id": id,
		}).Decode(&track); err != nil {
		return err
	}

	track.Favourited = !track.Favourited
	track.UpdatedAt = time.Now()

	r.Client.
		Database(r.Env.DB).
		Collection("tracks").
		FindOneAndUpdate(context.TODO(), bson.M{"_id": id}, track)

	return c.JSON(200, track)
}

// TracksFavourites - Get favourited tracks
// /tracks/favourites/
func (r *Router) TracksFavourites(c echo.Context) (err error) {
	pipe := []bson.M{
		bson.M{"$match": bson.M{
			"favourited": true,
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

	var trackCur *mongo.Cursor
	if trackCur, err = r.Client.
		Database(r.Env.DB).
		Collection("tracks").
		Aggregate(context.TODO(), pipe); err != nil {
		return c.JSON(500, err)
	}

	tracks := []map[string]interface{}{}
	trackCur.All(context.TODO(), &tracks)

	return c.JSON(200, tracks)
}

// TracksFavourites - Get new tracks
// /tracks/new/
func (r *Router) TracksNew(c echo.Context) (err error) {
	pipe := []bson.M{
		bson.M{"$group": bson.M{
			"_id":        "$album",
			"track_id":   bson.M{"$first": "$_id"},
			"plays":      bson.M{"$first": "$plays"},
			"favourited": bson.M{"$first": "$favourited"},
			"title":      bson.M{"$first": "$title"},
			"artists":    bson.M{"$first": "$artists"},
			"album":      bson.M{"$first": "$album"},
			"art":        bson.M{"$first": "$art"},
			"duration":   bson.M{"$first": "$duration"},
			"path":       bson.M{"$first": "$path"},
		},
		},
		bson.M{"$lookup": bson.M{
			"from":         "albums",
			"localField":   "album",
			"foreignField": "_id",
			"as":           "album",
		}},
		bson.M{"$lookup": bson.M{
			"from":         "artists",
			"localField":   "artists",
			"foreignField": "_id",
			"as":           "artists",
		}},
		bson.M{"$unwind": bson.M{"path": "$album"}},
		bson.M{"$project": bson.M{
			"_id":        "$track_id",
			"plays":      "$plays",
			"favourited": "$favourited",
			"title":      "$title",
			"artists":    "$artists",
			"album":      "$album",
			"art":        "$art",
			"duration":   "$duration",
			"path":       "$path",
		}},
		bson.M{"$limit": 15},
		bson.M{"$sort": bson.M{"created_at": -1}},
	}

	var trackCur *mongo.Cursor
	if trackCur, err = r.Client.
		Database(r.Env.DB).
		Collection("tracks").
		Aggregate(context.TODO(), pipe); err != nil {
		return c.JSON(500, err)
	}

	tracks := []map[string]interface{}{}
	trackCur.All(context.TODO(), &tracks)

	return c.JSON(200, tracks)
}
