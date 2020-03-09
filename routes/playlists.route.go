package routes

import (
	"context"
	"time"

	"github.com/brocococonut/waveline-server-go/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// PlaylistsIndex - Get your playlists
// GET /playlists/
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

	// Add Favourites playlist to the slice
	playlists = append(playlists, models.PlaylistExtended{
		Name:     "Favourites",
		ID:       "FAVOURITES",
		Tracks:   favCount,
		Pictures: favsImages,
		ReadOnly: true,
	})

	playlistTemps := []models.PlaylistExtended{}

	// Retrieve the other playlists
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
		bson.M{"$lookup": bson.M{
			"from":         "artists",
			"localField":   "tracks.artists",
			"foreignField": "_id",
			"as":           "tracks.artists",
		}},
		bson.M{
			"$project": bson.M{
				"_id":  1,
				"name": 1,
				"pictures": bson.M{
					"$cond": bson.M{
						"if": bson.M{"$isArray": "$tracks"},
						"then": bson.M{
							"$map": bson.M{
								"input": "$tracks",
								"as":    "track",
								"in":    "$$track.album",
							},
						},
						"else": "$tracks.album.picture",
					},
				},
				"tracks": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$isArray": "$tracks"},
						"then": bson.M{"$size": "$tracks"},
						"else": bson.M{
							"$cond": bson.M{
								"if": bson.M{
									"$not": "tracks",
								},
								"then": 0,
								"else": 1,
							},
						},
					},
				},
			},
		},
	}

	var playlistCur *mongo.Cursor
	if playlistCur, err = r.Client.Database(r.Env.DB).Collection("playlists").Aggregate(context.TODO(), playlistPipe); err != nil {
		return c.JSON(500, err)
	}
	if err := playlistCur.All(context.TODO(), &playlistTemps); err != nil {
		return c.JSON(500, err)
	}

	playlists = append(playlists, playlistTemps...)

	return c.JSON(200, playlists)
}

// PlaylistsIndexPost - Create a new playlist
// POST /playlists
func (r *Router) PlaylistsIndexPost(c echo.Context) (err error) {
	plist := models.PlaylistBasic{}

	plist.CreatedAt = time.Now()
	plist.UpdatedAt = time.Now()

	if err := c.Bind(&plist); err != nil {
		return c.JSON(500, err)
	}

	plist.ID = primitive.NewObjectID()

	if _, err := r.Client.
		Database(r.Env.DB).
		Collection("playlists").
		InsertOne(context.TODO(), plist); err != nil {
		return c.JSON(500, err)
	}

	return c.JSON(200, plist)
}

// PlaylistsPlaylistTrack - Add a track to an existing playlist
// GET /playlists/:id/:track
func (r *Router) PlaylistsPlaylistTrack(c echo.Context) (err error) {
	var (
		idStr     = c.Param("id")
		trackStr  = c.Param("track")
		id, track primitive.ObjectID
	)

	if id, err = primitive.ObjectIDFromHex(idStr); err != nil {
		return err
	}
	if track, err = primitive.ObjectIDFromHex(trackStr); err != nil {
		return err
	}

	if _, err = r.Client.
		Database(r.Env.DB).
		Collection("playlists").
		UpdateOne(context.TODO(), bson.M{
			"_id": id,
		}, bson.M{
			"$addToSet": bson.M{
				"tracks": track,
			},
		}); err != nil {
		return err
	}

	return c.JSON(200, bson.M{})
}

// PlaylistsPlaylist - Get a particular playlist
// GET /playlists/:id
func (r *Router) PlaylistsPlaylist(c echo.Context) (err error) {
	var (
		idStr = c.Param("id")
		id    primitive.ObjectID
	)

	if id, err = primitive.ObjectIDFromHex(idStr); err != nil {
		return err
	}

	var plist map[string]interface{}

	// Get the playlist and map it to a basic map
	if err := r.Client.
		Database(r.Env.DB).
		Collection("playlists").
		FindOne(context.TODO(), bson.M{
			"_id": id,
		}).
		Decode(&plist); err != nil {
		return err
	}

	if plist["tracks"] == nil {
		return c.JSON(400, bson.M{})
	}

	// Retrieve and populate the fields of the tracks
	tracksPipe := []bson.M{
		bson.M{"$match": bson.M{
			"_id": bson.M{
				"$in": plist["tracks"],
			},
		}},
		bson.M{"$lookup": bson.M{
			"from":         "albums",
			"localField":   "album",
			"foreignField": "_id",
			"as":           "album",
		}},
		bson.M{"$unwind": "$album"},

		bson.M{"$lookup": bson.M{
			"from":         "artists",
			"localField":   "artists",
			"foreignField": "_id",
			"as":           "artists",
		}},

		bson.M{"$lookup": bson.M{
			"from":         "genres",
			"localField":   "genre",
			"foreignField": "_id",
			"as":           "genre",
		}},
		bson.M{"$unwind": "$genre"},
	}

	var tracks []models.TrackExtended

	var tracksCur *mongo.Cursor
	if tracksCur, err = r.Client.
		Database(r.Env.DB).
		Collection("tracks").
		Aggregate(context.TODO(), tracksPipe); err != nil {
		return c.JSON(500, err)
	}

	if err := tracksCur.All(context.TODO(), &tracks); err != nil {
		return c.JSON(500, err)
	}

	plist["tracks"] = tracks

	return c.JSON(200, plist)
}

// PlaylistsPlaylistPut - Update a playlists name or tracks
// PUT /playlists/:id
func (r *Router) PlaylistsPlaylistPut(c echo.Context) (err error) {
	var (
		idStr = c.Param("id")
		id    primitive.ObjectID
	)

	if id, err = primitive.ObjectIDFromHex(idStr); err != nil {
		return err
	}

	plist := models.PlaylistBasic{}

	plist.UpdatedAt = time.Now()

	if err := c.Bind(&plist); err != nil {
		return c.JSON(500, err)
	}

	plist.ID = id

	if _, err := r.Client.
		Database(r.Env.DB).
		Collection("playlists").
		UpdateOne(context.TODO(), bson.M{
			"_id": id,
		}, plist); err != nil {
		return c.JSON(500, err)
	}

	return c.JSON(200, plist)
}

// PlaylistsPlaylistDelete - Remove a playlist
// DELETE /playlists/:id
func (r *Router) PlaylistsPlaylistDelete(c echo.Context) (err error) {
	var (
		idStr = c.Param("id")
		id    primitive.ObjectID
	)

	if id, err = primitive.ObjectIDFromHex(idStr); err != nil {
		return err
	}

	r.Client.
		Database(r.Env.DB).
		Collection("playlists").
		FindOneAndDelete(context.TODO(), bson.M{
			"_id": id,
		})

	return c.JSON(200, bson.M{})
}
