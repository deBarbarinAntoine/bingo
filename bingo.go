package bingo

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	
	"github.com/debarbarinantoine/bingo/internal/enum"
	"github.com/debarbarinantoine/bingo/jwtkit"
	"github.com/debarbarinantoine/bingo/middleware"
	"github.com/debarbarinantoine/bingo/sessions"
	
	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/mongodbstore"
	"github.com/alexedwards/scs/mssqlstore"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// Bingo represents the Bingo server.
type Bingo struct {
	Server         *http.Server
	Router         *Router
	Logger         zerolog.Logger
	SessionManager *scs.SessionManager
	JwtConfig      *jwtkit.Config
	environment    string
}

// Options represents the configuration options for the Bingo server.
type Options struct {
	
	// ServerAddr specifies the address to listen on for the server.
	//
	// For example: ":8080".
	ServerAddr string
	
	// Environment specifies the environment in which the server is running.
	// It can be "development", "production", or "test".
	//
	// Only "production" sets the logger to output JSON format.
	Environment string
	
	// UseRealIP specifies whether to use the RealIP middleware.
	//
	// If set to true, the RealIP middleware will be added to the middleware stack
	// to extract the real client IP from headers like X-Forwarded-For and X-Real-IP.
	//
	// Enable this when your server is behind proxies, load balancers, or CDNs.
	UseRealIP bool
}

