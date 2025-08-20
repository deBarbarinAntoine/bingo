// middlewares.go
package router

import (
	"net/http"
	"reflect"
	
	"BinGo/binder"
	
	"github.com/alexedwards/scs/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

type Middleware func(next http.Handler) http.Handler

func SessionMiddleware(sessionManager *scs.SessionManager) Middleware {
	return sessionManager.LoadAndSave
}

func LoggerMiddleware(logger zerolog.Logger) Middleware {
	logHandler := hlog.NewHandler(logger)
	remoteAddrHandler := hlog.RemoteAddrHandler("ip")
	userAgentHandler := hlog.UserAgentHandler("user_agent")
	refererHandler := hlog.RefererHandler("referer")
	requestIDHandler := hlog.RequestIDHandler("req_id", "Request-Id")
	methodHandler := hlog.MethodHandler("method")
	
	return func(next http.Handler) http.Handler {
		return logHandler(remoteAddrHandler(userAgentHandler(refererHandler(requestIDHandler(methodHandler(next))))))
	}
}

func BinderMiddleware(dstType reflect.Type, key string, binderOptions ...binder.MultiBinderOption) Middleware {
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
			r = SetCtxData(r, key, newDst)
			next.ServeHTTP(w, r)
		})
	}
}
