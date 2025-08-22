package main

import (
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"time"
	
	"github.com/debarbarinantoine/bingo"
	"github.com/debarbarinantoine/bingo/context"
	"github.com/debarbarinantoine/bingo/jwtkit"
	"github.com/debarbarinantoine/bingo/middleware"
	
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
	data, ok := context.GetCtxData(r.Context(), "data").(*Data)
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

func admin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	json.NewEncoder(w).Encode(map[string]string{
		"message": "hello admin!",
	})
}

//go:generate go run github.com/debarbarinantoine/go-enum-generate@latest --force
func main() {
	srv, err := bingo.New(bingo.Options{
		ServerAddr:  "localhost:8008",
		Environment: "development",
	}).
		WithLogMiddleware().
		WithJWT(jwtkit.NewConfigWithSecret(jwtkit.AlgorithmHS256, "|JwT53cr3T|", jwtkit.DefaultTokenResponseOptions()))
	if err != nil {
		panic(err)
	}
	
	srv.Router.Use(
		middleware.RealIP(),
		middleware.CleanPath(),
		middleware.RedirectSlashes(),
		middleware.Timeout(time.Minute),
		middleware.ThrottleBacklog(100, 500, time.Minute),
		middleware.RateLimiterByIP(30, time.Minute),
	)
	
	// public endpoint
	srv.Router.WithMultipartFormBindCtx("/:id", handler, &Data{}, "data", http.MethodPost)
	
	// restricted area
	srv.Router.Use(jwtkit.VerifyAndAuthenticateJWT())
	srv.Router.HandleFunc("/admin", admin, http.MethodGet)
	
	log.Fatal(srv.ListenAndServe())
}