// New returns a new Bingo instance with the default middleware stack:
// - middleware.CtxData()
// - middleware.Recoverer()
//
// Usage example:
//
// 	package main
//
// import (
//	"fmt"
//	"mime/multipart"
//	"net/http"
//	"time"
//
//	"github.com/debarbarinantoine/bingo"
//	"github.com/debarbarinantoine/bingo/jwtkit"
//	"github.com/debarbarinantoine/bingo/middleware"
//
//	"github.com/rs/zerolog/hlog"
// )
//
// // Publication represents a nested struct in the multipart data.
// type Publication struct {
//	Title string    `multipart:"title" json:"title"`
//	Score float64   `multipart:"score" json:"score"`
//	Date  time.Time `multipart:"date" json:"date"`
// }
//
// // Info represents another nested struct, containing a slice of Publication.
// type Info struct {
//	Id           uint          `multipart:"id" json:"id"`
//	Publications []Publication `multipart:"publications" json:"publications"`
// }
//
// // Data shows how to bind various data types from different request sources.
// // The `multipart` tag is used for form data.
// // The `cookie`, `header`, `param`, and `query` tags are for other request sources.
// type Data struct {
//	Session   string                `cookie:"session" json:"session"`
//	CSRFToken string                `header:"x-csrf-token" json:"x-csrf-token"`
//	Id        uint                  `param:"id" json:"id"`
//	Name      string                `query:"name" json:"name"`
//	Age       int                   `multipart:"age" json:"age"`
//	IsAdmin   bool                  `multipart:"is_admin" json:"is_admin"`
//	Score     float64               `multipart:"score" json:"score"`
//	Info      Info                  `multipart:"info" json:"info"`
//	File      *multipart.FileHeader `multipart:"file" json:"file"`
// }
//
// // handler processes the bound data and sends a JSON response.
// func handler(w http.ResponseWriter, r *http.Request) {
//	// Retrieve the bound data from the request context.
//	data, ok := bingo.GetCtxData(r.Context(), "data").(*Data)
//	if !ok {
//		bingo.ServerError(r, w, fmt.Errorf("no data found in context"), "no data")
//		return
//	}
//	hlog.FromRequest(r).Info().Msg("Request received successfully")
//
//	// Check if a file was uploaded and log its details.
//	if data.File != nil {
//		hlog.FromRequest(r).Info().
//			Int64("file_size", data.File.Size).
//			Str("file_name", data.File.Filename).
//			Msg("File uploaded successfully")
//	} else {
//		hlog.FromRequest(r).Error().Msg("No file uploaded")
//	}
//
//	// Respond to the client with the bound data in JSON format.
//	bingo.Json(r, w, data, http.StatusOK)
// }
//
// // ping is a handler for a common route without data.
// func ping(w http.ResponseWriter, r *http.Request) {
//	bingo.Json(r, w, bingo.H{"message": "pong"}, http.StatusOK)
// }
//
// // admin is a handler for a restricted route.
// func admin(w http.ResponseWriter, r *http.Request) {
//	bingo.Json(r, w, bingo.H{"message": "hello admin!"}, http.StatusOK)
// }
//
// func main() {
//	// Initialize JWT configuration with a secret key.
//	jwtConfig, err := jwtkit.NewConfigWithSecret(jwtkit.AlgorithmHS256, "|JwT53cr3T|", jwtkit.DefaultTokenResponseOptions())
//	if err != nil {
//		panic(err)
//	}
//
//	// Initialize the server and chain builder methods to add middleware.
//	// The "With..." methods use a builder pattern to configure the server.
//	srv, err := bingo.New(bingo.Options{
//		ServerAddr:  "localhost:8080",
//		Environment: "development",
//	}).
//		WithLogMiddleware().
//		WithJWT(jwtConfig)
//	if err != nil {
//		panic(err)
//	}
//
//	// Use global middlewares.
//	srv.Router.Use(
//		middleware.RealIP(),
//		middleware.CleanPath(),
//		middleware.RedirectSlashes(),
//		middleware.Timeout(time.Minute),
//		middleware.ThrottleBacklog(100, 500, time.Minute),
//		middleware.RateLimiterByIP(30, time.Minute),
//	)
//
//	// public endpoint with a URL parameter (id) that has a regex match, and data binding.
//	srv.Router.WithMultipartFormBindCtx("/:id|\\d+", handler, &Data{}, "data", http.MethodPost)
//
//	// grouped routes that may have specific middlewares applied.
//	srv.Router.Group(func(router *bingo.Router) {
//
//		// Add a middleware to the following routes onward
//		// using the `router` instance specific to this group.
//		router.Use(jwtkit.VerifyAndAuthenticateJWT())
//
//		// Restricted route, requires JWT authentication.
//		router.HandleFunc("/admin", admin, http.MethodGet)
//
//		// Any following route will be restricted by the JWT middleware.
//		// ...
//	})
//
//	// public endpoint, not affected by the group set before.
//	srv.Router.HandleFunc("/ping", ping, http.MethodGet)
//
//	// Start the server and listen for incoming requests.
//	if err := srv.ListenAndServe(); err != nil {
//		srv.Logger.Fatal().Err(err).Msg("Server error")
//	}
// }
//
func New(options Options) *Bingo {
	var output io.Writer = os.Stdout
	if options.Environment == "" {
		options.Environment = "development"
	}
	if options.Environment != "production" {
		output = zerolog.ConsoleWriter{Out: os.Stdout}
	}
	// Initialize the zerolog logger
	log := zerolog.New(output).With().
		Timestamp().
		Str("addr", options.ServerAddr).
		Logger()
	
	// Initialize a new router
	r := NewRouter()
	
	// Apply the RealIP middleware if UseRealIP is true
	if options.UseRealIP {
		r.Use(middleware.RealIP())
	}
	
	return &Bingo{
		environment: options.Environment,
		Logger:      log,
		Router:      r,
		Server: &http.Server{
			Addr:              options.ServerAddr,
			IdleTimeout:       time.Minute,
			ReadHeaderTimeout: 3 * time.Second,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			Handler:           r,
		},
	}
}

// ListenAndServe starts the server and listens for incoming requests.
func (b *Bingo) ListenAndServe() error {
	if b.environment != "production" {
		fmt.Println(":: [INFO] Registered routes:")
		for _, route := range b.Router.Routes {
			// GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD|CONNECT|TRACE (max length: 52)
			fmt.Printf("\t=> %52s  %s\n", strings.Join(route.Methods, "|"), route.Path)
		}
	}
	b.Logger.Info().Str("address", b.Server.Addr).Msg("Starting server")
	return b.Server.ListenAndServe()
}

