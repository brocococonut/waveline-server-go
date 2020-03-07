package routes

import (
	"context"

	"github.com/brocococonut/waveline-server-go/models"
	"github.com/brocococonut/waveline-server-go/wavelibrary"
	"github.com/labstack/echo/v4"
	"gopkg.in/mgo.v2/bson"
)

// SystemSync - Sync filesystem
func (r *Router) SystemSync(c echo.Context) (err error) {
	l := wavelibrary.Library{
		DB:  r.Client.Database(r.Env.DB),
		Env: r.Env,
	}

	go l.Sync(r.Env.AlbumArtPath, r.Env.MusicPath, []string{".mp3", ".flac", ".m4a"})
	// info, err := l.Sync(r.Env.AlbumArtPath, r.Env.MusicPath, []string{".mp3", ".flac", ".m4a"})

	return c.JSON(200, models.Info{
		Tracks:  0,
		Albums:  0,
		Artists: 0,
		Size:    0,
		Mount:   "",
	})
}

// SystemInfo - Waveline info
func (r *Router) SystemInfo(c echo.Context) (err error) {
	info := map[string]interface{}{}
	infoRes := r.Client.Database(r.Env.DB).Collection("info").FindOne(context.TODO(), bson.M{})
	infoRes.Decode(&info)

	return c.JSON(200, info)
}
