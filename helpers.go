package bingo

import (
	"net/http"
	
	"github.com/debarbarinantoine/bingo/internal/helpers"
)

// H is a map[string]any type alias for JSON response data.
type H map[string]any

// Json writes JSON response with the given data and status code.
func Json(r *http.Request, w http.ResponseWriter, data any, status int) {
	helpers.Json(r, w, data, status)
}

// ServerError writes a server error response with the given error message.
func ServerError(r *http.Request, w http.ResponseWriter, err error, msg string) {
	helpers.ServerError(r, w, err, msg)
}

// ClientError writes a client error response with the given status code and error message.
func ClientError(r *http.Request, w http.ResponseWriter, status int, err error, msg string) {
	helpers.ClientError(r, w, status, err, msg)
}

// NotFound writes a 404 Not Found response with a "route %s not found" message.
func NotFound(r *http.Request, w http.ResponseWriter) {
	helpers.NotFound(r, w)
}

// MethodNotAllowed writes a 405 Method Not Allowed response with a "method %s not allowed" message.
func MethodNotAllowed(r *http.Request, w http.ResponseWriter) {
	helpers.MethodNotAllowed(r, w)
}