// UseLogMiddleware adds a middleware to load logger in the request context.
func (b *Bingo) UseLogMiddleware() {
	b.Router.Use(middleware.Logger(b.Logger))
}

// UseSessionMiddleware adds a middleware to load the *scs.SessionManager in the request context.
func (b *Bingo) UseSessionMiddleware() {
	if b.SessionManager != nil {
		b.Router.Use(sessions.Session(b.SessionManager))
	}
}

// WithLogMiddleware adds a middleware to load logger in the request context.
//
// It can be used in conjunction with the bingo.New()
// function because it returns the same instance of Bingo.
//
// Example:
//
//	b := bingo.New().WithLogMiddleware()
func (b *Bingo) WithLogMiddleware() *Bingo {
	b.UseLogMiddleware()
	return b
}

// WithJWT adds a middleware to load the *jwtkit.Config in the request context.
//
// It can be used in conjunction with the bingo.New()
// function because it returns the same instance of Bingo.
//
// Example:
//
//	jwtConfig, err := jwtkit.NewConfigWithSecret(jwtkit.AlgorithmHS256, "|JwT53cr3T|", jwtkit.DefaultTokenResponseOptions())
//	if err != nil {
//		panic(err)
//	}
//
//	// Initialize the server and chain builder methods to add middleware.
//	// The "With..." methods use a builder pattern to configure the server.
//	srv, err := bingo.New(bingo.Options{
//		ServerAddr:  "localhost:8080",
//		Environment: "development",
//	}).
//		WithLogMiddleware().
//		WithJWT(jwtConfig)
//	if err != nil {
//		panic(err)
//	}
func (b *Bingo) WithJWT(config *jwtkit.Config) (*Bingo, error) {
	// Check if session manager or Config auth are already initialized
	if b.JwtConfig != nil || b.SessionManager != nil {
		return b, nil
	}
	
	// Initialize a new JWTAuth instance with the specified secret and algorithm
	b.JwtConfig = config
	b.Router.Use(jwtkit.SetJWT(config))
	
	return b, nil
}

// SessionOptions represents the options for configuring session management.
type SessionOptions struct {
	DBPool      any
	IdleTimeout time.Duration
	Lifetime    time.Duration
	
	// Store specifies the session store to use.
	// Use sessions.Stores struct to specify the session store to use.
	Store  enum.SessionStore
	Cookie scs.SessionCookie
}

// NewSessionOptions creates a new SessionOptions instance with default values.
func NewSessionOptions() SessionOptions {
	return SessionOptions{
		Lifetime: 24 * time.Hour,
		Store:    sessions.Stores.InMemory,
		Cookie: scs.SessionCookie{
			Name:        "session",
			HttpOnly:    true,
			Path:        "/",
			SameSite:    http.SameSiteLaxMode,
			Secure:      false,
			Partitioned: false,
			Persist:     true,
		},
	}
}

func assignSessionOptions(opt SessionOptions, sessionManager *scs.SessionManager) {
	if opt.IdleTimeout != 0 {
		sessionManager.IdleTimeout = opt.IdleTimeout
	}
	if opt.Cookie.Domain != "" {
		sessionManager.Cookie.Domain = opt.Cookie.Domain
	}
	sessionManager.Lifetime = opt.Lifetime
	sessionManager.Cookie.Name = opt.Cookie.Name
	sessionManager.Cookie.HttpOnly = opt.Cookie.HttpOnly
	sessionManager.Cookie.Path = opt.Cookie.Path
	sessionManager.Cookie.SameSite = opt.Cookie.SameSite
	sessionManager.Cookie.Secure = opt.Cookie.Secure
	sessionManager.Cookie.Partitioned = opt.Cookie.Partitioned
	sessionManager.Cookie.Persist = opt.Cookie.Persist
}

