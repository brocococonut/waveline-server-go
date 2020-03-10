package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/brocococonut/waveline-server-go/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if authEnabled := os.Getenv("AUTH_ENABLED"); authEnabled == "true" {
			currentPath := c.Path()

			if contains := any([]string{"art", "play"}, func(pathSnip string) bool {
				strings.Contains(currentPath, pathSnip)
				return false
			}); contains == true {
				return next(c)
			}

			var reqKey string
			keyAlt := c.QueryParam("x-api-key")
			key := c.QueryParam("key")

			if key == "" && keyAlt == "" {
				return c.String(403, "API key missing")
			}
			if key == "" {
				reqKey = keyAlt
			} else {
				reqKey = key
			}

			if apiKey := os.Getenv("API_KEY"); apiKey != reqKey {
				return c.String(403, "Invalid API key")
			}

			return next(c)
		}

		return next(c)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	e := echo.New()

	r := routes.Router{}
	r.InitEnv()

	// Set client options
	clientOptions := options.Client().ApplyURI(r.Env.DBString)

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

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	e.Use(auth)

	// Routes
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
	e.POST("/playlists/", r.PlaylistsIndexPost)
	e.POST("/playlists", r.PlaylistsIndexPost)
	e.GET("/playlists/:id", r.PlaylistsPlaylist)
	e.GET("/playlists/:id/:track", r.PlaylistsPlaylistTrack)

	e.GET("/search/", r.SearchIndex)
	e.GET("/search", r.SearchIndex)

	e.GET("/system/sync", r.SystemSync)
	e.GET("/system/info", r.SystemInfo)

	e.GET("/tracks/", r.TracksIndex)
	e.GET("/tracks", r.TracksIndex)
	e.GET("/tracks/play/:id", r.TracksPlay)
	e.GET("/tracks/like/:id", r.TracksLike)
	e.GET("/tracks/favourites/", r.TracksFavourites)
	e.GET("/tracks/favourites", r.TracksFavourites)
	e.GET("/tracks/new/", r.TracksFavourites)
	e.GET("/tracks/new", r.TracksFavourites)

	// Start generating host for qr code
	host, err := url.Parse(r.Env.Host)
	if err != nil {
		e.Logger.Fatal("Failed to parse HOST string")
	}

	// Add API key to the host
	q := host.Query()
	q.Add("key", os.Getenv("API_KEY"))
	host.RawQuery = q.Encode()

	// Generate the QR code and display it in the terminal
	qr, err := qrcode.New(host.String(), qrcode.High)
	if err != nil {
		e.Logger.Fatal("Failed to generate qr code")
	}
	println(qr.ToString(false))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", r.Env.Port)))
}
