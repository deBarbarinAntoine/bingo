package router

import (
	"net/http"
	
	"github.com/debarbarinantoine/bingo/binder"
	"github.com/debarbarinantoine/bingo/middleware"
	
	"github.com/alexedwards/flow"
)

// Router is a wrapper around flow.Mux that provides an API compatible with Bingo for registering routes.
type Router struct {
	*flow.Mux
}

// New returns a new Router instance with the default middleware stack:
//
// Middleware:
// - middleware.CtxData()
// - middleware.Recoverer()
func New() *Router {
	mux := flow.New()
	mux.Use(middleware.CtxData(), middleware.Recoverer())
	return &Router{
		Mux: mux,
	}
}

// Group is used to create 'groups' of routes in a Mux. Middleware registered
// inside the group will only be used on the routes in that group. See the
// example code at the start of the package documentation for how to use this
// feature.
//
// N.B.: it comes from alexedwards/flow library
func (r *Router) Group(fn func(*Router)) {
	// We call the original flow.Mux.Group method.
	// The key is the function we pass to it.
	r.Mux.Group(func(m *flow.Mux) {
		// Here, we take the *flow.Mux (m) provided by the original library,
		// wrap it in a new *Router, and then pass that to the user's function.
		// This is the crucial adapter/bridge step.
		fn(&Router{Mux: m})
	})
}

// RouteOption is a functional option that configures a route by adding
// middleware to its handler chain. Each option wraps the route's handler
// with specific middleware, such as for data binding or validation.
//
// For example, to bind JSON data to a struct and validate it:
//
//  router.Post("/users", userHandler,
//      WithBinderAndValidator(&user.CreateDTO{}, "user"),
//  )
type RouteOption func(*Router)

// Get is a shortcut for registering a GET route with the given pattern and handler.
//
// It may accept RouteOption to add middleware.Binder and middleware.Validator.
func (r *Router) Get(pattern string, handler http.HandlerFunc, opts ...RouteOption) {
	r.Group(func(router *Router) {
		for _, opt := range opts {
			opt(router)
		}
		router.HandleFunc(pattern, handler, http.MethodGet)
	})
}

// Post is a shortcut for registering a POST route with the given pattern and handler.
//
// It may accept RouteOption to add middleware.Binder and middleware.Validator.
func (r *Router) Post(pattern string, handler http.HandlerFunc, opts ...RouteOption) {
	r.Group(func(router *Router) {
		for _, opt := range opts {
			opt(router)
		}
		router.HandleFunc(pattern, handler, http.MethodPost)
	})
}

// Put is a shortcut for registering a PUT route with the given pattern and handler.
//
// It may accept RouteOption to add middleware.Binder and middleware.Validator.
func (r *Router) Put(pattern string, handler http.HandlerFunc, opts ...RouteOption) {
	r.Group(func(router *Router) {
		for _, opt := range opts {
			opt(router)
		}
		router.HandleFunc(pattern, handler, http.MethodPut)
	})
}

// Patch is a shortcut for registering a PATCH route with the given pattern and handler.
//
// It may accept RouteOption to add middleware.Binder and middleware.Validator.
func (r *Router) Patch(pattern string, handler http.HandlerFunc, opts ...RouteOption) {
	r.Group(func(router *Router) {
		for _, opt := range opts {
			opt(router)
		}
		router.HandleFunc(pattern, handler, http.MethodPatch)
	})
}

// Delete is a shortcut for registering a DELETE route with the given pattern and handler.
//
// It may accept RouteOption to add middleware.Binder and middleware.Validator.
func (r *Router) Delete(pattern string, handler http.HandlerFunc, opts ...RouteOption) {
	r.Group(func(router *Router) {
		for _, opt := range opts {
			opt(router)
		}
		router.HandleFunc(pattern, handler, http.MethodDelete)
	})
}

// WithBinder is a route option that provides a Binder middleware with the given destination and key.
func WithBinder(dst any, key string, options ...binder.MultiBinderOption) RouteOption {
	return func(r *Router) {
		r.Use(middleware.Binder(dst, key, options...))
	}
}

// WithValidator is a route option that provides a Validator middleware for the data bound to the request with the given key.
func WithValidator(key string) RouteOption {
	return func(r *Router) {
		r.Use(middleware.Validator(key))
	}
}

// WithBinderAndValidator is a route option that provides a Binder and Validator middleware with the given destination and key.
func WithBinderAndValidator(dst any, key string, options ...binder.MultiBinderOption) RouteOption {
	return func(r *Router) {
		r.Use(middleware.Binder(dst, key, options...))
		r.Use(middleware.Validator(key))
	}
}

// WithMiddleware is a route option that applies the given middleware(s) to the route.
func WithMiddleware(middleware ...func(http.Handler) http.Handler) RouteOption {
	return func(r *Router) {
		r.Use(middleware...)
	}
}
