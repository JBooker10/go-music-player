package controllers

import (
	"net/http"
	"database/sql"
	"bytes"
	"errors"
	"os"
	"mime"
	"fmt"
	"log"
	"io/ioutil"
	"github.com/go-music-player/utils"
	"github.com/go-music-player/audD"
	"github.com/rootsongjc/cloudinary-go"

)

// MimeType - Standard of classify files on Internet
type MimeType string

const (
	mimeTypeAudioMPEG  MimeType = "audio/mpeg"
	mimeTypeVideoMP4   MimeType = "video/mp4"
	maxUploadSize	   int64 	= 10 * 1024 * 1024
	errFileSize                 = "File size exceeds 10 MB"
	errInvalidFile              = "invalid file"
	errCannotReadFileType       = "cannot read file type"
	audioStoragePath            = "audio_files/"
	uploadFile					= "uploadFile"
	localStorage		bool    = false
)

// Upload -- Stores the audio locally or via cloudinary
func (l *Library) Upload(db *sql.DB, storage *cloudinary.Service) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {

		r.Body = http.MaxBytesReader(w, r.Body,maxUploadSize)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			utils.ErrorResponse(w, errFileSize, http.StatusBadRequest, l.Error)
			return
		}

		file, _, err := r.FormFile(uploadFile)
		if err != nil {
			utils.ErrorResponse(w, errInvalidFile, http.StatusBadRequest, l.Error)
			return
		}

		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			utils.ErrorResponse(w, errInvalidFile, http.StatusBadRequest, l.Error)
			return
		}

		filetype := http.DetectContentType(fileBytes)
		if MimeType(filetype) != mimeTypeAudioMPEG && MimeType(filetype) != mimeTypeVideoMP4 {
		    utils.ErrorResponse(w, errInvalidFile, http.StatusBadRequest, l.Error)
			return
		}
		fmt.Println(filetype)

		extension, err := mime.ExtensionsByType(filetype)

		// filename, _ := utils.NewUUID()
		// pathName := audioStoragePath + filename + extension[0]
		// aud, err := uploadAudToFS(fileBytes, filename, pathName, r.Host + filename + extension[0])

		aud, url, err := uploadAudToCloudinary(fileBytes, extension[0], storage)
		
		if err != nil {
			fmt.Println(err)
			utils.ErrorResponse(w, errInvalidFile, http.StatusBadRequest, l.Error)
			return
		}
	

			// Track
			// l.Track.Size = audiofile.Size()
			l.Track.SpotifyID = aud.Result.Spotify.ID
			l.Track.Title = aud.Result.Title
			l.Track.Duration = aud.Result.ITunes.TrackTimeMillis
			l.Track.Genre = aud.Result.ITunes.PrimaryGenreName
			l.Track.ReleaseDate = aud.Result.ITunes.ReleaseDate
			l.Track.Stream = url
			l.Track.Collection = aud.Result.ITunes.CollectionName
			l.Track.Country = aud.Result.ITunes.Country
			l.Track.AlbumTitle = aud.Result.Album
			l.Track.Explicit = aud.Result.Deezer.ExplicitLyrics
			l.Track.Lyrics = aud.Result.Lyrics.Lyrics
			l.Track.Popularity = aud.Result.Spotify.Popularity
	
			// Artist
			l.Artist.SpotifyID = aud.Result.Spotify.Album.Artists[0].ID
			l.Artist.Name = aud.Result.Artist
			l.Artist.Label = aud.Result.Label
			l.Artist.Picture = aud.Result.Deezer.Artist.Picture
			l.Artist.PictureXL = aud.Result.Deezer.Artist.PictureXL

			// Album
			l.Album.SpotifyID = aud.Result.Spotify.Album.ID
			l.Album.Title = aud.Result.ITunes.CollectionName
			l.Album.DiscCount = aud.Result.ITunes.DiscCount
			l.Album.ReleaseDate = aud.Result.ITunes.ReleaseDate
			l.Album.TrackCount = aud.Result.ITunes.TrackCount
			l.Album.TrackNumber = aud.Result.ITunes.TrackNumber
			l.Album.TrackList = aud.Result.Deezer.Album.TrackList
			l.Album.RecordLabel = aud.Result.Label
			l.Album.Cover = aud.Result.Deezer.Album.Cover
			l.Album.CoverBig = aud.Result.Deezer.Album.CoverBig
			l.Album.CoverXL = aud.Result.Deezer.Album.CoverXL
			l.Album.Type = aud.Result.Deezer.Album.Type

			queryArtist := "insert into artists (spotify_id, name, label, picture, picture_xl) values ($1, $2, $3, $4, $5) returning id"
			artistID := 0
			err = db.QueryRow(queryArtist, l.Artist.SpotifyID,l.Artist.Name, l.Artist.Label, l.Artist.Picture, l.Artist.PictureXL).Scan(&artistID)
			if err != nil {
				fmt.Println("artist")
				log.Fatal(err)
			}

			queryAlbum := "insert into albums (title, artist_id,spotify_id, disc_count, release_date, track_count, track_number, track_list, record_label, cover, cover_big, cover_xl, type) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) returning id"
			albumID := 0
			err = db.QueryRow(queryAlbum, l.Album.Title, artistID, l.Album.SpotifyID, l.Album.DiscCount, l.Album.ReleaseDate, l.Album.TrackCount, l.Album.TrackNumber, l.Album.TrackList, l.Album.RecordLabel, l.Album.Cover, l.Album.CoverBig, l.Album.CoverXL, l.Album.Type).Scan(&albumID)
			if err != nil {
				fmt.Println("albums")
				log.Fatal(err)
			}

			queryTrack := "insert into tracks (spotify_id, country, album_id, album_title, title, genre, duration, stream, release_date, lyrics, collection, explicit, popularity ) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) returning id"

			trackID := 0
			err = db.QueryRow(queryTrack, l.Track.SpotifyID, l.Track.Country,
			albumID,
			l.Track.AlbumTitle,
			l.Track.Title, l.Track.Genre, l.Track.Duration, l.Track.Stream, l.Track.ReleaseDate, l.Track.Lyrics,l.Track.Collection, l.Track.Explicit, l.Track.Popularity).Scan(&trackID)
			if err != nil {
				fmt.Println("tracks")
				log.Fatal(err)
			}
				
			fmt.Println("Success")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("SUCCESS"))
	}
}

