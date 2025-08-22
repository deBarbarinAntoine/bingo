package bingo

import (
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
