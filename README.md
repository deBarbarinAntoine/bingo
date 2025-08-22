### BinGo - A Go Web Library ⚡

BinGo is a highly flexible and extensible web library for Go, designed for rapid development of robust and scalable web services. It's built on idiomatic Go principles and provides a clear, modular architecture with sensible defaults.

It uses widely adopted libraries like `zerolog` for logging, `alexedwards/flow` for routing, `alexedwards/scs` for sessions and session stores, `justinas/nosurf` for CSRF and `go-chi` for middlewares, CORS, JWT and rate limiting.

Contrary to other web libraries (and like `go-chi`), it's completely compatible with the `net/http` standard library.

-----

### 🚧 Work in Progress 🚧

This project is under active development. While the core features are functional, some parts are not yet complete, specifically the `Binder` implementation.

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

Create a new `Bingo` instance and configure it with your desired options, such as the server address, environment, and Redis store for sessions.

```go
package main

import (
    "log"
    "time"

    "github.com/debarbarinantoine/bingo/bingo"
    "github.com/debarbarinantoine/bingo/middleware"
)

func main() {
    srv := bingo.New(bingo.Options{
        ServerAddr:   "localhost:8080",
        Environment:  "development",
        RedisAddr:    "localhost:6379",
    }).
    WithLogMiddleware()

    srv.Router.Use(
        middleware.RealIP(),
        middleware.Recoverer(),
        middleware.Timeout(time.Minute),
    )

    // Define your routes and handlers
    srv.Router.HandleFunc("/", myHandler, "GET")

    log.Fatal(srv.ListenAndServe())
}
```

#### 2\. Define a Struct for Binding

Define a struct and use tags to map fields to data from the request. BinGo's `Binder` middleware can handle multiple sources at once.

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
    "github.com/debarbarinantoine/bingo/bingo"
    "net/http"
    "log"
)

type User struct {
    Name string `json:"name"`
    Age  int    `param:"age"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    // Get the bound data from the context
    user, _ := bingo.GetCtxData(r.Context(), "user").(*User)

    log.Printf("User: %+v", user)
}

func main() {
    srv := bingo.New(...)

    // The middleware automatically binds the data from multiple sources
    srv.Mux.WithJsonBindCtx("/users/:age", handler, &User{}, "user", "POST")

    log.Fatal(srv.ListenAndServe())
}
```

-----

### 🧑‍💻 Author

**Thorgan**

-----

### 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.
