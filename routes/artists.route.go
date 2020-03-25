package routes

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// ArtistsNew - Get the 10 newest artists
func (r *Router) ArtistsNew(c echo.Context) (err error) {
	artists := []map[string]interface{}{}

	pipe := []bson.M{
		bson.M{"$match": bson.M{}},
		bson.M{"$limit": 15},
		bson.M{"$sort": bson.M{
			"created_at": -1,
		}},
		artistLookup,
		artistUnwind,
		artistLookup,
	}

	var artistCur *mongo.Cursor
	if artistCur, err = r.Client.Database(r.Env.DB).Collection("artists").Aggregate(context.TODO(), pipe); err != nil {
		return c.JSON(500, err)
	}

	artistCur.All(context.TODO(), &artists)

	return c.JSON(200, artists)
}

// ArtistsIndex - Get a stream of artists
func (r *Router) ArtistsIndex(c echo.Context) (err error) {
	artists := []map[string]interface{}{}

	var (
		limitQuery = c.QueryParam("limit")
		skipQuery  = c.QueryParam("skip")
		limit      = 20
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
		bson.M{"$sort": bson.M{
			"created_at": -1,
		}},
		bson.M{"$skip": skip},
		bson.M{"$limit": limit},
	}

	var artistCur *mongo.Cursor
	if artistCur, err = r.Client.Database(r.Env.DB).Collection("artists").Aggregate(context.TODO(), pipe); err != nil {
		return c.JSON(500, err)
	}

	artistCur.All(context.TODO(), &artists)

	artistCount, _ := r.Client.Database(r.Env.DB).Collection("artists").CountDocuments(context.TODO(), bson.M{})

	type resQuery struct {
		Skip  int `json:"skip"`
		Limit int `json:"limit"`
	}

	type res struct {
		Artists []map[string]interface{} `json:"artists"`
		Total   int64                    `json:"total"`
		Query   resQuery                 `json:"query"`
	}

	return c.JSON(200, res{
		Artists: artists,
		Total:   artistCount,
		Query: resQuery{
			Skip:  skip,
			Limit: limit,
		},
	})
}
