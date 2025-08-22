// router.go
package router

import (
	"net/http"
	"reflect"
	
	"github.com/debarbarinantoine/bingo/binder"
	"github.com/debarbarinantoine/bingo/context"
	"github.com/debarbarinantoine/bingo/enum"
	"github.com/debarbarinantoine/bingo/middleware"
	
	"github.com/alexedwards/flow"
)

type Router struct {
	*flow.Mux
}

func New() *Router {
	mux := flow.New()
	mux.Use(context.CtxData, middleware.Recoverer())
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
	mm := *r
	fn(&mm)
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
	// Get the type of dst
	dstType := reflect.TypeOf(dst)
	
	r.Group(func(router *Router) {
		router.Use(middleware.Binder(dstType, key, binderOptions...))
		router.HandleFunc(pattern, handler, methods...)
	})
}
