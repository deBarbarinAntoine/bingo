package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	
	"github.com/rs/zerolog/hlog"
)

// H is a map[string]any type alias for JSON response data.
type H map[string]any

// Json writes JSON response with the given data and status code.
func Json(r *http.Request, w http.ResponseWriter, data any, status int) {
	err := sendJson(w, data, status)
	if err != nil {
		switch {
		case errors.Is(err, ErrJsonEncode):
			hlog.FromRequest(r).Error().Err(err).Msg("Failed to encode JSON response")
		case errors.Is(err, ErrJsonWrite):
			hlog.FromRequest(r).Error().Err(err).Msg("Failed to write JSON response")
		default:
			hlog.FromRequest(r).Error().Err(err).Msg("An error occurred while writing JSON response")
		}
		err = sendJson(w, H{"error": http.StatusText(http.StatusInternalServerError), "status": http.StatusInternalServerError}, http.StatusInternalServerError)
		if err != nil {
			hlog.FromRequest(r).Error().Err(err).Msg("Failed to send JSON error to client")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func sendJson(w http.ResponseWriter, data any, status int) error {
	var byteBuffer bytes.Buffer
	if err := json.NewEncoder(&byteBuffer).Encode(data); err != nil {
		return ErrJsonEncode
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if _, err := w.Write(byteBuffer.Bytes()); err != nil {
		return ErrJsonWrite
	}
	
	return nil
}

// ServerError writes a server error response with the given error message.
func ServerError(r *http.Request, w http.ResponseWriter, err error, msg string) {
	hlog.FromRequest(r).Error().Err(err).Msg(msg)
	Json(r, w, H{
		"error":  msg,
		"status": http.StatusInternalServerError,
	}, http.StatusInternalServerError)
}

// ClientError writes a client error response with the given status code and error message.
func ClientError(r *http.Request, w http.ResponseWriter, status int, err error, msg string) {
	hlog.FromRequest(r).Info().Err(err).Msg(msg)
	Json(r, w, H{
		"error":  msg,
		"status": status,
	}, status)
}

// NotFound writes a 404 Not Found response with a "route %s not found" message.
func NotFound(r *http.Request, w http.ResponseWriter) {
	ClientError(r, w, http.StatusNotFound, fmt.Errorf("route %s not found", r.URL.Path), fmt.Sprintf("route %s not found", r.URL.Path))
}

// MethodNotAllowed writes a 405 Method Not Allowed response with a "method %s not allowed" message.
func MethodNotAllowed(r *http.Request, w http.ResponseWriter) {
	ClientError(r, w, http.StatusMethodNotAllowed, fmt.Errorf("method %s not allowed", r.Method), fmt.Sprintf("method %s not allowed", r.Method))
}
