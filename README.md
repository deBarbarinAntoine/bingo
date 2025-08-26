

### ![Bingo Logo](logox32.png) BinGo - A Go Web Library

BinGo is a highly flexible and extensible web library for Go, designed for rapid development of robust and scalable web services. It's built on idiomatic Go principles and provides a clear, modular architecture with sensible defaults.

> Requires **Go 1.24** or higher.

It uses widely adopted libraries like:
- `zerolog` for logging
- `alexedwards/flow` for routing
- `alexedwards/scs` for sessions and session stores
- `lestrrat-go/jwx/v3` for JWT handling
- `go-playground/validator/v10` for data validation
- `justinas/nosurf` for CSRF
- `go-chi` for additional middleware utilities (like CORS, JWT, rate limiting, etc.).

Contrary to some web frameworks (and like `go-chi`), it's completely compatible with the `net/http` standard library.

---

### 🚧 Work in Progress 🚧

This project is under active development. While the core features are functional, some parts are not yet complete, specifically the `Binder` implementation that may encounter issues with complex data types.

Use with caution.

---

### 💡 Features

* **Modular Architecture**: Built with a layered design, separating the server, router, and middleware for clean, maintainable code.
* **Request Binding**:  Binds data from various HTTP request sources (JSON, URL parameters, form data, headers, and cookies) directly to a Go struct using simple tags.
* **Automatic Type Conversion**: Automatically converts string values from requests to appropriate Go types (integers, floats, booleans, etc.).
* **Automatic Validation**: Validates data using the `go-playground/validator/v10` library, ensuring data integrity.
* **Secure by Default**: Integrates with battle-tested libraries for essential security features like CSRF protection.
* **Extensive Middleware**: A curated collection of powerful middlewares for common web tasks:
	* **Logging**: Structured logging via `zerolog`.
	* **Authentication**: Supports both session-based authentication with multiple stores (PostgreSQL, MySQL, MSSQL, SQLite3, GORM, Redis, etc.) and stateless JWT-based authentication (with secret or RSA/ECDSA/EdDSA encryption).
	* **Rate Limiting**: Throttling and rate limiting to protect against abuse.
	* **HTTP Helpers**: Clean URL paths, panic recovery, and timeout handling.
* **Useful Helpers**: JSON responses, Error responses, Context setter and getter, Session helpers, and more.
* **Extensible Design**: The core components are designed to be easily extensible, allowing developers to add custom binders, middlewares, or authentication methods.

---

### 🤔 Why Choose BinGo?

**BinGo** strikes the perfect balance between **developer productivity** and **Go's philosophy** of simplicity. Here's what sets it apart:

#### **🚀 Eliminates Boilerplate Without Magic**
Instead of writing repetitive parsing and validation code:
```go
// Traditional approach - lots of manual work
func handler(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }
    
    email := r.FormValue("email")
    if email == "" {
        http.Error(w, "Email required", http.StatusBadRequest)
        return
    }
    
    if !isValidEmail(email) {
        http.Error(w, "Invalid email", http.StatusBadRequest)
        return
    }
    
    age, err := strconv.Atoi(r.FormValue("age"))
    if err != nil || age < 0 {
        http.Error(w, "Invalid age", http.StatusBadRequest)
        return
    }
    
    // Finally, your business logic...
}
```

