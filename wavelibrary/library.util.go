package wavelibrary

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/brocococonut/waveline-server-go/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/dhowden/tag"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Library struct for additional functions
type Library struct {
	DB  *mongo.Database
	Env struct {
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

type fileType struct {
	Info os.FileInfo
	Path string
}

func (l *Library) ExtractMeta(files []fileType) (tracks []models.Track) {
	for _, file := range files {
		// var genre models.Genre

		// Check to see if file exists in tracks collection
		if count, err := l.DB.Collection("tracks").CountDocuments(context.TODO(), bson.M{"path": file.Path}); count > 0 || err != nil {
			if err != nil {
				spew.Dump(err)
			}

			continue
		}

		// Open file for meta reading
		fileCur, err := os.Open(file.Path)
		if err != nil {
			spew.Dump(err)
			continue
		}

		// Attempt to get Metadata from file
		meta, err := tag.ReadFrom(fileCur)
		if err != nil && err.Error() != "no tags found" {
			spew.Dump(err)
			continue
		}

		if meta != nil {
			// Create artist docs
			artist := models.Artist{}
			artists := []models.Artist{}
			if meta.Artist() != "" {
				if artists, err = artist.FindOrCreate(
					l.DB,
					strings.Split(meta.Artist(), "&,"),
					l.Env.SpotifyClient,
					l.Env.SpotifySecret,
				); err != nil {
					spew.Dump(err)
				}
			}

			// Start creating album search/creation data
			albumData := models.AlbumSearchData{
				Album:   meta.Album(),
				Year:    meta.Year(),
				Artists: artists,
			}

			// Add artists if possible
			if len(artists) > 0 {
				albumData.Artist = models.Artist{
					Name: meta.Artist(),
					ID:   artists[0].ID,
				}
			}

			// Add picture if it exists
			pic := meta.Picture()
			if pic != nil && len(pic.Data) > 0 {
				albumData.Picture = pic.Data
			}

			// Find/create album
			album := models.Album{}
			if album, err = album.FindOrCreate(albumData, l.DB.Collection("albums"), l.Env.Host, l.Env.AlbumArtPath, l.Env.SpotifyClient, l.Env.SpotifySecret); err != nil {
				spew.Dump(err)
			}

			// Find or create genre
			var genre models.Genre
			if meta.Genre() != "" {
				if genreCount, err := l.DB.Collection("genres").CountDocuments(context.TODO(), bson.M{
					"name": meta.Genre(),
				}); genreCount == 0 || err != nil {
					if err != nil {
						spew.Dump(err)
					}

					genre = models.Genre{
						ID:   primitive.NewObjectID(),
						Name: meta.Genre(),
					}

					if _, err := l.DB.Collection("genres").InsertOne(context.TODO(), genre); err != nil {
						spew.Dump(err)
					}
				}
			}

			artistIds := []primitive.ObjectID{}

			for _, artist := range artists {
				artistIds = append(artistIds, artist.ID)
			}

			track := models.Track{
				Name:      meta.Title(),
				Artists:   artistIds,
				Album:     album.ID,
				Artist:    meta.Artist(),
				Genre:     genre.ID,
				Path:      file.Path,
				Year:      meta.Year(),
				CreatedAt: time.Now(),
			}

			if _, err := l.DB.Collection("tracks").InsertOne(context.TODO(), track); err != nil {
				spew.Dump(err)
				continue
			}

			tracks = append(tracks, track)
		} else {
			continue
		}
	}

	return tracks
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func filterFiles(arr []fileType, cond func(fileType) bool) []fileType {
	result := []fileType{}
	for i := range arr {
		if cond(arr[i]) {
			result = append(result, arr[i])
		}
	}
	return result
}

func (l *Library) Download(url, file string) {

}

func (l *Library) Sync(artPath, path string, ext []string) (info models.Info, err error) {
	println("Starting new sync, this may take a while")

	// Get absolute paths
	path, err = filepath.Abs(path)
	artPath, err = filepath.Abs(path)
	if err != nil {
		return info, err
	}

	// Make directories if necessary
	if err := os.MkdirAll(path, 0704); err != nil {
		return info, err
	}
	if err := os.MkdirAll(artPath, 0704); err != nil {
		return info, err
	}

	// files, err := ioutil.ReadDir(path)
	files := []fileType{}

	// Get all files recursively
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			println(err)
			return nil
		}
		if info.IsDir() == true {
			return nil
		}
		files = append(files, fileType{
			Info: info,
			Path: path,
		})

		return nil
	})

	var byteCount int64 = 0

	// Filter out files that don't match the extension types
	files = filterFiles(files, func(file fileType) bool {
		if file.Info.IsDir() {
			return false
		}
		fileExt := filepath.Ext(file.Info.Name())

		if contains(ext, fileExt) {
			byteCount = byteCount + file.Info.Size()
			return true
		}
		return false
	})

	start := time.Now()
	l.ExtractMeta(files)
	end := time.Now()
	// spew.Dump(tracks)

	info.ID = primitive.NewObjectID()
	info.Start = start
	info.End = end
	info.Seconds = int(time.Since(start).Seconds())

	if tracksEst, err := l.DB.Collection("tracks").EstimatedDocumentCount(context.TODO()); err != nil {
		info.Tracks = 0
	} else {
		info.Tracks = int(tracksEst)
	}
	if albumsEst, err := l.DB.Collection("albums").EstimatedDocumentCount(context.TODO()); err != nil {
		info.Albums = 0
	} else {
		info.Albums = int(albumsEst)
	}
	if artistsEst, err := l.DB.Collection("artists").EstimatedDocumentCount(context.TODO()); err != nil {
		info.Artists = 0
	} else {
		info.Artists = int(artistsEst)
	}
	info.Size = int(byteCount)
	info.Mount = path
	info.LastScan = time.Now()

	l.DB.Collection("info").FindOneAndDelete(context.TODO(), bson.M{})
	l.DB.Collection("info").InsertOne(context.TODO(), info)

	println(" - Sync finished. Results:")
	spew.Dump(info)

	return info, err
}
