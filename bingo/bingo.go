// bingo.go
package bingo

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"time"
	
	"github.com/debarbarinantoine/bingo/enum"
	"github.com/debarbarinantoine/bingo/middleware"
	"github.com/debarbarinantoine/bingo/router"
	
	"github.com/alexedwards/scs/gormstore"
	"github.com/alexedwards/scs/mongodbstore"
	"github.com/alexedwards/scs/mssqlstore"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Bingo struct {
	SessionManager *scs.SessionManager
	Mux            *router.Mux
	Server         *http.Server
	Logger         zerolog.Logger
	JWT            *jwtauth.JWTAuth
}

type Options struct {
	ServerAddr  string
	Environment string
}

func New(options Options) *Bingo {
	var output io.Writer = os.Stdout
	if options.Environment != "production" {
		output = zerolog.ConsoleWriter{Out: os.Stdout}
	}
	// Initialize the zerolog logger
	log := zerolog.New(output).With().
		Timestamp().
		Str("addr", options.ServerAddr).
		Logger()
	
	// Initialize a new router
	mux := router.New()
	
	return &Bingo{
		Logger:    log,
		Mux:       mux,
		Server: &http.Server{
			Addr:              options.ServerAddr,
			IdleTimeout:       time.Minute,
			ReadHeaderTimeout: 3 * time.Second,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			Handler:           mux,
		},
	}
}

func (b *Bingo) ListenAndServe() error {
	return b.Server.ListenAndServe()
}

func (b *Bingo) UseLogMiddleware() {
	b.Mux.Use(middleware.Logger(b.Logger))
}

func (b *Bingo) UseSessionMiddleware() {
	if b.SessionManager != nil {
		b.Mux.Use(middleware.Session(b.SessionManager))
	}
}

func (b *Bingo) WithLogMiddleware() *Bingo {
	b.UseLogMiddleware()
	return b
}

func (b *Bingo) WithJWT(jwtAlgorithm, jwtSecret string) *Bingo {
	// Check if session manager or JWT auth are already initialized
	if b.JWT != nil || b.SessionManager != nil {
		return b
	}
	
	// Initialize a new JWTAuth instance with the specified secret and algorithm
	var jwt *jwtauth.JWTAuth = nil
	if jwtSecret != "" && jwtAlgorithm != "" {
		jwt = jwtauth.New(jwtAlgorithm, []byte(jwtSecret), nil)
		b.JWT = jwt
		b.Mux.Use(middleware.SetJWT(jwt))
	}
	return b
}

type SessionOptions struct {
	DBPool      any
	IdleTimeout time.Duration
	Lifetime    time.Duration
	Store       enum.SessionStore
	Cookie      scs.SessionCookie
}

func NewSessionOptions() SessionOptions {
	return SessionOptions{
		Lifetime: 24 * time.Hour,
		Store:    enum.SessionStores.InMemory,
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

func (b *Bingo) WithSessions(opt SessionOptions) (*Bingo, error) {
	var err error
	
	// Check if session manager or JWT auth are already initialized
	if b.SessionManager != nil || b.JWT != nil {
		return b, nil
	}
	
	// Initialize a new session manager and configure it to use redisstore as the session store.
	sessionManager := scs.New()
	assignSessionOptions(opt, sessionManager)
	
	switch opt.Store {
	
	case enum.SessionStores.GORM:
		gormConn, ok := opt.DBPool.(*gorm.DB)
		if !ok {
			b.Logger.Error().Err(ErrInvalidDBPool).Msg("opt.DBPool is not a valid database connection")
			return nil, ErrInvalidDBPool
		}
		if sessionManager.Store, err = gormstore.New(gormConn); err != nil {
			b.Logger.Error().Err(err).Msg("Failed to initialize session manager")
			return nil, err
		}
	
	case enum.SessionStores.Redis:
		redisConn, ok := opt.DBPool.(*redis.Pool)
		if !ok {
			b.Logger.Error().Err(ErrInvalidDBPool).Msg("opt.DBPool is not a valid database connection")
			return nil, ErrInvalidDBPool
		}
		sessionManager.Store = redisstore.New(redisConn)
	
	case enum.SessionStores.PostgreSQL:
		postgresConn, ok := opt.DBPool.(*sql.DB)
		if !ok {
			b.Logger.Error().Err(ErrInvalidDBPool).Msg("opt.DBPool is not a valid database connection")
			return nil, ErrInvalidDBPool
		}
		sessionManager.Store = postgresstore.New(postgresConn)
	
	case enum.SessionStores.MySQL:
		mysqlConn, ok := opt.DBPool.(*sql.DB)
		if !ok {
			b.Logger.Error().Err(ErrInvalidDBPool).Msg("opt.DBPool is not a valid database connection")
			return nil, ErrInvalidDBPool
		}
		sessionManager.Store = mysqlstore.New(mysqlConn)
	
	case enum.SessionStores.SQLite3:
		sqlite3Conn, ok := opt.DBPool.(*sql.DB)
		if !ok {
			b.Logger.Error().Err(ErrInvalidDBPool).Msg("opt.DBPool is not a valid database connection")
			return nil, ErrInvalidDBPool
		}
		sessionManager.Store = sqlite3store.New(sqlite3Conn)
	
	case enum.SessionStores.MongoDB:
		mongoConn, ok := opt.DBPool.(*mongo.Database)
		if !ok {
			b.Logger.Error().Err(ErrInvalidDBPool).Msg("opt.DBPool is not a valid database connection")
			return nil, ErrInvalidDBPool
		}
		sessionManager.Store = mongodbstore.New(mongoConn)
	
	case enum.SessionStores.MSSQL:
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
	b.Mux.Use(middleware.SetSessionManager(sessionManager))
	
	return b, nil
}
