// middlewares.go
package middleware

import (
	"net/http"
	"reflect"
	"time"
	
	"github.com/debarbarinantoine/bingo/binder"
	"github.com/debarbarinantoine/bingo/internal/ctx"
	
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/justinas/nosurf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

type Middleware func(next http.Handler) http.Handler

func CtxData() Middleware {
	return ctx.Data
}

// GenerateCSRF takes an HTTP request and returns
// the Csrf token for that request
// or an empty string if the token does not exist.
//
// Note that the token won't be available after
// CSRF finishes
// (that is, in another handler that wraps it,
// or after the request has been served)
//
// N.B.: it comes from justinas/nosurf library
func GenerateCSRF(r *http.Request) string {
	return nosurf.Token(r)
}

// CSRF constructs a new CSRFHandler that calls
// the specified handler if the CSRF check succeeds.
//
// N.B.: it comes from justinas/nosurf library
func CSRF(isCookieSecure bool) Middleware {
	return func(next http.Handler) http.Handler {
		csrfHandler := nosurf.New(next)
		csrfHandler.SetBaseCookie(http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Secure:   isCookieSecure,
		})
		return csrfHandler
	}
}

// Headers adds the specified headers to the response.
func Headers(headers map[string]string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for key, value := range headers {
				w.Header().Set(key, value)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// AllowContentType enforces a whitelist of request Content-Types otherwise responds
// with a 415 Unsupported Media Type status.
//
// N.B.: it comes from go-chi/chi/v5/middleware library
func AllowContentType(contentTypes ...string) Middleware {
	return middleware.AllowContentType(contentTypes...)
}

// CleanPath middleware will clean out double slash mistakes from a user's request path.
// For example, if a user requests /users//1 or //users////1 will both be treated as: /users/1
//
// N.B.: it comes from go-chi/chi/v5/middleware library
func CleanPath() Middleware {
	return middleware.CleanPath
}

// RealIP is a middleware that sets a http.Request's RemoteAddr to the results
// of parsing either the True-Client-IP, X-Real-IP or the X-Forwarded-For headers
// (in that order).
//
// This middleware should be inserted fairly early in the middleware stack to
// ensure that subsequent layers (e.g., request loggers) which examine the
// RemoteAddr will see the intended value.
//
// You should only use this middleware if you can trust the headers passed to
// you (in particular, the three headers this middleware uses), for example
// because you have placed a reverse proxy like HAProxy or nginx in front of
// chi. If your reverse proxies are configured to pass along arbitrary header
// values from the client, or if you use this middleware without a reverse
// proxy, malicious clients will be able to make you very sad (or, depending on
// how you're using RemoteAddr, vulnerable to an attack of some sort).
//
// N.B.: it comes from go-chi/chi/v5/middleware library
func RealIP() Middleware {
	return middleware.RealIP
}

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
//
// Alternatively, look at https://github.com/go-chi/httplog middleware pkgs.
//
// N.B.: it comes from go-chi/chi/v5/middleware library
func Recoverer() Middleware {
	return middleware.Recoverer
}

// RedirectSlashes is a middleware that will match request paths with a trailing
// slash and redirect to the same path, less the trailing slash.
//
// NOTE: RedirectSlashes middleware is *incompatible* with http.FileServer,
// see https://github.com/go-chi/chi/issues/343
//
// N.B.: it comes from go-chi/chi/v5/middleware library
func RedirectSlashes() Middleware {
	return middleware.RedirectSlashes
}

// CorsOptions is a configuration container to setup the CORS middleware.
//
// N.B.: it comes from go-chi/cors library
type CorsOptions struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
	// Only one wildcard can be used per origin.
	// Default value is ["*"]
	AllowedOrigins []string
	
	// AllowOriginFunc is a custom function to validate the origin. It takes the origin
	// as argument and returns true if allowed or false otherwise. If this option is
	// set, the content of AllowedOrigins is ignored.
	AllowOriginFunc func(r *http.Request, origin string) bool
	
	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string
	
	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string
	
	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string
	
	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool
	
	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge int
	
	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool
	
	// Debugging flag adds additional output to debug server side CORS issues
	Debug bool
}

// Cors creates a new CORS handler with passed options.
//
// N.B.: it comes from go-chi/cors library
func Cors(options CorsOptions) Middleware {
	opts := cors.Options{
		AllowedOrigins:     options.AllowedOrigins,
		AllowOriginFunc:    options.AllowOriginFunc,
		AllowedMethods:     options.AllowedMethods,
		AllowedHeaders:     options.AllowedHeaders,
		ExposedHeaders:     options.ExposedHeaders,
		AllowCredentials:   options.AllowCredentials,
		MaxAge:             options.MaxAge,
		OptionsPassthrough: options.OptionsPassthrough,
		Debug:              options.Debug,
	}
	return cors.Handler(opts)
}

// Throttle is a middleware that limits number of currently processed requests
// at a time across all users. Note: Throttle is not a rate-limiter per user,
// instead it just puts a ceiling on the number of current in-flight requests
// being processed from the point from where the Throttle middleware is mounted.
//
// N.B.: it comes from go-chi/chi/v5/middleware library
func Throttle(limit int) Middleware {
	return middleware.Throttle(limit)
}

// ThrottleBacklog is a middleware that limits number of currently processed
// requests at a time and provides a backlog for holding a finite number of
// pending requests.
//
// N.B.: it comes from go-chi/chi/v5/middleware library
func ThrottleBacklog(limit int, backlogLimit int, backlogTimeout time.Duration) Middleware {
	return middleware.ThrottleBacklog(limit, backlogLimit, backlogTimeout)
}

// Timeout is a middleware that cancels ctx after a given timeout and return
// a 504 Gateway Timeout error to the client.
//
// It's required that you select the ctx.Done() channel to check for the signal
// if the ctx has reached its deadline and return, otherwise the timeout
// signal will be just ignored.
//
// ie. a route/handler may look like:
//
//	r.Get("/long", func(w http.ResponseWriter, r *http.Request) {
//		ctx := r.Context()
//		processTime := time.Duration(rand.Intn(4)+1) * time.Second
//
//		select {
//		case <-ctx.Done():
//			return
//
//		case <-time.After(processTime):
//			// The above channel simulates some hard work.
//		}
//
//		w.Write([]byte("done"))
//	})
//
// N.B.: it comes from go-chi/chi/v5/middleware library
func Timeout(timeout time.Duration) Middleware {
	return middleware.Timeout(timeout)
}

// RateLimiterByIP is a middleware that limits number of requests per IP address in a given time window.
//
// N.B.: it comes from go-chi/chi/v5/middleware library
func RateLimiterByIP(requestLimit int, windowLength time.Duration) Middleware {
	return httprate.LimitByIP(requestLimit, windowLength)
}

// Logger is a middleware that loads a logger into the request's ctx, along
// with some useful data:
// 	- RemoteAddr
// 	- RemoteIP
// 	- UserAgent
// 	- Referer
// 	- RequestID
// 	- Method
//
// When standard output is a TTY, Logger will
// print in color, otherwise it will print in black and white.
//
// N.B.: it uses rs/zerolog logger
func Logger(logger zerolog.Logger) Middleware {
	logHandler := hlog.NewHandler(logger)
	remoteAddrHandler := hlog.RemoteAddrHandler("addr")
	remoteIPHandler := hlog.RemoteIPHandler("remote_ip")
	userAgentHandler := hlog.UserAgentHandler("user_agent")
	refererHandler := hlog.RefererHandler("referer")
	requestIDHandler := hlog.RequestIDHandler("req_id", "Request-Id")
	methodHandler := hlog.MethodHandler("method")
	
	return func(next http.Handler) http.Handler {
		return logHandler(
			remoteAddrHandler(
				remoteIPHandler(
					userAgentHandler(
						refererHandler(
							requestIDHandler(
								methodHandler(next),
							),
						),
					),
				),
			),
		)
	}
}

// Binder is a middleware that binds the request body to a struct using the
// provided binder. It supports:
// 	- JSON
// 	- form-urlencoded
// 	- multipart/form-data
// 	- query strings
// 	- url params
// 	- headers
// 	- cookies
func Binder(dst any, key string, binderOptions ...binder.MultiBinderOption) Middleware {
	
	// Get the type of dst
	dstType := reflect.TypeOf(dst)
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			
			// Create a new variable of the same type as dst
			newDst := reflect.New(dstType.Elem()).Interface()
			
			dataBinder, err := binder.NewMultiBinder(newDst, r, binderOptions...)
			
			err = dataBinder.Fetch()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			r = ctx.SetData(r, key, newDst)
			next.ServeHTTP(w, r)
		})
	}
}
