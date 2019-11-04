package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-music-player/models"
	"github.com/go-music-player/utils"
	"github.com/gorilla/mux"
)

type AlbumData struct {
	models.Album
	Tracks []models.Track
}

func (l *Library) GetAlbums(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := ""
		var audioData AlbumData
		var albums []AlbumData

		var track models.Track

		queryAlbums := "select * from tracks t LEFT JOIN albums a ON t.album_id = a.id"

		rows, err := db.Query(queryAlbums)
		if err != nil {
			utils.ErrorResponse(w, "Error retreving albums", http.StatusBadRequest, l.Error)
			return
		}

		for rows.Next() {

			err := rows.Scan(&track.ID, &track.Hash, &track.SpotifyID, &track.Country, &id, &track.AlbumTitle, &track.Title, &track.Genre, &track.Duration, &track.Stream, &track.ReleaseDate, &track.Lyrics, &track.Collection, &track.Explicit, &track.Popularity, &audioData.Album.ID, &audioData.Album.Hash, &audioData.Album.SpotifyID, &id, &audioData.Album.Title, &audioData.Album.DiscCount, &audioData.Album.Cover, &audioData.Album.CoverBig, &audioData.Album.CoverXL, &audioData.Album.ReleaseDate, &audioData.Album.TrackCount, &audioData.Album.TrackNumber, &audioData.Album.TrackList, &audioData.Album.RecordLabel, &audioData.Album.Type)

			if err != nil {
				log.Fatal(err)
			}
			audioData.Tracks = append(audioData.Tracks, track)
			albums = append(albums, audioData)
		}

		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(albums)
	}
}

func (l *Library) GetAlbum(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id string
		var audioData AlbumData
		var track models.Track
		params := mux.Vars(r)

		queryAlbums := "select * from tracks t LEFT JOIN albums a ON t.album_id = a.id where a.hash = $1"

		rows := db.QueryRow(queryAlbums, params["id"])

		err := rows.Scan(&track.ID, &track.Hash, &track.SpotifyID, &track.Country, &id, &track.AlbumTitle, &track.Title, &track.Genre, &track.Duration, &track.Stream, &track.ReleaseDate, &track.Lyrics, &track.Collection, &track.Explicit, &track.Popularity, &audioData.Album.ID, &audioData.Album.Hash, &audioData.Album.SpotifyID, &id, &audioData.Album.Title, &audioData.Album.DiscCount, &audioData.Album.Cover, &audioData.Album.CoverBig, &audioData.Album.CoverXL, &audioData.Album.ReleaseDate, &audioData.Album.TrackCount, &audioData.Album.TrackNumber, &audioData.Album.TrackList, &audioData.Album.RecordLabel, &audioData.Album.Type)

		audioData.Tracks = append(audioData.Tracks, track)

		if err != nil {
			fmt.Println(err)
			utils.ErrorResponse(w, fmt.Sprintf("No album found with id of %s", params["id"]), http.StatusBadRequest, l.Error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(audioData)
	}
}
