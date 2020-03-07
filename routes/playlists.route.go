package routes

import (
	"context"

	"github.com/brocococonut/waveline-server-go/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func (r *Router) PlaylistsIndex(c echo.Context) (err error) {
	playlists := []models.PlaylistExtended{}

	// Get user favourited tracks for a different playlist
	favs := []models.TrackExtended{}
	pipe := []bson.M{
		bson.M{"$match": bson.M{"favourited": true}},
		bson.M{"$limit": 4},
		bson.M{"$sort": bson.M{"created_at": -1}},
		bson.M{"$lookup": bson.M{
			"from":         "albums",
			"localField":   "album",
			"foreignField": "_id",
			"as":           "album",
		}},
		bson.M{"$unwind": "$album"},
	}

	// Run the query
	var favImageCur *mongo.Cursor
	if favImageCur, err = r.Client.Database(r.Env.DB).Collection("tracks").Aggregate(context.TODO(), pipe); err != nil {
		return c.JSON(500, err)
	}
	favImageCur.All(context.TODO(), &favs)

	// Get the album images for the tracks
	favsImages := []string{}
	for _, track := range favs {
		if track.Album.Picture != "" {
			favsImages = append(favsImages, track.Album.Picture)
		}
	}

	// Get the number of favourited tracks
	favCount, err := r.Client.Database(r.Env.DB).Collection("tracks").CountDocuments(context.TODO(), bson.M{
		"favourited": true,
	})

	playlists = append(playlists, models.PlaylistExtended{
		Name:     "Favourites",
		ID:       "FAVOURITES",
		Tracks:   favCount,
		Pictures: favsImages,
		ReadOnly: true,
	})

	playlistTemps := []models.PlaylistExtended{}
	// Get the other playlists
	playlistPipe := []bson.M{
		bson.M{"$match": bson.M{}},
		bson.M{"$lookup": bson.M{
			"from":         "tracks",
			"localField":   "tracks",
			"foreignField": "_id",
			"as":           "tracks",
		}},
		bson.M{"$lookup": bson.M{
			"from":         "albums",
			"localField":   "tracks.album",
			"foreignField": "_id",
			"as":           "tracks.album",
		}},
		bson.M{"$unwind": "$album"},
		bson.M{"$lookup": bson.M{
			"from":         "artists",
			"localField":   "tracks.artists",
			"foreignField": "_id",
			"as":           "tracks.artists",
		}},
	}

	var playlistCur *mongo.Cursor
	if playlistCur, err = r.Client.Database(r.Env.DB).Collection("playlists").Aggregate(context.TODO(), playlistPipe); err != nil {
		return c.JSON(500, err)
	}
	playlistCur.All(context.TODO(), &playlistTemps)

	return nil
}

// func (r *Router) PlaylistsIndexPost(c echo.Context) (err error) {}

// func (r *Router) PlaylistsPlaylistTrack(c echo.Context) (err error) {}

// func (r *Router) PlaylistsPlaylist(c echo.Context) (err error) {}
// func (r *Router) PlaylistsPlaylistPut(c echo.Context) (err error) {}
// func (r *Router) PlaylistsPlaylistDelete(c echo.Context) (err error) {}
// AlbumsNew - Get the 10 newest albums
// func (r *Router) AlbumsNew(c echo.Context) (err error) {
// 	albums := []map[string]interface{}{}

// 	pipe := []bson.M{
// 		bson.M{"$match": bson.M{}},
// 		bson.M{"$limit": 10},
// 		bson.M{"$sort": bson.M{
// 			"created_at": -1,
// 		}},
// 		bson.M{"$lookup": bson.M{
// 			"from":         "artists",
// 			"localField":   "artist",
// 			"foreignField": "_id",
// 			"as":           "artist",
// 		}},
// 		bson.M{"$unwind": "$artist"},
// 		bson.M{"$lookup": bson.M{
// 			"from":         "artists",
// 			"localField":   "artists",
// 			"foreignField": "_id",
// 			"as":           "artists",
// 		}},
// 		// bson.M{"$unwind": "$popArtists"},
// 	}

