// This example demonstrates:
// - Multi-source data binding (form, query, headers, cookies, URL params)
// - Automatic validation with user-friendly error messages
// - JWT authentication and route protection
// - Flexible middleware composition
// - Clean error handling and logging
package main

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"time"
	
	"github.com/debarbarinantoine/bingo"
	"github.com/debarbarinantoine/bingo/binder"
	"github.com/debarbarinantoine/bingo/jwtkit"
	"github.com/debarbarinantoine/bingo/middleware"
	
	"github.com/rs/zerolog/hlog"
)

// Publication represents a nested struct in the multipart data.
type Publication struct {
	Title string    `multipart:"title" json:"title" validate:"required"`
	Score float64   `multipart:"score" json:"score" validate:"required,min=0,max=10"`
	Date  time.Time `multipart:"date" json:"date"`
}

// Info represents another nested struct, containing a slice of Publication.
type Info struct {
	Id           uint          `multipart:"id" json:"id" validate:"required,min=1"`
	Publications []Publication `multipart:"publications" json:"publications"`
}

// Data shows how to bind various data types from different request sources.
//
// The `multipart` tag is used for multipart form data.
// The `cookie`, `header`, `param`, and `query` tags are for other request sources.
//
// The `validate` tag is used for data validation
// (see the go-playground/validator/v10 documentation for more details).
type Data struct {
	Session   string                `cookie:"session" json:"session" validate:"required,uuid"`
	CSRFToken string                `header:"x-csrf-token" json:"x-csrf-token"`
	Id        uint                  `param:"id" json:"id" validate:"required,min=1"`
	Name      string                `query:"name" json:"name" validate:"required"`
	Age       int                   `multipart:"age" json:"age" validate:"required,min=0,max=120"`
	IsAdmin   bool                  `multipart:"is_admin" json:"is_admin"`
	Score     float64               `multipart:"score" json:"score" validate:"required,min=0,max=10"`
	Info      Info                  `multipart:"info" json:"info"`
	File      *multipart.FileHeader `multipart:"file" json:"file"`
}

// ping is a handler for a common route without data.
func ping(w http.ResponseWriter, r *http.Request) {
	// Respond to the client with a JSON message using the bingo.Json helper and bingo.H type (shortcut for map[string]any).
	bingo.Json(r, w, bingo.H{"message": "pong"}, http.StatusOK)
}

// handler processes the bound and verified data and sends a JSON response.
func handler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the bound data from the request context.
	data, ok := bingo.GetCtxData(r.Context(), "data").(*Data)
	if !ok {
		// Return an `Internal Server Error` with the bingo.ServerError helper.
		bingo.ServerError(r, w, fmt.Errorf("no data found in context"), "no data")
		return
	}
	
	// Log at the info level with the data set by the Log middleware and the following message.
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
		UseRealIP:   true,
	}).
		WithLogMiddleware().
		WithJWT(jwtConfig)
	if err != nil {
		panic(err)
	}
	
	// Use global middlewares.
	srv.Router.Use(
		middleware.CleanPath(),
		middleware.RedirectSlashes(),
		middleware.Timeout(time.Minute),
		middleware.ThrottleBacklog(100, 500, time.Minute),
		middleware.RateLimiterByIP(30, time.Minute),
	)
	
	// public endpoint with a URL parameter (id) that has a regex match.
	//
	// Uses the *Router.Post() method with the `bingo.WithBinderAndValidator()` functional option to bind and validate data.
	srv.Router.Post("/:id|\\d+", handler, bingo.WithBinderAndValidator(&Data{}, "data", binder.WithoutJSONBinder()))
	
	// grouped routes that may have specific middlewares applied.
	srv.Router.Group(func(router *bingo.Router) {
		
		// Add a middleware to the following routes onward
		// using the `router` instance specific to this group.
		router.Use(jwtkit.VerifyAndAuthenticateJWT())
		
		// Restricted route, requires JWT authentication.
		// Uses the *Router.Get() method
		router.Get("/admin", admin)
		
		// Any following route will be restricted by the JWT middleware.
		// ...
	})
	
	// public endpoint, not affected by the group set before.
	// Uses the *Router.HandleFunc() method corresponding to the native http.HandleFunc method
	srv.Router.HandleFunc("/ping", ping, http.MethodGet)
	
	// Start the server and listen for incoming requests.
	if err := srv.ListenAndServe(); err != nil {
		srv.Logger.Fatal().Err(err).Msg("Server error")
	}
}
