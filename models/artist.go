package models

type Artist struct {
	ID        string `json:"-"`
	SpotifyID string
	Hash      string
	Name      string
	Label     string
	Picture   string
	PictureXL string
}
