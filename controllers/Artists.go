package controllers

import (
	"log"
	"fmt"
	"encoding/json"
	"net/http"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/go-music-player/models"
	"github.com/go-music-player/utils"
)

type ArtistData struct {
	models.Artist
	Albums  []models.Album
}

func (l *Library) GetArtists(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		id := ""
		var audioData ArtistData
		var artists []ArtistData

		var album models.Album

		queryArtists := "select * from albums a INNER JOIN artists r ON a.artist_id = r.id;"

		rows, err := db.Query(queryArtists)
		if err != nil {
			utils.ErrorResponse(w, "Error getting data", http.StatusBadRequest, l.Error)
			return	
		}

		for rows.Next() {
		err := rows.Scan(&album.ID, &album.Hash, &album.SpotifyID, &id, &album.Title, &album.DiscCount, &album.Cover, &album.CoverBig, &album.CoverXL, &album.ReleaseDate, &album.TrackCount, &album.TrackNumber, &album.TrackList, &album.RecordLabel, &album.Type,  &audioData.Artist.ID, &audioData.Artist.Hash, &audioData.Artist.SpotifyID, &audioData.Artist.Label, &audioData.Artist.Name, &audioData.Artist.Picture, &audioData.Artist.PictureXL)
		if err != nil {
			log.Fatal(err)
		}
		audioData.Albums = append(audioData.Albums, album)
		artists = append(artists, audioData)
		}

		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
	}
}

func (l *Library) GetArtist(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var id string
		var audioData ArtistData
		var album models.Album
		params := mux.Vars(r)

		queryArtists := "select * from albums a INNER JOIN artists r ON a.artist_id = r.id where r.hash = $1"

		rows := db.QueryRow(queryArtists, params["id"])

		err := rows.Scan(&album.ID, &album.Hash, &album.SpotifyID, &id, &album.Title, &album.DiscCount, &album.Cover, &album.CoverBig, &album.CoverXL, &album.ReleaseDate, &album.TrackCount, &album.TrackNumber, &album.TrackList, &album.RecordLabel, &album.Type,  &audioData.Artist.ID, &audioData.Artist.Hash, &audioData.Artist.SpotifyID, &audioData.Artist.Label, &audioData.Artist.Name, &audioData.Artist.Picture, &audioData.Artist.PictureXL)

		audioData.Albums = append(audioData.Albums, album)

		if err != nil {
			fmt.Println(err)
			utils.ErrorResponse(w, fmt.Sprintf("No artist found with id of %s", params["id"]), http.StatusBadRequest, l.Error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(audioData)
	}
}