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

type Mux struct {
	*flow.Mux
}

func New() *Mux {
	mux := flow.New()
	mux.Use(context.CtxData, middleware.Recoverer())
	return &Mux{
		Mux: mux,
	}
}

func (m *Mux) WithFormBindCtx(pattern string, handler http.HandlerFunc, dst any, key string, methods ...string) {
	m.withBindCtx(dst, key, binder.WithBodyBinder(enum.Tags.Form), pattern, handler, methods...)
}

func (m *Mux) WithMultipartFormBindCtx(pattern string, handler http.HandlerFunc, dst any, key string, methods ...string) {
	m.withBindCtx(dst, key, binder.WithBodyBinder(enum.Tags.MultipartForm), pattern, handler, methods...)
}

func (m *Mux) WithJsonBindCtx(pattern string, handler http.HandlerFunc, dst any, key string, methods ...string) {
	m.withBindCtx(dst, key, binder.WithBodyBinder(enum.Tags.Json), pattern, handler, methods...)
}

func (m *Mux) withBindCtx(dst any, key string, binderOptions []binder.MultiBinderOption, pattern string, handler http.HandlerFunc, methods ...string) {
	// Get the type of dst
	dstType := reflect.TypeOf(dst)
	
	m.Group(func(mux *flow.Mux) {
		mux.Use(middleware.Binder(dstType, key, binderOptions...))
		mux.HandleFunc(pattern, handler, methods...)
	})
}
