package bingo

import (
	"net/http"
	
	"github.com/debarbarinantoine/bingo/internal/helpers"
)

type H map[string]any

func Json(r *http.Request, w http.ResponseWriter, data any, status int) {
	helpers.Json(r, w, data, status)
}

func ServerError(r *http.Request, w http.ResponseWriter, err error, msg string) {
	helpers.ServerError(r, w, err, msg)
}

func ClientError(r *http.Request, w http.ResponseWriter, status int, err error, msg string) {
	helpers.ClientError(r, w, status, err, msg)
}
