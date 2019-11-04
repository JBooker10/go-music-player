package controllers

import (
	"database/sql"
	"net/http"
)

type Playlist struct {}

// Linked List Data Structure for Playlist

func (p *Playlist) GetPlaylist(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {

	}
}

func (p *Playlist) AddTrackToPlaylist(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
     
	}
}
