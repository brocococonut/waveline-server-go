package routes

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// AlbumsNew - Get the 10 newest albums
func (r *Router) AlbumsNew(c echo.Context) (err error) {
	albums := []map[string]interface{}{}

	pipe := []bson.M{
		bson.M{"$match": bson.M{}},
		bson.M{"$limit": 10},
		bson.M{"$sort": bson.M{
			"created_at": -1,
		}},
		artistLookup,
		artistUnwind,
		artistLookup,
	}

	var albumCur *mongo.Cursor
	if albumCur, err = r.Client.Database(r.Env.DB).Collection("albums").Aggregate(context.TODO(), pipe); err != nil {
		return c.JSON(500, err)
	}

	albumCur.All(context.TODO(), &albums)

	return c.JSON(200, albums)
}

// AlbumsArt - Find artwork from a file hash
func (r *Router) AlbumsArt(c echo.Context) (err error) {
	idStr := c.Param("id")

	idStr = filepath.Clean(idStr)

	path := fmt.Sprintf("%s/%s", r.Env.AlbumArtPath, idStr)

	return c.File(path)
}

// AlbumsIndex - Get a stream of albums
func (r *Router) AlbumsIndex(c echo.Context) (err error) {
	albums := []map[string]interface{}{}

	var (
		limitQuery = c.QueryParam("limit")
		skipQuery  = c.QueryParam("skip")
		limit      = 25
		skip       = 0
	)

	if limitQuery == "" {
		limit = 25
	} else {
		if limit, err = strconv.Atoi(limitQuery); err != nil {
			c.Logger().Error(err)
			return err
		}
	}
	if skipQuery == "" {
		skip = 0
	} else {
		if skip, err = strconv.Atoi(skipQuery); err != nil {
			c.Error(err)
		}
	}

	if limit > 200 || limit < 1 {
		limit = 25
	}
	if skip < 0 {
		skip = 0
	}

	pipe := []bson.M{
		bson.M{"$match": bson.M{}},
		bson.M{"$skip": skip},
		bson.M{"$limit": limit},
		bson.M{"$sort": bson.M{
			"created_at": -1,
		}},
		artistLookup,
		artistUnwind,
		artistLookup,
	}

	var albumCur *mongo.Cursor
	if albumCur, err = r.Client.Database(r.Env.DB).Collection("albums").Aggregate(context.TODO(), pipe); err != nil {
		return c.JSON(500, err)
	}

	albumCur.All(context.TODO(), &albums)

	return c.JSON(200, albums)
}

// AlbumsArtists - Get albums by a particular artist
func (r *Router) AlbumsArtists(c echo.Context) (err error) {
	var (
		albums = []map[string]interface{}{}
		idStr  = c.Param("id")
		// id     primitive.ObjectID
	)

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}

	pipe := []bson.M{
		bson.M{"$match": bson.M{
			"artist": id,
		}},
		artistLookup,
		artistUnwind,
		artistLookup,
	}

	var albumCur *mongo.Cursor
	if albumCur, err = r.Client.Database(r.Env.DB).Collection("albums").Aggregate(context.TODO(), pipe); err != nil {
		return c.JSON(500, err)
	}

	albumCur.All(context.TODO(), &albums)

	return c.JSON(200, albums)
}

// AlbumsRandom - Get 10 random albums
func (r *Router) AlbumsRandom(c echo.Context) (err error) {
	// albums := []map[string]interface{}{}
	type tempAlbum struct {
		Artists   interface{} `json:"artists" bson:"artists"`
		Year      interface{} `json:"year" bson:"year"`
		CreatedAt interface{} `json:"created_at" bson:"created_at"`
		ID        interface{} `json:"_id" bson:"_id"`
		Name      interface{} `json:"name" bson:"name"`
		Artist    interface{} `json:"artist" bson:"artist"`
	}
	albums := []tempAlbum{}

	var (
		albumCur *mongo.Cursor
		min      float64
		n        = []tempAlbum{}
		iterMax  = 10
	)

	pipe := []bson.M{
		artistLookup,
		artistUnwind,
		artistLookup,
	}

	if albumCur, err = r.Client.Database(r.Env.DB).Collection("albums").Aggregate(context.TODO(), pipe); err != nil {
		return c.JSON(500, err)
	}

	err = albumCur.All(context.TODO(), &albums)

	if len(albums) < iterMax {
		iterMax = len(albums)
	}

	for i := 0; i < 10; i++ {
		var ind = int(math.Floor(rand.Float64()*(float64(len(albums))-min+1)) + min)
		if ind >= len(albums) {
			ind = len(albums) - 1
		}
		n = append(n, albums[ind])
	}

	return c.JSON(200, n)
}

// AlbumsAlbum - Get a specific album
func (r *Router) AlbumsAlbum(c echo.Context) (err error) {
	var (
		albums = []map[string]interface{}{}
		idStr  = c.Param("id")
	)

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return err
	}

	pipe := []bson.M{
		bson.M{"$match": bson.M{
			"_id": id,
		}},
		artistLookup,
		artistUnwind,
		artistLookup,
	}

	var albumCur *mongo.Cursor
	if albumCur, err = r.Client.Database(r.Env.DB).Collection("albums").Aggregate(context.TODO(), pipe); err != nil {
		return c.JSON(500, err)
	}

	albumCur.All(context.TODO(), &albums)

	if len(albums) == 0 {
		return c.JSON(404, map[string]interface{}{})
	}

	return c.JSON(200, albums[0])
}
