package models

import (
	"time"
)

type Track struct {
	ID			int		 `json:"-"`
	SpotifyID	string
	Hash		string
	Country		string
	AlbumTitle	string  
	Title 		string
	Genre		string
	Size		int64
	Duration	int64
	Stream		string
	ReleaseDate	time.Time
	Collection	string
	Explicit	bool
	Lyrics		string
	Popularity	int
}