// 	var albumCur *mongo.Cursor
// 	if albumCur, err = r.Client.Database(r.Env.DB).Collection("albums").Aggregate(context.TODO(), pipe); err != nil {
// 		return c.JSON(500, err)
// 	}

// 	albumCur.All(context.TODO(), &albums)

// 	return c.JSON(200, albums)
// }

// // AlbumsArt - Find artwork from a file hash
// func (r *Router) AlbumsArt(c echo.Context) (err error) {
// 	idStr := c.Param("id")

// 	idStr = filepath.Clean(idStr)

// 	path := fmt.Sprintf("%s/%s", r.Env.AlbumArtPath, idStr)

// 	return c.File(path)
// }

// // AlbumsIndex - Get a stream of albums
// func (r *Router) AlbumsIndex(c echo.Context) (err error) {
// 	albums := []map[string]interface{}{}

// 	var (
// 		limitQuery = c.QueryParam("limit")
// 		skipQuery  = c.QueryParam("skip")
// 		limit      = 25
// 		skip       = 0
// 	)

// 	if limitQuery == "" {
// 		limit = 25
// 	} else {
// 		if limit, err = strconv.Atoi(limitQuery); err != nil {
// 			c.Logger().Error(err)
// 			return err
// 		}
// 	}
// 	if skipQuery == "" {
// 		skip = 0
// 	} else {
// 		if skip, err = strconv.Atoi(skipQuery); err != nil {
// 			c.Error(err)
// 		}
// 	}

// 	if limit > 200 || limit < 1 {
// 		limit = 25
// 	}
// 	if skip < 0 {
// 		skip = 0
// 	}

// 	pipe := []bson.M{
// 		bson.M{"$match": bson.M{}},
// 		bson.M{"$skip": skip},
// 		bson.M{"$limit": limit},
// 		bson.M{"$sort": bson.M{
// 			"created_at": -1,
// 		}},
// 		bson.M{"$lookup": bson.M{
// 			"from":         "artists",
// 			"localField":   "artist",
// 			"foreignField": "_id",
// 			"as":           "artist",
// 		}},
// 		bson.M{"$unwind": "$artist"},
// 		bson.M{"$lookup": bson.M{
// 			"from":         "artists",
// 			"localField":   "artists",
// 			"foreignField": "_id",
// 			"as":           "artists",
// 		}},
// 		// bson.M{"$unwind": "$popArtists"},
// 	}

// 	var albumCur *mongo.Cursor
// 	if albumCur, err = r.Client.Database(r.Env.DB).Collection("albums").Aggregate(context.TODO(), pipe); err != nil {
// 		return c.JSON(500, err)
// 	}

// 	albumCur.All(context.TODO(), &albums)

// 	return c.JSON(200, albums)
// }

// // AlbumsArtists - Get albums by a particular artist
// func (r *Router) AlbumsArtists(c echo.Context) (err error) {
// 	var (
// 		albums = []map[string]interface{}{}
// 		idStr  = c.Param("id")
// 		// id     primitive.ObjectID
// 	)

// 	id, err := primitive.ObjectIDFromHex(idStr)
// 	if err != nil {
// 		return err
// 	}

// 	pipe := []bson.M{
// 		bson.M{"$match": bson.M{
// 			"artist": id,
// 		}},
// 		bson.M{"$lookup": bson.M{
// 			"from":         "artists",
// 			"localField":   "artist",
// 			"foreignField": "_id",
// 			"as":           "artist",
// 		}},
// 		bson.M{"$unwind": "$artist"},
// 		bson.M{"$lookup": bson.M{
// 			"from":         "artists",
// 			"localField":   "artists",
// 			"foreignField": "_id",
// 			"as":           "artists",
// 		}},
// 	}

