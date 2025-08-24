package bingo

import (
	"net/http"
	
	"github.com/debarbarinantoine/bingo/binder"
	"github.com/debarbarinantoine/bingo/internal/router"
)

type Router struct {
	*router.Router
}

// NewRouter returns a new Router instance with the default middleware stack:
//
// Middlewares:
//  - middleware.CtxData()
//  - middleware.Recoverer()
//
// It's normally automatically called by bingo.New(),
// that creates a Bingo instance, containing a *Router
//
// Usage:
//
//	r := bingo.New()
//	r.HandleFunc("/home", HomeHandler, http.MethodGet)
//
func NewRouter() *Router {
	return &Router{
		Router: router.New(),
	}
}

// Group is used to create 'groups' of routes in a Mux. Middleware registered
// inside the group will only be used on the routes in that group. See the
// example code at the start of the package documentation for how to use this
// feature.
//
// N.B.: it comes from alexedwards/flow library
func (r *Router) Group(fn func(*Router)) {
	r.Router.Group(func(r *router.Router) {
		fn(&Router{Router: r})
	})
}

// RouteOption is a functional option that configures a route by adding middleware to its handler chain. Each option wraps the route's handler with specific middleware, such as for data binding or validation.
//
// For example, to bind JSON data to a struct and validate it:
// router.Post("/users", userHandler,
//    WithBinderAndValidator(&user.CreateDTO{}, "user"),
// )
type RouteOption router.RouteOption

// WithBinder is a route option that provides a Binder middleware with the given destination and key.
func WithBinder(dst any, key string, options ...binder.MultiBinderOption) router.RouteOption {
	return router.WithBinder(dst, key, options...)
}

// WithValidator is a route option that provides a Validator middleware for the data bound to the request with the given key.
func WithValidator(key string) router.RouteOption {
	return router.WithValidator(key)
}

// WithBinderAndValidator is a route option that provides a Binder and Validator middleware with the given destination and key.
func WithBinderAndValidator(dst any, key string, options ...binder.MultiBinderOption) router.RouteOption {
	return router.WithBinderAndValidator(dst, key, options...)
}

// WithMiddleware is a route option that applies the given middleware(s) to the route.
func WithMiddleware(middleware ...func(http.Handler) http.Handler) router.RouteOption {
	return router.WithMiddleware(middleware...)
}
