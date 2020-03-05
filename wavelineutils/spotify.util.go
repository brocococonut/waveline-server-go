package wavelineutils

import (
	"context"
	"log"
	"os"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/zmb3/spotify"
)

// Spotify - A spotify wrapper wrapper..
type Spotify struct {
	initialised bool
	client      spotify.Client
}

// Authorize - Get an auth token from the spotify API
func (s Spotify) Authorize() {
	config := &clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_ID"),
		ClientSecret: os.Getenv("SPOTIFY_SECRET"),
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	s.client = spotify.Authenticator{}.NewClient(token)
	s.initialised = true
}

// AlbumPicture - Retrieve the Album picture URL for a particular album query
func (s Spotify) AlbumPicture(query string) string {
	if s.initialised == false {
		s.Authorize()
	}

	res, err := s.client.Search(query, spotify.SearchTypeArtist)
	if err != nil {
		log.Print(err)
	}

	if len(res.Albums.Albums) > 0 {
		if len(res.Albums.Albums[0].Images) > 0 {
			return res.Albums.Albums[0].Images[0].URL
		}
	}

	return ""
}

// ArtistPicture - Retrieve the Artist picture URL for a particular album query
func (s Spotify) ArtistPicture(query string) string {
	if s.initialised == false {
		s.Authorize()
	}

	res, err := s.client.Search(query, spotify.SearchTypeArtist)
	if err != nil {
		log.Print(err)
	}

	if len(res.Artists.Artists) > 0 {
		if len(res.Artists.Artists[0].Images) > 0 {
			return res.Artists.Artists[0].Images[0].URL
		}
	}

	return ""
}