// WithSessions adds a middleware to load the *scs.SessionManager in the request context.
//
// It can be used in conjunction with the bingo.New()
// function because it returns the same instance of Bingo.
//
// Example:
//
//	srv, err := bingo.New(bingo.Options{
//		ServerAddr:  "localhost:8080",
//		Environment: "development",
//	}).
//		WithLogMiddleware().
//		WithSessions(NewSessionOptions())
//	if err != nil {
//		panic(err)
//	}
//
func (b *Bingo) WithSessions(opt SessionOptions) (*Bingo, error) {
	var err error
	
	// Check if session manager or Config auth are already initialized
	if b.SessionManager != nil || b.JwtConfig != nil {
		return b, nil
	}
	
	// Initialize a new session manager and configure it to use redisstore as the session store.
	sessionManager := scs.New()
	assignSessionOptions(opt, sessionManager)
	
	switch opt.Store {
	
	case sessions.Stores.GORM:
		gormConn, ok := opt.DBPool.(*gorm.DB)
		if !ok {
			b.Logger.Error().Err(ErrInvalidDBPool).Msg("opt.DBPool is not a valid database connection")
			return nil, ErrInvalidDBPool
		}
		if sessionManager.Store, err = gormstore.New(gormConn); err != nil {
			b.Logger.Error().Err(err).Msg("Failed to initialize session manager")
			return nil, err
		}
	
	case sessions.Stores.Redis:
		redisConn, ok := opt.DBPool.(*redis.Pool)
		if !ok {
			b.Logger.Error().Err(ErrInvalidDBPool).Msg("opt.DBPool is not a valid database connection")
			return nil, ErrInvalidDBPool
		}
		sessionManager.Store = redisstore.New(redisConn)
	
	case sessions.Stores.PostgreSQL:
		postgresConn, ok := opt.DBPool.(*sql.DB)
		if !ok {
			b.Logger.Error().Err(ErrInvalidDBPool).Msg("opt.DBPool is not a valid database connection")
			return nil, ErrInvalidDBPool
		}
		sessionManager.Store = postgresstore.New(postgresConn)
	
	case sessions.Stores.MySQL:
		mysqlConn, ok := opt.DBPool.(*sql.DB)
		if !ok {
			b.Logger.Error().Err(ErrInvalidDBPool).Msg("opt.DBPool is not a valid database connection")
			return nil, ErrInvalidDBPool
		}
		sessionManager.Store = mysqlstore.New(mysqlConn)
	
	case sessions.Stores.SQLite3:
		sqlite3Conn, ok := opt.DBPool.(*sql.DB)
		if !ok {
			b.Logger.Error().Err(ErrInvalidDBPool).Msg("opt.DBPool is not a valid database connection")
			return nil, ErrInvalidDBPool
		}
		sessionManager.Store = sqlite3store.New(sqlite3Conn)
	
	case sessions.Stores.MongoDB:
		mongoConn, ok := opt.DBPool.(*mongo.Database)
		if !ok {
			b.Logger.Error().Err(ErrInvalidDBPool).Msg("opt.DBPool is not a valid database connection")
			return nil, ErrInvalidDBPool
		}
		sessionManager.Store = mongodbstore.New(mongoConn)
	
	case sessions.Stores.MSSQL:
		mssqlConn, ok := opt.DBPool.(*sql.DB)
		if !ok {
			b.Logger.Error().Err(ErrInvalidDBPool).Msg("opt.DBPool is not a valid database connection")
			return nil, ErrInvalidDBPool
		}
		sessionManager.Store = mssqlstore.New(mssqlConn)
	
	// If no valid store is specified or if InMemory is specified, use the default in-memory store.
	default:
		sessionManager.Store = memstore.New()
	}
	
	// Set the session manager on the Bingo instance and apply the dependency injection middleware.
	b.SessionManager = sessionManager
	b.Router.Use(sessions.SetSessionManager(sessionManager))
	
	return b, nil
}
