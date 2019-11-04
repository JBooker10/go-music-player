package models

import (
	"time"
)


type Album struct {
	ID			string		`json:"-"`
	SpotifyID	string		
	Title		string
	Hash		string
	DiscCount	int			
	Cover 		string
	CoverBig	string		
	CoverXL		string	
	ReleaseDate	time.Time	
	TrackCount  int			
	TrackNumber	int			
	TrackList	string		
	RecordLabel	string		
	Type 		string
}