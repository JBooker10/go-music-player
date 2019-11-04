package models

import (
	"time"
)

type User struct {
	ID				int			`json:"-"`
	Hash			string
	FirstName		string
	LastName		string
	Avatar			string
	DOB				string
	Email			string
	Password		string
	ConfirmPassword string
	CreatedAt		time.Time
}

type JWT struct {
	Token 	string	`json:"token"`
}