BinGo reduces this to:
```go
// BinGo approach - focus on business logic
type UserData struct {
    ID    int    `param:"id" validate:"required,gt=0"`
    Email string `form:"email" validate:"required,email"`
    Age   int    `form:"age" validate:"required,gte=0"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    user := bingo.GetCtxData(r.Context(), "user").(*UserData)
    // Your business logic starts here - data is already validated!
}
```

#### **🛡️ Security & Reliability First**
- **Battle-tested dependencies**: Built on proven libraries like `zerolog`, `go-playground/validator`, and `alexedwards/flow`
- **Secure defaults**: CSRF protection, rate limiting, and proper error handling out of the box
- **Type safety**: Automatic type conversion with validation prevents runtime errors
- **Panic Recovery**: Built-in panic recovery middleware and graceful server shutdown
- **Concurrency helper**: Built-in function to spawn goroutines with recovery feature and wait group when shutting server down.

#### **🧩 True Modularity**
Unlike monolithic frameworks, BinGo is **composable**:
- Use only what you need - each middleware is independent
- Mix with existing `net/http` code seamlessly
- Add your own middlewares without framework lock-in
- Compatible with the entire Go HTTP ecosystem

#### **📈 Scales With Your Needs**
- **Prototype quickly** with automatic binding and validation
- **Grow confidently** with enterprise-grade middleware (JWT, sessions, rate limiting)
- **Deploy anywhere** - it's still just a standard Go HTTP server

#### **💡 Developer Experience That Just Works**
- **Intuitive API**: If you know `net/http`, you already understand BinGo
- **Clear error messages**: Validation errors are automatically translated to user-friendly messages
- **Comprehensive logging**: Request tracing and structured logging built-in
- **Zero configuration**: Sensible defaults get you started immediately

#### **🆚 How It Compares**

| Feature                     | BinGo    | Gin      | Echo     | Chi     | net/http  |
|-----------------------------|----------|----------|----------|---------|-----------|
| Learning Curve              | Low      | Medium   | Medium   | Low     | Low       |
| Boilerplate                 | Minimal  | Low      | Low      | High    | Very High |
| Data Binding                | Built-in | Built-in | Built-in | Manual  | Manual    |
| Validation                  | Built-in | Manual   | Manual   | Manual  | Manual    |
| Standard Library Compatible | ✅        | ❌        | ❌        | ✅       | ✅         |
| Middleware Ecosystem        | Curated  | Large    | Large    | Large   | DIY       |
| Performance Overhead        | Minimal  | Low      | Low      | Minimal | None      |

**Perfect for:**
- REST APIs that need robust input validation
- Microservices requiring consistent error handling
- Projects transitioning from manual `net/http` code
- Teams that value Go's simplicity but want modern DX

**Not ideal for:**
- Simple static file servers
- Projects requiring maximum performance (use `net/http` directly)
- Teams preferring full-framework approaches like Gin/Echo

---

### 📦 Installation

To use BinGo in your project, install it with `go get`:

```sh
go get github.com/debarbarinantoine/bingo
```

---

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

#### 2\. Define a Struct for Binding and Validation

Define a struct and use tags to map fields to data from the request. BinGo's `Binder` middleware can handle multiple sources at once.

But know that only one source placed in the body can be sent at a time (e.g., JSON, form data and multipart form data).

To know more about the validation tags, see the [go-playground/validator/v10](https://github.com/go-playground/validator/v10) documentation.

```go
import (
    "mime/multipart"
    "time"
)

type MyRequestData struct {
    UserID  int     `param:"user_id" validate:"required,gt=0"`
    Query   string  `query:"q"`
    Email   string  `json:"email" validate:"required,email"`
    File    *multipart.FileHeader `multipart:"file"`
    AuthToken string `header:"Authorization" validate:"required"`
    SessionID string `cookie:"session_id" validate:"required,len=32"`
}
```

#### 3\. Handle the Request

Use BinGo's `Get`, `Post`, `Put`, `Patch`, `Delete` helpers with the `bingo.WithBinderAndValidator()` functional option to automatically bind the data to your struct, validate it, and make it available in the handler's context.

```go
package main

import (
	"log"
	"net/http"
	
    "github.com/debarbarinantoine/bingo"
)

type User struct {
    Name string `json:"name" validate:"required"`
    Age  int    `param:"age" validate:"required,lte=16,gte=120"`
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

    // The middleware automatically binds the data from multiple sources and validates it
    srv.Router.Post("/users/:age", handler, bingo.WithBinderAndValidator(&User{}, "user"))

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
	"github.com/debarbarinantoine/bingo/binder"
	"github.com/debarbarinantoine/bingo/jwtkit"
	"github.com/debarbarinantoine/bingo/middleware"
	
	"github.com/rs/zerolog/hlog"
)

// Publication represents a nested struct in the multipart data.
type Publication struct {
	Title string    `multipart:"title" json:"title" validate:"required,min=2,max=255"`
	Score float64   `multipart:"score" json:"score" validate:"required,gte=0,lte=10"`
	Date  time.Time `multipart:"date" json:"date" validate:"required"`
}

// Info represents another nested struct, containing a slice of Publication.
type Info struct {
	Id           uint          `multipart:"id" json:"id" validate:"required,gt=0"`
	Publications []Publication `multipart:"publications" json:"publications" validate:"required,dive"`
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
	Id        uint                  `param:"id" json:"id" validate:"required,gt=0"`
	Name      string                `query:"name" json:"name" validate:"required,min=2,max=100"`
	Age       int                   `multipart:"age" json:"age" validate:"required,gt=0,lte=120"`
	IsAdmin   bool                  `multipart:"is_admin" json:"is_admin"`
	Score     float64               `multipart:"score" json:"score" validate:"required,gte=0,lte=10"`
	Info      Info                  `multipart:"info" json:"info" validate:"required,dive"`
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
```

---

### 🧑‍💻 Author

**Thorgan**

---

### 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.
