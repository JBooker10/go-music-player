package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-music-player/controllers"
	"github.com/go-music-player/persistence"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rootsongjc/cloudinary-go"
	"github.com/subosito/gotenv"
)

var (
	cdn *cloudinary.Service
	db  *sql.DB
)

func init() {
	gotenv.Load()
}

func main() {
	r := mux.NewRouter()
	db = persistence.PostgreSQLConnection()
	cdn = persistence.CloudinaryConnection()
	port := flag.String("port", "8000", "application is listening on port")

	auth := controllers.Authentication{}
	library := controllers.Library{}
	user := controllers.User{}

	r.HandleFunc("/v1/auth/login", auth.Login(db)).Methods("POST")
	r.HandleFunc("/v1/auth/register", auth.Register(db)).Methods("POST")
	r.HandleFunc("/v1/home", auth.AuthMiddleware(user.Home(db))).Methods("GET")
	r.HandleFunc("/v1/library/upload", library.Upload(db, cdn)).Methods("POST")
	r.HandleFunc("/v1/library/tracks", auth.AuthMiddleware(library.GetTracks(db))).Methods("GET")
	r.HandleFunc("/v1/library/tracks/{id}", auth.AuthMiddleware(library.GetTrack(db))).Methods("GET")
	r.HandleFunc("/v1/library/artists", auth.AuthMiddleware(library.GetArtists(db))).Methods("GET")
	r.HandleFunc("/v1/library/artists/{id}", auth.AuthMiddleware(library.GetArtist(db))).Methods("GET")
	r.HandleFunc("/v1/library/albums", auth.AuthMiddleware(library.GetAlbums(db))).Methods("GET")
	r.HandleFunc("/v1/library/albums/{id}", auth.AuthMiddleware(library.GetAlbum(db))).Methods("GET")

	// Serve Static Audio Files
	r.Handle("/v1/audio/{id}", http.StripPrefix("/v1/audio/", http.FileServer(http.Dir("./audio_files"))))

	srv := &http.Server{
		Handler:      handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r),
		Addr:         "0.0.0.0:" + *port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("Server started on %s\n", *port)
	log.Fatal(srv.ListenAndServe())
}