func createAudFile(path string) {
	_, err  := os.Stat(path)
	if err != nil {
		fmt.Println(err)
	}

	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
	}
}

func uploadAudToFS(file []byte, filename string, fs string, url string ) (audD.AudioRecognitionAPI, error ) {
	createAudFile(fs)
	err := ioutil.WriteFile(fs, file, 0644)
	if err != nil {
		 var aud audD.AudioRecognitionAPI
		 return aud, errors.New("Could not write file to disk")	
	}
	
	aud, err := audD.AudRecognition(url)
	if aud.Status == "error" {
			return aud, errors.New("cannot recognize audio file")
	}
	if err != nil {
		return  aud, errors.New("Could not upload file")
	}
	return aud, err 
}

func uploadAudToCloudinary(file []byte, ext string, storage *cloudinary.Service) (audD.AudioRecognitionAPI, string, error) {
	audiofile := bytes.NewReader(file)
	audio, err := storage.Upload("audio", audiofile, ext, true, 2)
	if err != nil {
		fmt.Println(err)
	}

	audFile := storage.Url(audio, 2)
	aud, err := audD.AudRecognition(audFile)
	if aud.Status == "error" {
		return aud, "", errors.New("cannot recognize audio file")
	}
	if err != nil {
		return  aud, "", errors.New("Could not upload file")
	}
		return aud, audFile, err
}



