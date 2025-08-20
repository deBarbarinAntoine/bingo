package main

import (
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"time"
	
	"BinGo/bingo"
	"BinGo/router"
	
	"github.com/rs/zerolog/hlog"
)

type Publication struct {
	Title string    `multipart:"title" json:"title"`
	Score float64   `multipart:"score" json:"score"`
	Date  time.Time `multipart:"date" json:"date"`
}

type Info struct {
	Id           uint          `multipart:"id" json:"id"`
	Publications []Publication `multipart:"publications" json:"publications"`
}

type Data struct {
	Session   string                `cookie:"session" json:"session"`
	CSRFToken string                `header:"x-csrf-token" json:"x-csrf-token"`
	Id        uint                  `param:"id" json:"id"`
	Name      string                `query:"name" json:"name"`
	Age       int                   `multipart:"age" json:"age"`
	IsAdmin   bool                  `multipart:"is_admin" json:"is_admin"`
	Score     float64               `multipart:"score" json:"score"`
	Info      Info                  `multipart:"info" json:"info"`
	File      *multipart.FileHeader `multipart:"file" json:"file"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	data, ok := router.GetCtxData(r.Context(), "data").(*Data)
	if !ok {
		http.Error(w, "no data", http.StatusInternalServerError)
		return
	}
	hlog.FromRequest(r).Info().Msg("Request received successfully")
	
	if data.File != nil {
		hlog.FromRequest(r).Info().
			Int64("file_size", data.File.Size).
			Str("file_name", data.File.Filename).
			Msg("File uploaded successfully")
	} else {
		hlog.FromRequest(r).Error().Msg("No file uploaded")
	}
	
	w.Header().Set("Content-Type", "application/json")
	
	json.NewEncoder(w).Encode(data)
}

//go:generate go run github.com/debarbarinantoine/go-enum-generate
func main() {
	srv := bingo.New(bingo.Options{
		ServerAddr:  "localhost:8008",
		RedisAddr:   "localhost:6379",
		Environment: "development",
	}).
		WithLogMiddleware().
		WithSessionMiddleware()
	
	srv.Mux.WithMultipartFormBindCtx("/:id", handler, &Data{}, "data", http.MethodPost)
	log.Fatal(srv.ListenAndServe())
}
