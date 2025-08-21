// bingo.go
package bingo

import (
	"io"
	"net/http"
	"os"
	"time"
	
	"github.com/debarbarinantoine/bingo/middleware"
	"github.com/debarbarinantoine/bingo/router"
	
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog"
)

type Bingo struct {
	SessionManager *scs.SessionManager
	Mux            *router.Mux
	Server         *http.Server
	Logger         zerolog.Logger
	JWT            *jwtauth.JWTAuth
	RedisAddr      string
}

type Options struct {
	ServerAddr  string
	RedisAddr   string
	Environment string
}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
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
		RedisAddr: options.RedisAddr,
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

func (b *Bingo) WithSessions() *Bingo {
	// Check if session manager or JWT auth are already initialized
	if b.SessionManager != nil || b.JWT != nil {
		return b
	}
	
	// Initialize a new session manager and configure it to use redisstore as the session store.
	sessionManager := scs.New()
	if b.RedisAddr != "" {
		sessionManager.Store = redisstore.New(newPool(b.RedisAddr))
	}
	b.SessionManager = sessionManager
	b.Mux.Use(middleware.SetSessionManager(sessionManager))
	
	return b
}
