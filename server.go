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

	e.GET("/albums/new/", r.AlbumsNew)
	e.GET("/albums/new", r.AlbumsNew)
	e.GET("/albums/random/", r.AlbumsRandom)
	e.GET("/albums/random", r.AlbumsRandom)
	e.GET("/albums/artists/:id", r.AlbumsArtists)
	e.GET("/albums/art/:id", r.AlbumsArt)
	e.GET("/albums/:id", r.AlbumsAlbum)
	e.GET("/albums/", r.AlbumsIndex)
	e.GET("/albums", r.AlbumsIndex)

	e.GET("/artists/new/", r.ArtistsNew)
	e.GET("/artists/new", r.ArtistsNew)
	e.GET("/artists/", r.ArtistsIndex)
	e.GET("/artists", r.ArtistsIndex)

	e.GET("/genres/", r.GenresIndex)
	e.GET("/genres", r.GenresIndex)

	e.GET("/playlists/", r.PlaylistsIndex)
	e.GET("/playlists", r.PlaylistsIndex)

	e.GET("/system/sync", r.SystemSync)
	e.GET("/system/info", r.SystemInfo)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s%s", r.Env.Host, r.Env.Port)))
}
