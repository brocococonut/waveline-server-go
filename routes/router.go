package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// Router structure
	Router struct {
		Client *mongo.Client
		Env    struct {
			Host          string
			DB            string
			DBHost        string
			DBPort        string
			AlbumArtPath  string
			MusicPath     string
			SpotifyClient string
			SpotifySecret string
		}
	}
)

// InitEnv - Initialise the Env based variabls on a Router struct
func (r *Router) InitEnv() {
	cwd, _ := os.Getwd()
	if providedHost := os.Getenv("HOST"); providedHost != "" {
		r.Env.Host = providedHost
	} else {
		r.Env.Host = "http://127.0.0.1"
	}

	if providedPath := os.Getenv("ART_PATH"); providedPath != "" {
		r.Env.AlbumArtPath = providedPath
	} else {
		r.Env.AlbumArtPath = fmt.Sprintf("%s/albumart", cwd)
	}

	if providedPath := os.Getenv("MUSIC_PATH"); providedPath != "" {
		r.Env.MusicPath = providedPath
	} else {
		r.Env.MusicPath = fmt.Sprintf("%s/music", cwd)
	}

	if providedStr := os.Getenv("DB"); providedStr != "" {
		r.Env.DB = providedStr
	} else {
		r.Env.DB = "waveline"
	}
	if providedStr := os.Getenv("DB_HOST"); providedStr != "" {
		r.Env.DBHost = providedStr
	} else {
		r.Env.DBHost = "127.0.0.1"
	}
	if providedStr := os.Getenv("DB_PORT"); providedStr != "" {
		r.Env.DBPort = providedStr
	} else {
		r.Env.DBPort = "27017"
	}

	if providedStr := os.Getenv("SPOTIFY_CLIENT"); providedStr != "" {
		r.Env.SpotifyClient = providedStr
	} else {
		log.Fatal("Missing \"SPOTIFY_CLIENT\" env variable")
	}
	if providedStr := os.Getenv("SPOTIFY_SECRET"); providedStr != "" {
		r.Env.SpotifySecret = providedStr
	} else {
		log.Fatal("Missing \"SPOTIFY_SECRET\" env variable")
	}

	stringified, _ := json.Marshal(r.Env)

	println(string(stringified))
}
