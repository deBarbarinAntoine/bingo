// router.go
package router

import (
	"net/http"
	
	"github.com/debarbarinantoine/bingo/binder"
	"github.com/debarbarinantoine/bingo/internal/enum"
	"github.com/debarbarinantoine/bingo/middleware"
	
	"github.com/alexedwards/flow"
)

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

func (r *Router) WithFormBindCtx(pattern string, handler http.HandlerFunc, dst any, key string, methods ...string) {
	r.withBindCtx(dst, key, binder.WithBodyBinder(enum.Tags.Form), pattern, handler, methods...)
}

func (r *Router) WithMultipartFormBindCtx(pattern string, handler http.HandlerFunc, dst any, key string, methods ...string) {
	r.withBindCtx(dst, key, binder.WithBodyBinder(enum.Tags.MultipartForm), pattern, handler, methods...)
}

func (r *Router) WithJsonBindCtx(pattern string, handler http.HandlerFunc, dst any, key string, methods ...string) {
	r.withBindCtx(dst, key, binder.WithBodyBinder(enum.Tags.Json), pattern, handler, methods...)
}

func (r *Router) WithBindCtx(pattern string, handler http.HandlerFunc, dst any, key string, methods ...string) {
	r.withBindCtx(dst, key, nil, pattern, handler, methods...)
}

func (r *Router) withBindCtx(dst any, key string, binderOptions []binder.MultiBinderOption, pattern string, handler http.HandlerFunc, methods ...string) {
	r.Group(func(router *Router) {
		router.Use(middleware.Binder(dst, key, binderOptions...))
		router.HandleFunc(pattern, handler, methods...)
	})
}
