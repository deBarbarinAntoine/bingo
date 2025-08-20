// bingo.go
package bingo

import (
	"io"
	"net/http"
	"os"
	"time"
	
	"BinGo/router"
	
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog"
)

type Bingo struct {
	SessionManager *scs.SessionManager
	Mux            *router.Mux
	Server         *http.Server
	Logger         zerolog.Logger
}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		},
	}
}

func New(serverAddr, redisAddr, environment string) *Bingo {
	var output io.Writer = os.Stdout
	if environment != "production" {
		output = zerolog.ConsoleWriter{Out: os.Stdout}
	}
	// Initialize the zerolog logger
	log := zerolog.New(output).With().
		Timestamp().
		Str("addr", serverAddr).
		Logger()
	
	// Initialize a new session manager and configure it to use redisstore as the session store.
	sessionManager := scs.New()
	sessionManager.Store = redisstore.New(newPool(redisAddr))
	mux := router.New()
	return &Bingo{
		SessionManager: sessionManager,
		Logger:         log,
		Mux:            mux,
		Server: &http.Server{
			Addr:              serverAddr,
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
	b.Mux.Use(router.LoggerMiddleware(b.Logger))
}

func (b *Bingo) UseSessionMiddleware() {
	b.Mux.Use(router.SessionMiddleware(b.SessionManager))
}

func (b *Bingo) WithLogMiddleware() *Bingo {
	b.UseLogMiddleware()
	return b
}

func (b *Bingo) WithSessionMiddleware() *Bingo {
	b.UseSessionMiddleware()
	return b
}
