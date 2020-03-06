package routes

import (
	"github.com/brocococonut/waveline-server-go/models"
	"github.com/labstack/echo/v4"
)

// NewAlbum - Create a new album
func (r *Router) NewAlbum(c echo.Context) (err error) {

	album := models.Album{}
	artist := models.Artist{
		Name: "test",
	}

	album, err = album.FindOrCreate(
		models.AlbumSearchData{
			Artist: &artist,
			Album:  "Test",
		},
		r.Client.Database("waveline").Collection("albums"),
		r.Env.Host,
		r.Env.AlbumArtPath,
		r.Env.SpotifyClient,
		r.Env.SpotifySecret,
	)

	// healthy := true

	// startPingTime := ms()
	// if err = h.Client.Ping(context.TODO(), readpref.Primary()); err != nil {
	// 	res.AddErrorT(0, err, true)
	// 	healthy = false
	// } else {
	// 	res.AddMessageT(startPingTime, "DB Ping successful", true, false)
	// }

	// if healthy == false {
	// 	res.SetStatus(500)
	// }

	return c.JSON(200, struct {
		album models.Album
		err   error
	}{
		album: album,
		err:   err,
	})
}
