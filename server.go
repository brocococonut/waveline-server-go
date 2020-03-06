package main

import (
	"context"
	"fmt"
	"log"

	"github.com/brocococonut/waveline-server-go/routes"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	e := echo.New()

	r := routes.Router{}
	r.InitEnv()

	// Set client options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s/%s", r.Env.DBHost, r.Env.DBPort, r.Env.DB))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	r.Client = client

	e.GET("/albums/new", r.NewAlbum)
	e.Logger.Fatal(e.Start(":1323"))
}
