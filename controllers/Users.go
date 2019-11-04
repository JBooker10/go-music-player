package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-music-player/models"
	"github.com/go-music-player/utils"
	"github.com/gorilla/context"
)

type User struct{}

// Home -- Retreives user profile
func (u User) Home(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error utils.Error

		token := context.Get(r, "token")
		claim := token.(*jwt.Token).Claims.(jwt.MapClaims)

		rows := db.QueryRow("select * from users where email=$1", claim["email"])
		err := rows.Scan(&user.ID, &user.Hash, &user.FirstName, &user.LastName, &user.Avatar, &user.DOB, &user.Email, &user.Password, &user.CreatedAt)
		user.ConfirmPassword = ""
		user.Password = ""

		if err != nil {
			fmt.Println(err)
			utils.ErrorResponse(w, fmt.Sprintf("No user found with email of %s", claim["email"]), http.StatusBadRequest, error)
			return
		}

		json.NewEncoder(w).Encode(user)
	}
}
