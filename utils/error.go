package utils

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type LoginError struct {
	Code     int    `json:"code"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func ErrorResponse(w http.ResponseWriter, message string, status int, error Error) {
	error.Message = message
	error.Code = status
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
	return
}