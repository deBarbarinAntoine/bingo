### BinGo - A Go Web Library ⚡

BinGo is a highly flexible and extensible web library for Go, designed for rapid development of robust and scalable web services. It's built on idiomatic Go principles and provides a clear, modular architecture with sensible defaults.

It uses widely adopted libraries like `zerolog` for logging, `alexedwards/flow` for routing, `alexedwards/scs` for sessions and session stores, `justinas/nosurf` for CSRF and `go-chi` for middlewares, CORS, JWT and rate limiting.

Contrary to other web libraries (and like `go-chi`), it's completely compatible with the `net/http` standard library.

-----

### 🚧 Work in Progress 🚧

This project is under active development. While the core features are functional, some parts are not yet complete, specifically the `Binder` implementation that may encounter issues with complex data types.

Use with caution.

-----

### 💡 Features

* **Modular Architecture**: Built with a layered design, separating the server, router, and middleware for clean, maintainable code.
* **Request Binding**:  Binds data from various HTTP request sources (JSON, URL parameters, form data, headers, and cookies) directly to a Go struct using simple tags.
* **Automatic Type Conversion**: Automatically converts string values from requests to appropriate Go types (integers, floats, booleans, etc.).
* **Secure by Default**: Integrates with battle-tested libraries for essential security features like CSRF protection.
* **Extensive Middleware**: A curated collection of powerful middlewares for common web tasks:
	* **Logging**: Structured logging via `zerolog`.
	* **Authentication**: Supports both session-based authentication with multiple stores (PostgreSQL, MySQL, MSSQL, SQLite3, GORM, Redis, etc.) and stateless JWT-based authentication (with secret or RSA/ECDSA/EdDSA encryption).
	* **Rate Limiting**: Throttling and rate limiting to protect against abuse.
	* **HTTP Helpers**: Clean URL paths, panic recovery, and timeout handling.
* **Extensible Design**: The core components are designed to be easily extensible, allowing developers to add custom binders, middlewares, or authentication methods.

-----

### 📦 Installation

To use BinGo in your project, install it with `go get`:

```sh
go get github.com/debarbarinantoine/bingo
```

-----

### 📚 Usage

#### 1\. Initialize Your Server

Create a new `Bingo` instance and configure it with your desired options, such as the server address, environment, and authentication configuration.

```go
package main

import (
	"log"
	"time"
	
	"github.com/debarbarinantoine/bingo"
	"github.com/debarbarinantoine/bingo/jwtkit"
	"github.com/debarbarinantoine/bingo/middleware"
)

func main() {
	// Initialize JWT configuration with a secret key.
	jwtConfig, err := jwtkit.NewConfigWithSecret(jwtkit.AlgorithmHS256, "|JwT53cr3T|", jwtkit.DefaultTokenResponseOptions())
	if err != nil {
		panic(err)
	}
	
	// Initialize the server and chain builder methods to add middleware.
	// The "With..." methods use a builder pattern to configure the server.
	srv, err := bingo.New(bingo.Options{
		ServerAddr:  "localhost:8080",
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
	
	// The routing will go here...
	
	// Start the server.
	log.Fatal(srv.ListenAndServe())
}
```

#### 2\. Define a Struct for Binding

Define a struct and use tags to map fields to data from the request. BinGo's `Binder` middleware can handle multiple sources at once.

But know that only one source placed in the body can be sent at a time (e.g., JSON, form data and multipart form data).

```go
import (
    "mime/multipart"
    "time"
)

type MyRequestData struct {
    UserID  int     `param:"user_id"`
    Query   string  `query:"q"`
    Email   string  `json:"email"`
    File    *multipart.FileHeader `multipart:"file"`
    AuthToken string `header:"Authorization"`
    SessionID string `cookie:"session_id"`
}
```

#### 3\. Handle the Request

Use BinGo's `With...BindCtx` helpers to automatically bind the data to your struct and make it available in the handler's context.

```go
package main

import (
	"log"
	"net/http"
	
    "github.com/debarbarinantoine/bingo"
)

type User struct {
    Name string `json:"name"`
    Age  int    `param:"age"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Get the bound data from the ctx
    user, ok := bingo.GetCtxData(r.Context(), "user").(*User)
	if !ok {
		// Handle error if no user found in context
		bingo.ServerError(r, w, fmt.Errorf("no user found in context"), "no user")
		return
	}

	// Respond with the data in JSON format
    bingo.Json(r, w, bingo.H{"user": bingo.H{"name": user.Name, "age": user.Age}}, http.StatusOK)
}

func main() {
    srv := bingo.New(...)

    // The middleware automatically binds the data from multiple sources
    srv.Router.WithJsonBindCtx("/users/:age", handler, &User{}, "user", "POST")

    log.Fatal(srv.ListenAndServe())
}
```

---

#### Complete example

```go
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
		ServerAddr:  "localhost:8080",
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
```

-----

### 🧑‍💻 Author

**Thorgan**

-----

### 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.
