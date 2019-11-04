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


type Library struct {
	Track	models.Track
	Album	models.Album
	Artist	models.Artist
	Error	utils.Error
}

type TrackData struct {
	models.Track
	Album	models.Album
	Artist  models.Artist
}

func (l *Library) GetTracks(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {		
		id := ""
		var audioData TrackData
		var tracks []TrackData

		queryTracks := "select * from tracks t INNER JOIN albums a ON t.album_id = a.id INNER JOIN artists r ON a.artist_id = r.id"

		rows, err := db.Query(queryTracks)
		if err != nil {
			utils.ErrorResponse(w, "Error getting data", http.StatusBadRequest, l.Error)
			return	
		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&audioData.Track.ID,  &audioData.Track.Hash, &audioData.Track.SpotifyID, &audioData.Track.Country, &id, &audioData.Track.AlbumTitle, &audioData.Track.Title, &audioData.Track.Genre, &audioData.Track.Duration, &audioData.Track.Stream, &audioData.Track.ReleaseDate, &audioData.Track.Lyrics, &audioData.Track.Collection, &audioData.Track.Explicit,  &audioData.Track.Popularity,  &audioData.Album.ID, &audioData.Album.Hash, &audioData.Album.SpotifyID, &id, &audioData.Album.Title, &audioData.Album.DiscCount, &audioData.Album.Cover, &audioData.Album.CoverBig, &audioData.Album.CoverXL, &audioData.Album.ReleaseDate, &audioData.Album.TrackCount, &audioData.Album.TrackNumber, &audioData.Album.TrackList, &audioData.Album.RecordLabel, &audioData.Album.Type,  &audioData.Artist.ID, &audioData.Artist.Hash, &audioData.Artist.SpotifyID, &audioData.Artist.Label, &audioData.Artist.Name, &audioData.Artist.Picture, &audioData.Artist.PictureXL)

			if err != nil {
				log.Fatal(err)
			}

			tracks = append(tracks, audioData)
		}

		err = rows.Err()
		if err != nil {
		log.Fatal(err)
		}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tracks)
	}
}

func (l *Library) GetTrack(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var id string
		var audioData TrackData
		params := mux.Vars(r)

		queryTracks := "select * from tracks t INNER JOIN albums a ON t.album_id = a.id INNER JOIN artists r ON a.artist_id = r.id where t.hash = $1"

		rows :=  db.QueryRow(queryTracks, params["id"])

		err := rows.Scan(&audioData.Track.ID,  &audioData.Track.Hash, &audioData.Track.SpotifyID, &audioData.Track.Country, &id, &audioData.Track.AlbumTitle, &audioData.Track.Title, &audioData.Track.Genre, &audioData.Track.Duration, &audioData.Track.Stream, &audioData.Track.ReleaseDate, &audioData.Track.Lyrics, &audioData.Track.Collection, &audioData.Track.Explicit,  &audioData.Track.Popularity,  &audioData.Album.ID, &audioData.Album.Hash, &audioData.Album.SpotifyID, &id, &audioData.Album.Title, &audioData.Album.DiscCount, &audioData.Album.Cover, &audioData.Album.CoverBig, &audioData.Album.CoverXL, &audioData.Album.ReleaseDate, &audioData.Album.TrackCount, &audioData.Album.TrackNumber, &audioData.Album.TrackList, &audioData.Album.RecordLabel, &audioData.Album.Type,  &audioData.Artist.ID, &audioData.Artist.Hash, &audioData.Artist.SpotifyID, &audioData.Artist.Label, &audioData.Artist.Name, &audioData.Artist.Picture,&audioData.Artist.PictureXL)
	

	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(w, fmt.Sprintf("No track found with id of %s", params["id"]), http.StatusBadRequest, l.Error)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(audioData)

}
}


