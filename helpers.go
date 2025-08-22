package bingo

import (
	"bytes"
	"encoding/json"
	"net/http"
	
	"github.com/rs/zerolog/hlog"
)

type H map[string]any

func Json(w http.ResponseWriter, data any, status int) error {
	
	var byteBuffer bytes.Buffer
	if err := json.NewEncoder(&byteBuffer).Encode(data); err != nil {
		return ErrJsonResponseWith(err)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	
	if _, err := w.Write(byteBuffer.Bytes()); err != nil {
		return ErrJsonResponseWith(err)
	}
	
	return nil
}

func ServerError(r *http.Request, w http.ResponseWriter, err error, msg string) error {
	hlog.FromRequest(r).Error().Err(err).Msg(msg)
	return Json(w, H{
		"error":  msg,
		"status": http.StatusInternalServerError,
	}, http.StatusInternalServerError)
}

func ClientError(r *http.Request, w http.ResponseWriter, status int, err error, msg string) error {
	hlog.FromRequest(r).Error().Err(err).Msg(msg)
	return Json(w, H{
		"error":  msg,
		"status": status,
	}, status)
}
