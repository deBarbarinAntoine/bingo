package main

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"time"
	
	"github.com/debarbarinantoine/bingo"
	"github.com/debarbarinantoine/bingo/jwtkit"
	"github.com/debarbarinantoine/bingo/middleware"
	
	"github.com/rs/zerolog/hlog"
)

// Publication represents a nested struct in the multipart data.
type Publication struct {
	Title string    `multipart:"title" json:"title"`
	Score float64   `multipart:"score" json:"score"`
	Date  time.Time `multipart:"date" json:"date"`
}

// Info represents another nested struct, containing a slice of Publication.
type Info struct {
	Id           uint          `multipart:"id" json:"id"`
	Publications []Publication `multipart:"publications" json:"publications"`
}

// Data shows how to bind various data types from different request sources.
// The `multipart` tag is used for form data.
// The `cookie`, `header`, `param`, and `query` tags are for other request sources.
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

// handler processes the bound data and sends a JSON response.
func handler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the bound data from the request context.
	data, ok := bingo.GetCtxData(r.Context(), "data").(*Data)
	if !ok {
		bingo.ServerError(r, w, fmt.Errorf("no data found in context"), "no data")
		return
	}
	hlog.FromRequest(r).Info().Msg("Request received successfully")
	
	// Check if a file was uploaded and log its details.
	if data.File != nil {
		hlog.FromRequest(r).Info().
			Int64("file_size", data.File.Size).
			Str("file_name", data.File.Filename).
			Msg("File uploaded successfully")
	} else {
		hlog.FromRequest(r).Error().Msg("No file uploaded")
	}
	
	// Respond to the client with the bound data in JSON format.
	bingo.Json(r, w, data, http.StatusOK)
}

// ping is a handler for a common route without data.
func ping(w http.ResponseWriter, r *http.Request) {
	bingo.Json(r, w, bingo.H{"message": "pong"}, http.StatusOK)
}

// admin is a handler for a restricted route.
func admin(w http.ResponseWriter, r *http.Request) {
	bingo.Json(r, w, bingo.H{"message": "hello admin!"}, http.StatusOK)
}

func main() {
	// Initialize JWT configuration with a secret key.
	jwtConfig, err := jwtkit.NewConfigWithSecret(jwtkit.AlgorithmHS256, "|JwT53cr3T|", jwtkit.DefaultTokenResponseOptions())
	if err != nil {
		panic(err)
	}
	
	// Initialize the server and chain builder methods to add middleware.
	// The "With..." methods use a builder pattern to configure the server.
	srv, err := bingo.New(bingo.Options{
		ServerAddr:  "localhost:8008",
		Environment: "development",
	}).
		WithLogMiddleware().
		WithJWT(jwtConfig)
	if err != nil {
		panic(err)
	}
	
	// Use global middlewares.
	srv.Router.Use(
		middleware.RealIP(),
		middleware.CleanPath(),
		middleware.RedirectSlashes(),
		middleware.Timeout(time.Minute),
		middleware.ThrottleBacklog(100, 500, time.Minute),
		middleware.RateLimiterByIP(30, time.Minute),
	)
	
	// public endpoint with a URL parameter (id) that has a regex match, and data binding.
	srv.Router.WithMultipartFormBindCtx("/:id|\\d+", handler, &Data{}, "data", http.MethodPost)
	
	// grouped routes that may have specific middlewares applied.
	srv.Router.Group(func(router *bingo.Router) {
		
		// Add a middleware to the following routes onward
		// using the `router` instance specific to this group.
		router.Use(jwtkit.VerifyAndAuthenticateJWT())
		
		// Restricted route, requires JWT authentication.
		router.HandleFunc("/admin", admin, http.MethodGet)
		
		// Any following route will be restricted by the JWT middleware.
		// ...
	})
	
	// public endpoint, not affected by the group set before.
	srv.Router.HandleFunc("/ping", ping, http.MethodGet)
	
	// Start the server and listen for incoming requests.
	if err := srv.ListenAndServe(); err != nil {
		srv.Logger.Fatal().Err(err).Msg("Server error")
	}
}
