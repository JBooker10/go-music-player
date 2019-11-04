package controllers

import (	
	"fmt"
	"log"
	"os"

	"strings"
	"encoding/json"
	"database/sql"
	"net/http"
	"github.com/go-music-player/models"
	"github.com/go-music-player/utils"
	"github.com/gorilla/context"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	REDIRECT_URL =  "http://localhost:8000"
	AUTHORIZE_URL = "https://github.com/login/oauth/authorize"
    TOKEN_URL     = "https://github.com/login/oauth/access_token"
)

type Authentication struct {
	User		models.User
	JWT			models.JWT
	LoginError	utils.LoginError
	Error 		utils.Error
}


func (a *Authentication) Login(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&a.User)
		
		if a.User.Email == "" {
		  a.LoginError.Email = "Email is required"
		  a.LoginError.Code = http.StatusUnauthorized
		  w.WriteHeader(http.StatusUnauthorized)
		  json.NewEncoder(w).Encode(a.LoginError)
		}

		if a.User.Password == "" {
		  a.LoginError.Password = "Password is required"
		  a.LoginError.Code = http.StatusUnauthorized
		  w.WriteHeader(http.StatusUnauthorized)
		  json.NewEncoder(w).Encode(a.LoginError)
		}

		password := a.User.Password

		row :=  db.QueryRow("select * from users where email=$1", a.User.Email)
		err := row.Scan(&a.User.ID, &a.User.Hash, &a.User.FirstName, &a.User.LastName, &a.User.Avatar, &a.User.DOB, &a.User.Email, &a.User.Password, &a.User.CreatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				a.LoginError.Email = "User does not exist"
				a.LoginError.Code = http.StatusUnauthorized
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(a.LoginError)
				return 
			} 
				fmt.Println(err)
			
		}

		hashedPassword := a.User.Password

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			a.LoginError.Password = "Invalid Password"
			a.LoginError.Code = http.StatusUnauthorized
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(a.LoginError)
			return
		}

		token, err := utils.GenerateUserToken(a.User)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		a.JWT.Token = token
		json.NewEncoder(w).Encode(a.JWT)

	}
} 

func (a *Authentication) Register(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&a.User)

		if a.User.FirstName == "" {
			utils.ErrorResponse(w, "First name is required", http.StatusBadRequest, a.Error)
			return
		}

		if a.User.LastName == "" {
			utils.ErrorResponse(w, "Last name is required", http.StatusBadRequest, a.Error)
			return
		}

		if a.User.Email == "" {
			utils.ErrorResponse(w, "Email is required", http.StatusBadRequest, a.Error)
			return
		}

		if a.User.Password == "" {
			utils.ErrorResponse(w, "Password is required", http.StatusBadRequest, a.Error)
			return
		}

		if a.User.ConfirmPassword == "" {
			utils.ErrorResponse(w, "Please Confirm Password", http.StatusBadRequest, a.Error)
			return
		}

		if a.User.Password != a.User.ConfirmPassword {
			utils.ErrorResponse(w, "Password does not match", http.StatusBadRequest, a.Error)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(a.User.Password), 10)
		if err != nil {
			log.Fatal(err)
		}

		a.User.Password = string(hash)
		a.User.ConfirmPassword = ""

		queryParam := "insert into users (first_name, last_name, avatar, dob, email, password, created_at) values($1, $2, $3, $4, $5, $6, $7) returning id"
		err = db.QueryRow(queryParam, a.User.FirstName, a.User.LastName, a.User.Avatar, a.User.DOB, a.User.Email, a.User.Password, a.User.CreatedAt).Scan(&a.User.ID)

		if err != nil {
			fmt.Println(err)
			utils.ErrorResponse(w, "Error registration failed", http.StatusInternalServerError, a.Error)
			return
		}

		token, err := utils.GenerateUserToken(a.User)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		a.JWT.Token = token
		json.NewEncoder(w).Encode(a.JWT)
	}
}

func (a *Authentication) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return []byte(os.Getenv("SECRET")), nil
			})

			if err != nil {
				utils.ErrorResponse(w, fmt.Sprintf("Error: %s", err), http.StatusUnauthorized, a.Error)
				return
			}

			if token.Valid {
				context.Set(r, "token", token)
				next.ServeHTTP(w, r)
			} else {
				utils.ErrorResponse(w, fmt.Sprintf("Error: %s", a.Error), http.StatusUnauthorized, a.Error)
				return
			}
		} else {
			utils.ErrorResponse(w, "Invalid token", http.StatusUnauthorized, a.Error)
			return
		}
	}
}

// GithubLogin not completed 
// func (a *Authentication) GithubLogin(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		oauthConfig := &oauth2.Config{
// 			ClientID: os.Getenv("GITHUB_CLIENT_ID"),
// 			ClientSecret: os.Getenv("GITHUB_SECRET"),
// 			Endpoint: oauth2.Endpoint{
// 				AuthURL: AUTHORIZE_URL,
// 				TokenURL: TOKEN_URL,
// 			},
// 			RedirectURL: REDIRECT_URL,
// 			Scopes: []string{"user:email"},
// 		}

// 		fmt.Println(r.URL.Query())
	
// 		token, err := oauthConfig.Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
// 		if err != nil {
// 			fmt.Println(err)
// 			utils.ErrorResponse(w, "Error Retreiving token", http.StatusBadRequest, a.Error)
// 			return
// 		}

// 		if !token.Valid() {
// 			utils.ErrorResponse(w, "Invalid Token", http.StatusBadRequest, a.Error)
// 			return
// 		}

// 		client := github.NewClient(oauthConfig.Client(oauth2.NoContext, token))
// 		user, _, err := client.Users.Get(context.Background(), "")
// 		if err != nil {
// 			utils.ErrorResponse(w, "Could not connect to client", http.StatusBadRequest, a.Error)
// 			return
// 		}

// 		fmt.Printf("Name: %s\n", *user.Name)
// 		http.Redirect(w, r, REDIRECT_URL, http.StatusPermanentRedirect)
// 	}
// }