// 	var albumCur *mongo.Cursor
// 	if albumCur, err = r.Client.Database(r.Env.DB).Collection("albums").Aggregate(context.TODO(), pipe); err != nil {
// 		return c.JSON(500, err)
// 	}

// 	albumCur.All(context.TODO(), &albums)

// 	return c.JSON(200, albums)
// }

// // AlbumsRandom - Get 10 random albums
// func (r *Router) AlbumsRandom(c echo.Context) (err error) {
// 	// albums := []map[string]interface{}{}
// 	type tempAlbum struct {
// 		Artists    interface{} `json:"artists" bson:"artists"`
// 		Year       interface{} `json:"year" bson:"year"`
// 		CreatedAt  interface{} `json:"created_at" bson:"created_at"`
// 		ID         interface{} `json:"_id" bson:"_id"`
// 		Name       interface{} `json:"name" bson:"name"`
// 		Artist     interface{} `json:"artist" bson:"artist"`
// 		ArtistsPop interface{} `json:"-" bson:"artistsPop"`
// 		ArtistPop  interface{} `json:"-" bson:"artistPop"`
// 	}
// 	albums := []tempAlbum{}

// 	var (
// 		albumCur *mongo.Cursor
// 		min      float64
// 		n        = []tempAlbum{}
// 		iterMax  = 10
// 	)

// 	pipe := []bson.M{
// 		bson.M{"$lookup": bson.M{
// 			"from":         "artists",
// 			"localField":   "artist",
// 			"foreignField": "_id",
// 			"as":           "artist",
// 		}},
// 		bson.M{"$unwind": "$artist"},
// 		bson.M{"$lookup": bson.M{
// 			"from":         "artists",
// 			"localField":   "artists",
// 			"foreignField": "_id",
// 			"as":           "artists",
// 		}},
// 	}

// 	if albumCur, err = r.Client.Database(r.Env.DB).Collection("albums").Aggregate(context.TODO(), pipe); err != nil {
// 		return c.JSON(500, err)
// 	}

// 	err = albumCur.All(context.TODO(), &albums)

// 	if len(albums) < iterMax {
// 		iterMax = len(albums)
// 	}

// 	for i := 0; i < 10; i++ {
// 		var ind = int(math.Floor(rand.Float64()*(float64(len(albums))-min+1)) + min)
// 		if ind >= len(albums) {
// 			ind = len(albums) - 1
// 		}
// 		n = append(n, albums[ind])
// 	}

// 	return c.JSON(200, n)
// }

// // AlbumsAlbum - Get a specific album
// func (r *Router) AlbumsAlbum(c echo.Context) (err error) {
// 	var (
// 		albums = []map[string]interface{}{}
// 		idStr  = c.Param("id")
// 	)

// 	id, err := primitive.ObjectIDFromHex(idStr)
// 	if err != nil {
// 		return err
// 	}

// 	pipe := []bson.M{
// 		bson.M{"$match": bson.M{
// 			"_id": id,
// 		}},
// 		bson.M{"$lookup": bson.M{
// 			"from":         "artists",
// 			"localField":   "artist",
// 			"foreignField": "_id",
// 			"as":           "artist",
// 		}},
// 		bson.M{"$unwind": "$artist"},
// 		bson.M{"$lookup": bson.M{
// 			"from":         "artists",
// 			"localField":   "artists",
// 			"foreignField": "_id",
// 			"as":           "artists",
// 		}},
// 	}

// 	var albumCur *mongo.Cursor
// 	if albumCur, err = r.Client.Database(r.Env.DB).Collection("albums").Aggregate(context.TODO(), pipe); err != nil {
// 		return c.JSON(500, err)
// 	}

// 	albumCur.All(context.TODO(), &albums)

// 	if len(albums) == 0 {
// 		return c.JSON(404, map[string]interface{}{})
// 	}

// 	return c.JSON(200, albums[0])
// }
