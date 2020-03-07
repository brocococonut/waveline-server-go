package routes

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// GenresIndex - Get all the genres
func (r *Router) GenresIndex(c echo.Context) (err error) {
	genres := []map[string]interface{}{}

	pipe := []bson.M{
		bson.M{"$match": bson.M{}},
	}

	var genreCur *mongo.Cursor
	if genreCur, err = r.Client.Database(r.Env.DB).Collection("artists").Aggregate(context.TODO(), pipe); err != nil {
		return c.JSON(500, err)
	}

	genreCur.All(context.TODO(), &genres)

	return c.JSON(200, genres)
}
