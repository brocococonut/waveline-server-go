package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// Router structure
	Router struct {
		Client *mongo.Client
		Env    struct {
			Host          string
			Port          string
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
		r.Env.Host = "127.0.0.1"
	}

	if providedPort := os.Getenv("PORT"); providedPort != "" {
		r.Env.Port = providedPort
	} else {
		r.Env.Port = ":1323"
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

	r.Env.MusicPath, _ = filepath.Abs(r.Env.MusicPath)
	r.Env.AlbumArtPath, _ = filepath.Abs(r.Env.AlbumArtPath)

	os.MkdirAll(r.Env.MusicPath, 0704)
	os.MkdirAll(r.Env.AlbumArtPath, 0704)

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

	if providedStr := os.Getenv("SPOTIFY_ID"); providedStr != "" {
		r.Env.SpotifyClient = providedStr
	} else {
		log.Fatal("Missing \"SPOTIFY_ID\" env variable")
	}
	if providedStr := os.Getenv("SPOTIFY_SECRET"); providedStr != "" {
		r.Env.SpotifySecret = providedStr
	} else {
		log.Fatal("Missing \"SPOTIFY_SECRET\" env variable")
	}

	stringified, _ := json.Marshal(r.Env)

	println(string(stringified))
}
