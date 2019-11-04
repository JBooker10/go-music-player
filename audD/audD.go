package audD

import (
	"time"
	"os"
	// "fmt"
	"github.com/levigross/grequests"
	"encoding/json"
)

type AudioRecognitionAPI struct {
	Status string
	Result struct {
		Artist      string
		Title       string
		Album       string
		ReleaseDate string `json:"release_date"`
		Label       string
		TimeCode    string
		Lyrics      Lyrics
		ITunes      ITunes
		Deezer      Deezer
		Spotify		Spotify
	}
}

type ITunes struct {
	WrapperType            string
	Kind                   string
	ArtistID               int64
	CollectionID           int64
	TrackID                int64
	ArtistName             string
	CollectionName         string
	TrackName              string
	CollectionCensoredName string
	TrackCensoredName      string
	ArtistViewURL          string
	CollectionViewURL      string
	TrackViewURL           string
	PreviewURL             string
	ArtworkURL30           string
	ArtworkURL60           string
	ArtworkURL100          string
	CollectionPrice        float32
	TrackPrice             float32
	ReleaseDate            time.Time
	CollectionExplicitness string
	TrackExplicitness      string
	DiscCount              int
	DiscNumber             int
	TrackCount             int
	TrackNumber            int
	TrackTimeMillis        int64
	Country                string
	Currency               string
	PrimaryGenreName       string
	IsStreamable           bool
}

type Lyrics struct {
	SongID            string `json:"song_id"`
	ArtistID          string `json:"artist_id"`
	Title             string
	TitleWithFeatured string `json:"title_with_featured"`
	FullTitle         string `json:"full_title"`
	Artist            string
	Lyrics            string
}

type Deezer struct {
	ID             int64
	Readable       bool
	Title          string
	TitleShort     string `json:"title_short"`
	TitleVersion   string `json:"title_version"`
	Link           string
	Duration       int
	Rank           int32
	ExplicitLyrics bool `json:"explicit_lyrics"`
	Preview        string
	Artist         Artist
	Album          Album
	Type           string
}

type Spotify struct {
	Album struct {
		AlbumType	string `json:"album_type"`
		Artists		[]struct {
			ExternalURLs struct  {
				Spotify	string
			} `json:"external_urls"`
			Href	string
			ID		string
			Name	string
			Type	string
			URI		string
		}
	AvaliableMarket	string `json:"available_market"`
	ExternalURLs 	struct {
		Spotify 	string
	} `json:"external_urls"`
	Href	string
	ID		string
	Images 	[]struct{
		Height	int
		URL		string
		Width	int
	}
	Name	string
	ReleaseDate	string 	`json:"release_date"`
	ReleaseDatePrecision string `json:"release_data_precision"`
	TotalTracks	int `json:"total_tracks"`
	Type string `json:"type"`
	URI	string `json:"uri"`
	}
	

	AvaliableMarket	string
	DiscNumber		int
	DurationMS		int32
	Explicit		bool
	ExternalIDs		struct {
		IRSC	string
	}
	ExternalURLs  struct {
		Spotify	string
	}
	Href		string
	ID			string
	IsLocal 	bool
	Name		string
	Popularity	int
	TrackNumber	int
	URI  		string

}

type Artist struct {
	ID            int32
	Name          string
	Link          string
	Picture       string
	PictureSmall  string `json:"picture_small"`
	PictureMedium string `json:"picture_medium"`
	PictureBig    string `json:"picture_big"`
	PictureXL     string `json:"picture_xl"`
	TrackList     string
	Type          string
}

type Album struct {
	ID         int32
	Title      string
	Cover      string
	CoverSmall string `json:"cover_small"`
	CoverBig   string `json:"cover_big"`
	CoverXL    string `json:"cover_xl"`
	TrackList  string
	Type       string
}


func AudRecognition(audioURL string) (AudioRecognitionAPI, error) {
	var audioData AudioRecognitionAPI
	pathname := "https://audd.p.rapidapi.com/"
	options := &grequests.RequestOptions{
		Headers: map[string]string{"X-RapidAPI-Key": os.Getenv("AUDD_KEY")},
	}
	resp, err := grequests.Get(pathname+"?return=timecode%2Citunes%2Clyrics%2Cdeezer%2Cspotify&itunes_country=us&url="+audioURL, options)
	json.NewDecoder(resp.RawResponse.Body).Decode(&audioData)
	return audioData, err
}
