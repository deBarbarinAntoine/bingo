package bingo

import (
	"bytes"
	"encoding/json"
	"net/http"
	
	"github.com/rs/zerolog/hlog"
)

type H map[string]any

func Json(r *http.Request, w http.ResponseWriter, data any, status int) {
	
	var byteBuffer bytes.Buffer
	if err := json.NewEncoder(&byteBuffer).Encode(data); err != nil {
		hlog.FromRequest(r).Error().Err(ErrJsonResponseWith(err)).Msg("Failed to encode JSON response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if _, err := w.Write(byteBuffer.Bytes()); err != nil {
		hlog.FromRequest(r).Error().Err(ErrJsonResponseWith(err)).Msg("Failed to write JSON response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func ServerError(r *http.Request, w http.ResponseWriter, err error, msg string) {
	hlog.FromRequest(r).Error().Err(err).Msg(msg)
	Json(r, w, H{
		"error":  msg,
		"status": http.StatusInternalServerError,
	}, http.StatusInternalServerError)
}

func ClientError(r *http.Request, w http.ResponseWriter, status int, err error, msg string) {
	hlog.FromRequest(r).Info().Err(err).Msg(msg)
	Json(r, w, H{
		"error":  msg,
		"status": status,
	}, status)
}
